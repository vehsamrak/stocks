package main

import "time"

type StableBuyTradeStrategy struct {
	MoneyAmount int // money amount to open trade in dollars
	DayOfTrade  int // day to open trade every month
}

func (s *StableBuyTradeStrategy) Name() string {
	return "stable buy every month"
}

func (s *StableBuyTradeStrategy) Conditions() []TradeCondition {
	return []TradeCondition{
		{
			MoneyAmount: s.MoneyAmount,
			Direction:   buy,
			OpenOnDays: []TradeDay{
				{Month: time.January, Day: s.DayOfTrade},
				{Month: time.February, Day: s.DayOfTrade},
				{Month: time.March, Day: s.DayOfTrade},
				{Month: time.April, Day: s.DayOfTrade},
				{Month: time.May, Day: s.DayOfTrade},
				{Month: time.June, Day: s.DayOfTrade},
				{Month: time.July, Day: s.DayOfTrade},
				{Month: time.August, Day: s.DayOfTrade},
				{Month: time.September, Day: s.DayOfTrade},
				{Month: time.October, Day: s.DayOfTrade},
				{Month: time.November, Day: s.DayOfTrade},
				{Month: time.December, Day: s.DayOfTrade},
			},
		},
	}
}
