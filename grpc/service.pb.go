// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

package grpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type RunRequest struct {
	Ping                 string   `protobuf:"bytes,1,opt,name=ping,proto3" json:"ping,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunRequest) Reset()         { *m = RunRequest{} }
func (m *RunRequest) String() string { return proto.CompactTextString(m) }
func (*RunRequest) ProtoMessage()    {}
func (*RunRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0}
}

func (m *RunRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunRequest.Unmarshal(m, b)
}
func (m *RunRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunRequest.Marshal(b, m, deterministic)
}
func (m *RunRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunRequest.Merge(m, src)
}
func (m *RunRequest) XXX_Size() int {
	return xxx_messageInfo_RunRequest.Size(m)
}
func (m *RunRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RunRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RunRequest proto.InternalMessageInfo

func (m *RunRequest) GetPing() string {
	if m != nil {
		return m.Ping
	}
	return ""
}

type RunResponse struct {
	Pong                 string   `protobuf:"bytes,1,opt,name=pong,proto3" json:"pong,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunResponse) Reset()         { *m = RunResponse{} }
func (m *RunResponse) String() string { return proto.CompactTextString(m) }
func (*RunResponse) ProtoMessage()    {}
func (*RunResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{1}
}

func (m *RunResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunResponse.Unmarshal(m, b)
}
func (m *RunResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunResponse.Marshal(b, m, deterministic)
}
func (m *RunResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunResponse.Merge(m, src)
}
func (m *RunResponse) XXX_Size() int {
	return xxx_messageInfo_RunResponse.Size(m)
}
func (m *RunResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RunResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RunResponse proto.InternalMessageInfo

func (m *RunResponse) GetPong() string {
	if m != nil {
		return m.Pong
	}
	return ""
}

func init() {
	proto.RegisterType((*RunRequest)(nil), "grpc.RunRequest")
	proto.RegisterType((*RunResponse)(nil), "grpc.RunResponse")
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_a0b84a42fa06f626) }

var fileDescriptor_a0b84a42fa06f626 = []byte{
	// 137 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x49, 0x2f, 0x2a, 0x48, 0x56,
	0x52, 0xe0, 0xe2, 0x0a, 0x2a, 0xcd, 0x0b, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x12, 0xe2,
	0x62, 0x29, 0xc8, 0xcc, 0x4b, 0x97, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x95, 0x14,
	0xb9, 0xb8, 0xc1, 0x2a, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0xc1, 0x4a, 0xf2, 0x91, 0x94, 0xe4,
	0xe7, 0xa5, 0x1b, 0x59, 0x70, 0x71, 0x78, 0x65, 0xe5, 0xa7, 0x27, 0xa6, 0xa6, 0x67, 0x0a, 0xe9,
	0x70, 0x31, 0x07, 0x95, 0xe6, 0x09, 0x09, 0xe8, 0x81, 0x8c, 0xd7, 0x43, 0x98, 0x2d, 0x25, 0x88,
	0x24, 0x02, 0x31, 0x4b, 0x89, 0x21, 0x89, 0x0d, 0xec, 0x16, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xaa, 0xcf, 0xc1, 0xfc, 0x9c, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// JjogaegiClient is the client API for Jjogaegi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type JjogaegiClient interface {
	Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (*RunResponse, error)
}

type jjogaegiClient struct {
	cc *grpc.ClientConn
}

func NewJjogaegiClient(cc *grpc.ClientConn) JjogaegiClient {
	return &jjogaegiClient{cc}
}

func (c *jjogaegiClient) Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (*RunResponse, error) {
	out := new(RunResponse)
	err := c.cc.Invoke(ctx, "/grpc.Jjogaegi/Run", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JjogaegiServer is the server API for Jjogaegi service.
type JjogaegiServer interface {
	Run(context.Context, *RunRequest) (*RunResponse, error)
}

// UnimplementedJjogaegiServer can be embedded to have forward compatible implementations.
type UnimplementedJjogaegiServer struct {
}

func (*UnimplementedJjogaegiServer) Run(ctx context.Context, req *RunRequest) (*RunResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Run not implemented")
}

func RegisterJjogaegiServer(s *grpc.Server, srv JjogaegiServer) {
	s.RegisterService(&_Jjogaegi_serviceDesc, srv)
}

func _Jjogaegi_Run_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JjogaegiServer).Run(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Jjogaegi/Run",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JjogaegiServer).Run(ctx, req.(*RunRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Jjogaegi_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Jjogaegi",
	HandlerType: (*JjogaegiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Run",
			Handler:    _Jjogaegi_Run_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
