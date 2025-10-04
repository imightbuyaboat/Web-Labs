package main

import (
	pb "Web-Labs/23/proto"
	"context"
	"fmt"
)

type Server struct {
	pb.UnimplementedCalculatorServer
}

func (s *Server) Add(ctx context.Context, r *pb.AddRequest) (*pb.AddResponse, error) {
	return &pb.AddResponse{
		Res: Add(r.X, r.Y),
	}, nil
}

func (s *Server) Sub(ctx context.Context, r *pb.SubRequest) (*pb.SubResponse, error) {
	return &pb.SubResponse{
		Res: Sub(r.X, r.Y),
	}, nil
}

func (s *Server) Mult(ctx context.Context, r *pb.MultRequest) (*pb.MultResponse, error) {
	return &pb.MultResponse{
		Res: Mult(r.X, r.Y),
	}, nil
}

func (s *Server) Div(ctx context.Context, r *pb.DivRequest) (*pb.DivResponse, error) {
	res, err := Div(r.X, r.Y)
	if err != nil {
		return nil, err
	}

	return &pb.DivResponse{
		Res: res,
	}, nil
}

func (s *Server) Sqrt(ctx context.Context, r *pb.SqrtRequest) (*pb.SqrtResponse, error) {
	res, err := Sqrt(r.X)
	if err != nil {
		return nil, err
	}

	return &pb.SqrtResponse{
		Res: res,
	}, nil
}

func (s *Server) Percent(ctx context.Context, r *pb.PercentRequest) (*pb.PercentResponse, error) {
	res, err := Percent(r.X, r.Percent)
	if err != nil {
		return nil, err
	}

	return &pb.PercentResponse{
		Res: res,
	}, nil
}

func (s *Server) Round(ctx context.Context, r *pb.RoundRequest) (*pb.RoundResponse, error) {
	return &pb.RoundResponse{
		Res: Round(r.X, r.Y),
	}, nil
}

func (s *Server) Pow(ctx context.Context, r *pb.PowRequest) (*pb.PowResponse, error) {
	return &pb.PowResponse{
		Res: Pow(r.X, r.Y),
	}, nil
}

func (s *Server) Sequence(ctx context.Context, r *pb.SequenceRequest) (*pb.SequenceResponse, error) {
	res := r.Start
	var err error

	for _, op := range r.Operations {
		switch op.Op {
		case "+":
			res = Add(res, op.Value)

		case "-":
			res = Sub(res, op.Value)

		case "*":
			res = Mult(res, op.Value)

		case "/":
			res, err = Div(res, op.Value)
			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("unknown operation: %s", op.Op)
		}
	}

	return &pb.SequenceResponse{
		Res: res,
	}, nil
}
