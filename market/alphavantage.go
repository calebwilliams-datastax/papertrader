package market

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//AVConfig -
type AVConfig struct {
	APIKey string
	Base   string
}

//TimeSeriesRaw -
type TimeSeriesRaw struct {
	Symbol    string
	Timestamp string
	Open      string
	High      string
	Low       string
	Close     string
	Volume    string
}

//CurrentPrice -
func CurrentPrice(symbol string, av AVConfig) (float64, error) {
	log.Printf("alphavantage : CurrentPrice symbol:%s\n", symbol)
	current, err := TimeSeriesIntraday(symbol, 5, av)
	if err != nil {
		return 0, err
	}
	actual, _ := strconv.ParseFloat(current[len(current)-1].Close, 64)
	fmt.Printf(" pice:%v", actual)
	return actual, nil
}

//TimeSeriesIntraday - csv output
func TimeSeriesIntraday(symbol string, interval int, av AVConfig) ([]TimeSeriesRaw, error) {
	resSlice := make([]TimeSeriesRaw, 0)
	/*
		https://www.alphavantage.co/
		query?function=TIME_SERIES_INTRADAY
		&symbol=GME
		&interval=5min
		&apikey=HR9QB2RM5GWOO0IX
	*/
	//TODO : Break the http functionality out into a generic function?
	url := fmt.Sprintf("%squery?function=TIME_SERIES_INTRADAY&symbol=%s&interval=%vmin&apikey=%s&datatype=csv", av.Base, symbol, interval, av.APIKey)
	fmt.Printf("alphavantage: TimeSeriesIntraday: \n%s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return resSlice, err
	}
	//Begin handling the responsebody
	dataBytes, err := ioutil.ReadAll(resp.Body)
	fmt.Printf(" response body length: %v", len(dataBytes))
	if len(dataBytes) > 0 {
		data := string(dataBytes)
		for i, tsi := range strings.Split(data, "\n") {
			if i != 0 {
				/*
					timestamp,open,high,low,close,volume
					2021-02-08 20:00:00,59.3400,59.5000,59.2000,59.5000,11046
				*/
				row := strings.Split(tsi, ",")
				if len(row) > 1 {
					av := TimeSeriesRaw{
						Symbol:    symbol,
						Timestamp: row[0],
						Open:      row[1],
						High:      row[2],
						Low:       row[3],
						Close:     row[4],
						Volume:    row[5][:len(row[5])-2],
					}
					resSlice = append(resSlice, av)
				}
			}
		}
		fmt.Printf(" time series obj count: %v\n", len(resSlice))
		return resSlice, nil
	}
	return resSlice, fmt.Errorf("no response body returned from %s", url)
}
