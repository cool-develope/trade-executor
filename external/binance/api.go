package binance

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
	"github.com/cool-develope/trade-executor/utils"
	"github.com/gorilla/websocket"
)

const (
	endpoint = "wss://stream.binance.com:9443"
)

// OrderBook struct
type OrderBook struct {
	OrderBookID uint64 `json:"u,omitempty"`
	Symbol      string `json:"s,omitempty"`
	BidPrice    string `json:"b,omitempty"`
	BidQty      string `json:"B,omitempty"`
	AskPrice    string `json:"a,omitempty"`
	AskQty      string `json:"A,omitempty"`
}

// Subscribe opens the websocket connection.
func Subscribe(symbol string, orderBook chan<- *pb.OrderBook) {
	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws/%s@bookTicker", endpoint, symbol), nil)
	if err != nil {
		log.Fatalf("erro websocket dial: %v", err)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatalf("error read message: %v", err)
			break
		}
		var ob OrderBook

		err = json.Unmarshal(message, &ob)
		if err != nil {
			log.Fatalf("error json unmarshal: %v", err)
			break
		}

		orderBook <- &pb.OrderBook{
			OrderBookId: ob.OrderBookID,
			Symbol:      ob.Symbol,
			BidPrice:    utils.ParseFloat(ob.BidPrice),
			BidQty:      utils.ParseFloat(ob.BidQty),
			AskPrice:    utils.ParseFloat(ob.AskPrice),
			AskQty:      utils.ParseFloat(ob.AskQty),
		}
	}
}
