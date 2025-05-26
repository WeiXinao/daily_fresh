package main

import (
	"os"
	"runtime"

	"github.com/WeiXinao/daily_your_go/app/user"
)

func main() {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	
	user.NewApp("user-server").Run()
}