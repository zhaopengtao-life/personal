package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

func Decimal() {
	var ifHCOutOctets_start, ifHCOutOctets_end uint64
	ifHCOutOctets_start = 172252719015828
	ifHCOutOctets_end = 172249605304388
	numa := fmt.Sprintf("%v", ifHCOutOctets_start)
	fmt.Println("numa: ", numa)
	numb := fmt.Sprintf("%v", ifHCOutOctets_end)
	fmt.Println("numb: ", numb)
	num1, err := decimal.NewFromString(numa)
	if err != nil {
		panic(err)
	}
	num2, err := decimal.NewFromString(numb)
	if err != nil {
		panic(err)
	}
	num3 := num1.Sub(num2)
	fmt.Println("num3: ", num3)
	numc, _ := decimal.NewFromString("56")
	fmt.Println("numc: ", numc)
	num4 := num3.Div(numc)
	fmt.Println("num4: ", num4)
	numd, _ := decimal.NewFromString("8")
	fmt.Println("numd: ", numd)
	num5 := num4.Mul(numd)
	fmt.Println("num5: ", num5)

	nume := fmt.Sprintf("%v", num5)
	fl, err := strconv.ParseFloat(nume, 64)
	if err != nil {
		fmt.Println("转换错误：", err)
	} else {
		fmt.Println("nume: ", fl)
	}
}

func DemoDecimal() {
	// float 四舍五入取小数点后两位
	value := 0.123434
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)

	// 字符串转decimal
	divisor := "0.45"
	price, _ := decimal.NewFromString(divisor) //获取数字
	fmt.Println("字符串转decimal: ", price)

	// float64转decimal
	values := decimal.NewFromFloat(value)
	fmt.Println("float64转decimal: ", values)

	// decimal转float64
	float, _ := strconv.ParseFloat(price.String(), 64)
	fmt.Println("decimal转float64: ", float)
}
