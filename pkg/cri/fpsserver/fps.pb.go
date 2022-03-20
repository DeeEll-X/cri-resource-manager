//
//Copyright 2019 Intel Corporation
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: fps.proto

package fpsserver

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

type FPSDropRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *FPSDropRequest) Reset() {
	*x = FPSDropRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fps_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FPSDropRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FPSDropRequest) ProtoMessage() {}

func (x *FPSDropRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fps_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FPSDropRequest.ProtoReflect.Descriptor instead.
func (*FPSDropRequest) Descriptor() ([]byte, []int) {
	return file_fps_proto_rawDescGZIP(), []int{0}
}

func (x *FPSDropRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type FPSDropReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FPSDropReply) Reset() {
	*x = FPSDropReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fps_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FPSDropReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FPSDropReply) ProtoMessage() {}

func (x *FPSDropReply) ProtoReflect() protoreflect.Message {
	mi := &file_fps_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FPSDropReply.ProtoReflect.Descriptor instead.
func (*FPSDropReply) Descriptor() ([]byte, []int) {
	return file_fps_proto_rawDescGZIP(), []int{1}
}

var File_fps_proto protoreflect.FileDescriptor

var file_fps_proto_rawDesc = []byte{
	0x0a, 0x09, 0x66, 0x70, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x66, 0x70, 0x73,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x22, 0x20, 0x0a, 0x0e, 0x46, 0x50, 0x53, 0x44, 0x72, 0x6f,
	0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x0e, 0x0a, 0x0c, 0x46, 0x50, 0x53, 0x44,
	0x72, 0x6f, 0x70, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x32, 0x4d, 0x0a, 0x0a, 0x46, 0x50, 0x53, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3f, 0x0a, 0x07, 0x46, 0x50, 0x53, 0x44, 0x72, 0x6f,
	0x70, 0x12, 0x19, 0x2e, 0x66, 0x70, 0x73, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x46, 0x50,
	0x53, 0x44, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x66,
	0x70, 0x73, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x46, 0x50, 0x53, 0x44, 0x72, 0x6f, 0x70,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x0c, 0x5a, 0x0a, 0x2f, 0x66, 0x70, 0x73, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fps_proto_rawDescOnce sync.Once
	file_fps_proto_rawDescData = file_fps_proto_rawDesc
)

func file_fps_proto_rawDescGZIP() []byte {
	file_fps_proto_rawDescOnce.Do(func() {
		file_fps_proto_rawDescData = protoimpl.X.CompressGZIP(file_fps_proto_rawDescData)
	})
	return file_fps_proto_rawDescData
}

var file_fps_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_fps_proto_goTypes = []interface{}{
	(*FPSDropRequest)(nil), // 0: fpsserver.FPSDropRequest
	(*FPSDropReply)(nil),   // 1: fpsserver.FPSDropReply
}
var file_fps_proto_depIdxs = []int32{
	0, // 0: fpsserver.FPSService.FPSDrop:input_type -> fpsserver.FPSDropRequest
	1, // 1: fpsserver.FPSService.FPSDrop:output_type -> fpsserver.FPSDropReply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_fps_proto_init() }
func file_fps_proto_init() {
	if File_fps_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fps_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FPSDropRequest); i {
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
		file_fps_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FPSDropReply); i {
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
			RawDescriptor: file_fps_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_fps_proto_goTypes,
		DependencyIndexes: file_fps_proto_depIdxs,
		MessageInfos:      file_fps_proto_msgTypes,
	}.Build()
	File_fps_proto = out.File
	file_fps_proto_rawDesc = nil
	file_fps_proto_goTypes = nil
	file_fps_proto_depIdxs = nil
}