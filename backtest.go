package main

import (
	"fmt"
	"time"
)

type TradeResult struct {
	initialBalance int
	balance        int
	deposits       int
	profit         int
	openedTrades   []Trade
	tradedMonths   map[time.Time]bool
}

type Trade struct {
	Direction   TradeDirection
	Price       Price
	MoneyAmount int
}

func NewTradeResult(initialBalance int) *TradeResult {
	return &TradeResult{
		initialBalance: initialBalance,
		balance:        initialBalance,
		tradedMonths:   make(map[time.Time]bool),
	}
}

func (tr *TradeResult) AddBalance(balanceAddition int) int {
	tr.balance = tr.balance + balanceAddition
	return tr.balance
}

func (tr *TradeResult) Balance() int {
	return tr.balance
}

func (tr *TradeResult) HasOpenedTrades() bool {
	return len(tr.openedTrades) > 0
}

func (tr *TradeResult) OpenTrade(moneyAmount int, direction TradeDirection, price Price) {
	tradeMonth := time.Date(
		price.Date.Year(),
		price.Date.Month(),
		0,
		0,
		0,
		0,
		0,
		time.Now().Location(),
	)

	// TODO[petr]: refactor this to use YEAR+MONTH structure instead of time.Time in map
	_, ok := tr.tradedMonths[tradeMonth]
	if ok {
		return
	}

	tr.tradedMonths[tradeMonth] = true
	tr.openedTrades = append(
		tr.openedTrades, Trade{
			Direction:   direction,
			Price:       price,
			MoneyAmount: moneyAmount,
		},
	)
}

func (tr *TradeResult) CloseTrade() {
	tr.openedTrades = nil
}

func (tr *TradeResult) Profit() int {
	return tr.profit
}

func (tr *TradeResult) Deposits() int {
	return tr.deposits
}

type Backtest struct {
	From time.Time
	Till time.Time
}

func (b *Backtest) Run(initialBalance int, strategy TradeStrategy, prices []Price) *TradeResult {
	tradeResult := NewTradeResult(initialBalance)

	var lastPrice Price
	for _, price := range prices {
		if price.Date.Before(b.From) || price.Date.After(b.Till) {
			continue
		}

		for _, condition := range strategy.Conditions() {
			if tradeResult.HasOpenedTrades() {
				if b.canCloseTrade(price, condition) {
					b.closeTrade(tradeResult)
				}
			}

			if b.canOpenTrade(tradeResult, price, condition) {
				b.openTrade(tradeResult, condition, price)
			}
		}

		lastPrice = price
	}

	b.calculateProfit(tradeResult, lastPrice)

	return tradeResult
}

func (b *Backtest) canOpenTrade(tradeResult *TradeResult, price Price, condition TradeCondition) bool {
	if canOpenTradeOnDays(tradeResult, price, condition.OpenOnDays) {
		return true
	}

	// if canOpenTradeOnChange(price, condition.OpenOnChange) {
	// 	return true
	// }

	// if canOpenTradeOnChangePercent(price, condition.OpenOnChangePercent) {
	// 	return true
	// }

	return false
}

func canOpenTradeOnDays(tradeResult *TradeResult, price Price, tradeDays map[time.Time]int) bool {
	for tradeMonth, tradeDay := range tradeDays {
		_, isMonthAlreadyBeenTraded := tradeResult.tradedMonths[tradeMonth]
		if isMonthAlreadyBeenTraded {
			return false
		}

		if tradeMonth.Month() == price.Date.Month() && tradeDay <= price.Date.Day() {
			return true
		}
	}

	return false
}

func (b *Backtest) canCloseTrade(price Price, condition TradeCondition) bool {
	// TODO[petr]: implement trade closing processing
	return false
}

func (b *Backtest) openTrade(tradeResult *TradeResult, condition TradeCondition, price Price) {
	direction := condition.Direction
	tradeResult.OpenTrade(condition.MoneyAmount, direction, price)
}

func (b *Backtest) closeTrade(tradeResult *TradeResult) {
	// TODO[petr]: implement trade closing
	tradeResult.CloseTrade()
}

func (b *Backtest) calculateProfit(tradeResult *TradeResult, lastPrice Price) {
	fmt.Printf("date | last | price | profit\n")
	for _, trade := range tradeResult.openedTrades {
		tradeResult.deposits = tradeResult.deposits + trade.MoneyAmount
		saldo := lastPrice.Price - trade.Price.Price
		profitPercent := saldo / (trade.Price.Price / 100)
		tradeResult.profit += int(float64(trade.MoneyAmount) * profitPercent / 100)

		fmt.Printf(
			"%s | %d | %d | %.02f%%| %d\n",
			trade.Price.Date.Format("2006-01-02"),
			int(lastPrice.Price),
			int(trade.Price.Price),
			profitPercent,
			tradeResult.profit,
		)
	}

	tradeResult.balance += tradeResult.deposits + tradeResult.profit
}
