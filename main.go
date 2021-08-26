package main

import "fmt"

const (
	initialBalance = 100
)

func main() {
	prices := FetchPrices()
	backtest := &Backtest{}
	strategy := &StableBuyTradeStrategy{
		MoneyAmount: 2000,
		DayOfTrade:  6,
	}

	tradeResult := backtest.Run(initialBalance, strategy, prices)

	fmt.Printf("%#v\n", len(tradeResult.openedTrades))
	fmt.Printf("%#v\n", tradeResult.Profit())
}
