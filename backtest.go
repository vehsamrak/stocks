package main

type TradeResult struct {
	initialBalance int
	balance        int
}

func NewTradeResult(initialBalance int) *TradeResult {
	return &TradeResult{
		initialBalance: initialBalance,
		balance:        initialBalance,
	}
}

func (tr *TradeResult) AddBalance(balanceAddition int) int {
	tr.balance = tr.balance + balanceAddition
	return tr.balance
}

func (tr *TradeResult) Balance() int {
	return tr.balance
}

type Backtest struct {
}

func (b *Backtest) Calculate(initialBalance int, prices []Price) *TradeResult {
	tradeResult := NewTradeResult(initialBalance)

	return tradeResult
}
