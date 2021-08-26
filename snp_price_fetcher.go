package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/extrame/xls"
	"go.uber.org/ratelimit"
)

const (
	rateLimitInSeconds = 1
	// Export link from https://www.spglobal.com/spdji/en/indices/equity/sp-500/#overview
	urlSNP500              = "https://www.spglobal.com/spdji/en/idsexport/file.xls?redesignExport=true&languageId=1&selectedModule=PerformanceGraphView&selectedSubModule=Graph&yearFlag=tenYearFlag&indexId=340"
	urlTimeout             = time.Second * 30
	urlTLSHandshakeTimeout = time.Second * 5
	stockName              = "S&P 500"
	SNP500PricesFileName   = "spx_prices.xls"
	rowStockName           = 6
	columnDate             = 0
	columnPrice            = 1
	maxLoopCount           = 1000000
)

type Price struct {
	Date  time.Time
	Price float64
}

// note: this is off by two days on the real epoch (1/1/1900) because
// - the days are 1 indexed so 1/1/1900 is 1 not 0
// - Excel pretends that Feb 29, 1900 existed even though it did not
// The following function will fail for dates before March 1st 1900
// Before that date the Julian calendar was used so a conversion would be necessary
var excelEpoch = time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)

func FetchPrices() []Price {
	rateLimiter := ratelimit.New(rateLimitInSeconds)
	prices, err := downloadPrices(rateLimiter, urlSNP500)
	if err != nil {
		panic(err)
	}

	return prices
}

func downloadPrices(rateLimiter ratelimit.Limiter, pricesUrl string) ([]Price, error) {
	// TODO[petr]: save file with DATE. Do not download file if date of the last file is today
	filename := SNP500PricesFileName

	// check if file exists
	if _, err := os.Stat(filename); err != nil {
		rateLimiter.Take()

		var netClient = &http.Client{
			Timeout: urlTimeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   urlTimeout,
					KeepAlive: urlTimeout,
				}).DialContext,
				TLSHandshakeTimeout: urlTLSHandshakeTimeout,
			},
		}
		response, _ := netClient.Get(pricesUrl)
		defer response.Body.Close()

		responseBodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(filename, responseBodyBytes, 0644)
		if err != nil {
			return nil, err
		}
	}

	prices, err := fetchPricesFromXLS(filename)
	if err != nil {
		return nil, err
	}

	return prices, nil
}

func fetchPricesFromXLS(xlsFilePath string) ([]Price, error) {
	table, err := xls.Open(xlsFilePath, "UTF8")
	if err != nil {
		return nil, err
	}

	sheet := table.GetSheet(0)

	err = validateSheet(sheet)
	if err != nil {
		return nil, err
	}

	currentRow := rowStockName + 1
	var prices []Price
	for {
		if currentRow > maxLoopCount {
			return nil, fmt.Errorf("maximum loop size reached")
		}

		dateValue := sheet.Row(currentRow).Col(columnDate)
		priceValue := sheet.Row(currentRow).Col(columnPrice)

		if dateValue == "" {
			break
		}

		price, err := strconv.ParseFloat(priceValue, 64)
		prices = append(
			prices, Price{
				Date:  convertExcelDate(dateValue),
				Price: price,
			},
		)

		if err != nil {
			return nil, err
		}

		currentRow++
	}

	return prices, nil
}

func validateSheet(sheet *xls.WorkSheet) error {
	row6Column1 := sheet.Row(rowStockName).Col(columnPrice)
	if row6Column1 != stockName {
		return fmt.Errorf(
			"row %d column %d must be \"%s\". Current is \"%s\"",
			rowStockName,
			columnPrice,
			stockName,
			row6Column1,
		)
	}

	return nil
}

func printPrices(prices *[]Price) {
	for _, price := range *prices {
		fmt.Printf("%s - %.2f\n", price.Date.Format("02.01.2006"), price.Price)
	}
}

func convertExcelDate(excelDate string) time.Time {
	var days, _ = strconv.Atoi(excelDate)
	return excelEpoch.Add(time.Second * time.Duration(days*86400))
}
