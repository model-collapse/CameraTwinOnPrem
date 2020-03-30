package main

import (
	"flag"
	"log"
	"os"

	rt "github.com/qiangxue/fasthttp-routing"
	h "github.com/valyala/fasthttp"
)

var fsHandler h.RequestHandler

func handleFSGet(ctx *rt.Context) (reterr error) {
	fsHandler(ctx.RequestCtx)
	return
}

func main() {
	servicePort := os.Getenv("SERVICE_PORT")
	flag.Parse()
	fsHandler = h.FSHandler("/store", 2)

	router := rt.New()
	router.Get("/files/*")

	log.Printf("Starting listening on %s....", servicePort)
	h.ListenAndServe("0.0.0.0:"+servicePort, router.HandleRequest)
}
