package main

import "C"
import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server("S1")
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("Hello World,bazhe")
	})

	s.BindHandler("/test", abc)
	s.SetAddr("120.23.23.1")
	s.SetPort(7090, 9090)
	s.Run()

}

func abc(r *ghttp.Request) {
	r.Response.Write("test")
}
