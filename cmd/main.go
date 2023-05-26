package main

import (
	"fmt"
	"github.com/vsavritsky/go-currency-rate/pkg/common/service/cbr"
)

func main() {
	rates := cbr.GetCurrencyRates()
	fmt.Printf("%+v\n", rates["USD"])
}
