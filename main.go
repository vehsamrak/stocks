package main

import (
	"fmt"
	"time"
)

const (
	initialBalance = 23000
	fromDate       = "2020-08-01"
)

func main() {
	prices := FetchPrices()
	fromTime, _ := time.Parse("2006-01-02", fromDate)
	backtest := &Backtest{From: fromTime}
	strategy := &StableBuyTradeStrategy{
		MoneyAmount: 2000,
		DayOfTrade:  6,
	}

	tradeResult := backtest.Run(initialBalance, strategy, prices)

	fmt.Printf("%#v - opened trades\n", len(tradeResult.openedTrades))
	fmt.Printf("%#v - deposits\n", tradeResult.Deposits())
	profit := tradeResult.Profit()
	profitPercent := profit / (tradeResult.Deposits() / 100)
	fmt.Printf("%#v - profit (%d%%)\n", profit, profitPercent)
	fmt.Printf("%#v - end capital\n", tradeResult.Balance())
}
