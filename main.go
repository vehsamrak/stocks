package main

import (
	"fmt"
	"time"
)

const (
	initialBalance = 100
	fromDate       = "2011-08-01"
	tillDate       = "2022-01-01"
)

func main() {
	prices := FetchPrices()
	fromTime, _ := time.Parse("2006-01-02", fromDate)
	tillTime, _ := time.Parse("2006-01-02", tillDate)
	backtest := &Backtest{From: fromTime, Till: tillTime}
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
