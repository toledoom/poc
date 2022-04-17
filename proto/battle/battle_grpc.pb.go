// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: proto/battle/battle.proto

package battle

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

// BattleClient is the client API for Battle service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BattleClient interface {
	StartBattle(ctx context.Context, in *StartBattleRequest, opts ...grpc.CallOption) (*StartBattleResponse, error)
	FinishBattle(ctx context.Context, in *FinishBattleRequest, opts ...grpc.CallOption) (*FinishBattleResponse, error)
}

type battleClient struct {
	cc grpc.ClientConnInterface
}

func NewBattleClient(cc grpc.ClientConnInterface) BattleClient {
	return &battleClient{cc}
}

func (c *battleClient) StartBattle(ctx context.Context, in *StartBattleRequest, opts ...grpc.CallOption) (*StartBattleResponse, error) {
	out := new(StartBattleResponse)
	err := c.cc.Invoke(ctx, "/battle.Battle/StartBattle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *battleClient) FinishBattle(ctx context.Context, in *FinishBattleRequest, opts ...grpc.CallOption) (*FinishBattleResponse, error) {
	out := new(FinishBattleResponse)
	err := c.cc.Invoke(ctx, "/battle.Battle/FinishBattle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BattleServer is the server API for Battle service.
// All implementations must embed UnimplementedBattleServer
// for forward compatibility
type BattleServer interface {
	StartBattle(context.Context, *StartBattleRequest) (*StartBattleResponse, error)
	FinishBattle(context.Context, *FinishBattleRequest) (*FinishBattleResponse, error)
	mustEmbedUnimplementedBattleServer()
}

// UnimplementedBattleServer must be embedded to have forward compatible implementations.
type UnimplementedBattleServer struct {
}

func (UnimplementedBattleServer) StartBattle(context.Context, *StartBattleRequest) (*StartBattleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartBattle not implemented")
}
func (UnimplementedBattleServer) FinishBattle(context.Context, *FinishBattleRequest) (*FinishBattleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FinishBattle not implemented")
}
func (UnimplementedBattleServer) mustEmbedUnimplementedBattleServer() {}

// UnsafeBattleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BattleServer will
// result in compilation errors.
type UnsafeBattleServer interface {
	mustEmbedUnimplementedBattleServer()
}

func RegisterBattleServer(s grpc.ServiceRegistrar, srv BattleServer) {
	s.RegisterService(&Battle_ServiceDesc, srv)
}

func _Battle_StartBattle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartBattleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BattleServer).StartBattle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/battle.Battle/StartBattle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BattleServer).StartBattle(ctx, req.(*StartBattleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Battle_FinishBattle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FinishBattleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BattleServer).FinishBattle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/battle.Battle/FinishBattle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BattleServer).FinishBattle(ctx, req.(*FinishBattleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Battle_ServiceDesc is the grpc.ServiceDesc for Battle service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Battle_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "battle.Battle",
	HandlerType: (*BattleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartBattle",
			Handler:    _Battle_StartBattle_Handler,
		},
		{
			MethodName: "FinishBattle",
			Handler:    _Battle_FinishBattle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/battle/battle.proto",
}
