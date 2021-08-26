package main

import "time"

type TradeResult struct {
	initialBalance int
	balance        int
	openedTrades   []Trade
	tradedMonths   map[time.Time]bool
}

type Trade struct {
	Direction TradeDirection
	Price     Price
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

func (tr *TradeResult) OpenTrade(direction TradeDirection, price Price) {
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
			Direction: direction,
			Price:     price,
		},
	)
}

func (tr *TradeResult) CloseTrade() {
	tr.openedTrades = nil
}

func (tr *TradeResult) Profit() int {
	return tr.balance - tr.initialBalance
}

type Backtest struct {
}

func (b *Backtest) Run(initialBalance int, strategy TradeStrategy, prices []Price) *TradeResult {
	tradeResult := NewTradeResult(initialBalance)

	for _, price := range prices {
		for _, condition := range strategy.Conditions() {
			if tradeResult.HasOpenedTrades() {
				if b.canCloseTrade(price, condition) {
					b.closeTrade(tradeResult)
				}
			}

			if b.canOpenTrade(tradeResult, price, condition) {
				b.openTrade(tradeResult, condition.Direction, price)
			}
		}
	}

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

func (b *Backtest) openTrade(tradeResult *TradeResult, direction TradeDirection, price Price) {
	tradeResult.OpenTrade(direction, price)
}

func (b *Backtest) closeTrade(tradeResult *TradeResult) {
	// TODO[petr]: implement trade closing
	tradeResult.CloseTrade()
}
