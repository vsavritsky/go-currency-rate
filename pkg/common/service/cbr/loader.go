package cbr

import (
	"encoding/xml"
	"github.com/vsavritsky/go-currency-rate/pkg/common/model"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const URL = "https://cbr.ru/scripts/XML_daily.asp"
const DF = "02/01/2006"

var mu sync.Mutex

type XmlResult struct {
	ValCurs xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"Name,attr"`
	Valute  []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string  `xml:"ID,attr"`
	NumCode  int64   `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	Nominal  float64 `xml:"Nominal"`
	Name     string  `xml:"Name"`
	Value    string  `xml:"Value"`
}

type CurrencyRate struct {
	ID      string
	NumCode int64
	ISOCode string
	Name    string
	Value   float64
}

func GetCurrencyRates() map[string]model.Rate {
	mu.Lock()
	results := make(map[string]model.Rate)
	var rates = FetchCurrencyRates(time.Time{})

	for _, el := range rates {
		results[el.ISOCode] = model.Rate{
			CurrencyCode: el.ISOCode,
			Provider:     "cbr",
			Value:        el.Value,
		}
	}

	defer mu.Unlock()

	return results
}

func FetchCurrencyRates(d time.Time) map[string]CurrencyRate {
	url := URL
	if !d.IsZero() {
		url = url + "?date_req=" + d.Format(DF)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error of get currency: %v", err.Error())
		return nil
	}

	var data XmlResult

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&data)

	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}

	rates := make(map[string]CurrencyRate)

	for _, el := range data.Valute {
		value, _ := strconv.ParseFloat(strings.Replace(el.Value, ",", ".", -1), 64)

		rates[el.CharCode] = CurrencyRate{
			ID:      el.ID,
			NumCode: el.NumCode,
			ISOCode: el.CharCode,
			Name:    el.Name,
			Value:   value / el.Nominal,
		}
	}

	return rates
}
