package main

import (
	"os"
	"runtime"

	"github.com/WeiXinao/daily_fresh/app/product/server"
)

func main() {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	
	server.NewApp("product-rpc").Run()
}