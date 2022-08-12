package mock

import (
	"math/rand"
	"time"

	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
)

const (
	priceSeed = 100.0
	qtySeed   = 10.0
	noise     = 0.1
	diff      = 0.5
)

var bookID = 1000000

func generateOrderBook(symbol string) *pb.OrderBook {
	dev := (rand.Float64() - 0.5) * noise * priceSeed
	askQty := qtySeed * (1 + rand.Float64()*noise)
	bidQty := qtySeed * (1 + rand.Float64()*noise)
	bookID++

	return &pb.OrderBook{
		OrderBookId: uint64(bookID),
		Symbol:      symbol,
		BidPrice:    priceSeed + dev - diff,
		BidQty:      bidQty,
		AskPrice:    priceSeed + dev,
		AskQty:      askQty,
	}
}

// Subscribe opens the websocket connection.
func Subscribe(symbol string, orderBook chan<- *pb.OrderBook) {
	for {
		orderBook <- generateOrderBook(symbol)
		time.Sleep(100 * time.Millisecond)
	}
}
