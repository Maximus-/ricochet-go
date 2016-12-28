package core

import (
	"errors"
	"fmt"
	"github.com/ricochet-im/ricochet-go/core/utils"
	"github.com/ricochet-im/ricochet-go/rpc"
	"strconv"
	"sync"
	"time"
)

type ContactList struct {
	core *Ricochet

	mutex  sync.RWMutex
	events *utils.Publisher

	contacts map[int]*Contact
}

func LoadContactList(core *Ricochet) (*ContactList, error) {
	list := &ContactList{
		core:   core,
		events: utils.CreatePublisher(),
	}

	config := core.Config.OpenRead()
	defer config.Close()

	list.contacts = make(map[int]*Contact, len(config.Contacts))
	for idStr, data := range config.Contacts {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("Invalid contact id '%s'", idStr)
		}
		if _, exists := list.contacts[id]; exists {
			return nil, fmt.Errorf("Duplicate contact id '%d'", id)
		}

		contact, err := ContactFromConfig(core, id, data, list.events)
		if err != nil {
			return nil, err
		}
		list.contacts[id] = contact
	}

	// XXX Requests aren't implemented
	return list, nil
}

func (this *ContactList) EventMonitor() utils.Subscribable {
	return this.events
}

func (this *ContactList) Contacts() []*Contact {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	re := make([]*Contact, 0, len(this.contacts))
	for _, contact := range this.contacts {
		re = append(re, contact)
	}
	return re
}

func (this *ContactList) ContactById(id int) *Contact {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.contacts[id]
}

func (this *ContactList) ContactByAddress(address string) *Contact {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	for _, contact := range this.contacts {
		if contact.Address() == address {
			return contact
		}
	}
	return nil
}

func (this *ContactList) AddContactRequest(address, name, fromName, text string) (*Contact, error) {
	if !IsAddressValid(address) {
		return nil, errors.New("Invalid ricochet address")
	}
	if len(fromName) > 0 && !IsNicknameAcceptable(fromName) {
		return nil, errors.New("Invalid nickname")
	}
	if len(text) > 0 && !IsMessageAcceptable(text) {
		return nil, errors.New("Invalid message")
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	for _, contact := range this.contacts {
		if contact.Address() == address {
			return nil, errors.New("Contact already exists with this address")
		}
		if contact.Nickname() == name {
			return nil, errors.New("Contact already exists with this nickname")
		}
	}

	// XXX check inbound requests

	// Write new contact into config
	config := this.core.Config.OpenWrite()

	maxContactId := 0
	for idstr, _ := range config.Contacts {
		if id, err := strconv.Atoi(idstr); err == nil {
			if maxContactId < id {
				maxContactId = id
			}
		}
	}

	contactId := maxContactId + 1
	onion, _ := OnionFromAddress(address)
	configContact := ConfigContact{
		Hostname:    onion,
		Nickname:    name,
		WhenCreated: time.Now().Format(time.RFC3339),
		Request: ConfigContactRequest{
			Pending:    true,
			MyNickname: fromName,
			Message:    text,
		},
	}

	config.Contacts[strconv.Itoa(contactId)] = configContact
	if err := config.Save(); err != nil {
		return nil, err
	}

	// Create Contact
	contact, err := ContactFromConfig(this.core, contactId, configContact, this.events)
	if err != nil {
		return nil, err
	}
	this.contacts[contactId] = contact

	event := ricochet.ContactEvent{
		Type: ricochet.ContactEvent_ADD,
		Subject: &ricochet.ContactEvent_Contact{
			Contact: contact.Data(),
		},
	}
	this.events.Publish(event)

	contact.StartConnection()
	return contact, nil
}

func (this *ContactList) RemoveContact(contact *Contact) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.contacts[contact.Id()] != contact {
		return errors.New("Not in contact list")
	}

	contact.StopConnection()

	config := this.core.Config.OpenWrite()
	delete(config.Contacts, strconv.Itoa(contact.Id()))
	if err := config.Save(); err != nil {
		return err
	}

	delete(this.contacts, contact.Id())

	event := ricochet.ContactEvent{
		Type: ricochet.ContactEvent_DELETE,
		Subject: &ricochet.ContactEvent_Contact{
			Contact: &ricochet.Contact{
				Id:      int32(contact.Id()),
				Address: contact.Address(),
			},
		},
	}
	this.events.Publish(event)

	return nil
}

func (this *ContactList) StartConnections() {
	for _, contact := range this.Contacts() {
		contact.StartConnection()
	}
}

func (this *ContactList) StopConnections() {
	for _, contact := range this.Contacts() {
		contact.StopConnection()
	}
}
