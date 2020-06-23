package main

//go build -ldflags "-w -s" main.go
import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
	"xxnet/middleware"
	"xxnet/modstation"

	"github.com/gin-gonic/gin"
)

var (
	gPch *modstation.ChannelManager
	gErr error
)

func init() {
	logFile, err := os.OpenFile("./error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lmicroseconds | log.Ldate)

}

//process ---
func process(pch *modstation.ChannelManager) {
	t1 := time.NewTimer(time.Second * time.Duration(pch.RequestGap))
	for {
		select {
		case <-t1.C:
			pch.RequestData()
			t1.Reset(time.Second * time.Duration(pch.RequestGap))
		}
		runtime.Gosched()
	}
}

//fatalError --
func fatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	gPch, gErr = modstation.NewChannelManager("./cfg/config.json")
	fatalError(gErr)
	go process(gPch)
	fatalError(gErr)

	APIPort := strconv.Itoa(int(gPch.APIPort))
	modstation.Debug("Server v1.21 Start At :", APIPort)
	go func() {
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		router.Use(middleware.Cors()) //跨域支持
		router.POST("/setalarm", gPch.PostParams)
		router.GET("/data", gPch.GetData)
		router.GET("/json", gPch.GetJSON)
		router.GET("/reload", gPch.ReloadModbus)
		router.Run(":" + APIPort)
		fmt.Println("end")
		log.Println("end")
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown Server ...")
	log.Println("Shutdown Server ...")
	gPch.Close()

}
