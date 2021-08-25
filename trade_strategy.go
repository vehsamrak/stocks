package main

import "time"

type TradeDirection string

type TradeDay struct {
	Month time.Month
	Day   int
}

const (
	buy  TradeDirection = "buy"
	sell TradeDirection = "sell"
)

type TradeCondition struct {
	MoneyAmount         int            // amount of money in dollars to trade
	MoneyAmountPercent  float64        // amount of money in percents to trade
	Direction           TradeDirection // trade direction: "buy" or "sell"
	OpenOnChange        int            // open trade when price changes to certain dollars
	OpenOnChangePercent float64        // open trade when price changes to certain percents
	OpenOnDays          []TradeDay     // days to open trades
	TakeProfit          int            // close trade when profit in dollars reached
	TakeProfitPercent   float64        // close trade when profit in percents reached
	StopLoss            int            // close trade when loss in dollars reached
	StopLossPercent     float64        // close trade when loss in percents reached
}

type TradeStrategy interface {
	Name() string
	Conditions() []TradeCondition
}
