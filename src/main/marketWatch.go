package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/resty.v1"
)

type OuterLayer struct {
	Msg []RSIData
}

type RSIData struct {
	Data [][]string
}

type MarketOrder struct {
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"`
	T        string  `json:"type"`
	Quantity float64 `json:"quantity"`
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

func MarketWatch() {

	// A counter for the number of times we print
	// printed := 0

	// We call set interval to print Hello World forever
	// every 1 second
	// clear := make(chan float64)
	// var gainList []float64
	// var lossList []float64
	period := 14
	isOrderExecuted := false

	SetInterval(func() {

		go func() {
			resp, _ := resty.R().Get("https://api.binance.com/api/v1/klines?symbol=BTCUSDT&interval=1m&limit=" + strconv.Itoa(period))
			// fmt.Println("\nhttps://api.binance.com/api/v1/klines?symbol=BTCUSDT&interval=1m&limit=" + strconv.Itoa(period))
			// explore response object
			// fmt.Printf("\nError: %v", err)
			fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
			fmt.Printf("\nResponse Status: %v", resp.Status())
			fmt.Printf("\nResponse Time: %v", resp.Time())
			fmt.Println("\nResponse Received At: %v", resp.ReceivedAt())
			// k := &[]RSIData{}
			// _ = json.Unmarshal(resp.Body(), &k)
			// // fmt.Println(string(resp.Body()))
			rt := string(resp.Body())
			rt = rt[1 : len(rt)-2]
			rt = strings.Replace(rt, "[", "", 10)
			rt = strings.Replace(rt, "]", "", 10)
			rt = strings.Replace(rt, `"`, "", 100)
			arrRt := strings.Split(rt, ",")
			arrRt = append(arrRt, " ")
			allCandles := make([][]string, 0)
			singleCandle := make([]string, 0)
			for index, element := range arrRt {
				if index%12 != 0 {
					singleCandle = append(singleCandle, element)
				} else if index != 0 && index%12 == 0 {

					allCandles = append(allCandles, singleCandle)

					singleCandle = make([]string, 0)
					singleCandle = append(singleCandle, element)

				} else if index == 0 && index%12 == 0 {
					singleCandle = append(singleCandle, element)

				}

			}

			// fmt.Println(allCandles)
			gainList := make([]float64, 0)
			lossList := make([]float64, 0)
			totalGain := 0.0
			totalLoss := 0.0
			for _, elem := range allCandles {
				startPrice, _ := strconv.ParseFloat(elem[1], 64)
				endPrice, _ := strconv.ParseFloat(elem[4], 64)
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
			fmt.Println("   RSI : ", RSI)
			if isOrderExecuted == true {
				fmt.Println("Order alread executed")
				return
			}
			if RSI >= 70 {
				isOrderExecuted = true
				//sell
			} else if RSI <= 30 {
				isOrderExecuted = true
				//buy
			}
			// fmt.Printf("\nResponse Body: %v", string(resp.Body())) // or resp.String() or string(resp.Body())
			// respString := string(resp.Body())
			// rawResp := string(resp.Body())
			// respData := rawResp[2 : len(rawResp)-2]
			// dataSlice := strings.Split(respData, ",")
			// rawOpeningStr := dataSlice[1][1 : len(dataSlice[1])-1]
			// rawClosingStr := dataSlice[4][1 : len(dataSlice[1])-1]
			// fmt.Println("rawOpeningStr : ", rawOpeningStr)

			// openingPrice, err := strconv.ParseFloat(rawOpeningStr, 32)
			// closingPrice, err := strconv.ParseFloat(rawClosingStr, 32)
			// diff := closingPrice - openingPrice
			// fmt.Println("Diff : ", diff)
			// clear <- diff
			////////////////////////////////////////////////////////////

		}()
		// currentValue := <-clear
		// fmt.Println("current value ", currentValue)
		// if currentValue > 0 {
		// 	fmt.Println("current value at gain")

		// 	gainList = append(gainList, currentValue)
		// } else {
		// 	fmt.Println("current value at loss")

		// 	lossList = append(lossList, currentValue)

		// }
		// fmt.Println("len(gainList) :", len(gainList))
		// fmt.Println("len(lossList) : ", len(lossList))

		// if len(gainList)+len(lossList) == period {
		// 	gainSum := 0.0
		// 	lossSum := 0.0
		// 	for i, s := range gainList {
		// 		fmt.Println(i, s)
		// 		gainSum += s
		// 	}
		// 	if len(gainList) == 0 {
		// 		gainList = append(gainList, 0)

		// 	}
		// 	avgGain := gainSum / float64(len(gainList))
		// 	for i, s := range lossList {
		// 		fmt.Println(i, s)
		// 		lossSum += s
		// 	}
		// 	if len(lossList) == 0 {
		// 		lossList = append(lossList, 0)

		// 	}
		// 	avgLoss := lossSum / float64(len(lossList))
		// 	avgLoss *= -1
		// 	fmt.Println("Average gain : ", avgGain)
		// 	fmt.Println("Average loss : ", avgLoss)
		// 	if avgLoss == 0 {
		// 		fmt.Println("avg loss is zero.Caculation voided")
		// 	} else {
		// 		RSI := 100 - (100 / (1 + (avgGain / avgLoss)))
		// 		fmt.Println("RSI : ", RSI)
		// 	}
		// 	gainList = gainList[:0]
		// 	lossList = lossList[:0]
		// }
	}, 60000, false)

	// If we wanted to we had a long running task (i.e. network call)
	// we could pass in true as the last argument to run the function
	// as a goroutine

	// Some artificial work here to wait till we've printed
	// 5 times
	// for {
	// 	if printed == 5 {
	// 		// Stop the ticket, ending the interval go routine
	// 		stop <- true
	// 		return
	// 	}
	// }

}
