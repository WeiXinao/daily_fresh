package main

import (
	"os"
	"runtime"

	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api"
)

func main() {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	
	api.NewApp("api-server").Run()
}