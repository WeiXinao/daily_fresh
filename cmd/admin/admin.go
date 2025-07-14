package main

import (
	"os"
	"runtime"

	"github.com/WeiXinao/daily_fresh/app/bff/admin"
)

func main() {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	
	admin.NewApp("admin-server").Run()
}