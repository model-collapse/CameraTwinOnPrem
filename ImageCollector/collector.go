package main

import (
	fmt "fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func toURL(u string) (ret *url.URL) {
	ret, err := url.Parse(u)
	if err != nil {
		log.Fatalf("Error in parsing url, %s, %v", u, err)
	}

	return
}

func main() {
	mqttHost := os.Getenv("MQTT_HOST")
	mqttPort := os.Getenv("MQTT_PORT")
	mclient := mqtt.NewClient((&mqtt.ClientOptions{
		Servers:  []*url.URL{toURL(fmt.Sprintf("tcp://%s:%s", mqttHost, mqttPort))},
		ClientID: "image collector",
	}))

	if token := mclient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	token := mclient.Subscribe("camera/capture/dump", 0, func(cl mqtt.Client, msg mqtt.Message) {
		pck := ImgPack{}
		if err := pck.Unmarshal(msg.Payload()); err != nil {
			log.Panicf("Invalid data format of image")
			return
		}

		dir := fmt.Sprintf("/store/%s", pck.DeviceId)
		if stat, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				os.MkdirAll(dir, os.ModePerm)
			} else {
				log.Printf("ERROR: %v", err)
			}
		} else if !stat.IsDir() {
			log.Panicf("This dir is occupied by a file")
		}

		tm := time.Unix(pck.TimestampSend, 0)
		timeStr := tm.Format("20060102150405")

		path := fmt.Sprintf("%s/%s_%d.jpg", dir, pck.DeviceId, timeStr)
		ioutil.WriteFile(path, pck.Image, os.ModePerm)
	})

	token.Wait()
	if token.Error() != nil {
		log.Panicf("Error in publush! %v", token.Error())
	} else {
		log.Printf("Command Subscribed!")
	}

	c chan os.Signal
	signal.Notify(c, syscall.SIGTERM)
	<-c
	
	return
}
