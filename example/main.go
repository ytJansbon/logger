package main

import (
	"fmt"
	"logger"
	"os/user"
	"time"
)

func main() {
	user, err := user.Current()
	config := "{" +
		"\"prefix\":\"mylogger\"," +
		"\"bufferSize\":1024," +
		"\"fileDir\":\"%s\"," +
		"\"fileName\":\"example.log\"," +
		"\"fileSize\":%d," +
		"\"fileCount\":5," +
		"\"console\":false" +
		"}"
	config = fmt.Sprintf(config, user.HomeDir+"/log", 1024*1024*1)
	logger, err := logger.NewLogger(config)
	if err != nil {
		fmt.Println("logger new error:", err.Error())
		return
	}

	logger.Info("test logger info")
	time.Sleep(time.Duration(5) * time.Second)
	fmt.Println("test end")
}
