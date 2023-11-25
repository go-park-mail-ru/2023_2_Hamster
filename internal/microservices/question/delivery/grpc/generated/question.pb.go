// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.24.4
// source: question.proto

package __

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

type AnswerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Rating int32  `protobuf:"varint,3,opt,name=rating,proto3" json:"rating,omitempty"`
}

func (x *AnswerRequest) Reset() {
	*x = AnswerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnswerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnswerRequest) ProtoMessage() {}

func (x *AnswerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_question_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnswerRequest.ProtoReflect.Descriptor instead.
func (*AnswerRequest) Descriptor() ([]byte, []int) {
	return file_question_proto_rawDescGZIP(), []int{0}
}

func (x *AnswerRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AnswerRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AnswerRequest) GetRating() int32 {
	if x != nil {
		return x.Rating
	}
	return 0
}

type AverageRatingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AverageRating int32 `protobuf:"varint,1,opt,name=averageRating,proto3" json:"averageRating,omitempty"`
}

func (x *AverageRatingResponse) Reset() {
	*x = AverageRatingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AverageRatingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AverageRatingResponse) ProtoMessage() {}

func (x *AverageRatingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_question_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AverageRatingResponse.ProtoReflect.Descriptor instead.
func (*AverageRatingResponse) Descriptor() ([]byte, []int) {
	return file_question_proto_rawDescGZIP(), []int{1}
}

func (x *AverageRatingResponse) GetAverageRating() int32 {
	if x != nil {
		return x.AverageRating
	}
	return 0
}

type CalculateAverageRatingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QuestionName string `protobuf:"bytes,2,opt,name=questionName,proto3" json:"questionName,omitempty"`
}

func (x *CalculateAverageRatingRequest) Reset() {
	*x = CalculateAverageRatingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CalculateAverageRatingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CalculateAverageRatingRequest) ProtoMessage() {}

func (x *CalculateAverageRatingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_question_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CalculateAverageRatingRequest.ProtoReflect.Descriptor instead.
func (*CalculateAverageRatingRequest) Descriptor() ([]byte, []int) {
	return file_question_proto_rawDescGZIP(), []int{2}
}

func (x *CalculateAverageRatingRequest) GetQuestionName() string {
	if x != nil {
		return x.QuestionName
	}
	return ""
}

type CheckUserAnswerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	QuestionName string `protobuf:"bytes,2,opt,name=questionName,proto3" json:"questionName,omitempty"`
}

func (x *CheckUserAnswerRequest) Reset() {
	*x = CheckUserAnswerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckUserAnswerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckUserAnswerRequest) ProtoMessage() {}

func (x *CheckUserAnswerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_question_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckUserAnswerRequest.ProtoReflect.Descriptor instead.
func (*CheckUserAnswerRequest) Descriptor() ([]byte, []int) {
	return file_question_proto_rawDescGZIP(), []int{3}
}

func (x *CheckUserAnswerRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CheckUserAnswerRequest) GetQuestionName() string {
	if x != nil {
		return x.QuestionName
	}
	return ""
}

type CheckUserAnswerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Average bool `protobuf:"varint,1,opt,name=average,proto3" json:"average,omitempty"`
}

func (x *CheckUserAnswerResponse) Reset() {
	*x = CheckUserAnswerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckUserAnswerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckUserAnswerResponse) ProtoMessage() {}

func (x *CheckUserAnswerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_question_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckUserAnswerResponse.ProtoReflect.Descriptor instead.
func (*CheckUserAnswerResponse) Descriptor() ([]byte, []int) {
	return file_question_proto_rawDescGZIP(), []int{4}
}

func (x *CheckUserAnswerResponse) GetAverage() bool {
	if x != nil {
		return x.Average
	}
	return false
}

type AverageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Average int32 `protobuf:"varint,1,opt,name=average,proto3" json:"average,omitempty"`
}

func (x *AverageResponse) Reset() {
	*x = AverageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AverageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AverageResponse) ProtoMessage() {}

func (x *AverageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_question_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AverageResponse.ProtoReflect.Descriptor instead.
func (*AverageResponse) Descriptor() ([]byte, []int) {
	return file_question_proto_rawDescGZIP(), []int{5}
}

func (x *AverageResponse) GetAverage() int32 {
	if x != nil {
		return x.Average
	}
	return 0
}

var File_question_proto protoreflect.FileDescriptor

var file_question_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4b, 0x0a, 0x0d, 0x41, 0x6e, 0x73, 0x77, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x72, 0x61,
	0x74, 0x69, 0x6e, 0x67, 0x22, 0x3d, 0x0a, 0x15, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x52,
	0x61, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a,
	0x0d, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74,
	0x69, 0x6e, 0x67, 0x22, 0x43, 0x0a, 0x1d, 0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x65,
	0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x4c, 0x0a, 0x16, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x55, 0x73, 0x65, 0x72, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69,
	0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x33, 0x0a, 0x17, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x55,
	0x73, 0x65, 0x72, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x22, 0x2b, 0x0a, 0x0f, 0x41,
	0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x07, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x32, 0x88, 0x02, 0x0a, 0x0f, 0x51, 0x75, 0x65,
	0x73, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3f, 0x0a, 0x0c,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12, 0x17, 0x2e, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x56, 0x0a,
	0x0f, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x55, 0x73, 0x65, 0x72, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72,
	0x12, 0x20, 0x2e, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x55, 0x73, 0x65, 0x72, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x21, 0x2e, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x55, 0x73, 0x65, 0x72, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5c, 0x0a, 0x16, 0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61,
	0x74, 0x65, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x12,
	0x27, 0x2e, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x61, 0x6c, 0x63, 0x75,
	0x6c, 0x61, 0x74, 0x65, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_question_proto_rawDescOnce sync.Once
	file_question_proto_rawDescData = file_question_proto_rawDesc
)

func file_question_proto_rawDescGZIP() []byte {
	file_question_proto_rawDescOnce.Do(func() {
		file_question_proto_rawDescData = protoimpl.X.CompressGZIP(file_question_proto_rawDescData)
	})
	return file_question_proto_rawDescData
}

var file_question_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_question_proto_goTypes = []interface{}{
	(*AnswerRequest)(nil),                 // 0: question.AnswerRequest
	(*AverageRatingResponse)(nil),         // 1: question.AverageRatingResponse
	(*CalculateAverageRatingRequest)(nil), // 2: question.CalculateAverageRatingRequest
	(*CheckUserAnswerRequest)(nil),        // 3: question.CheckUserAnswerRequest
	(*CheckUserAnswerResponse)(nil),       // 4: question.CheckUserAnswerResponse
	(*AverageResponse)(nil),               // 5: question.AverageResponse
	(*emptypb.Empty)(nil),                 // 6: google.protobuf.Empty
}
var file_question_proto_depIdxs = []int32{
	0, // 0: question.QuestionService.CreateAnswer:input_type -> question.AnswerRequest
	3, // 1: question.QuestionService.CheckUserAnswer:input_type -> question.CheckUserAnswerRequest
	2, // 2: question.QuestionService.CalculateAverageRating:input_type -> question.CalculateAverageRatingRequest
	6, // 3: question.QuestionService.CreateAnswer:output_type -> google.protobuf.Empty
	4, // 4: question.QuestionService.CheckUserAnswer:output_type -> question.CheckUserAnswerResponse
	5, // 5: question.QuestionService.CalculateAverageRating:output_type -> question.AverageResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_question_proto_init() }
func file_question_proto_init() {
	if File_question_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_question_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnswerRequest); i {
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
		file_question_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AverageRatingResponse); i {
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
		file_question_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CalculateAverageRatingRequest); i {
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
		file_question_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckUserAnswerRequest); i {
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
		file_question_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckUserAnswerResponse); i {
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
		file_question_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AverageResponse); i {
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
			RawDescriptor: file_question_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_question_proto_goTypes,
		DependencyIndexes: file_question_proto_depIdxs,
		MessageInfos:      file_question_proto_msgTypes,
	}.Build()
	File_question_proto = out.File
	file_question_proto_rawDesc = nil
	file_question_proto_goTypes = nil
	file_question_proto_depIdxs = nil
}
