package main

import (
	"os"
	"runtime"

	"github.com/WeiXinao/daily_fresh/app/user/server"
)

func main() {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	
	srv.NewApp("user-server").Run()
}