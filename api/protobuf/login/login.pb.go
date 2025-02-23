// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: login.proto

package loginpb

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

type ReqRegisterUser struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Email           string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Username        string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Password        string                 `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	ConfirmPassword string                 `protobuf:"bytes,4,opt,name=confirm_password,json=confirmPassword,proto3" json:"confirm_password,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *ReqRegisterUser) Reset() {
	*x = ReqRegisterUser{}
	mi := &file_login_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReqRegisterUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqRegisterUser) ProtoMessage() {}

func (x *ReqRegisterUser) ProtoReflect() protoreflect.Message {
	mi := &file_login_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqRegisterUser.ProtoReflect.Descriptor instead.
func (*ReqRegisterUser) Descriptor() ([]byte, []int) {
	return file_login_proto_rawDescGZIP(), []int{0}
}

func (x *ReqRegisterUser) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *ReqRegisterUser) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *ReqRegisterUser) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *ReqRegisterUser) GetConfirmPassword() string {
	if x != nil {
		return x.ConfirmPassword
	}
	return ""
}

type RspRegisterUser struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ErrorCode     int32                  `protobuf:"varint,1,opt,name=error_code,json=errorCode,proto3" json:"error_code,omitempty"`
	UserId        int32                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RspRegisterUser) Reset() {
	*x = RspRegisterUser{}
	mi := &file_login_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RspRegisterUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspRegisterUser) ProtoMessage() {}

func (x *RspRegisterUser) ProtoReflect() protoreflect.Message {
	mi := &file_login_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspRegisterUser.ProtoReflect.Descriptor instead.
func (*RspRegisterUser) Descriptor() ([]byte, []int) {
	return file_login_proto_rawDescGZIP(), []int{1}
}

func (x *RspRegisterUser) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *RspRegisterUser) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type ReqLoginUser struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReqLoginUser) Reset() {
	*x = ReqLoginUser{}
	mi := &file_login_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReqLoginUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqLoginUser) ProtoMessage() {}

func (x *ReqLoginUser) ProtoReflect() protoreflect.Message {
	mi := &file_login_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqLoginUser.ProtoReflect.Descriptor instead.
func (*ReqLoginUser) Descriptor() ([]byte, []int) {
	return file_login_proto_rawDescGZIP(), []int{2}
}

func (x *ReqLoginUser) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *ReqLoginUser) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type RspLoginUser struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ErrorCode     int32                  `protobuf:"varint,1,opt,name=error_code,json=errorCode,proto3" json:"error_code,omitempty"`
	AccessToken   string                 `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken  string                 `protobuf:"bytes,3,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RspLoginUser) Reset() {
	*x = RspLoginUser{}
	mi := &file_login_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RspLoginUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspLoginUser) ProtoMessage() {}

func (x *RspLoginUser) ProtoReflect() protoreflect.Message {
	mi := &file_login_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspLoginUser.ProtoReflect.Descriptor instead.
func (*RspLoginUser) Descriptor() ([]byte, []int) {
	return file_login_proto_rawDescGZIP(), []int{3}
}

func (x *RspLoginUser) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *RspLoginUser) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *RspLoginUser) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

type ReqUpdateUser struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        uint32                 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Password      string                 `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	Age           uint32                 `protobuf:"varint,4,opt,name=age,proto3" json:"age,omitempty"`
	PhoneNumber   string                 `protobuf:"bytes,5,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	Address       string                 `protobuf:"bytes,6,opt,name=address,proto3" json:"address,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReqUpdateUser) Reset() {
	*x = ReqUpdateUser{}
	mi := &file_login_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReqUpdateUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqUpdateUser) ProtoMessage() {}

func (x *ReqUpdateUser) ProtoReflect() protoreflect.Message {
	mi := &file_login_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqUpdateUser.ProtoReflect.Descriptor instead.
func (*ReqUpdateUser) Descriptor() ([]byte, []int) {
	return file_login_proto_rawDescGZIP(), []int{4}
}

func (x *ReqUpdateUser) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ReqUpdateUser) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *ReqUpdateUser) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *ReqUpdateUser) GetAge() uint32 {
	if x != nil {
		return x.Age
	}
	return 0
}

func (x *ReqUpdateUser) GetPhoneNumber() string {
	if x != nil {
		return x.PhoneNumber
	}
	return ""
}

func (x *ReqUpdateUser) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type RspUpdateUser struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ErrorCode     int32                  `protobuf:"varint,1,opt,name=error_code,json=errorCode,proto3" json:"error_code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RspUpdateUser) Reset() {
	*x = RspUpdateUser{}
	mi := &file_login_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RspUpdateUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspUpdateUser) ProtoMessage() {}

func (x *RspUpdateUser) ProtoReflect() protoreflect.Message {
	mi := &file_login_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspUpdateUser.ProtoReflect.Descriptor instead.
func (*RspUpdateUser) Descriptor() ([]byte, []int) {
	return file_login_proto_rawDescGZIP(), []int{5}
}

func (x *RspUpdateUser) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

var File_login_proto protoreflect.FileDescriptor

var file_login_proto_rawDesc = string([]byte{
	0x0a, 0x0b, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6c,
	0x6f, 0x67, 0x69, 0x6e, 0x22, 0x8a, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x71, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x29, 0x0a, 0x10, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72,
	0x6d, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x22, 0x49, 0x0a, 0x0f, 0x52, 0x73, 0x70, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x55, 0x73, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x40, 0x0a, 0x0c,
	0x52, 0x65, 0x71, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x75,
	0x0a, 0x0c, 0x52, 0x73, 0x70, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1d,
	0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x21, 0x0a,
	0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0xaf, 0x01, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x67, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x61, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x68,
	0x6f, 0x6e, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x18, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x2e, 0x0a, 0x0d, 0x52, 0x73, 0x70, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x32, 0xc5, 0x01, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x40, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x6e,
	0x2e, 0x52, 0x65, 0x71, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72,
	0x1a, 0x16, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x2e, 0x52, 0x73, 0x70, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x09, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x13, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x2e,
	0x52, 0x65, 0x71, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x13, 0x2e, 0x6c,
	0x6f, 0x67, 0x69, 0x6e, 0x2e, 0x52, 0x73, 0x70, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65,
	0x72, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65,
	0x72, 0x12, 0x14, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x14, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x2e,
	0x52, 0x73, 0x70, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x22, 0x00, 0x42,
	0x1e, 0x5a, 0x1c, 0x2e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x3b, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_login_proto_rawDescOnce sync.Once
	file_login_proto_rawDescData []byte
)

func file_login_proto_rawDescGZIP() []byte {
	file_login_proto_rawDescOnce.Do(func() {
		file_login_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_login_proto_rawDesc), len(file_login_proto_rawDesc)))
	})
	return file_login_proto_rawDescData
}

var file_login_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_login_proto_goTypes = []any{
	(*ReqRegisterUser)(nil), // 0: login.ReqRegisterUser
	(*RspRegisterUser)(nil), // 1: login.RspRegisterUser
	(*ReqLoginUser)(nil),    // 2: login.ReqLoginUser
	(*RspLoginUser)(nil),    // 3: login.RspLoginUser
	(*ReqUpdateUser)(nil),   // 4: login.ReqUpdateUser
	(*RspUpdateUser)(nil),   // 5: login.RspUpdateUser
}
var file_login_proto_depIdxs = []int32{
	0, // 0: login.LoginService.RegisterUser:input_type -> login.ReqRegisterUser
	2, // 1: login.LoginService.LoginUser:input_type -> login.ReqLoginUser
	4, // 2: login.LoginService.UpdateUser:input_type -> login.ReqUpdateUser
	1, // 3: login.LoginService.RegisterUser:output_type -> login.RspRegisterUser
	3, // 4: login.LoginService.LoginUser:output_type -> login.RspLoginUser
	5, // 5: login.LoginService.UpdateUser:output_type -> login.RspUpdateUser
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_login_proto_init() }
func file_login_proto_init() {
	if File_login_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_login_proto_rawDesc), len(file_login_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_login_proto_goTypes,
		DependencyIndexes: file_login_proto_depIdxs,
		MessageInfos:      file_login_proto_msgTypes,
	}.Build()
	File_login_proto = out.File
	file_login_proto_goTypes = nil
	file_login_proto_depIdxs = nil
}
