// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.3
// source: casservice.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CasServiceClient is the client API for CasService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CasServiceClient interface {
	CasLogin(ctx context.Context, in *CasLoginReq, opts ...grpc.CallOption) (*CasLoginRes, error)
}

type casServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCasServiceClient(cc grpc.ClientConnInterface) CasServiceClient {
	return &casServiceClient{cc}
}

func (c *casServiceClient) CasLogin(ctx context.Context, in *CasLoginReq, opts ...grpc.CallOption) (*CasLoginRes, error) {
	out := new(CasLoginRes)
	err := c.cc.Invoke(ctx, "/casservice.CasService/casLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CasServiceServer is the server API for CasService service.
// All implementations must embed UnimplementedCasServiceServer
// for forward compatibility
type CasServiceServer interface {
	CasLogin(context.Context, *CasLoginReq) (*CasLoginRes, error)
	mustEmbedUnimplementedCasServiceServer()
}

// UnimplementedCasServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCasServiceServer struct {
}

func (UnimplementedCasServiceServer) CasLogin(context.Context, *CasLoginReq) (*CasLoginRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CasLogin not implemented")
}
func (UnimplementedCasServiceServer) mustEmbedUnimplementedCasServiceServer() {}

// UnsafeCasServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CasServiceServer will
// result in compilation errors.
type UnsafeCasServiceServer interface {
	mustEmbedUnimplementedCasServiceServer()
}

func RegisterCasServiceServer(s grpc.ServiceRegistrar, srv CasServiceServer) {
	s.RegisterService(&CasService_ServiceDesc, srv)
}

func _CasService_CasLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CasLoginReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CasServiceServer).CasLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/casservice.CasService/casLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CasServiceServer).CasLogin(ctx, req.(*CasLoginReq))
	}
	return interceptor(ctx, in, info, handler)
}

// CasService_ServiceDesc is the grpc.ServiceDesc for CasService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CasService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "casservice.CasService",
	HandlerType: (*CasServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "casLogin",
			Handler:    _CasService_CasLogin_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "casservice.proto",
}