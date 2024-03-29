// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: ssh.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SshAccount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsDel     bool   `protobuf:"varint,1,opt,name=is_del,json=isDel,proto3" json:"is_del,omitempty"`
	IsKill    bool   `protobuf:"varint,2,opt,name=is_kill,json=isKill,proto3" json:"is_kill,omitempty"`
	Username  string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	PublicKey string `protobuf:"bytes,4,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
}

func (x *SshAccount) Reset() {
	*x = SshAccount{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ssh_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SshAccount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SshAccount) ProtoMessage() {}

func (x *SshAccount) ProtoReflect() protoreflect.Message {
	mi := &file_ssh_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SshAccount.ProtoReflect.Descriptor instead.
func (*SshAccount) Descriptor() ([]byte, []int) {
	return file_ssh_proto_rawDescGZIP(), []int{0}
}

func (x *SshAccount) GetIsDel() bool {
	if x != nil {
		return x.IsDel
	}
	return false
}

func (x *SshAccount) GetIsKill() bool {
	if x != nil {
		return x.IsKill
	}
	return false
}

func (x *SshAccount) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *SshAccount) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

type AccountStream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsInit      bool          `protobuf:"varint,1,opt,name=is_init,json=isInit,proto3" json:"is_init,omitempty"`
	IsHeartBeat bool          `protobuf:"varint,2,opt,name=is_heart_beat,json=isHeartBeat,proto3" json:"is_heart_beat,omitempty"`
	Accounts    []*SshAccount `protobuf:"bytes,3,rep,name=accounts,proto3" json:"accounts,omitempty"`
}

func (x *AccountStream) Reset() {
	*x = AccountStream{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ssh_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccountStream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountStream) ProtoMessage() {}

func (x *AccountStream) ProtoReflect() protoreflect.Message {
	mi := &file_ssh_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountStream.ProtoReflect.Descriptor instead.
func (*AccountStream) Descriptor() ([]byte, []int) {
	return file_ssh_proto_rawDescGZIP(), []int{1}
}

func (x *AccountStream) GetIsInit() bool {
	if x != nil {
		return x.IsInit
	}
	return false
}

func (x *AccountStream) GetIsHeartBeat() bool {
	if x != nil {
		return x.IsHeartBeat
	}
	return false
}

func (x *AccountStream) GetAccounts() []*SshAccount {
	if x != nil {
		return x.Accounts
	}
	return nil
}

var File_ssh_proto protoreflect.FileDescriptor

var file_ssh_proto_rawDesc = []byte{
	0x0a, 0x09, 0x73, 0x73, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x77, 0x0a, 0x0a, 0x53, 0x73, 0x68, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x15, 0x0a,
	0x06, 0x69, 0x73, 0x5f, 0x64, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x69,
	0x73, 0x44, 0x65, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x6b, 0x69, 0x6c, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x4b, 0x69, 0x6c, 0x6c, 0x12, 0x1a, 0x0a,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x22, 0x7b, 0x0a, 0x0d, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f,
	0x69, 0x6e, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x49, 0x6e,
	0x69, 0x74, 0x12, 0x22, 0x0a, 0x0d, 0x69, 0x73, 0x5f, 0x68, 0x65, 0x61, 0x72, 0x74, 0x5f, 0x62,
	0x65, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x48, 0x65, 0x61,
	0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x12, 0x2d, 0x0a, 0x08, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x73, 0x68, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x08, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x73, 0x32, 0x48, 0x0a, 0x0b, 0x53, 0x73, 0x68, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x73, 0x12, 0x39, 0x0a, 0x05, 0x57, 0x61, 0x74, 0x63, 0x68, 0x12, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x22, 0x00, 0x30, 0x01, 0x42,
	0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_ssh_proto_rawDescOnce sync.Once
	file_ssh_proto_rawDescData = file_ssh_proto_rawDesc
)

func file_ssh_proto_rawDescGZIP() []byte {
	file_ssh_proto_rawDescOnce.Do(func() {
		file_ssh_proto_rawDescData = protoimpl.X.CompressGZIP(file_ssh_proto_rawDescData)
	})
	return file_ssh_proto_rawDescData
}

var file_ssh_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_ssh_proto_goTypes = []interface{}{
	(*SshAccount)(nil),    // 0: proto.SshAccount
	(*AccountStream)(nil), // 1: proto.AccountStream
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_ssh_proto_depIdxs = []int32{
	0, // 0: proto.AccountStream.accounts:type_name -> proto.SshAccount
	2, // 1: proto.SshAccounts.Watch:input_type -> google.protobuf.Empty
	1, // 2: proto.SshAccounts.Watch:output_type -> proto.AccountStream
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ssh_proto_init() }
func file_ssh_proto_init() {
	if File_ssh_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ssh_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SshAccount); i {
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
		file_ssh_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccountStream); i {
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
			RawDescriptor: file_ssh_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ssh_proto_goTypes,
		DependencyIndexes: file_ssh_proto_depIdxs,
		MessageInfos:      file_ssh_proto_msgTypes,
	}.Build()
	File_ssh_proto = out.File
	file_ssh_proto_rawDesc = nil
	file_ssh_proto_goTypes = nil
	file_ssh_proto_depIdxs = nil
}
