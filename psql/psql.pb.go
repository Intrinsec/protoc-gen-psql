// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.14.0
// source: psql.proto

package psql

import (
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Table type enum
type TableType int32

const (
	TableType_DATA     TableType = 0
	TableType_RELATION TableType = 1
)

// Enum value maps for TableType.
var (
	TableType_name = map[int32]string{
		0: "DATA",
		1: "RELATION",
	}
	TableType_value = map[string]int32{
		"DATA":     0,
		"RELATION": 1,
	}
)

func (x TableType) Enum() *TableType {
	p := new(TableType)
	*p = x
	return p
}

func (x TableType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TableType) Descriptor() protoreflect.EnumDescriptor {
	return file_psql_proto_enumTypes[0].Descriptor()
}

func (TableType) Type() protoreflect.EnumType {
	return &file_psql_proto_enumTypes[0]
}

func (x TableType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TableType.Descriptor instead.
func (TableType) EnumDescriptor() ([]byte, []int) {
	return file_psql_proto_rawDescGZIP(), []int{0}
}

type CascadeUpdateOnRelatedTable struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field string `protobuf:"bytes,1,opt,name=field,proto3" json:"field,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *CascadeUpdateOnRelatedTable) Reset() {
	*x = CascadeUpdateOnRelatedTable{}
	if protoimpl.UnsafeEnabled {
		mi := &file_psql_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CascadeUpdateOnRelatedTable) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CascadeUpdateOnRelatedTable) ProtoMessage() {}

func (x *CascadeUpdateOnRelatedTable) ProtoReflect() protoreflect.Message {
	mi := &file_psql_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CascadeUpdateOnRelatedTable.ProtoReflect.Descriptor instead.
func (*CascadeUpdateOnRelatedTable) Descriptor() ([]byte, []int) {
	return file_psql_proto_rawDescGZIP(), []int{0}
}

func (x *CascadeUpdateOnRelatedTable) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *CascadeUpdateOnRelatedTable) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type RelayCascadeUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SourceForeignKey string                            `protobuf:"bytes,1,opt,name=source_foreign_key,json=sourceForeignKey,proto3" json:"source_foreign_key,omitempty"`
	Destinations     []*RelayCascadeUpdate_Destination `protobuf:"bytes,2,rep,name=destinations,proto3" json:"destinations,omitempty"`
}

func (x *RelayCascadeUpdate) Reset() {
	*x = RelayCascadeUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_psql_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RelayCascadeUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelayCascadeUpdate) ProtoMessage() {}

func (x *RelayCascadeUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_psql_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelayCascadeUpdate.ProtoReflect.Descriptor instead.
func (*RelayCascadeUpdate) Descriptor() ([]byte, []int) {
	return file_psql_proto_rawDescGZIP(), []int{1}
}

func (x *RelayCascadeUpdate) GetSourceForeignKey() string {
	if x != nil {
		return x.SourceForeignKey
	}
	return ""
}

func (x *RelayCascadeUpdate) GetDestinations() []*RelayCascadeUpdate_Destination {
	if x != nil {
		return x.Destinations
	}
	return nil
}

type RelayCascadeUpdate_Destination struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ForeignKey string `protobuf:"bytes,1,opt,name=foreign_key,json=foreignKey,proto3" json:"foreign_key,omitempty"`
	Field      string `protobuf:"bytes,2,opt,name=field,proto3" json:"field,omitempty"`
	Value      string `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *RelayCascadeUpdate_Destination) Reset() {
	*x = RelayCascadeUpdate_Destination{}
	if protoimpl.UnsafeEnabled {
		mi := &file_psql_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RelayCascadeUpdate_Destination) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelayCascadeUpdate_Destination) ProtoMessage() {}

func (x *RelayCascadeUpdate_Destination) ProtoReflect() protoreflect.Message {
	mi := &file_psql_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelayCascadeUpdate_Destination.ProtoReflect.Descriptor instead.
func (*RelayCascadeUpdate_Destination) Descriptor() ([]byte, []int) {
	return file_psql_proto_rawDescGZIP(), []int{1, 0}
}

func (x *RelayCascadeUpdate_Destination) GetForeignKey() string {
	if x != nil {
		return x.ForeignKey
	}
	return ""
}

func (x *RelayCascadeUpdate_Destination) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *RelayCascadeUpdate_Destination) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var file_psql_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.FileOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         1091,
		Name:          "psql.initialization",
		Tag:           "bytes,1091,rep,name=initialization",
		Filename:      "psql.proto",
	},
	{
		ExtendedType:  (*descriptor.FileOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         1092,
		Name:          "psql.finalization",
		Tag:           "bytes,1092,rep,name=finalization",
		Filename:      "psql.proto",
	},
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         1091,
		Name:          "psql.disabled",
		Tag:           "varint,1091,opt,name=disabled",
		Filename:      "psql.proto",
	},
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         1092,
		Name:          "psql.prefix",
		Tag:           "bytes,1092,rep,name=prefix",
		Filename:      "psql.proto",
	},
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         1093,
		Name:          "psql.suffix",
		Tag:           "bytes,1093,rep,name=suffix",
		Filename:      "psql.proto",
	},
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         1094,
		Name:          "psql.constraint",
		Tag:           "bytes,1094,rep,name=constraint",
		Filename:      "psql.proto",
	},
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: (*TableType)(nil),
		Field:         1095,
		Name:          "psql.tableType",
		Tag:           "varint,1095,opt,name=tableType,enum=psql.TableType",
		Filename:      "psql.proto",
	},
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: ([]*RelayCascadeUpdate)(nil),
		Field:         1096,
		Name:          "psql.relay_cascade_update",
		Tag:           "bytes,1096,rep,name=relay_cascade_update",
		Filename:      "psql.proto",
	},
	{
		ExtendedType:  (*descriptor.FieldOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         1091,
		Name:          "psql.column",
		Tag:           "bytes,1091,opt,name=column",
		Filename:      "psql.proto",
	},
	{
		ExtendedType:  (*descriptor.FieldOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         1092,
		Name:          "psql.auto_fill_on_update",
		Tag:           "bytes,1092,opt,name=auto_fill_on_update",
		Filename:      "psql.proto",
	},
	{
		ExtendedType:  (*descriptor.FieldOptions)(nil),
		ExtensionType: ([]*CascadeUpdateOnRelatedTable)(nil),
		Field:         1093,
		Name:          "psql.cascade_update_on_related_table",
		Tag:           "bytes,1093,rep,name=cascade_update_on_related_table",
		Filename:      "psql.proto",
	},
}

// Extension fields to descriptor.FileOptions.
var (
	// repeated string initialization = 1091;
	E_Initialization = &file_psql_proto_extTypes[0]
	// repeated string finalization = 1092;
	E_Finalization = &file_psql_proto_extTypes[1]
)

// Extension fields to descriptor.MessageOptions.
var (
	// optional bool disabled = 1091;
	E_Disabled = &file_psql_proto_extTypes[2]
	// repeated string prefix = 1092;
	E_Prefix = &file_psql_proto_extTypes[3]
	// repeated string suffix = 1093;
	E_Suffix = &file_psql_proto_extTypes[4]
	// repeated string constraint = 1094;
	E_Constraint = &file_psql_proto_extTypes[5]
	// optional psql.TableType tableType = 1095;
	E_TableType = &file_psql_proto_extTypes[6]
	// repeated psql.RelayCascadeUpdate relay_cascade_update = 1096;
	E_RelayCascadeUpdate = &file_psql_proto_extTypes[7]
)

// Extension fields to descriptor.FieldOptions.
var (
	// optional string column = 1091;
	E_Column = &file_psql_proto_extTypes[8]
	// optional string auto_fill_on_update = 1092;
	E_AutoFillOnUpdate = &file_psql_proto_extTypes[9]
	// repeated psql.CascadeUpdateOnRelatedTable cascade_update_on_related_table = 1093;
	E_CascadeUpdateOnRelatedTable = &file_psql_proto_extTypes[10]
)

var File_psql_proto protoreflect.FileDescriptor

var file_psql_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x70, 0x73, 0x71, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x70, 0x73,
	0x71, 0x6c, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x49, 0x0a, 0x1b, 0x43, 0x61, 0x73, 0x63, 0x61, 0x64, 0x65, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x6e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x54, 0x61,
	0x62, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22,
	0xe8, 0x01, 0x0a, 0x12, 0x52, 0x65, 0x6c, 0x61, 0x79, 0x43, 0x61, 0x73, 0x63, 0x61, 0x64, 0x65,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x2c, 0x0a, 0x12, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5f, 0x66, 0x6f, 0x72, 0x65, 0x69, 0x67, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x10, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x46, 0x6f, 0x72, 0x65, 0x69, 0x67,
	0x6e, 0x4b, 0x65, 0x79, 0x12, 0x48, 0x0a, 0x0c, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x70, 0x73, 0x71,
	0x6c, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x79, 0x43, 0x61, 0x73, 0x63, 0x61, 0x64, 0x65, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0c, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x5a,
	0x0a, 0x0b, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a,
	0x0b, 0x66, 0x6f, 0x72, 0x65, 0x69, 0x67, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x66, 0x6f, 0x72, 0x65, 0x69, 0x67, 0x6e, 0x4b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2a, 0x23, 0x0a, 0x09, 0x54, 0x61,
	0x62, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x41, 0x54, 0x41, 0x10,
	0x00, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x4c, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x01, 0x3a,
	0x45, 0x0a, 0x0e, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0xc3, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x41, 0x0a, 0x0c, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0xc4, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x69, 0x6e,
	0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x3f, 0x0a, 0x08, 0x64, 0x69, 0x73,
	0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xc3, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64,
	0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x88, 0x01, 0x01, 0x3a, 0x38, 0x0a, 0x06, 0x70, 0x72,
	0x65, 0x66, 0x69, 0x78, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xc4, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72,
	0x65, 0x66, 0x69, 0x78, 0x3a, 0x38, 0x0a, 0x06, 0x73, 0x75, 0x66, 0x66, 0x69, 0x78, 0x12, 0x1f,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0xc5, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x73, 0x75, 0x66, 0x66, 0x69, 0x78, 0x3a, 0x40,
	0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x12, 0x1f, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xc6, 0x08,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74,
	0x3a, 0x52, 0x0a, 0x09, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xc7,
	0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x70, 0x73, 0x71, 0x6c, 0x2e, 0x54, 0x61, 0x62,
	0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x88, 0x01, 0x01, 0x3a, 0x6c, 0x0a, 0x14, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x63, 0x61,
	0x73, 0x63, 0x61, 0x64, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xc8, 0x08,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x73, 0x71, 0x6c, 0x2e, 0x52, 0x65, 0x6c, 0x61,
	0x79, 0x43, 0x61, 0x73, 0x63, 0x61, 0x64, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x12,
	0x72, 0x65, 0x6c, 0x61, 0x79, 0x43, 0x61, 0x73, 0x63, 0x61, 0x64, 0x65, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x3a, 0x39, 0x0a, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x12, 0x1d, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xc3, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x88, 0x01, 0x01, 0x3a, 0x50, 0x0a,
	0x13, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x66, 0x69, 0x6c, 0x6c, 0x5f, 0x6f, 0x6e, 0x5f, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0xc4, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x61, 0x75, 0x74, 0x6f,
	0x46, 0x69, 0x6c, 0x6c, 0x4f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x3a,
	0x87, 0x01, 0x0a, 0x1f, 0x63, 0x61, 0x73, 0x63, 0x61, 0x64, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x5f, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0xc5, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x70, 0x73, 0x71, 0x6c,
	0x2e, 0x43, 0x61, 0x73, 0x63, 0x61, 0x64, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x6e,
	0x52, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x1b, 0x63, 0x61,
	0x73, 0x63, 0x61, 0x64, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x6e, 0x52, 0x65, 0x6c,
	0x61, 0x74, 0x65, 0x64, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x74, 0x72, 0x69, 0x6e, 0x73, 0x65,
	0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x70, 0x73, 0x71,
	0x6c, 0x2f, 0x70, 0x73, 0x71, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_psql_proto_rawDescOnce sync.Once
	file_psql_proto_rawDescData = file_psql_proto_rawDesc
)

func file_psql_proto_rawDescGZIP() []byte {
	file_psql_proto_rawDescOnce.Do(func() {
		file_psql_proto_rawDescData = protoimpl.X.CompressGZIP(file_psql_proto_rawDescData)
	})
	return file_psql_proto_rawDescData
}

var file_psql_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_psql_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_psql_proto_goTypes = []interface{}{
	(TableType)(0),                         // 0: psql.TableType
	(*CascadeUpdateOnRelatedTable)(nil),    // 1: psql.CascadeUpdateOnRelatedTable
	(*RelayCascadeUpdate)(nil),             // 2: psql.RelayCascadeUpdate
	(*RelayCascadeUpdate_Destination)(nil), // 3: psql.RelayCascadeUpdate.Destination
	(*descriptor.FileOptions)(nil),         // 4: google.protobuf.FileOptions
	(*descriptor.MessageOptions)(nil),      // 5: google.protobuf.MessageOptions
	(*descriptor.FieldOptions)(nil),        // 6: google.protobuf.FieldOptions
}
var file_psql_proto_depIdxs = []int32{
	3,  // 0: psql.RelayCascadeUpdate.destinations:type_name -> psql.RelayCascadeUpdate.Destination
	4,  // 1: psql.initialization:extendee -> google.protobuf.FileOptions
	4,  // 2: psql.finalization:extendee -> google.protobuf.FileOptions
	5,  // 3: psql.disabled:extendee -> google.protobuf.MessageOptions
	5,  // 4: psql.prefix:extendee -> google.protobuf.MessageOptions
	5,  // 5: psql.suffix:extendee -> google.protobuf.MessageOptions
	5,  // 6: psql.constraint:extendee -> google.protobuf.MessageOptions
	5,  // 7: psql.tableType:extendee -> google.protobuf.MessageOptions
	5,  // 8: psql.relay_cascade_update:extendee -> google.protobuf.MessageOptions
	6,  // 9: psql.column:extendee -> google.protobuf.FieldOptions
	6,  // 10: psql.auto_fill_on_update:extendee -> google.protobuf.FieldOptions
	6,  // 11: psql.cascade_update_on_related_table:extendee -> google.protobuf.FieldOptions
	0,  // 12: psql.tableType:type_name -> psql.TableType
	2,  // 13: psql.relay_cascade_update:type_name -> psql.RelayCascadeUpdate
	1,  // 14: psql.cascade_update_on_related_table:type_name -> psql.CascadeUpdateOnRelatedTable
	15, // [15:15] is the sub-list for method output_type
	15, // [15:15] is the sub-list for method input_type
	12, // [12:15] is the sub-list for extension type_name
	1,  // [1:12] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_psql_proto_init() }
func file_psql_proto_init() {
	if File_psql_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_psql_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CascadeUpdateOnRelatedTable); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_psql_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RelayCascadeUpdate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_psql_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RelayCascadeUpdate_Destination); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_psql_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 11,
			NumServices:   0,
		},
		GoTypes:           file_psql_proto_goTypes,
		DependencyIndexes: file_psql_proto_depIdxs,
		EnumInfos:         file_psql_proto_enumTypes,
		MessageInfos:      file_psql_proto_msgTypes,
		ExtensionInfos:    file_psql_proto_extTypes,
	}.Build()
	File_psql_proto = out.File
	file_psql_proto_rawDesc = nil
	file_psql_proto_goTypes = nil
	file_psql_proto_depIdxs = nil
}
