/*
Copyright 2024. The NeverDown Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/kzz45/neverdown/pkg/jingx/proto/generated.proto

package proto

import (
	fmt "fmt"

	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"

	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func (m *Context) Reset()      { *m = Context{} }
func (*Context) ProtoMessage() {}
func (*Context) Descriptor() ([]byte, []int) {
	return fileDescriptor_bad5b002d2a47b83, []int{0}
}
func (m *Context) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Context) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *Context) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Context.Merge(m, src)
}
func (m *Context) XXX_Size() int {
	return m.Size()
}
func (m *Context) XXX_DiscardUnknown() {
	xxx_messageInfo_Context.DiscardUnknown(m)
}

var xxx_messageInfo_Context proto.InternalMessageInfo

func (m *Request) Reset()      { *m = Request{} }
func (*Request) ProtoMessage() {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_bad5b002d2a47b83, []int{1}
}
func (m *Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return m.Size()
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Response) Reset()      { *m = Response{} }
func (*Response) ProtoMessage() {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_bad5b002d2a47b83, []int{2}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return m.Size()
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *WatchEvent) Reset()      { *m = WatchEvent{} }
func (*WatchEvent) ProtoMessage() {}
func (*WatchEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_bad5b002d2a47b83, []int{3}
}
func (m *WatchEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WatchEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *WatchEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WatchEvent.Merge(m, src)
}
func (m *WatchEvent) XXX_Size() int {
	return m.Size()
}
func (m *WatchEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_WatchEvent.DiscardUnknown(m)
}

var xxx_messageInfo_WatchEvent proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Context)(nil), "github.com.kzz45.neverdown.pkg.jingx.proto.Context")
	proto.RegisterType((*Request)(nil), "github.com.kzz45.neverdown.pkg.jingx.proto.Request")
	proto.RegisterType((*Response)(nil), "github.com.kzz45.neverdown.pkg.jingx.proto.Response")
	proto.RegisterType((*WatchEvent)(nil), "github.com.kzz45.neverdown.pkg.jingx.proto.WatchEvent")
}

func init() {
	proto.RegisterFile("github.com/kzz45/neverdown/pkg/jingx/proto/generated.proto", fileDescriptor_bad5b002d2a47b83)
}

var fileDescriptor_bad5b002d2a47b83 = []byte{
	// 532 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x94, 0x4d, 0x6b, 0x53, 0x4f,
	0x14, 0xc6, 0x33, 0x79, 0x21, 0xc9, 0xa4, 0x7f, 0xfe, 0xf1, 0xea, 0x22, 0xba, 0xb8, 0x09, 0x71,
	0x13, 0x8b, 0xce, 0xa5, 0x45, 0x41, 0x14, 0x17, 0x9d, 0x50, 0xa4, 0x08, 0x22, 0x43, 0xa9, 0xe0,
	0xee, 0xe6, 0xde, 0xe3, 0xed, 0x35, 0xcd, 0xcc, 0x38, 0x77, 0xf2, 0xd2, 0xae, 0x04, 0x71, 0xef,
	0xc7, 0xf0, 0xa3, 0x64, 0xd9, 0x65, 0x57, 0xc1, 0x5c, 0xbf, 0x84, 0xb8, 0x92, 0x4c, 0x26, 0x2f,
	0xa4, 0x85, 0xb6, 0x2b, 0x57, 0x61, 0xce, 0x79, 0x9e, 0xe7, 0x9c, 0xf3, 0x83, 0x5c, 0xfc, 0x22,
	0x8a, 0xf5, 0x71, 0xbf, 0x43, 0x02, 0xd1, 0xf3, 0xba, 0x67, 0x67, 0x4f, 0x9f, 0x79, 0x1c, 0x06,
	0xa0, 0x42, 0x31, 0xe4, 0x9e, 0xec, 0x46, 0xde, 0xa7, 0x98, 0x47, 0x23, 0x4f, 0x2a, 0xa1, 0x85,
	0x17, 0x01, 0x07, 0xe5, 0x6b, 0x08, 0x89, 0x79, 0x3b, 0xdb, 0x2b, 0x2f, 0x31, 0x5e, 0xb2, 0xf4,
	0x12, 0xd9, 0x8d, 0x88, 0xf1, 0xce, 0xb5, 0x0f, 0x9e, 0xac, 0xcd, 0x89, 0x44, 0x24, 0xe6, 0x91,
	0x9d, 0xfe, 0x47, 0xf3, 0xb2, 0xf9, 0x22, 0x12, 0x56, 0xfe, 0xf2, 0x9a, 0xb5, 0x7c, 0x19, 0x27,
	0x9e, 0xea, 0xf8, 0x81, 0x37, 0xd8, 0xd9, 0xdc, 0xab, 0xf9, 0x1b, 0xe1, 0x62, 0x5b, 0x70, 0x0d,
	0x23, 0xed, 0x3c, 0xc4, 0x05, 0x2d, 0xba, 0xc0, 0x6b, 0xa8, 0x81, 0x5a, 0x65, 0xfa, 0xdf, 0x78,
	0x52, 0xcf, 0xa4, 0x93, 0x7a, 0xe1, 0x70, 0x56, 0x64, 0xf3, 0x9e, 0xf3, 0x08, 0x17, 0xe3, 0x64,
	0x2f, 0xec, 0xc5, 0xbc, 0x96, 0x6d, 0xa0, 0x56, 0x89, 0xfe, 0x6f, 0x65, 0xc5, 0x83, 0x79, 0x99,
	0x2d, 0xfa, 0xce, 0x63, 0x5c, 0x82, 0x91, 0x8c, 0x15, 0xec, 0xe9, 0x5a, 0xae, 0x81, 0x5a, 0x05,
	0x5a, 0xb5, 0xda, 0xd2, 0xbe, 0xad, 0xb3, 0xa5, 0xc2, 0x11, 0xb8, 0x12, 0x9c, 0xf4, 0x13, 0x0d,
	0x8a, 0x89, 0x13, 0xa8, 0xe5, 0x1b, 0xa8, 0x55, 0xd9, 0x7d, 0x4e, 0xae, 0xe1, 0x36, 0x3b, 0x8e,
	0xcc, 0x8e, 0x23, 0x83, 0x1d, 0xd2, 0x5e, 0xf9, 0xe9, 0x5d, 0x3b, 0xaa, 0xb2, 0x56, 0x64, 0xeb,
	0x13, 0x9a, 0xdf, 0xb2, 0xb8, 0xc8, 0xe0, 0x73, 0x1f, 0x12, 0xed, 0x7c, 0x45, 0xb8, 0x1a, 0x29,
	0xd1, 0x97, 0x47, 0xa0, 0x92, 0x58, 0xf0, 0x37, 0x31, 0x0f, 0x0d, 0x86, 0xca, 0xee, 0xab, 0x5b,
	0xad, 0xf0, 0x7a, 0x23, 0x84, 0xd6, 0xec, 0x1e, 0xd5, 0xcd, 0x0e, 0xbb, 0x34, 0xd0, 0xf1, 0x70,
	0x99, 0xfb, 0x3d, 0x48, 0xa4, 0x1f, 0x80, 0xa1, 0x5b, 0xa6, 0x77, 0xac, 0xbd, 0xfc, 0x76, 0xd1,
	0x60, 0x2b, 0x8d, 0xd3, 0xc2, 0xf9, 0x01, 0xa8, 0x8e, 0xa1, 0x5b, 0xa6, 0xf7, 0xac, 0x36, 0x7f,
	0x04, 0xaa, 0xf3, 0xc7, 0xfe, 0x32, 0xa3, 0x70, 0xee, 0xe3, 0x9c, 0xf2, 0x87, 0x86, 0xea, 0x16,
	0x2d, 0xa6, 0x93, 0x7a, 0x8e, 0xf9, 0x43, 0x36, 0xab, 0x35, 0x7f, 0x64, 0x71, 0x89, 0x41, 0x22,
	0x05, 0x4f, 0xc0, 0x69, 0xe0, 0x7c, 0x20, 0x42, 0x30, 0xb7, 0x17, 0xe8, 0xd6, 0x22, 0xb1, 0x2d,
	0x42, 0x60, 0xa6, 0x73, 0x35, 0xaa, 0xec, 0x3f, 0x45, 0x95, 0xbb, 0x05, 0xaa, 0xfc, 0x4d, 0x51,
	0x15, 0xae, 0x40, 0x75, 0x80, 0xf1, 0x7b, 0x5f, 0x07, 0xc7, 0xfb, 0x03, 0xe0, 0x7a, 0xc6, 0x4a,
	0x9f, 0x4a, 0xb0, 0x7f, 0x97, 0x25, 0xab, 0xc3, 0x53, 0x09, 0xcc, 0x74, 0x16, 0x51, 0xd9, 0xcb,
	0x51, 0xf4, 0xdd, 0x78, 0xea, 0x66, 0xce, 0xa7, 0x6e, 0xe6, 0x62, 0xea, 0x66, 0xbe, 0xa4, 0x2e,
	0x1a, 0xa7, 0x2e, 0x3a, 0x4f, 0x5d, 0x74, 0x91, 0xba, 0xe8, 0x67, 0xea, 0xa2, 0xef, 0xbf, 0xdc,
	0xcc, 0x87, 0xed, 0x9b, 0x7f, 0x72, 0xfe, 0x06, 0x00, 0x00, 0xff, 0xff, 0xb7, 0x47, 0x48, 0xff,
	0x9f, 0x04, 0x00, 0x00,
}

func (m *Context) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Context) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Context) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.ClusterRole.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	i = encodeVarintGenerated(dAtA, i, uint64(m.ExpireAt))
	i--
	dAtA[i] = 0x18
	i--
	if m.IsAdmin {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i--
	dAtA[i] = 0x10
	i -= len(m.Token)
	copy(dAtA[i:], m.Token)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Token)))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Raw != nil {
		i -= len(m.Raw)
		copy(dAtA[i:], m.Raw)
		i = encodeVarintGenerated(dAtA, i, uint64(len(m.Raw)))
		i--
		dAtA[i] = 0x22
	}
	i -= len(m.Verb)
	copy(dAtA[i:], m.Verb)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Verb)))
	i--
	dAtA[i] = 0x1a
	i -= len(m.Namespace)
	copy(dAtA[i:], m.Namespace)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Namespace)))
	i--
	dAtA[i] = 0x12
	{
		size, err := m.GroupVersionKind.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Response) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Response) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Response) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Raw != nil {
		i -= len(m.Raw)
		copy(dAtA[i:], m.Raw)
		i = encodeVarintGenerated(dAtA, i, uint64(len(m.Raw)))
		i--
		dAtA[i] = 0x2a
	}
	i -= len(m.Verb)
	copy(dAtA[i:], m.Verb)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Verb)))
	i--
	dAtA[i] = 0x22
	i -= len(m.Namespace)
	copy(dAtA[i:], m.Namespace)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Namespace)))
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.GroupVersionKind.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	i = encodeVarintGenerated(dAtA, i, uint64(m.Code))
	i--
	dAtA[i] = 0x8
	return len(dAtA) - i, nil
}

func (m *WatchEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WatchEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WatchEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Raw != nil {
		i -= len(m.Raw)
		copy(dAtA[i:], m.Raw)
		i = encodeVarintGenerated(dAtA, i, uint64(len(m.Raw)))
		i--
		dAtA[i] = 0x12
	}
	i -= len(m.Type)
	copy(dAtA[i:], m.Type)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Type)))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenerated(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenerated(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Context) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Token)
	n += 1 + l + sovGenerated(uint64(l))
	n += 2
	n += 1 + sovGenerated(uint64(m.ExpireAt))
	l = m.ClusterRole.Size()
	n += 1 + l + sovGenerated(uint64(l))
	return n
}

func (m *Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.GroupVersionKind.Size()
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.Namespace)
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.Verb)
	n += 1 + l + sovGenerated(uint64(l))
	if m.Raw != nil {
		l = len(m.Raw)
		n += 1 + l + sovGenerated(uint64(l))
	}
	return n
}

func (m *Response) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovGenerated(uint64(m.Code))
	l = m.GroupVersionKind.Size()
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.Namespace)
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.Verb)
	n += 1 + l + sovGenerated(uint64(l))
	if m.Raw != nil {
		l = len(m.Raw)
		n += 1 + l + sovGenerated(uint64(l))
	}
	return n
}

func (m *WatchEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Type)
	n += 1 + l + sovGenerated(uint64(l))
	if m.Raw != nil {
		l = len(m.Raw)
		n += 1 + l + sovGenerated(uint64(l))
	}
	return n
}

func sovGenerated(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenerated(x uint64) (n int) {
	return sovGenerated(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Context) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Context{`,
		`Token:` + fmt.Sprintf("%v", this.Token) + `,`,
		`IsAdmin:` + fmt.Sprintf("%v", this.IsAdmin) + `,`,
		`ExpireAt:` + fmt.Sprintf("%v", this.ExpireAt) + `,`,
		`ClusterRole:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.ClusterRole), "ClusterRole", "v1.ClusterRole", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Request) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Request{`,
		`GroupVersionKind:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.GroupVersionKind), "GroupVersionKind", "v1.GroupVersionKind", 1), `&`, ``, 1) + `,`,
		`Namespace:` + fmt.Sprintf("%v", this.Namespace) + `,`,
		`Verb:` + fmt.Sprintf("%v", this.Verb) + `,`,
		`Raw:` + valueToStringGenerated(this.Raw) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Response) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Response{`,
		`Code:` + fmt.Sprintf("%v", this.Code) + `,`,
		`GroupVersionKind:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.GroupVersionKind), "GroupVersionKind", "v1.GroupVersionKind", 1), `&`, ``, 1) + `,`,
		`Namespace:` + fmt.Sprintf("%v", this.Namespace) + `,`,
		`Verb:` + fmt.Sprintf("%v", this.Verb) + `,`,
		`Raw:` + valueToStringGenerated(this.Raw) + `,`,
		`}`,
	}, "")
	return s
}
func (this *WatchEvent) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&WatchEvent{`,
		`Type:` + fmt.Sprintf("%v", this.Type) + `,`,
		`Raw:` + valueToStringGenerated(this.Raw) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGenerated(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Context) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Context: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Context: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsAdmin", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsAdmin = bool(v != 0)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpireAt", wireType)
			}
			m.ExpireAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExpireAt |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterRole", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ClusterRole.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Request) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Request: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Request: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupVersionKind", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GroupVersionKind.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Namespace", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Namespace = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Verb", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Verb = Verb(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Raw", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Raw = append(m.Raw[:0], dAtA[iNdEx:postIndex]...)
			if m.Raw == nil {
				m.Raw = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Response) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Response: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Response: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupVersionKind", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GroupVersionKind.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Namespace", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Namespace = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Verb", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Verb = Verb(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Raw", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Raw = append(m.Raw[:0], dAtA[iNdEx:postIndex]...)
			if m.Raw == nil {
				m.Raw = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *WatchEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: WatchEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WatchEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Raw", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Raw = append(m.Raw[:0], dAtA[iNdEx:postIndex]...)
			if m.Raw == nil {
				m.Raw = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenerated(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenerated
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenerated
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenerated
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenerated        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenerated          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenerated = fmt.Errorf("proto: unexpected end of group")
)
