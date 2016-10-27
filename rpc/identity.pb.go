// Code generated by protoc-gen-go.
// source: identity.proto
// DO NOT EDIT!

package ricochet

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Identity struct {
	Address string `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
}

func (m *Identity) Reset()                    { *m = Identity{} }
func (m *Identity) String() string            { return proto.CompactTextString(m) }
func (*Identity) ProtoMessage()               {}
func (*Identity) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

type IdentityRequest struct {
}

func (m *IdentityRequest) Reset()                    { *m = IdentityRequest{} }
func (m *IdentityRequest) String() string            { return proto.CompactTextString(m) }
func (*IdentityRequest) ProtoMessage()               {}
func (*IdentityRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func init() {
	proto.RegisterType((*Identity)(nil), "ricochet.Identity")
	proto.RegisterType((*IdentityRequest)(nil), "ricochet.IdentityRequest")
}

func init() { proto.RegisterFile("identity.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 95 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcb, 0x4c, 0x49, 0xcd,
	0x2b, 0xc9, 0x2c, 0xa9, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x28, 0xca, 0x4c, 0xce,
	0x4f, 0xce, 0x48, 0x2d, 0x51, 0x52, 0xe1, 0xe2, 0xf0, 0x84, 0xca, 0x09, 0x49, 0x70, 0xb1, 0x27,
	0xa6, 0xa4, 0x14, 0xa5, 0x16, 0x17, 0x4b, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xb8, 0x4a,
	0x82, 0x5c, 0xfc, 0x30, 0x55, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x49, 0x6c, 0x60, 0x93,
	0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xcc, 0xc7, 0x44, 0x6d, 0x5b, 0x00, 0x00, 0x00,
}
