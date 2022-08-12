package main

import (
	"github.com/cool-develope/trade-executor/config"
	"github.com/cool-develope/trade-executor/external/binance"
	"github.com/cool-develope/trade-executor/internal/orderctrl"
	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
	"github.com/cool-develope/trade-executor/internal/server"
	"github.com/cool-develope/trade-executor/storage"
	"github.com/urfave/cli/v2"
)

func serve(ctx *cli.Context) error {
	configFilePath := ctx.String(flagCfg)
	c, err := config.Load(configFilePath)
	if err != nil {
		return err
	}

	storage, err := storage.NewSqlite3Storage()
	if err != nil {
		return err
	}

	orderBook := make(chan *pb.OrderBook)

	go registerExchange(orderBook, c.Exchange)

	orderCtrl := orderctrl.NewOrderCtrl(storage, orderBook)

	go func() {
		err = orderCtrl.Execute()
	}()
	if err != nil {
		return err
	}

	return server.RunServer(storage, orderCtrl, c.Server)
}

func registerExchange(orderBook chan *pb.OrderBook, config config.ExchangeConfig) {
	if config.Name == "binance" {
		for _, symbol := range config.Symbols {
			binance.Subscribe(symbol, orderBook)
		}
	}
}
