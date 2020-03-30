package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	rt "github.com/qiangxue/fasthttp-routing"
)

func handleLogPost(ctx *rt.Context) (reterr error) {
	log.Printf("Getting post for logs")
	deviceName := string(ctx.QueryArgs().Peek("device_id"))
	log.Printf("Device Name = %s", deviceName)

	if deviceName == "" {
		log.Printf("Empty device name")
		ctx.SetStatusCode(500)
		return
	}

	timeStr := string(ctx.QueryArgs().Peek("timestamp"))
	if timeStr == "" {
		timeStr = fmt.Sprintf("%d", time.Now().Unix())
	}
	timestamp, _ := strconv.ParseInt(timeStr, 10, 64)
	tm := time.Unix(timestamp, 0)
	timeStrf := tm.Format("20060102150405")

	logTerm := struct {
		DeviceName string `json:"device_id"`
		Timestamp  int64  `json:"timestamp"`
		TimeStr    string `json:"time_str"`
		Log        string `json:"log"`
	}{
		DeviceName: deviceName,
		Timestamp:  timestamp,
		TimeStr:    timeStrf,
		Log:        string(ctx.PostBody()),
	}

	dt, _ := json.Marshal(logTerm)
	token := mclient.Publish("camera/log/scratched", 0, false, dt)
	token.Wait()
	if token.Error() != nil {
		log.Panicf("Error in publush! %v", token.Error())
	} else {
		log.Printf("log published!")
	}

	ctx.SetStatusCode(200)
	return
}
