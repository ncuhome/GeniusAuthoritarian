// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: ssh.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	SshAccounts_Watch_FullMethodName = "/proto.SshAccounts/Watch"
)

// SshAccountsClient is the client API for SshAccounts service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SshAccountsClient interface {
	Watch(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (SshAccounts_WatchClient, error)
}

type sshAccountsClient struct {
	cc grpc.ClientConnInterface
}

func NewSshAccountsClient(cc grpc.ClientConnInterface) SshAccountsClient {
	return &sshAccountsClient{cc}
}

func (c *sshAccountsClient) Watch(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (SshAccounts_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &SshAccounts_ServiceDesc.Streams[0], SshAccounts_Watch_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &sshAccountsWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SshAccounts_WatchClient interface {
	Recv() (*AccountStream, error)
	grpc.ClientStream
}

type sshAccountsWatchClient struct {
	grpc.ClientStream
}

func (x *sshAccountsWatchClient) Recv() (*AccountStream, error) {
	m := new(AccountStream)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SshAccountsServer is the server API for SshAccounts service.
// All implementations must embed UnimplementedSshAccountsServer
// for forward compatibility
type SshAccountsServer interface {
	Watch(*emptypb.Empty, SshAccounts_WatchServer) error
	mustEmbedUnimplementedSshAccountsServer()
}

// UnimplementedSshAccountsServer must be embedded to have forward compatible implementations.
type UnimplementedSshAccountsServer struct {
}

func (UnimplementedSshAccountsServer) Watch(*emptypb.Empty, SshAccounts_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedSshAccountsServer) mustEmbedUnimplementedSshAccountsServer() {}

// UnsafeSshAccountsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SshAccountsServer will
// result in compilation errors.
type UnsafeSshAccountsServer interface {
	mustEmbedUnimplementedSshAccountsServer()
}

func RegisterSshAccountsServer(s grpc.ServiceRegistrar, srv SshAccountsServer) {
	s.RegisterService(&SshAccounts_ServiceDesc, srv)
}

func _SshAccounts_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SshAccountsServer).Watch(m, &sshAccountsWatchServer{stream})
}

type SshAccounts_WatchServer interface {
	Send(*AccountStream) error
	grpc.ServerStream
}

type sshAccountsWatchServer struct {
	grpc.ServerStream
}

func (x *sshAccountsWatchServer) Send(m *AccountStream) error {
	return x.ServerStream.SendMsg(m)
}

// SshAccounts_ServiceDesc is the grpc.ServiceDesc for SshAccounts service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SshAccounts_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.SshAccounts",
	HandlerType: (*SshAccountsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _SshAccounts_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "ssh.proto",
}
