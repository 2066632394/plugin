// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hacker.proto

package types

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type HackerAction struct {
	// Types that are valid to be assigned to Value:
	//	*HackerAction_AddBill
	Value                isHackerAction_Value `protobuf_oneof:"value"`
	Ty                   int32                `protobuf:"varint,50,opt,name=Ty,proto3" json:"Ty,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *HackerAction) Reset()         { *m = HackerAction{} }
func (m *HackerAction) String() string { return proto.CompactTextString(m) }
func (*HackerAction) ProtoMessage()    {}
func (*HackerAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_3aa5634b99676f1c, []int{0}
}

func (m *HackerAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HackerAction.Unmarshal(m, b)
}
func (m *HackerAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HackerAction.Marshal(b, m, deterministic)
}
func (m *HackerAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HackerAction.Merge(m, src)
}
func (m *HackerAction) XXX_Size() int {
	return xxx_messageInfo_HackerAction.Size(m)
}
func (m *HackerAction) XXX_DiscardUnknown() {
	xxx_messageInfo_HackerAction.DiscardUnknown(m)
}

var xxx_messageInfo_HackerAction proto.InternalMessageInfo

type isHackerAction_Value interface {
	isHackerAction_Value()
}

type HackerAction_AddBill struct {
	AddBill *HackerAddBill `protobuf:"bytes,1,opt,name=addBill,proto3,oneof"`
}

func (*HackerAction_AddBill) isHackerAction_Value() {}

func (m *HackerAction) GetValue() isHackerAction_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *HackerAction) GetAddBill() *HackerAddBill {
	if x, ok := m.GetValue().(*HackerAction_AddBill); ok {
		return x.AddBill
	}
	return nil
}

func (m *HackerAction) GetTy() int32 {
	if m != nil {
		return m.Ty
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*HackerAction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _HackerAction_OneofMarshaler, _HackerAction_OneofUnmarshaler, _HackerAction_OneofSizer, []interface{}{
		(*HackerAction_AddBill)(nil),
	}
}

func _HackerAction_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*HackerAction)
	// value
	switch x := m.Value.(type) {
	case *HackerAction_AddBill:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.AddBill); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("HackerAction.Value has unexpected type %T", x)
	}
	return nil
}

func _HackerAction_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*HackerAction)
	switch tag {
	case 1: // value.addBill
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(HackerAddBill)
		err := b.DecodeMessage(msg)
		m.Value = &HackerAction_AddBill{msg}
		return true, err
	default:
		return false, nil
	}
}

func _HackerAction_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*HackerAction)
	// value
	switch x := m.Value.(type) {
	case *HackerAction_AddBill:
		s := proto.Size(x.AddBill)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type HackerAddBill struct {
	StockNumber          string   `protobuf:"bytes,1,opt,name=stockNumber,proto3" json:"stockNumber,omitempty"`
	StockName            string   `protobuf:"bytes,2,opt,name=stockName,proto3" json:"stockName,omitempty"`
	Brand                string   `protobuf:"bytes,3,opt,name=brand,proto3" json:"brand,omitempty"`
	BatchRequest         string   `protobuf:"bytes,4,opt,name=batchRequest,proto3" json:"batchRequest,omitempty"`
	PledgeRate           string   `protobuf:"bytes,5,opt,name=pledgeRate,proto3" json:"pledgeRate,omitempty"`
	BasicUnit            string   `protobuf:"bytes,6,opt,name=basicUnit,proto3" json:"basicUnit,omitempty"`
	CommodityCode        string   `protobuf:"bytes,7,opt,name=commodityCode,proto3" json:"commodityCode,omitempty"`
	ExpirationDate       string   `protobuf:"bytes,8,opt,name=expirationDate,proto3" json:"expirationDate,omitempty"`
	PledgePrice          string   `protobuf:"bytes,9,opt,name=pledgePrice,proto3" json:"pledgePrice,omitempty"`
	EarlyWarningDate     string   `protobuf:"bytes,10,opt,name=earlyWarningDate,proto3" json:"earlyWarningDate,omitempty"`
	Specification        string   `protobuf:"bytes,11,opt,name=specification,proto3" json:"specification,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HackerAddBill) Reset()         { *m = HackerAddBill{} }
func (m *HackerAddBill) String() string { return proto.CompactTextString(m) }
func (*HackerAddBill) ProtoMessage()    {}
func (*HackerAddBill) Descriptor() ([]byte, []int) {
	return fileDescriptor_3aa5634b99676f1c, []int{1}
}

func (m *HackerAddBill) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HackerAddBill.Unmarshal(m, b)
}
func (m *HackerAddBill) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HackerAddBill.Marshal(b, m, deterministic)
}
func (m *HackerAddBill) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HackerAddBill.Merge(m, src)
}
func (m *HackerAddBill) XXX_Size() int {
	return xxx_messageInfo_HackerAddBill.Size(m)
}
func (m *HackerAddBill) XXX_DiscardUnknown() {
	xxx_messageInfo_HackerAddBill.DiscardUnknown(m)
}

var xxx_messageInfo_HackerAddBill proto.InternalMessageInfo

func (m *HackerAddBill) GetStockNumber() string {
	if m != nil {
		return m.StockNumber
	}
	return ""
}

func (m *HackerAddBill) GetStockName() string {
	if m != nil {
		return m.StockName
	}
	return ""
}

func (m *HackerAddBill) GetBrand() string {
	if m != nil {
		return m.Brand
	}
	return ""
}

func (m *HackerAddBill) GetBatchRequest() string {
	if m != nil {
		return m.BatchRequest
	}
	return ""
}

func (m *HackerAddBill) GetPledgeRate() string {
	if m != nil {
		return m.PledgeRate
	}
	return ""
}

func (m *HackerAddBill) GetBasicUnit() string {
	if m != nil {
		return m.BasicUnit
	}
	return ""
}

func (m *HackerAddBill) GetCommodityCode() string {
	if m != nil {
		return m.CommodityCode
	}
	return ""
}

func (m *HackerAddBill) GetExpirationDate() string {
	if m != nil {
		return m.ExpirationDate
	}
	return ""
}

func (m *HackerAddBill) GetPledgePrice() string {
	if m != nil {
		return m.PledgePrice
	}
	return ""
}

func (m *HackerAddBill) GetEarlyWarningDate() string {
	if m != nil {
		return m.EarlyWarningDate
	}
	return ""
}

func (m *HackerAddBill) GetSpecification() string {
	if m != nil {
		return m.Specification
	}
	return ""
}

type HackerBillIndex struct {
	StockNumber          string   `protobuf:"bytes,1,opt,name=stockNumber,proto3" json:"stockNumber,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Height               int64    `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	Ty                   int32    `protobuf:"varint,4,opt,name=ty,proto3" json:"ty,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HackerBillIndex) Reset()         { *m = HackerBillIndex{} }
func (m *HackerBillIndex) String() string { return proto.CompactTextString(m) }
func (*HackerBillIndex) ProtoMessage()    {}
func (*HackerBillIndex) Descriptor() ([]byte, []int) {
	return fileDescriptor_3aa5634b99676f1c, []int{2}
}

func (m *HackerBillIndex) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HackerBillIndex.Unmarshal(m, b)
}
func (m *HackerBillIndex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HackerBillIndex.Marshal(b, m, deterministic)
}
func (m *HackerBillIndex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HackerBillIndex.Merge(m, src)
}
func (m *HackerBillIndex) XXX_Size() int {
	return xxx_messageInfo_HackerBillIndex.Size(m)
}
func (m *HackerBillIndex) XXX_DiscardUnknown() {
	xxx_messageInfo_HackerBillIndex.DiscardUnknown(m)
}

var xxx_messageInfo_HackerBillIndex proto.InternalMessageInfo

func (m *HackerBillIndex) GetStockNumber() string {
	if m != nil {
		return m.StockNumber
	}
	return ""
}

func (m *HackerBillIndex) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *HackerBillIndex) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *HackerBillIndex) GetTy() int32 {
	if m != nil {
		return m.Ty
	}
	return 0
}

type HackerNfcCodeIndexList struct {
	List                 []*HackerBillIndex `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *HackerNfcCodeIndexList) Reset()         { *m = HackerNfcCodeIndexList{} }
func (m *HackerNfcCodeIndexList) String() string { return proto.CompactTextString(m) }
func (*HackerNfcCodeIndexList) ProtoMessage()    {}
func (*HackerNfcCodeIndexList) Descriptor() ([]byte, []int) {
	return fileDescriptor_3aa5634b99676f1c, []int{3}
}

func (m *HackerNfcCodeIndexList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HackerNfcCodeIndexList.Unmarshal(m, b)
}
func (m *HackerNfcCodeIndexList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HackerNfcCodeIndexList.Marshal(b, m, deterministic)
}
func (m *HackerNfcCodeIndexList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HackerNfcCodeIndexList.Merge(m, src)
}
func (m *HackerNfcCodeIndexList) XXX_Size() int {
	return xxx_messageInfo_HackerNfcCodeIndexList.Size(m)
}
func (m *HackerNfcCodeIndexList) XXX_DiscardUnknown() {
	xxx_messageInfo_HackerNfcCodeIndexList.DiscardUnknown(m)
}

var xxx_messageInfo_HackerNfcCodeIndexList proto.InternalMessageInfo

func (m *HackerNfcCodeIndexList) GetList() []*HackerBillIndex {
	if m != nil {
		return m.List
	}
	return nil
}

func init() {
	proto.RegisterType((*HackerAction)(nil), "types.HackerAction")
	proto.RegisterType((*HackerAddBill)(nil), "types.HackerAddBill")
	proto.RegisterType((*HackerBillIndex)(nil), "types.HackerBillIndex")
	proto.RegisterType((*HackerNfcCodeIndexList)(nil), "types.HackerNfcCodeIndexList")
}

func init() { proto.RegisterFile("hacker.proto", fileDescriptor_3aa5634b99676f1c) }

var fileDescriptor_3aa5634b99676f1c = []byte{
	// 399 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xc9, 0x1f, 0x27, 0x64, 0x9c, 0x16, 0xb4, 0xaa, 0xa2, 0x3d, 0x20, 0x64, 0x59, 0x08,
	0x45, 0x3d, 0x44, 0x28, 0x3c, 0x01, 0xa5, 0x87, 0x22, 0xa1, 0x0a, 0xad, 0x8a, 0x10, 0xc7, 0xf5,
	0xee, 0x34, 0x5e, 0xd5, 0xf6, 0x9a, 0xf5, 0x06, 0xd5, 0x8f, 0xc0, 0x5b, 0xa3, 0x9d, 0x0d, 0x21,
	0x86, 0x4b, 0x6f, 0x9e, 0xdf, 0x7c, 0xde, 0x6f, 0x34, 0xdf, 0xc0, 0xb2, 0x94, 0xea, 0x01, 0xdd,
	0xa6, 0x75, 0xd6, 0x5b, 0x96, 0xf8, 0xbe, 0xc5, 0x2e, 0xff, 0x0e, 0xcb, 0x1b, 0xc2, 0x1f, 0x94,
	0x37, 0xb6, 0x61, 0xef, 0x60, 0x2e, 0xb5, 0xbe, 0x32, 0x55, 0xc5, 0x47, 0xd9, 0x68, 0x9d, 0x6e,
	0x2f, 0x36, 0x24, 0xdc, 0x1c, 0x54, 0xb1, 0x77, 0xf3, 0x4c, 0xfc, 0x91, 0xb1, 0x73, 0x18, 0xdf,
	0xf5, 0x7c, 0x9b, 0x8d, 0xd6, 0x89, 0x18, 0xdf, 0xf5, 0x57, 0x73, 0x48, 0x7e, 0xca, 0x6a, 0x8f,
	0xf9, 0xaf, 0x09, 0x9c, 0x0d, 0xfe, 0x62, 0x19, 0xa4, 0x9d, 0xb7, 0xea, 0xe1, 0x76, 0x5f, 0x17,
	0xe8, 0xc8, 0x60, 0x21, 0x4e, 0x11, 0x7b, 0x05, 0x8b, 0x58, 0xca, 0x1a, 0xf9, 0x98, 0xfa, 0x7f,
	0x01, 0xbb, 0x80, 0xa4, 0x70, 0xb2, 0xd1, 0x7c, 0x42, 0x9d, 0x58, 0xb0, 0x1c, 0x96, 0x85, 0xf4,
	0xaa, 0x14, 0xf8, 0x63, 0x8f, 0x9d, 0xe7, 0x53, 0x6a, 0x0e, 0x18, 0x7b, 0x0d, 0xd0, 0x56, 0xa8,
	0x77, 0x28, 0xa4, 0x47, 0x9e, 0x90, 0xe2, 0x84, 0x04, 0xdf, 0x42, 0x76, 0x46, 0x7d, 0x6d, 0x8c,
	0xe7, 0xb3, 0xe8, 0x7b, 0x04, 0xec, 0x0d, 0x9c, 0x29, 0x5b, 0xd7, 0x56, 0x1b, 0xdf, 0x7f, 0xb4,
	0x1a, 0xf9, 0x9c, 0x14, 0x43, 0xc8, 0xde, 0xc2, 0x39, 0x3e, 0xb6, 0xc6, 0xc9, 0xb0, 0xc8, 0xeb,
	0xe0, 0xf3, 0x9c, 0x64, 0xff, 0xd0, 0xb0, 0x85, 0xe8, 0xfc, 0xc5, 0x19, 0x85, 0x7c, 0x11, 0xb7,
	0x70, 0x82, 0xd8, 0x25, 0xbc, 0x44, 0xe9, 0xaa, 0xfe, 0x9b, 0x74, 0x8d, 0x69, 0x76, 0xf4, 0x16,
	0x90, 0xec, 0x3f, 0x1e, 0x66, 0xeb, 0x5a, 0x54, 0xe6, 0xde, 0x28, 0xb2, 0xe0, 0x69, 0x9c, 0x6d,
	0x00, 0x73, 0x0b, 0x2f, 0x62, 0x14, 0x21, 0x87, 0x4f, 0x8d, 0xc6, 0xc7, 0x27, 0x84, 0xc1, 0x60,
	0x2a, 0xb5, 0x76, 0x87, 0x1c, 0xe8, 0x9b, 0xad, 0x60, 0x56, 0xa2, 0xd9, 0x95, 0x9e, 0x32, 0x98,
	0x88, 0x43, 0x15, 0xae, 0xc0, 0xf7, 0xb4, 0xfa, 0x44, 0x8c, 0x7d, 0x9f, 0x5f, 0xc3, 0x2a, 0x1a,
	0xde, 0xde, 0xab, 0xb0, 0x1d, 0xf2, 0xfc, 0x6c, 0x3a, 0xcf, 0x2e, 0x61, 0x5a, 0x99, 0xce, 0xf3,
	0x51, 0x36, 0x59, 0xa7, 0xdb, 0xd5, 0xe0, 0xbc, 0x8e, 0xd3, 0x09, 0xd2, 0x14, 0x33, 0xba, 0xd5,
	0xf7, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x26, 0xfb, 0xca, 0xdf, 0xbb, 0x02, 0x00, 0x00,
}