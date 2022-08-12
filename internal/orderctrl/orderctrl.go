package orderctrl

import (
	"fmt"

	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
)

const (
	buyType  = "BUY"
	sellType = "SELL"
	eps      = 1e-10
)

type storageInterface interface {
	SetOrderBook(ob *pb.OrderBook) error
	SetAppliedOrder(ao *pb.Order) (uint64, error)
	SetExecutedOrder(po *pb.PartialOrder, orderID uint64) error
}

// OrderCtrl is a controller to execute orders.
type OrderCtrl struct {
	storage   storageInterface
	orderBook <-chan *pb.OrderBook
	orderPool map[uint64]*pb.Order
}

// NewOrderCtrl creates new order controller.
func NewOrderCtrl(storage storageInterface, orderBook chan *pb.OrderBook) *OrderCtrl {
	return &OrderCtrl{
		storage:   storage,
		orderBook: orderBook,
		orderPool: make(map[uint64]*pb.Order),
	}
}

// SetOrder adds new order to the pool.
func (o *OrderCtrl) SetOrder(ao *pb.Order) (uint64, error) {
	orderID, err := o.storage.SetAppliedOrder(ao)
	if err != nil {
		return orderID, err
	}
	o.orderPool[orderID] = ao

	return orderID, nil
}

func (o *OrderCtrl) matchOrders(ob *pb.OrderBook) ([]uint64, error) {
	executedOrders := make([]uint64, 0)
	for orderID, order := range o.orderPool {
		executedAmount := 0.0
		executedPrice := 0.0
		if order.OrderType == buyType {
			if ob.AskPrice <= order.Price {
				executedAmount = ob.AskQty
				if ob.AskQty > order.Qty {
					executedAmount = order.Qty
				}
				ob.AskQty -= executedAmount
				executedPrice = ob.AskPrice
			}
		} else if order.OrderType == sellType {
			if ob.BidPrice >= order.Price {
				executedAmount = ob.BidQty
				if ob.BidQty > order.Qty {
					executedAmount = order.Qty
				}
				ob.AskQty -= executedAmount
				executedPrice = ob.BidPrice
			}
		} else {
			return nil, fmt.Errorf("error unregistered order type: %d %v", orderID, order)
		}

		if executedAmount > eps {
			err := o.storage.SetExecutedOrder(&pb.PartialOrder{
				OrderBookId: ob.OrderBookId,
				Price:       executedPrice,
				Qty:         executedAmount,
			}, orderID)
			if err != nil {
				return nil, err
			}
		}
		if order.Qty <= executedAmount {
			executedOrders = append(executedOrders, orderID)
		}
		order.Qty -= executedAmount
		o.orderPool[orderID] = order
	}

	return executedOrders, nil
}

// Execute runs the executor.
func (o *OrderCtrl) Execute() error {
	for {
		ob := <-o.orderBook
		err := o.storage.SetOrderBook(ob)
		if err != nil {
			return err
		}

		executedOrders, err := o.matchOrders(ob)
		if err != nil {
			return err
		}

		for _, orderID := range executedOrders {
			delete(o.orderPool, orderID)
		}
	}
}
