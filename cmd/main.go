package main

import "github.com/urfave/cli/v2"

const (
	// App name
	appName = "trade-executor"
	// version represents the program based on the git tag
	version = "v0.1.0"
)

const (
	flagSymbol    = "symbol"
	flagOrderType = "type"
	flagAmount    = "amount"
	flagPrice     = "price"
	flagOrderID   = "id"
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = version

	orderApplyFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     flagSymbol,
			Aliases:  []string{"s"},
			Usage:    "Symbol Name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     flagOrderType,
			Aliases:  []string{"t"},
			Usage:    "Order Type: buy, sell.",
			Required: true,
		},
		&cli.Float64Flag{
			Name:     flagAmount,
			Aliases:  []string{"a"},
			Usage:    "Order Amount",
			Required: true,
		},
		&cli.Float64Flag{
			Name:     flagPrice,
			Aliases:  []string{"p"},
			Usage:    "Order Price",
			Required: true,
		},
	}

	priceFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     flagSymbol,
			Aliases:  []string{"s"},
			Usage:    "Symbol Name",
			Required: true,
		},
	}

	orderResultFlags := []cli.Flag{
		&cli.Int64Flag{
			Name:     flagOrderID,
			Aliases:  []string{"i"},
			Usage:    "Order ID",
			Required: true,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "serve",
			Usage:  "Serve the executor service",
			Action: serve,
		},
		{
			Name:  "order",
			Usage: "Apply an order or get the response",
			Subcommands: []*cli.Command{
				{
					Name:   "apply",
					Usage:  "Execute an order",
					Action: orderApply,
					Flags:  orderApplyFlags,
				},
				{
					Name:   "result",
					Usage:  "Get the order executed response",
					Action: getOrder,
					Flags:  orderResultFlags,
				},
			},
		},
		{
			Name:   "price",
			Usage:  "Get the symbol price",
			Action: getPrice,
			Flags:  priceFlags,
		},
	}
}
