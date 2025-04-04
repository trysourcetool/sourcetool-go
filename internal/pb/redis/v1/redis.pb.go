// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: redis/v1/redis.proto

package redisv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RedisMessage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Payload       []byte                 `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RedisMessage) Reset() {
	*x = RedisMessage{}
	mi := &file_redis_v1_redis_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RedisMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedisMessage) ProtoMessage() {}

func (x *RedisMessage) ProtoReflect() protoreflect.Message {
	mi := &file_redis_v1_redis_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedisMessage.ProtoReflect.Descriptor instead.
func (*RedisMessage) Descriptor() ([]byte, []int) {
	return file_redis_v1_redis_proto_rawDescGZIP(), []int{0}
}

func (x *RedisMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RedisMessage) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_redis_v1_redis_proto protoreflect.FileDescriptor

const file_redis_v1_redis_proto_rawDesc = "" +
	"\n" +
	"\x14redis/v1/redis.proto\x12\bredis.v1\"8\n" +
	"\fRedisMessage\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x18\n" +
	"\apayload\x18\x02 \x01(\fR\apayloadB\xa0\x01\n" +
	"\fcom.redis.v1B\n" +
	"RedisProtoP\x01ZCgithub.com/trysourcetool/sourcetool-go/internal/pb/redis/v1;redisv1\xa2\x02\x03RXX\xaa\x02\bRedis.V1\xca\x02\bRedis\\V1\xe2\x02\x14Redis\\V1\\GPBMetadata\xea\x02\tRedis::V1b\x06proto3"

var (
	file_redis_v1_redis_proto_rawDescOnce sync.Once
	file_redis_v1_redis_proto_rawDescData []byte
)

func file_redis_v1_redis_proto_rawDescGZIP() []byte {
	file_redis_v1_redis_proto_rawDescOnce.Do(func() {
		file_redis_v1_redis_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_redis_v1_redis_proto_rawDesc), len(file_redis_v1_redis_proto_rawDesc)))
	})
	return file_redis_v1_redis_proto_rawDescData
}

var file_redis_v1_redis_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_redis_v1_redis_proto_goTypes = []any{
	(*RedisMessage)(nil), // 0: redis.v1.RedisMessage
}
var file_redis_v1_redis_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_redis_v1_redis_proto_init() }
func file_redis_v1_redis_proto_init() {
	if File_redis_v1_redis_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_redis_v1_redis_proto_rawDesc), len(file_redis_v1_redis_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_redis_v1_redis_proto_goTypes,
		DependencyIndexes: file_redis_v1_redis_proto_depIdxs,
		MessageInfos:      file_redis_v1_redis_proto_msgTypes,
	}.Build()
	File_redis_v1_redis_proto = out.File
	file_redis_v1_redis_proto_goTypes = nil
	file_redis_v1_redis_proto_depIdxs = nil
}
