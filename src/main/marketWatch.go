package main

import (
	"context"
	"fmt"
	"go-contacts/src/models"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"

	"github.com/adshao/go-binance"
)

type MarketOrder struct {
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"`
	T        string  `json:"type"`
	Quantity float64 `json:"quantity"`
}
type Platform struct {
	Symbols []Symbol
}

type Symbol struct {
	Symbol string `json:"symbol"`
	Status string `json:"status"`
}
type CurrentIndicator struct {
	Symbol string  `json:"symbol"`
	Value  float64 `json:"value"`
}

func ExcecuteMarketOrder(symbol string, side string, quantity string) {
	var (
		apiKey    = "9cBb2iRgQOKzQMH9PxWALRXuwU6Ey0mNBKKzHmZUO9cLUTgCd7KtgTKLP4CDYJnR"
		secretKey = "pBHk1F2vSUfIkCVQQMd4ds5suA93tm5wPj3jPCmXSHulhy2rhHp0pFSoEWdvIrE4"
	)
	marketOrderSide := binance.SideTypeBuy
	if side == "buy" {
		marketOrderSide = binance.SideTypeBuy
	} else if side == "sell" {
		marketOrderSide = binance.SideTypeSell
	}

	client := binance.NewClient(apiKey, secretKey)
	order, err := client.NewCreateOrderService().Symbol(symbol).TimeInForce(binance.TimeInForceGTC).Price("0.0255").
		Side(marketOrderSide).Type(binance.OrderTypeLimit).Quantity(quantity).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Order API response")
	fmt.Println(order)
}

func SetInterval(someFunc func(), milliseconds int, async bool) chan bool {

	// How often to fire the passed in function
	// in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	// Put the selection in a go routine
	// so that the for loop is none blocking
	go func() {
		for {

			select {
			case <-ticker.C:
				if async {
					// This won't block
					go someFunc()
				} else {
					// This will block
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				return
			}

		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clear

}

func UpdateRSI(symbol string, period int) CurrentIndicator {
	var (
		apiKey    = "9cBb2iRgQOKzQMH9PxWALRXuwU6Ey0mNBKKzHmZUO9cLUTgCd7KtgTKLP4CDYJnR"
		secretKey = "pBHk1F2vSUfIkCVQQMd4ds5suA93tm5wPj3jPCmXSHulhy2rhHp0pFSoEWdvIrE4"
	)
	client := binance.NewClient(apiKey, secretKey)

	klines, err := client.NewKlinesService().Symbol(symbol).Limit(period).Interval("1m").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return CurrentIndicator{}
	}

	gainList := make([]float64, 0)
	lossList := make([]float64, 0)
	totalGain := 0.0
	totalLoss := 0.0
	// fmt.Println("kline length ", len(klines))
	for _, k := range klines {

		startPrice, _ := strconv.ParseFloat(k.Open, 64)
		endPrice, _ := strconv.ParseFloat(k.Close, 64)
		diff := endPrice - startPrice

		if diff > 0 {
			gainList = append(gainList, diff)
		} else {
			lossList = append(lossList, diff)

		}

	}

	for _, elem := range gainList {
		totalGain += elem
	}
	for _, elem := range lossList {
		totalLoss += elem
	}

	avgGain := totalGain / float64(period)
	avgLoss := totalLoss / float64(period)
	avgLoss *= -1
	if avgLoss == 0 {
		avgLoss = 1
	}
	RSI := 100 - 100/(1+(avgGain/avgLoss))
	if RSI < 0 || RSI > 100 {
		fmt.Println(" Error  RSI : ", RSI)

	}
	//fmt.Println("   Symbol  : ", symbol)
	//fmt.Println("   RSI : ", RSI)
	output := CurrentIndicator{Symbol: symbol, Value: RSI}
	return output
}



func MarketWatch() {

	period := 14
	// isOrderExecuted := false

	SetInterval(func() {
		fmt.Println("New Period =======================================")
		allSymbols := models.GetAllSymbols("binance")
		//allRSIVal := []CurrentIndicator{}
		allRSIVal := make (gin.H)
		for _, symbol := range allSymbols {
			output := UpdateRSI(symbol.Symbol, period)
			gotSymbol := symbol.Symbol
			allRSIVal[gotSymbol ] = output.Value;
		}
		// a ,_ :=  strconv.ParseFloat(allRSIVal[],64)
		// fmt.Println(a < 20)
	
		//class
	
		_,allBuyOrder := models.GetAllMarketOrderBySideAndIndicator("buy" , "RSI")
		_,allSellOrder := models.GetAllMarketOrderBySideAndIndicator("sell" , "RSI")
		for _,buyOrder := range(allBuyOrder){
			k := allRSIVal[ buyOrder.Symbol].(float64)	
		
			if k < buyOrder.Threshold{
				fmt.Println("Buy order executed")
				fmt.Println("current RSI :" , k)
				fmt.Println("buyOrder Threshold ", buyOrder.Threshold)
				fmt.Println(buyOrder)
			}
			
		}

		for _,sellOrder := range (allSellOrder){
			k := allRSIVal[ sellOrder.Symbol].(float64)	
		
			if k > sellOrder.Threshold{
				fmt.Println("Sell order executed")
				fmt.Println("current RSI :" , k)
				fmt.Println("sellOrder Threshold ", sellOrder.Threshold)
				fmt.Println(sellOrder)
			}
		}
		// go func() {
		// 	symbol := "MCOETH"
		// 	var (
		// 		apiKey    = "9cBb2iRgQOKzQMH9PxWALRXuwU6Ey0mNBKKzHmZUO9cLUTgCd7KtgTKLP4CDYJnR"
		// 		secretKey = "pBHk1F2vSUfIkCVQQMd4ds5suA93tm5wPj3jPCmXSHulhy2rhHp0pFSoEWdvIrE4"
		// 	)
		// 	client := binance.NewClient(apiKey, secretKey)

		// 	klines, err := client.NewKlinesService().Symbol(symbol).Limit(period).Interval("1m").Do(context.Background())
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		return
		// 	}

		// 	gainList := make([]float64, 0)
		// 	lossList := make([]float64, 0)
		// 	totalGain := 0.0
		// 	totalLoss := 0.0
		// 	// fmt.Println("kline length ", len(klines))
		// 	for _, k := range klines {

		// 		startPrice, _ := strconv.ParseFloat(k.Open, 64)
		// 		endPrice, _ := strconv.ParseFloat(k.Close, 64)
		// 		diff := endPrice - startPrice

		// 		if diff > 0 {
		// 			gainList = append(gainList, diff)
		// 		} else {
		// 			lossList = append(lossList, diff)

		// 		}

		// 	}

		// 	for _, elem := range gainList {
		// 		totalGain += elem
		// 	}
		// 	for _, elem := range lossList {
		// 		totalLoss += elem
		// 	}

		// 	avgGain := totalGain / float64(period)
		// 	avgLoss := totalLoss / float64(period)
		// 	avgLoss *= -1
		// 	if avgLoss == 0 {
		// 		avgLoss = 1
		// 	}
		// 	RSI := 100 - 100/(1+(avgGain/avgLoss))
		// 	if RSI < 0 || RSI > 100 {
		// 		fmt.Println(" Error  RSI : ", RSI)

		// 	}
		// 	fmt.Println("   RSI : ", RSI)
		// 	if isOrderExecuted == true {
		// 		fmt.Println("Order alread executed. No repeated action will be done.")
		// 		return
		// 	}
		// 	if RSI >= 70 {
		// 		isOrderExecuted = true
		// 		//sell
		// 	} else if RSI <= 63 {
		// 		isOrderExecuted = true
		// 		//buy
		// 	}

		// }()

	}, 60000, false)

}
