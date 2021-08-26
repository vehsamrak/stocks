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
	tradeDays := make(map[time.Time]int)
	tradeDays[s.createTimeForMonth(time.January)] = s.DayOfTrade
	tradeDays[s.createTimeForMonth(time.February)] = s.DayOfTrade
	tradeDays[s.createTimeForMonth(time.March)] = s.DayOfTrade
	tradeDays[s.createTimeForMonth(time.April)] = s.DayOfTrade
	tradeDays[s.createTimeForMonth(time.May)] = s.DayOfTrade
	tradeDays[s.createTimeForMonth(time.June)] = s.DayOfTrade
	tradeDays[s.createTimeForMonth(time.July)] = s.DayOfTrade
	tradeDays[s.createTimeForMonth(time.August)] = s.DayOfTrade
	tradeDays[s.createTimeForMonth(time.September)] = s.DayOfTrade
	tradeDays[s.createTimeForMonth(time.October)] = s.DayOfTrade
	tradeDays[s.createTimeForMonth(time.November)] = s.DayOfTrade
	tradeDays[s.createTimeForMonth(time.December)] = s.DayOfTrade

	return []TradeCondition{
		{
			MoneyAmount: s.MoneyAmount,
			Direction:   buy,
			OpenOnDays:  tradeDays,
		},
	}
}

func (s *StableBuyTradeStrategy) createTimeForMonth(month time.Month) time.Time {
	return time.Date(
		0,
		month,
		0,
		0,
		0,
		0,
		0,
		time.Now().Location(),
	)
}
