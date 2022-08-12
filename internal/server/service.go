package server

import (
	"context"

	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
)

type storageInterface interface {
	GetLastOrderBook() (*pb.OrderBook, error)
	GetExecutedResults(orderID uint64) ([]*pb.PartialOrder, error)
}

type orderCtrlInterface interface {
	SetOrder(ao *pb.Order) (uint64, error)
}

type executorService struct {
	storage   storageInterface
	orderCtrl orderCtrlInterface
	pb.UnsafeExecutorServiceServer
}

// NewExecutorService creates new executor service.
func NewExecutorService(storage storageInterface, orderCtrl orderCtrlInterface) pb.ExecutorServiceServer {
	return &executorService{
		storage:   storage,
		orderCtrl: orderCtrl,
	}
}

// Get the best price for the dedicated symbol.
func (e *executorService) GetPrice(ctx context.Context, req *pb.GetPriceRequest) (*pb.GetPriceResponse, error) {
	ob, err := e.storage.GetLastOrderBook()

	return &pb.GetPriceResponse{
		Price: ob,
	}, err
}

// Apply new order.
func (e *executorService) SetOrder(ctx context.Context, req *pb.SetOrderRequest) (*pb.SetOrderRespnose, error) {
	orderID, err := e.orderCtrl.SetOrder(req.Order)

	return &pb.SetOrderRespnose{
		OrderId: orderID,
	}, err
}

// Get the order result for the order id.
func (e *executorService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	pos, err := e.storage.GetExecutedResults(req.OrderId)
	return &pb.GetOrderResponse{
		Result: pos,
	}, err
}
