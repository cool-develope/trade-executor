package main

import (
	"fmt"
	"os"

	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultServerURL = "http://127.0.0.1:9090"
	envName          = "EXECUTOR_SERVER_URL"
)

func newExecutorClient() pb.ExecutorServiceClient {
	serverURL := getServerURL()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	proverConn, err := grpc.Dial(serverURL, opts...)
	if err != nil {
		panic(fmt.Errorf("fail to dial: %v", err))
	}

	return pb.NewExecutorServiceClient(proverConn)
}

func getServerURL() string {
	serverURL := os.Getenv(envName)
	if len(serverURL) == 0 {
		serverURL = defaultServerURL
	}
	return serverURL
}

func orderApply(ctx *cli.Context) error {
	client := newExecutorClient()

	order := &pb.Order{
		Symbol:    ctx.String(flagSymbol),
		OrderType: ctx.String(flagOrderType),
		Qty:       ctx.Float64(flagAmount),
		Price:     ctx.Float64(flagPrice),
	}
	res, err := client.SetOrder(ctx.Context, &pb.SetOrderRequest{Order: order})
	if err != nil {
		return err
	}

	fmt.Printf("Order received: %d", res.OrderId)
	return nil
}

func getOrder(ctx *cli.Context) error {
	client := newExecutorClient()

	orderID := ctx.Int64(flagOrderID)

	res, err := client.GetOrder(ctx.Context, &pb.GetOrderRequest{OrderId: uint64(orderID)})
	if err != nil {
		return err
	}

	fmt.Printf("Executed results: %v", res.Result)
	return nil
}

func getPrice(ctx *cli.Context) error {
	client := newExecutorClient()

	symbol := ctx.String(flagSymbol)

	res, err := client.GetPrice(ctx.Context, &pb.GetPriceRequest{Symbol: symbol})
	if err != nil {
		return err
	}

	fmt.Printf("Order Book: %v", res.Price)
	return nil
}
