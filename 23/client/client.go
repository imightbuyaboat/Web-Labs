package main

import (
	pb "Web-Labs/23/proto"
	"context"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pb.CalculatorClient
}

func NewClient() (*Client, error) {
	clientConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		pb.NewCalculatorClient(clientConn),
	}, nil
}

func (c *Client) Calculate(op string, vars ...string) (float64, error) {
	if len(vars) < 1 {
		return .0, fmt.Errorf("invalid expression")
	}

	var res float64

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch op {
	case "add":
		if len(vars) != 2 {
			return .0, fmt.Errorf("invalid expression")
		}
		x, _ := strconv.ParseFloat(vars[0], 64)
		y, _ := strconv.ParseFloat(vars[1], 64)
		r, err := c.Add(ctx, &pb.AddRequest{X: x, Y: y})
		if err != nil {
			return .0, err
		}
		res = r.Res

	case "sub":
		if len(vars) != 2 {
			return .0, fmt.Errorf("invalid expression")
		}
		x, _ := strconv.ParseFloat(vars[0], 64)
		y, _ := strconv.ParseFloat(vars[1], 64)
		r, err := c.Sub(ctx, &pb.SubRequest{X: x, Y: y})
		if err != nil {
			return .0, err
		}
		res = r.Res

	case "mult":
		if len(vars) != 2 {
			return .0, fmt.Errorf("invalid expression")
		}
		x, _ := strconv.ParseFloat(vars[0], 64)
		y, _ := strconv.ParseFloat(vars[1], 64)
		r, err := c.Mult(ctx, &pb.MultRequest{X: x, Y: y})
		if err != nil {
			return .0, err
		}
		res = r.Res

	case "div":
		if len(vars) != 2 {
			return .0, fmt.Errorf("invalid expression")
		}
		x, _ := strconv.ParseFloat(vars[0], 64)
		y, _ := strconv.ParseFloat(vars[1], 64)
		r, err := c.Div(ctx, &pb.DivRequest{X: x, Y: y})
		if err != nil {
			return .0, err
		}
		res = r.Res

	case "sqrt":
		if len(vars) != 1 {
			return .0, fmt.Errorf("invalid expression")
		}
		x, _ := strconv.ParseFloat(vars[0], 64)
		r, err := c.Sqrt(ctx, &pb.SqrtRequest{X: x})
		if err != nil {
			return .0, err
		}
		res = r.Res

	case "percent":
		if len(vars) != 2 {
			return .0, fmt.Errorf("invalid expression")
		}
		x, _ := strconv.ParseFloat(vars[0], 64)
		y, _ := strconv.ParseFloat(vars[1], 64)
		r, err := c.Percent(ctx, &pb.PercentRequest{X: x, Percent: y})
		if err != nil {
			return .0, err
		}
		res = r.Res

	case "round":
		if len(vars) != 2 {
			return .0, fmt.Errorf("invalid expression")
		}
		x, _ := strconv.ParseFloat(vars[0], 64)
		y, _ := strconv.Atoi(vars[1])
		r, err := c.Round(ctx, &pb.RoundRequest{X: x, Y: int64(y)})
		if err != nil {
			return .0, err
		}
		res = r.Res

	case "pow":
		if len(vars) != 2 {
			return .0, fmt.Errorf("invalid expression")
		}
		x, _ := strconv.ParseFloat(vars[0], 64)
		y, _ := strconv.ParseFloat(vars[1], 64)
		r, err := c.Pow(ctx, &pb.PowRequest{X: x, Y: y})
		if err != nil {
			return .0, err
		}
		res = r.Res

	case "seq":
		if len(vars) < 3 || len(vars)%2 == 0 {
			return .0, fmt.Errorf("invalid expression")
		}
		start, _ := strconv.ParseFloat(vars[0], 64)

		var operations []*pb.Operation
		for i := 1; i < len(vars); i += 2 {
			value, _ := strconv.ParseFloat(vars[i+1], 64)
			operations = append(operations, &pb.Operation{
				Op:    vars[i],
				Value: value,
			})
		}

		r, err := c.Sequence(ctx, &pb.SequenceRequest{
			Start:      start,
			Operations: operations,
		})
		if err != nil {
			return .0, err
		}
		res = r.Res

	default:
		return .0, fmt.Errorf("unknown operations: %s", op)
	}

	return res, nil
}
