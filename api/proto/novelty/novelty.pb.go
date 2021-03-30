// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ug/api/proto/novelty/novelty.proto

package novelty

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Novelty struct {
	Click                *Click   `protobuf:"bytes,1,opt,name=click,proto3" json:"click,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Novelty) Reset()         { *m = Novelty{} }
func (m *Novelty) String() string { return proto.CompactTextString(m) }
func (*Novelty) ProtoMessage()    {}
func (*Novelty) Descriptor() ([]byte, []int) {
	return fileDescriptor_c58844a2ef16a5c0, []int{0}
}

func (m *Novelty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Novelty.Unmarshal(m, b)
}
func (m *Novelty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Novelty.Marshal(b, m, deterministic)
}
func (m *Novelty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Novelty.Merge(m, src)
}
func (m *Novelty) XXX_Size() int {
	return xxx_messageInfo_Novelty.Size(m)
}
func (m *Novelty) XXX_DiscardUnknown() {
	xxx_messageInfo_Novelty.DiscardUnknown(m)
}

var xxx_messageInfo_Novelty proto.InternalMessageInfo

func (m *Novelty) GetClick() *Click {
	if m != nil {
		return m.Click
	}
	return nil
}

type Click struct {
	DayClickCount        uint64   `protobuf:"varint,1,opt,name=day_click_count,json=dayClickCount,proto3" json:"day_click_count,omitempty"`
	LastClickDay         int64    `protobuf:"varint,2,opt,name=last_click_day,json=lastClickDay,proto3" json:"last_click_day,omitempty"`
	AccumulateClickCount uint64   `protobuf:"varint,3,opt,name=accumulate_click_count,json=accumulateClickCount,proto3" json:"accumulate_click_count,omitempty"`
	LastClickTimestamp   int64    `protobuf:"varint,4,opt,name=last_click_timestamp,json=lastClickTimestamp,proto3" json:"last_click_timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Click) Reset()         { *m = Click{} }
func (m *Click) String() string { return proto.CompactTextString(m) }
func (*Click) ProtoMessage()    {}
func (*Click) Descriptor() ([]byte, []int) {
	return fileDescriptor_c58844a2ef16a5c0, []int{1}
}

func (m *Click) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Click.Unmarshal(m, b)
}
func (m *Click) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Click.Marshal(b, m, deterministic)
}
func (m *Click) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Click.Merge(m, src)
}
func (m *Click) XXX_Size() int {
	return xxx_messageInfo_Click.Size(m)
}
func (m *Click) XXX_DiscardUnknown() {
	xxx_messageInfo_Click.DiscardUnknown(m)
}

var xxx_messageInfo_Click proto.InternalMessageInfo

func (m *Click) GetDayClickCount() uint64 {
	if m != nil {
		return m.DayClickCount
	}
	return 0
}

func (m *Click) GetLastClickDay() int64 {
	if m != nil {
		return m.LastClickDay
	}
	return 0
}

func (m *Click) GetAccumulateClickCount() uint64 {
	if m != nil {
		return m.AccumulateClickCount
	}
	return 0
}

func (m *Click) GetLastClickTimestamp() int64 {
	if m != nil {
		return m.LastClickTimestamp
	}
	return 0
}

func init() {
	proto.RegisterType((*Novelty)(nil), "novelty.Novelty")
	proto.RegisterType((*Click)(nil), "novelty.Click")
}

func init() { proto.RegisterFile("ug/api/proto/novelty/novelty.proto", fileDescriptor_c58844a2ef16a5c0) }

var fileDescriptor_c58844a2ef16a5c0 = []byte{
	// 248 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x89, 0x6d, 0x2d, 0x8c, 0x5a, 0x61, 0x29, 0x92, 0x83, 0x87, 0x12, 0x8a, 0xf6, 0xb4,
	0x2b, 0xea, 0x13, 0x58, 0xc1, 0x9b, 0x87, 0xe0, 0xc9, 0x4b, 0x19, 0x27, 0x4b, 0x08, 0x66, 0xb3,
	0xd1, 0xce, 0x8a, 0xfb, 0x70, 0xbe, 0x9b, 0x64, 0xd2, 0x68, 0x05, 0x4f, 0xbb, 0xf3, 0xfd, 0xff,
	0xfc, 0xc3, 0x0c, 0x64, 0xa1, 0x34, 0xd8, 0x56, 0xa6, 0x7d, 0xf7, 0xec, 0x4d, 0xe3, 0x3f, 0x6c,
	0xcd, 0x71, 0x78, 0xb5, 0x50, 0x35, 0xdd, 0x95, 0x99, 0x81, 0xe9, 0x63, 0xff, 0x55, 0x4b, 0x98,
	0x50, 0x5d, 0xd1, 0x6b, 0x9a, 0x2c, 0x92, 0xd5, 0xd1, 0xf5, 0x4c, 0x0f, 0x2d, 0xeb, 0x8e, 0xe6,
	0xbd, 0x98, 0x7d, 0x25, 0x30, 0x11, 0xa0, 0x2e, 0xe0, 0xb4, 0xc0, 0xb8, 0x11, 0xbc, 0x21, 0x1f,
	0x1a, 0x96, 0xce, 0x71, 0x7e, 0x52, 0x60, 0x14, 0xcb, 0xba, 0x83, 0x6a, 0x09, 0xb3, 0x1a, 0xb7,
	0xbc, 0x33, 0x16, 0x18, 0xd3, 0x83, 0x45, 0xb2, 0x1a, 0xe5, 0xc7, 0x1d, 0x15, 0xdf, 0x3d, 0x46,
	0x75, 0x0b, 0x67, 0x48, 0x14, 0x5c, 0xa8, 0x91, 0xed, 0x9f, 0xd0, 0x91, 0x84, 0xce, 0x7f, 0xd5,
	0xbd, 0xec, 0x2b, 0x98, 0xef, 0x65, 0x73, 0xe5, 0xec, 0x96, 0xd1, 0xb5, 0xe9, 0x58, 0x26, 0xa8,
	0x9f, 0x09, 0x4f, 0x83, 0x72, 0xf7, 0x00, 0xe7, 0xe4, 0x9d, 0x66, 0xdb, 0x90, 0x6d, 0x58, 0x87,
	0xb2, 0xbf, 0xc8, 0xb0, 0xec, 0xf3, 0x65, 0x59, 0xb1, 0x26, 0x5f, 0x58, 0xed, 0x51, 0x93, 0x77,
	0xe6, 0xf3, 0xcd, 0xfc, 0x77, 0xd0, 0x97, 0x43, 0x29, 0x6f, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff,
	0xdd, 0xa4, 0x04, 0x1f, 0x6f, 0x01, 0x00, 0x00,
}
