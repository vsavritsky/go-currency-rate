package main

import (
	"fmt"
	"github.com/vsavritsky/currencyRate/pkg/common/service/cbr"
)

func main() {
	rates := cbr.GetCurrencyRates()
	fmt.Printf("%+v\n", rates["USD"])
}
