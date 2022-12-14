package mock_test

import (
	"testing"

	"github.com/cool-develope/trade-executor/external/mock"
	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
	"github.com/stretchr/testify/require"
)

func TestMockAPI(t *testing.T) {
	orderBook := make(chan *pb.OrderBook)
	symbol := "BNBUSDT"
	go mock.Subscribe(symbol, orderBook)
	for i := 0; i < 10; i++ {
		ob := <-orderBook
		require.Equal(t, ob.Symbol, symbol)
		require.Greater(t, ob.BidPrice, 0.0)
		require.Greater(t, ob.AskPrice, ob.BidPrice)
		require.Greater(t, ob.BidQty, 0.0)
		require.Greater(t, ob.AskQty, 0.0)
	}
}
