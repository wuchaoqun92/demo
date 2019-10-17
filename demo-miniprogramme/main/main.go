package main

import (
	"awesomeProject/demo-miniprogramme/method"
	"fmt"
	"log"
	"net/http"
)


func main() {

	http.HandleFunc("/checkIn",Panics(method.CheckIn))
	http.HandleFunc("/GetContentList",Panics(method.ChoiceMethodToList))
	err := http.ListenAndServe(":8888",nil);
	fmt.Println("shshshshshs",err)
	if err == nil {
		fmt.Println("shshshshshs", err)
	}


}


//panic预处理
func Panics(handle http.HandlerFunc)http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		defer func(){
			if x:=recover();x!=nil{
				log.Printf("[%v]caught panic:%v",r.RemoteAddr,x)
			}
		}()
		handle(w,r)
	}
}