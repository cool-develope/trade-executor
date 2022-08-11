package storage_test

import (
	"testing"

	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
	"github.com/cool-develope/trade-executor/storage"
	"github.com/stretchr/testify/require"
)

const symbol = "BNBUSDT"

func TestStorage(t *testing.T) {
	s, err := storage.NewSqlite3Storage()
	require.NoError(t, err)
	err = s.InitDB()
	require.NoError(t, err)

	var orderBookID = uint64(1000)
	ob := pb.OrderBook{
		OrderBookId: orderBookID,
		Symbol:      symbol,
		AskPrice:    41.5,
		AskQty:      10.0,
		BidPrice:    40.5,
		BidQty:      12,
	}

	err = s.SetOrderBook(&ob)
	require.NoError(t, err)

	gob, err := s.GetLastOrderBook()
	require.NoError(t, err)
	require.Equal(t, gob.Symbol, ob.Symbol)
	require.Equal(t, gob.OrderBookId, ob.OrderBookId)

	ao := pb.Order{
		Symbol:    symbol,
		Price:     40.0,
		Qty:       25.0,
		OrderType: "SELL",
	}

	orderID, err := s.SetAppliedOrder(&ao)
	require.NoError(t, err)
	require.Equal(t, orderID, uint64(1))

	ob.BidPrice = 40.0
	ob.BidQty = 20
	ob.OrderBookId += 1
	err = s.SetOrderBook(&ob)
	require.NoError(t, err)

	pos := []pb.PartialOrder{
		pb.PartialOrder{
			OrderBookId: orderBookID,
			Price:       40.5,
			Qty:         10.0,
		},
		pb.PartialOrder{
			OrderBookId: orderBookID + 1,
			Price:       40.0,
			Qty:         15.0,
		},
	}

	err = s.SetExecutedOrder(&pos[0], orderID)
	require.NoError(t, err)
	err = s.SetExecutedOrder(&pos[1], orderID)
	require.NoError(t, err)

	gpos, err := s.GetExecutedResults(orderID)
	require.NoError(t, err)
	require.Equal(t, len(gpos), len(pos))
	require.Equal(t, gpos[0].OrderBookId, pos[0].OrderBookId)
	require.Equal(t, gpos[1].OrderBookId, pos[1].OrderBookId)
	require.Equal(t, gpos[1].Qty, pos[1].Qty)
}
