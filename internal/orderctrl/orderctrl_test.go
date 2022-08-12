package orderctrl

import (
	"sort"
	"testing"

	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
	"github.com/stretchr/testify/require"
)

type MockStorage struct {
	Calls      map[string]int
	orderCount uint64
}

func (m *MockStorage) SetAppliedOrder(ao *pb.Order) (uint64, error) {
	m.Calls["SetAppliedOrder"]++
	m.orderCount++
	return m.orderCount, nil
}

func (m *MockStorage) SetExecutedOrder(po *pb.PartialOrder, orderID uint64) error {
	m.Calls["SetExecutedOrder"]++
	return nil
}

func (m *MockStorage) SetOrderBook(ob *pb.OrderBook) error {
	m.Calls["SetOrderBook"]++
	return nil
}

func TestOrderMatch(t *testing.T) {
	testCases := []*struct {
		name         string
		orders       []*pb.Order
		orderBook    *pb.OrderBook
		expectResult []uint64
		expectErr    bool
		expectCalls  int
	}{
		{
			"invalid order",
			[]*pb.Order{
				{
					Symbol:    "BNBUSDT",
					OrderType: "HOLD",
				},
			},
			&pb.OrderBook{},
			[]uint64{},
			true,
			0,
		},
		{
			"matched order",
			[]*pb.Order{
				{
					Symbol:    "BNBUSDT",
					OrderType: "BUY",
					Qty:       10.0,
					Price:     40.5,
				},
			},
			&pb.OrderBook{
				Symbol:   "BNBUSDT",
				AskPrice: 40.0,
				AskQty:   20.0,
				BidPrice: 39.5,
				BidQty:   10.0,
			},
			[]uint64{1},
			false,
			1,
		},
		{
			"unmatched order",
			[]*pb.Order{
				{
					Symbol:    "BNBUSDT",
					OrderType: "BUY",
					Qty:       10.0,
					Price:     39,
				},
			},
			&pb.OrderBook{
				Symbol:   "BNBUSDT",
				AskPrice: 40.0,
				AskQty:   20.0,
				BidPrice: 39.5,
				BidQty:   10.0,
			},
			[]uint64{},
			false,
			0,
		},
		{
			"partial matched order",
			[]*pb.Order{
				{
					Symbol:    "BNBUSDT",
					OrderType: "BUY",
					Qty:       30.0,
					Price:     40.0,
				},
			},
			&pb.OrderBook{
				Symbol:   "BNBUSDT",
				AskPrice: 40.0,
				AskQty:   20.0,
				BidPrice: 39.5,
				BidQty:   10.0,
			},
			[]uint64{},
			false,
			1,
		},
		{
			"full matched order",
			[]*pb.Order{
				{
					Symbol:    "BNBUSDT",
					OrderType: "BUY",
					Qty:       10.0,
					Price:     40.5,
				},
				{
					Symbol:    "BNBUSDT",
					OrderType: "SELL",
					Qty:       10.0,
					Price:     39.0,
				},
			},
			&pb.OrderBook{
				Symbol:   "BNBUSDT",
				AskPrice: 40.0,
				AskQty:   20.0,
				BidPrice: 39.5,
				BidQty:   10.0,
			},
			[]uint64{1, 2},
			false,
			2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ob := make(chan *pb.OrderBook)
			m := &MockStorage{}
			m.Calls = make(map[string]int)

			ctrl := NewOrderCtrl(m, ob)
			for _, order := range tc.orders {
				_, err := ctrl.SetOrder(order)
				require.NoError(t, err)
			}

			res, err := ctrl.matchOrders(tc.orderBook)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
				require.Equal(t, res, tc.expectResult)
				require.Equal(t, m.Calls["SetExecutedOrder"], tc.expectCalls)
			}
		})
	}
}
