package main

type DrawdownTradeStrategy struct {
	MoneyAmount int // money amount to open trade in dollars
	DayOfTrade  int // day to open trade every month
}

func (s *DrawdownTradeStrategy) Name() string {
	return "buy on 2 percent drawdown"
}

func (s *DrawdownTradeStrategy) Conditions() []TradeCondition {
	return []TradeCondition{
		{
			MoneyAmount:         s.MoneyAmount,
			Direction:           buy,
			OpenOnChangePercent: -2,
		},
	}
}
