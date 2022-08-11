package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
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

func parseFloat(value string) float64 {
	fValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("can't parse float: %s", value)
		return 0.0
	}
	return fValue
}

// Subscribe opens the websocket connection.
func Subscribe(symbol string, orderBook chan<- *pb.OrderBook, quit <-chan bool) {
	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws/%s@bookTicker", endpoint, symbol), nil)
	if err != nil {
		log.Fatalf("erro websocket dial: %v", err)
		panic(err)
	}

	for {
		select {
		case <-quit:
			return
		default:
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
				BidPrice:    parseFloat(ob.BidPrice),
				BidQty:      parseFloat(ob.BidQty),
				AskPrice:    parseFloat(ob.AskPrice),
				AskQty:      parseFloat(ob.AskQty),
			}
		}
	}
}
