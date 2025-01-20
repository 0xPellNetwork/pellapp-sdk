// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/pelldvs/pellapp-sdk/pelldvs/avsi.proto

package types

import (
	fmt "fmt"
	types "github.com/0xPellNetwork/pelldvs/avsi/types"
	_ "github.com/cosmos/cosmos-proto"
	types1 "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
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

// StringEvent defines an Event object wrapper where all the attributes
// contain key/value pairs that are strings instead of raw bytes.
type StringEvent struct {
	EventType  string      `protobuf:"bytes,1,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
	Attributes []Attribute `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes"`
}

func (m *StringEvent) Reset()      { *m = StringEvent{} }
func (*StringEvent) ProtoMessage() {}
func (*StringEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba553b13051d5021, []int{0}
}
func (m *StringEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StringEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StringEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StringEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StringEvent.Merge(m, src)
}
func (m *StringEvent) XXX_Size() int {
	return m.Size()
}
func (m *StringEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_StringEvent.DiscardUnknown(m)
}

var xxx_messageInfo_StringEvent proto.InternalMessageInfo

func (m *StringEvent) GetEventType() string {
	if m != nil {
		return m.EventType
	}
	return ""
}

func (m *StringEvent) GetAttributes() []Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

// Attribute defines an attribute wrapper where the key and value are
// strings instead of raw bytes.
type Attribute struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *Attribute) Reset()      { *m = Attribute{} }
func (*Attribute) ProtoMessage() {}
func (*Attribute) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba553b13051d5021, []int{1}
}
func (m *Attribute) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Attribute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Attribute.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Attribute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Attribute.Merge(m, src)
}
func (m *Attribute) XXX_Size() int {
	return m.Size()
}
func (m *Attribute) XXX_DiscardUnknown() {
	xxx_messageInfo_Attribute.DiscardUnknown(m)
}

var xxx_messageInfo_Attribute proto.InternalMessageInfo

func (m *Attribute) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Attribute) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// TODO buf import pelldvs.avsi.Event
// Result is the union of ResponseFormat and ResponseCheckTx.
type Result struct {
	// Data is any data returned from message or handler execution. It MUST be
	// length prefixed in order to separate data from multiple message executions.
	// Deprecated. This field is still populated, but prefer msg_response instead
	// because it also contains the Msg response typeURL.
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"` // Deprecated: Do not use.
	// Log contains the log information from message or handler execution.
	Log string `protobuf:"bytes,2,opt,name=log,proto3" json:"log,omitempty"`
	// Events contains a slice of Event objects that were emitted during message
	// or handler execution.
	Events []types.Event `protobuf:"bytes,3,rep,name=events,proto3" json:"events"`
	// msg_responses contains the Msg handler responses type packed in Anys.
	MsgResponses []*types1.Any `protobuf:"bytes,4,rep,name=msg_responses,json=msgResponses,proto3" json:"msg_responses,omitempty"`
}

func (m *Result) Reset()      { *m = Result{} }
func (*Result) ProtoMessage() {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba553b13051d5021, []int{2}
}
func (m *Result) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Result.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(m, src)
}
func (m *Result) XXX_Size() int {
	return m.Size()
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

func init() {
	proto.RegisterType((*StringEvent)(nil), "intellix.pelldvs.StringEvent")
	proto.RegisterType((*Attribute)(nil), "intellix.pelldvs.Attribute")
	proto.RegisterType((*Result)(nil), "intellix.pelldvs.Result")
}

func init() { proto.RegisterFile("github.com/pelldvs/pellapp-sdk/pelldvs/avsi.proto", fileDescriptor_ba553b13051d5021) }

var fileDescriptor_ba553b13051d5021 = []byte{
	// 396 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x52, 0xb1, 0x6e, 0xe2, 0x40,
	0x10, 0xf5, 0x82, 0x0f, 0x89, 0x85, 0xd3, 0xa1, 0x05, 0x9d, 0x7c, 0xa0, 0x33, 0x16, 0x15, 0x0d,
	0xeb, 0xbb, 0xe3, 0x74, 0x3a, 0xd1, 0x61, 0x29, 0x6d, 0x0a, 0x27, 0x55, 0x1a, 0x64, 0xc2, 0x66,
	0x65, 0xb1, 0x78, 0x2d, 0xef, 0xda, 0x8a, 0xbb, 0x94, 0x29, 0xf3, 0x09, 0xf9, 0x88, 0x74, 0xf9,
	0x01, 0x4a, 0x4a, 0x14, 0x45, 0x51, 0x02, 0x3f, 0x12, 0x79, 0x6d, 0x13, 0x94, 0xee, 0xcd, 0xcc,
	0x7b, 0x9e, 0xe7, 0xb7, 0x03, 0x7b, 0x7e, 0x20, 0x09, 0x63, 0xfe, 0xb5, 0x1d, 0x12, 0xc6, 0x16,
	0x89, 0xb0, 0xbd, 0x44, 0xf8, 0x38, 0x8c, 0xb8, 0xe4, 0xa8, 0x55, 0x0e, 0x71, 0x31, 0xec, 0x1a,
	0xc7, 0x2c, 0x5b, 0xa6, 0x21, 0x11, 0x39, 0xb7, 0xdb, 0xa1, 0x9c, 0x72, 0x05, 0xed, 0x0c, 0x15,
	0xdd, 0x1f, 0x94, 0x73, 0xca, 0x88, 0xad, 0xaa, 0x79, 0x7c, 0x65, 0x7b, 0x41, 0x5a, 0x8e, 0x2e,
	0xb9, 0x58, 0x71, 0x31, 0xcb, 0x35, 0x79, 0x91, 0x8f, 0x06, 0x31, 0x6c, 0x9c, 0xc9, 0xc8, 0x0f,
	0xe8, 0x49, 0x42, 0x02, 0x89, 0x7e, 0x42, 0x48, 0x32, 0x30, 0xcb, 0xf6, 0x19, 0xc0, 0x02, 0xc3,
	0xba, 0x5b, 0x57, 0x9d, 0xf3, 0x34, 0x24, 0x68, 0x0a, 0xa1, 0x27, 0x65, 0xe4, 0xcf, 0x63, 0x49,
	0x84, 0x51, 0xb1, 0xaa, 0xc3, 0xc6, 0x9f, 0x1e, 0xfe, 0x6c, 0x1d, 0x4f, 0x4b, 0x8e, 0xa3, 0xaf,
	0x5f, 0xfa, 0x9a, 0x7b, 0x24, 0x9a, 0xe8, 0x37, 0xcf, 0x16, 0x18, 0x8c, 0x61, 0xfd, 0x40, 0x42,
	0x2d, 0x58, 0x5d, 0x92, 0xb4, 0xd8, 0x96, 0x41, 0xd4, 0x81, 0x5f, 0x12, 0x8f, 0xc5, 0xc4, 0xa8,
	0xa8, 0x5e, 0x5e, 0x0c, 0x1e, 0x01, 0xac, 0xb9, 0x44, 0xc4, 0x4c, 0xa2, 0xef, 0x50, 0x5f, 0x78,
	0xd2, 0x53, 0x9a, 0xa6, 0x53, 0x31, 0x80, 0xab, 0xea, 0xec, 0x53, 0x8c, 0xd3, 0x42, 0x96, 0x41,
	0xf4, 0x1b, 0xd6, 0x94, 0x7f, 0x61, 0x54, 0x95, 0xdd, 0xf6, 0xc1, 0xa5, 0x4a, 0x5f, 0xfd, 0x76,
	0x61, 0xb3, 0x20, 0xa2, 0x53, 0xf8, 0x75, 0x25, 0xe8, 0x2c, 0x22, 0x22, 0xe4, 0x81, 0x20, 0xc2,
	0xd0, 0x95, 0xb2, 0x83, 0xf3, 0x84, 0x71, 0x99, 0x30, 0x9e, 0x06, 0xa9, 0xd3, 0x7e, 0x7a, 0x18,
	0x7d, 0xcb, 0x23, 0x1d, 0x89, 0xc5, 0xd2, 0xfa, 0x85, 0xff, 0xfe, 0x73, 0x9b, 0x2b, 0x41, 0xdd,
	0x52, 0x3e, 0xd1, 0x6f, 0xef, 0xfb, 0x9a, 0xf3, 0x7f, 0xfb, 0x66, 0x6a, 0xeb, 0x9d, 0x09, 0x36,
	0x3b, 0x13, 0xbc, 0xee, 0x4c, 0x70, 0xb7, 0x37, 0xb5, 0xcd, 0xde, 0xd4, 0xb6, 0x7b, 0x53, 0xbb,
	0xe8, 0x7e, 0x1c, 0xc7, 0x92, 0x1e, 0x0e, 0x44, 0xbd, 0xfa, 0xbc, 0xa6, 0x16, 0x8e, 0xdf, 0x03,
	0x00, 0x00, 0xff, 0xff, 0x9b, 0x7b, 0xd8, 0xef, 0x41, 0x02, 0x00, 0x00,
}

func (m *StringEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StringEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StringEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Attributes) > 0 {
		for iNdEx := len(m.Attributes) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Attributes[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintAvsi(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.EventType) > 0 {
		i -= len(m.EventType)
		copy(dAtA[i:], m.EventType)
		i = encodeVarintAvsi(dAtA, i, uint64(len(m.EventType)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Attribute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Attribute) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Attribute) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintAvsi(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintAvsi(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Result) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Result) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Result) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.MsgResponses) > 0 {
		for iNdEx := len(m.MsgResponses) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MsgResponses[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintAvsi(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Events) > 0 {
		for iNdEx := len(m.Events) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Events[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintAvsi(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Log) > 0 {
		i -= len(m.Log)
		copy(dAtA[i:], m.Log)
		i = encodeVarintAvsi(dAtA, i, uint64(len(m.Log)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintAvsi(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAvsi(dAtA []byte, offset int, v uint64) int {
	offset -= sovAvsi(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *StringEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.EventType)
	if l > 0 {
		n += 1 + l + sovAvsi(uint64(l))
	}
	if len(m.Attributes) > 0 {
		for _, e := range m.Attributes {
			l = e.Size()
			n += 1 + l + sovAvsi(uint64(l))
		}
	}
	return n
}

func (m *Attribute) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovAvsi(uint64(l))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovAvsi(uint64(l))
	}
	return n
}

func (m *Result) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovAvsi(uint64(l))
	}
	l = len(m.Log)
	if l > 0 {
		n += 1 + l + sovAvsi(uint64(l))
	}
	if len(m.Events) > 0 {
		for _, e := range m.Events {
			l = e.Size()
			n += 1 + l + sovAvsi(uint64(l))
		}
	}
	if len(m.MsgResponses) > 0 {
		for _, e := range m.MsgResponses {
			l = e.Size()
			n += 1 + l + sovAvsi(uint64(l))
		}
	}
	return n
}

func sovAvsi(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAvsi(x uint64) (n int) {
	return sovAvsi(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *StringEvent) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForAttributes := "[]Attribute{"
	for _, f := range this.Attributes {
		repeatedStringForAttributes += fmt.Sprintf("%v", f) + ","
	}
	repeatedStringForAttributes += "}"
	s := strings.Join([]string{`&StringEvent{`,
		`EventType:` + fmt.Sprintf("%v", this.EventType) + `,`,
		`Attributes:` + repeatedStringForAttributes + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringAvsi(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *StringEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAvsi
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
			return fmt.Errorf("proto: StringEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StringEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EventType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAvsi
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
				return ErrInvalidLengthAvsi
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAvsi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EventType = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Attributes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAvsi
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
				return ErrInvalidLengthAvsi
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAvsi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Attributes = append(m.Attributes, Attribute{})
			if err := m.Attributes[len(m.Attributes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAvsi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAvsi
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
func (m *Attribute) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAvsi
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
			return fmt.Errorf("proto: Attribute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Attribute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAvsi
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
				return ErrInvalidLengthAvsi
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAvsi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAvsi
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
				return ErrInvalidLengthAvsi
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAvsi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAvsi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAvsi
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
func (m *Result) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAvsi
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
			return fmt.Errorf("proto: Result: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Result: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAvsi
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
				return ErrInvalidLengthAvsi
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthAvsi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Log", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAvsi
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
				return ErrInvalidLengthAvsi
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAvsi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Log = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Events", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAvsi
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
				return ErrInvalidLengthAvsi
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAvsi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Events = append(m.Events, types.Event{})
			if err := m.Events[len(m.Events)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MsgResponses", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAvsi
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
				return ErrInvalidLengthAvsi
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAvsi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MsgResponses = append(m.MsgResponses, &types1.Any{})
			if err := m.MsgResponses[len(m.MsgResponses)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAvsi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAvsi
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
func skipAvsi(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAvsi
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
					return 0, ErrIntOverflowAvsi
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
					return 0, ErrIntOverflowAvsi
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
				return 0, ErrInvalidLengthAvsi
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAvsi
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAvsi
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAvsi        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAvsi          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAvsi = fmt.Errorf("proto: unexpected end of group")
)

func (m *Result) String() string {
	if m == nil {
		return "nil"
	}
	return fmt.Sprintf("&Result{Data: %v, Log: %s, Events: %v, MsgResponses: %v}", m.Data, m.Log, m.Events, m.MsgResponses)
}
