package crontab

import (
	"demo-person/demo-miniprogramme/common"
	"fmt"
	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	spec := "*/5 * * * * *"
	c.AddFunc(spec, func() {
		fmt.Println("every 5 seconds executing")
	})

	c.AddFunc(spec, func() {
		a := common.TransTimeToChineseWords()
		s := a[0] + "年" + a[1] + "月" + a[2] + "日"
		fmt.Println(s)
		//common.CreateList(s)
		//common.DeleteList()
	})

	go c.Start()
	defer c.Stop()

	select {
	//case <-time.After(time.Second * 10):
	//	return
	}

}
