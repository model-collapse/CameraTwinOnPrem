package main

import (
	//"fmt"
	"log"
	"time"

	rt "github.com/qiangxue/fasthttp-routing"
)

func handleFSPost(ctx *rt.Context) (err error) {
	log.Printf("Getting post...")

	bodyData := ctx.PostBody()
	timeStr := string(ctx.QueryArgs().Peek("timestamp"))
	tm, _ := time.Parse("20060102150405", timeStr)
	tmr := time.Now()

	deviceName := string(ctx.QueryArgs().Peek("device_id"))
	log.Printf("device_id = %s", deviceName)
	log.Printf("timestamp = %s", timeStr)

	if timeStr == "" && deviceName == "" {
		log.Printf("empty file name")
		ctx.SetStatusCode(500)
		return
	}

	if isBrokenJPEG(bodyData) {
		ctx.SetStatusCode(406)
		log.Printf("Broken JPG, returning 406...")
		return
	}

	pk := ImgPack{
		DeviceId:      deviceName,
		TimestampSend: tm.Unix(),
		TimestampRecv: tmr.Unix(),
		Image:         bodyData,
	}

	sd, _ := pk.Marshal()
	token := mclient.Publish("camera/capture/people_analytics", 0, false, sd)
	token.Wait()
	if token.Error() != nil {
		log.Panicf("Error in publush! %v", token.Error())
	} else {
		log.Printf("Command published!")
	}

	token = mclient.Publish("camera/capture/dump", 0, false, sd)
	token.Wait()
	if token.Error() != nil {
		log.Panicf("Error in publush! %v", token.Error())
	} else {
		log.Printf("Command published!")
	}

	ctx.SetStatusCode(200)

	return
}
