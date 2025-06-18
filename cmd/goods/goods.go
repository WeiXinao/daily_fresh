package main

import (
	"os"
	"runtime"

	"github.com/WeiXinao/daily_your_go/app/goods/srv"
)

func main() {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	
	srv.NewApp("goods-server").Run()
}