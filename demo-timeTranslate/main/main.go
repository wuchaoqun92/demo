package main

import (
	"demo-person/demo-timeTranslate/toChinese"
	"fmt"
	"time"
)

func main() {
	a := toChinese.TransTimeToChineseWords(time.Now())
	b := a[0]+"年"+a[1]+"月"+a[2]+"日"+a[3]+"时"+a[4]+"分"+a[5]+"秒"
	fmt.Printf("%T,%s\n",b,b)


}
