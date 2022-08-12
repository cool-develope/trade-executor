package main

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

const dbPath = "./storage.db"

// e2e tests
func TestE2E(t *testing.T) {
	var (
		pFlagCfg, pFlagSymbol, pFlagOrderType string
		pFlagAmount, pFlagPrice               float64
		PFlagOrderID                          int64
	)

	if _, err := os.Stat(dbPath); err == nil {
		os.Remove(dbPath) //nolint
	} else {
		t.Log(err)
	}

	flagSet := flag.NewFlagSet("serve", 0)
	flagSet.StringVar(&pFlagCfg, flagCfg, "../config/configlocal.toml", "")
	ctx := cli.NewContext(cli.NewApp(), flagSet, nil)
	err := serve(ctx)
	require.NoError(t, err)

	time.Sleep(time.Second)

	rescueStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	flagSet = flag.NewFlagSet("get price", 0)
	flagSet.StringVar(&pFlagSymbol, flagSymbol, "BNBUSDT", "")
	ctx = cli.NewContext(cli.NewApp(), flagSet, nil)
	err = getPrice(ctx)
	require.NoError(t, err)

	w.Close() //nolint
	out, err := ioutil.ReadAll(r)
	require.NoError(t, err)

	orderBookResult := string(out)
	require.Contains(t, orderBookResult, "Order Book:")
	t.Log(orderBookResult)

	r, w, err = os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	flagSet = flag.NewFlagSet("apply order", 0)
	flagSet.StringVar(&pFlagSymbol, flagSymbol, "BNBUSDT", "")
	flagSet.StringVar(&pFlagOrderType, flagOrderType, "BUY", "")
	flagSet.Float64Var(&pFlagAmount, flagAmount, 55.0, "")
	flagSet.Float64Var(&pFlagPrice, flagPrice, 305.0, "")
	ctx = cli.NewContext(cli.NewApp(), flagSet, nil)
	err = orderApply(ctx)
	require.NoError(t, err)

	w.Close() //nolint
	out, err = ioutil.ReadAll(r)
	require.NoError(t, err)

	orderApplyResult := string(out)
	require.Contains(t, orderApplyResult, "Order received:")
	t.Log(orderApplyResult)

	time.Sleep(time.Second)

	r, w, err = os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	flagSet = flag.NewFlagSet("get order", 0)
	flagSet.Int64Var(&PFlagOrderID, flagOrderID, 1, "")
	ctx = cli.NewContext(cli.NewApp(), flagSet, nil)
	err = getOrder(ctx)
	require.NoError(t, err)

	w.Close() //nolint
	out, err = ioutil.ReadAll(r)
	require.NoError(t, err)

	orderExecutedResult := string(out)
	require.Contains(t, orderExecutedResult, "Executed results:")
	t.Log(orderExecutedResult)

	os.Stdout = rescueStdout
}
