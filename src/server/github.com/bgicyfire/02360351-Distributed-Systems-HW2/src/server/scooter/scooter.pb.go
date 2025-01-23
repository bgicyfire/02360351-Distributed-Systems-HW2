// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v3.21.12
// source: scooter.proto

package scooter

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

type ScooterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ScooterId     string                 `protobuf:"bytes,1,opt,name=scooter_id,json=scooterId,proto3" json:"scooter_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ScooterRequest) Reset() {
	*x = ScooterRequest{}
	mi := &file_scooter_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScooterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScooterRequest) ProtoMessage() {}

func (x *ScooterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scooter_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScooterRequest.ProtoReflect.Descriptor instead.
func (*ScooterRequest) Descriptor() ([]byte, []int) {
	return file_scooter_proto_rawDescGZIP(), []int{0}
}

func (x *ScooterRequest) GetScooterId() string {
	if x != nil {
		return x.ScooterId
	}
	return ""
}

type ScooterResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Hostname      string                 `protobuf:"bytes,2,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Myleader      string                 `protobuf:"bytes,3,opt,name=myleader,proto3" json:"myleader,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ScooterResponse) Reset() {
	*x = ScooterResponse{}
	mi := &file_scooter_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScooterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScooterResponse) ProtoMessage() {}

func (x *ScooterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scooter_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScooterResponse.ProtoReflect.Descriptor instead.
func (*ScooterResponse) Descriptor() ([]byte, []int) {
	return file_scooter_proto_rawDescGZIP(), []int{1}
}

func (x *ScooterResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ScooterResponse) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *ScooterResponse) GetMyleader() string {
	if x != nil {
		return x.Myleader
	}
	return ""
}

var File_scooter_proto protoreflect.FileDescriptor

var file_scooter_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x22, 0x2f, 0x0a, 0x0e, 0x53, 0x63, 0x6f, 0x6f,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x63,
	0x6f, 0x6f, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x49, 0x64, 0x22, 0x61, 0x0a, 0x0f, 0x53, 0x63, 0x6f,
	0x6f, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x79, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6d, 0x79, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x32, 0x57, 0x0a, 0x0e,
	0x53, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x17, 0x2e, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x53, 0x63, 0x6f,
	0x6f, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x63,
	0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2e, 0x53, 0x63, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x4a, 0x5a, 0x48, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x67, 0x69, 0x63, 0x79, 0x66, 0x69, 0x72, 0x65, 0x2f, 0x30, 0x32,
	0x33, 0x36, 0x30, 0x33, 0x35, 0x31, 0x2d, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x64, 0x2d, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2d, 0x48, 0x57, 0x32, 0x2f, 0x73,
	0x72, 0x63, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x73, 0x63, 0x6f, 0x6f, 0x74, 0x65,
	0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_scooter_proto_rawDescOnce sync.Once
	file_scooter_proto_rawDescData = file_scooter_proto_rawDesc
)

func file_scooter_proto_rawDescGZIP() []byte {
	file_scooter_proto_rawDescOnce.Do(func() {
		file_scooter_proto_rawDescData = protoimpl.X.CompressGZIP(file_scooter_proto_rawDescData)
	})
	return file_scooter_proto_rawDescData
}

var file_scooter_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_scooter_proto_goTypes = []any{
	(*ScooterRequest)(nil),  // 0: scooter.ScooterRequest
	(*ScooterResponse)(nil), // 1: scooter.ScooterResponse
}
var file_scooter_proto_depIdxs = []int32{
	0, // 0: scooter.ScooterService.GetScooterStatus:input_type -> scooter.ScooterRequest
	1, // 1: scooter.ScooterService.GetScooterStatus:output_type -> scooter.ScooterResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_scooter_proto_init() }
func file_scooter_proto_init() {
	if File_scooter_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_scooter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_scooter_proto_goTypes,
		DependencyIndexes: file_scooter_proto_depIdxs,
		MessageInfos:      file_scooter_proto_msgTypes,
	}.Build()
	File_scooter_proto = out.File
	file_scooter_proto_rawDesc = nil
	file_scooter_proto_goTypes = nil
	file_scooter_proto_depIdxs = nil
}
