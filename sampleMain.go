package main

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

func main() {
	a := "2/13"
	b := strings.Split(a, "/")
	c, _ := decimal.NewFromString(b[0])
	d, _ := decimal.NewFromString(b[1])
	e, _ := c.Div(d).Float64()
	fmt.Println(e)
	fmt.Println(decimal.NewFromFloat(0.0))

	// decimal.NewFromInt(1)
	// price := decimal.NewFromInt(1)
	// quantity := decimal.NewFromInt(13)
	// subtotal := price.Div(quantity)
	// newPrice := decimal.NewFromInt(100)
	// val := subtotal.Mul(newPrice)
	// fmt.Println(fmt.Println(big.NewFloat(subtotal).Text('f', 2)))
}
