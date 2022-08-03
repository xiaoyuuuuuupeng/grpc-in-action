// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.17.3
// source: easyresponse.proto

package pb

import (
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

type CasLoginReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserName string `protobuf:"bytes,1,opt,name=userName,proto3" json:"userName,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *CasLoginReq) Reset() {
	*x = CasLoginReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_casservice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CasLoginReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CasLoginReq) ProtoMessage() {}

func (x *CasLoginReq) ProtoReflect() protoreflect.Message {
	mi := &file_casservice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CasLoginReq.ProtoReflect.Descriptor instead.
func (*CasLoginReq) Descriptor() ([]byte, []int) {
	return file_casservice_proto_rawDescGZIP(), []int{0}
}

func (x *CasLoginReq) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *CasLoginReq) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type CasLoginRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Sex string `protobuf:"bytes,2,opt,name=sex,proto3" json:"sex,omitempty"`
}

func (x *CasLoginRes) Reset() {
	*x = CasLoginRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_casservice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CasLoginRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CasLoginRes) ProtoMessage() {}

func (x *CasLoginRes) ProtoReflect() protoreflect.Message {
	mi := &file_casservice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CasLoginRes.ProtoReflect.Descriptor instead.
func (*CasLoginRes) Descriptor() ([]byte, []int) {
	return file_casservice_proto_rawDescGZIP(), []int{1}
}

func (x *CasLoginRes) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *CasLoginRes) GetSex() string {
	if x != nil {
		return x.Sex
	}
	return ""
}

var File_casservice_proto protoreflect.FileDescriptor

var file_casservice_proto_rawDesc = []byte{
	0x0a, 0x10, 0x63, 0x61, 0x73, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x63, 0x61, 0x73, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x45,
	0x0a, 0x0b, 0x43, 0x61, 0x73, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x31, 0x0a, 0x0b, 0x43, 0x61, 0x73, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x52, 0x65, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65, 0x78, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x65, 0x78, 0x32, 0x4a, 0x0a, 0x0a, 0x43, 0x61, 0x73, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3c, 0x0a, 0x08, 0x63, 0x61, 0x73, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x12, 0x17, 0x2e, 0x63, 0x61, 0x73, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x43, 0x61, 0x73, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x63, 0x61,
	0x73, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x61, 0x73, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x52, 0x65, 0x73, 0x42, 0x35, 0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x2e, 0x63, 0x61, 0x73, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x42, 0x12,
	0x43, 0x61, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x62, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x80, 0x01, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_casservice_proto_rawDescOnce sync.Once
	file_casservice_proto_rawDescData = file_casservice_proto_rawDesc
)

func file_casservice_proto_rawDescGZIP() []byte {
	file_casservice_proto_rawDescOnce.Do(func() {
		file_casservice_proto_rawDescData = protoimpl.X.CompressGZIP(file_casservice_proto_rawDescData)
	})
	return file_casservice_proto_rawDescData
}

var file_casservice_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_casservice_proto_goTypes = []interface{}{
	(*CasLoginReq)(nil), // 0: casservice.CasLoginReq
	(*CasLoginRes)(nil), // 1: casservice.CasLoginRes
}
var file_casservice_proto_depIdxs = []int32{
	0, // 0: casservice.CasService.casLogin:input_type -> casservice.CasLoginReq
	1, // 1: casservice.CasService.casLogin:output_type -> casservice.CasLoginRes
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_casservice_proto_init() }
func file_casservice_proto_init() {
	if File_casservice_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_casservice_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CasLoginReq); i {
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
		file_casservice_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CasLoginRes); i {
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
			RawDescriptor: file_casservice_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_casservice_proto_goTypes,
		DependencyIndexes: file_casservice_proto_depIdxs,
		MessageInfos:      file_casservice_proto_msgTypes,
	}.Build()
	File_casservice_proto = out.File
	file_casservice_proto_rawDesc = nil
	file_casservice_proto_goTypes = nil
	file_casservice_proto_depIdxs = nil
}
