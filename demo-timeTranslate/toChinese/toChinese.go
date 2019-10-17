package toChinese

import (
	"strconv"
	"strings"
	"time"
)

//将时间转化为数组
func transToSlice(time time.Time)([]string) {

	tNow := time.Format("2006-01-02-15-04-05")

	return strings.Split(tNow, "-")
}

//判断时间数组是否需要加“十”，如十三，四十五
func decimalCheck(timeSlice []string)([]string){
	for i,v := range timeSlice{
		v1 ,_ :=strconv.Atoi(v)
		if  i>0 && i<4 {
			x:= v1/10
			y:= v1%10
			z :=strconv.Itoa(x)+"+"+strconv.Itoa(y)
			if x == 0 {
				z =strconv.Itoa(y)
			}else if x == 1 && y == 0 {
				z ="+"
			}else if y == 0 {
				z =strconv.Itoa(x)+"+"
			}else if x == 1 {
				z ="+"+strconv.Itoa(y)
			}

			timeSlice[i] = z
		}
		if i>=4 {
			x:= v1/10
			y:= v1%10
			z :=strconv.Itoa(x)+"+"+strconv.Itoa(y)
			if x == 0 {
				z =strconv.Itoa(x)+strconv.Itoa(y)
			}else if y==0 {
				z =strconv.Itoa(x)+"+"
			}else if x==1 {
				z ="+"+strconv.Itoa(y)
			}
			timeSlice[i] = z
		}

	}
	return timeSlice
}

//将时间切片转换为map，便于区分年月日信息
func transToMap(timeSlice []string)(map[int]([]string)){
	var timeMap = make(map[int]([]string),4)
	for i, v := range timeSlice {
		var timeArray1 []string = make([]string,0)
		for _,v1:=range v{
			timeArray1 = append(timeArray1,string(v1))
		}
		timeMap[i]=timeArray1
	}
	return timeMap
}




func TransTimeToChineseWords(time time.Time)([]string) {

	timeslice := transToSlice(time)
	timeslice1 := decimalCheck(timeslice)
	timeMap := transToMap(timeslice1)

	//对map中的字符进行转换
	for i,v:=range timeMap{
		for i1,v1:=range v {
			switch v1 {
			case "0":
				timeMap[i][i1]="零"
			case "1":
				timeMap[i][i1]="一";
			case "2":
				timeMap[i][i1]="二";
			case "3":
				timeMap[i][i1]="三";
			case "4":
				timeMap[i][i1]="四";
			case "5":
				timeMap[i][i1]="五";
			case "6":
				timeMap[i][i1]="六";
			case "7":
				timeMap[i][i1]="七";
			case "8":
				timeMap[i][i1]="八";
			case "9":
				timeMap[i][i1]="九";
			case "+":
				timeMap[i][i1]="十"
			default:
				continue
			}
		}

	}

	var a []string = make([]string,6)
	for i,v:=range timeMap{ //将map拼接为切片返回，切片便于使用
		a[i] = strings.Join(v,"")
	}

	//result :=a[0]+"年"+a[1]+"月"+a[2]+"日"+a[3]+"时"+a[4]+"分"
	//fmt.Println(result)
	return a
}
