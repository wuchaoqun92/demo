package main

import (
	"bytes"
	"crypto/md5"
	"demo-person/demo-httpCommon/netRequest"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	//test1()
	netRequest.GetRequest("http://www.baidu.com")
}

func test1() {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Handle("GET", "/welcome/{id:int}", func(ctx iris.Context) {
		ctx.HTML("<h1>Welcome</h1>")
		fmt.Println(ctx.Params().Get("id"))
		app.Logger().Warn(ctx.Params().Get("id"))
	})

	app.Get("/hello", func(ctx iris.Context) {
		n, err := ctx.JSON(iris.Map{"message": "Hello Iris!"})
		if err != nil || n != 25 {
			fmt.Println("err", n)
			return
		}
	})
	app.RegisterView(iris.HTML("/Users/wuchaoqun/Desktop/codeMaterial", ".html"))
	app.Get("/up", func(ctx iris.Context) {
		now := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(now, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		ctx.View("2019-12-06.html", token)

	})

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func test2() {
	app := iris.New()
	//请在参数化路径部分
	users := app.Party("/users", myAuthMiddlewareHandler)
	// http://localhost:8080/users/42/profile
	users.Get("/{id:int}/profile", userProfileHandler)
	// http://localhost:8080/users/inbox/1?4
	users.Get("/inbox/{id:int}/{uid:int}", userMessageHandler)

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
func myAuthMiddlewareHandler(ctx iris.Context) {
	ctx.WriteString("Authentication failed")
	ctx.Next() //继续执行后续的handler
}
func userProfileHandler(ctx iris.Context) { //
	id := ctx.Params().Get("id")
	ctx.WriteString(id)
}
func userMessageHandler(ctx iris.Context) {

	id := ctx.Params().Get("id")
	uid := ctx.Params().Get("uid")
	ctx.WriteString(id)
	ctx.WriteString(uid)
}

const (
	accessToken = "24.0d310a2fe08577bcc94ccf8520d611ae.2592000.1575602067.282335-17706172"
	recogUrl    = "https://aip.baidubce.com/rest/2.0/face/v3/detect"
	requestUrl  = recogUrl + "?access_token=" + accessToken
)

func test3() {
	app := iris.New()

	app.Get("/h", func(ctx iris.Context) {
		ctx.Writef("Hello from the server")
	})

	app.Post("/mypath", func(ctx iris.Context) {
		conType := ctx.Request().Header
		urlpa := ctx.URLParams()
		b := urlpa["b"]
		a, _ := ctx.GetBody()

		ctx.Writef("Hello from %s;conType=%s;urlpa=%s;a=%s ", ctx.Path(), conType, b, string(a))
	})

	// Note: It's not needed if the first action is "go app.Run".
	if err := app.Build(); err != nil {
		panic(err)
	}

	srv1 := &http.Server{Addr: ":9090", Handler: app}
	go srv1.ListenAndServe()
	println("Start a server listening on http://localhost:9090")

	srv2 := &http.Server{Addr: ":5050", Handler: app}
	go srv2.ListenAndServe()
	println("Start a server listening on http://localhost:5050")

	app.Run(iris.Addr(":8080"))

}

type logtest struct {
	time    string
	creater string
	msg     string
}

func (a *logtest) Levels() []log.Level {
	back := make([]log.Level, 3)
	back[0] = log.TraceLevel
	back[1] = log.DebugLevel
	back[2] = log.ErrorLevel
	return back
}

func (a *logtest) Fire(entry *log.Entry) error {
	//entry.Data["time-re"] = a.time
	data := make(log.Fields)
	data["Level"] = entry.Level.String()
	data["Time"] = entry.Time
	data["Message"] = entry.Message
	entry.Data = data
	return nil
}

type myFormatter struct {
}

func (my *myFormatter) Format(entry *log.Entry) ([]byte, error) {
	data := make(log.Fields)
	data["Level"] = "hahahahahahah"
	data["Time"] = entry.Time

	var b *bytes.Buffer

	b = entry.Buffer

	return b.Bytes(), nil
}

func logTest() {
	file, _ := os.OpenFile("/Users/wuchaoqun/Desktop/codeMaterial/1234.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer file.Close()
	log.SetFormatter(&myFormatter{})
	log.SetLevel(log.TraceLevel)
	log.SetOutput(file)

	b := logtest{
		time:    time.Now().Format("2006-1-2 15:04:05"),
		creater: "me",
		msg:     "test",
	}
	log.AddHook(&b)

	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Trace("A walrus appears")
}

func testMd5() {

	str1 := "123456"
	str2 := "w123456"

	//方法一
	start := time.Now()
	data := []byte(str1)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	cost := time.Since(start)
	fmt.Println("result1:", md5str1)
	fmt.Println("cost1:", cost)

	//方法二

	start = time.Now()
	w := md5.New()
	io.WriteString(w, str2)                  //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil)) //w.Sum(nil)将w的hash转成[]byte格式,sprintf 将 byte 转成 16 进制
	cost = time.Since(start)

	fmt.Println("result1:", md5str2)
	fmt.Println("cost1:", cost)

	fmt.Println(md5str1 == md5str2)

}
