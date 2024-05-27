package asakamiplugins

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	hbpath = engine.DataFolder() + "heartbeat"
)

func init() {
	// Check if heartbeat file exists
	if _, err := os.Stat(hbpath); os.IsNotExist(err) {
		// Create heartbeat file
		_, err := os.Create(hbpath)
		if err != nil {
			log.Print(err)
		}
	}

	// Read the heartbeat file
	hbtime, err := ioutil.ReadFile(hbpath)
	if err != nil {
		log.Print(err)
	}

	//将时间转换为时间戳
	var t time.Time
	t, err = time.Parse(time.RFC3339, string(hbtime))
	if err != nil {
		log.Print(err)
		t = time.Now()
	}

	// Get current time
	currentTime := time.Now()

	//获取离线时长（秒）
	offlineTime := currentTime.Sub(t).Seconds()

	//添加到统计
	statistics.lastofflinetime(offlineTime)
	statistics.sumofflinetime(offlineTime)
	runtimeitem, err := backpack.GetItem(103, 0)
	if err != nil {
		log.Print(err)
	}
	runtimeitem.Quantity = 0
	backpack.UpdateItem(runtimeitem, 0)

	// Start updating heartbeat file every 30 seconds
	go func() {
		var t = 0
		for {
			// Get current time
			currentTime := time.Now().Format(time.RFC3339)

			// Write current time to heartbeat file every 30 seconds
			if t == 30 {
				err := ioutil.WriteFile(hbpath, []byte(currentTime), 0644)
				if err != nil {
					log.Print(err)
				}
				t = 0
			}

			// Sleep for 1 seconds
			time.Sleep(1 * time.Second)
			t++
			statistics.nowruntime(1)
			statistics.sumruntime(1)
		}
	}()

}
