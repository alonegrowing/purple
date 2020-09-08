// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.13.0
// source: purple.proto

package purple

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type HomePageParam struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *HomePageParam) Reset() {
	*x = HomePageParam{}
	if protoimpl.UnsafeEnabled {
		mi := &file_purple_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HomePageParam) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HomePageParam) ProtoMessage() {}

func (x *HomePageParam) ProtoReflect() protoreflect.Message {
	mi := &file_purple_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HomePageParam.ProtoReflect.Descriptor instead.
func (*HomePageParam) Descriptor() ([]byte, []int) {
	return file_purple_proto_rawDescGZIP(), []int{0}
}

func (x *HomePageParam) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type HomePageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *HomePageResponse) Reset() {
	*x = HomePageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_purple_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HomePageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HomePageResponse) ProtoMessage() {}

func (x *HomePageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_purple_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HomePageResponse.ProtoReflect.Descriptor instead.
func (*HomePageResponse) Descriptor() ([]byte, []int) {
	return file_purple_proto_rawDescGZIP(), []int{1}
}

func (x *HomePageResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *HomePageResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_purple_proto protoreflect.FileDescriptor

var file_purple_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x75, 0x72, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x70, 0x75, 0x72, 0x70, 0x6c, 0x65, 0x22, 0x1f, 0x0a, 0x0d, 0x48, 0x6f, 0x6d, 0x65, 0x50, 0x61,
	0x67, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x36, 0x0a, 0x10, 0x48, 0x6f, 0x6d, 0x65, 0x50,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x32,
	0x4a, 0x0a, 0x06, 0x50, 0x75, 0x72, 0x70, 0x6c, 0x65, 0x12, 0x40, 0x0a, 0x0b, 0x47, 0x65, 0x74,
	0x48, 0x6f, 0x6d, 0x65, 0x50, 0x61, 0x67, 0x65, 0x12, 0x15, 0x2e, 0x70, 0x75, 0x72, 0x70, 0x6c,
	0x65, 0x2e, 0x48, 0x6f, 0x6d, 0x65, 0x50, 0x61, 0x67, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x1a,
	0x18, 0x2e, 0x70, 0x75, 0x72, 0x70, 0x6c, 0x65, 0x2e, 0x48, 0x6f, 0x6d, 0x65, 0x50, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2f,
	0x70, 0x75, 0x72, 0x70, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_purple_proto_rawDescOnce sync.Once
	file_purple_proto_rawDescData = file_purple_proto_rawDesc
)

func file_purple_proto_rawDescGZIP() []byte {
	file_purple_proto_rawDescOnce.Do(func() {
		file_purple_proto_rawDescData = protoimpl.X.CompressGZIP(file_purple_proto_rawDescData)
	})
	return file_purple_proto_rawDescData
}

var file_purple_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_purple_proto_goTypes = []interface{}{
	(*HomePageParam)(nil),    // 0: purple.HomePageParam
	(*HomePageResponse)(nil), // 1: purple.HomePageResponse
}
var file_purple_proto_depIdxs = []int32{
	0, // 0: purple.Purple.GetHomePage:input_type -> purple.HomePageParam
	1, // 1: purple.Purple.GetHomePage:output_type -> purple.HomePageResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_purple_proto_init() }
func file_purple_proto_init() {
	if File_purple_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_purple_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HomePageParam); i {
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
		file_purple_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HomePageResponse); i {
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
			RawDescriptor: file_purple_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_purple_proto_goTypes,
		DependencyIndexes: file_purple_proto_depIdxs,
		MessageInfos:      file_purple_proto_msgTypes,
	}.Build()
	File_purple_proto = out.File
	file_purple_proto_rawDesc = nil
	file_purple_proto_goTypes = nil
	file_purple_proto_depIdxs = nil
}