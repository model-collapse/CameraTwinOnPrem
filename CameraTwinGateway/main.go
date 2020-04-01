package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	rt "github.com/qiangxue/fasthttp-routing"
	h "github.com/valyala/fasthttp"
)

var mclient mqtt.Client
var mqttPort string
var mqttHost string
var cleanSession string
var servicePort string

func toURL(u string) (ret *url.URL) {
	ret, err := url.Parse(u)
	if err != nil {
		log.Fatalf("Error in parsing url, %s, %v", u, err)
	}

	return
}

func initMQTT() {
	mqttHost = os.Getenv("MQTT_HOST")
	mqttPort = os.Getenv("MQTT_PORT")
	cleanSession = strings.ToLower(os.Getenv("MQTT_CLEAN_SESSION"))

	mclient = mqtt.NewClient((&mqtt.ClientOptions{
		Servers:      []*url.URL{toURL(fmt.Sprintf("tcp://emq:8883", mqttHost, mqttPort))},
		ClientID:     "camera_twin",
		CleanSession: cleanSession == "true" || cleanSession == "yes",
	}))

	if token := mclient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

func main() {
	initMQTT()
	servicePort = os.Getenv("SERVICE_PORT")

	router := rt.New()
	router.Post("/camera_twin/upload", handleFSPost)
	router.Post("/camera_twin/log", handleLogPost)

	log.Printf("Starting listening on %s....", servicePort)
	h.ListenAndServe("0.0.0.0:"+servicePort, router.HandleRequest)
}
