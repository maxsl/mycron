package main

import (
	"fmt"
	"git.oschina.net/wida/mycron/src/cron"
	"git.oschina.net/wida/mycron/src/mydb"
	"git.oschina.net/wida/mycron/src/myexec"
	"time"
)

func main() {
	jobs, _ := mydb.GetCronList()
	c := cron.New()
	defer func() {
		c.Stop()
	}()

	//添加jobs
	for i := 0; i < len(jobs); i++ {
		job := jobs[i]
		c.AddFunc(job.Time,
			func() {
				job.ChangeRunningStatu(1)
				s, e := myexec.ExecWithTimeout(time.Second*10, job.Cmd)
				if e != nil {
					fmt.Print(e)
				}
				job.ChangeRunningStatu(0)
				fmt.Println(s)
			}, int(job.Status), int(job.ID), int64(job.STime), int64(job.ETime))
	}
	//start
	c.Start()

	//监听更新事件
	for {
		select {
		case <-time.After(time.Second):
			jobs, _ := mydb.GetModifyList()
			for i := 0; i < len(jobs); i++ {
				job := jobs[i]
				c.AddFunc(job.Time,
					func() {
						job.ChangeRunningStatu(1)
						s, e := myexec.ExecWithTimeout(time.Second*10, job.Cmd)
						if e != nil {
							fmt.Print(e)
						}
						job.ChangeRunningStatu(0)
						fmt.Println(s)
					}, int(job.Status), int(job.ID), int64(job.STime), int64(job.ETime))
			}
			mydb.UpdateModifyList()
			continue
		}
	}
}
func printfEntry(c *cron.Cron) {
	for _, v := range c.Entries() {
		fmt.Println(v.ID, v.Status, v.Start, v.Ending)
	}
}