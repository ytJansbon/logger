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
		"\"outputType\":%d," +
		"\"level\":%d}}"
	config = fmt.Sprintf(config, user.HomeDir+"/log", 1024*1024*1, logger.OutputType_console, logger.Level_All)

	err = logger.InsLogger().Init(config)
	if err != nil {
		fmt.Println("logger new error:", err.Error())
		return
	}

	logger.InsLogger().Info("test logger info")
	time.Sleep(time.Duration(5) * time.Second)
	fmt.Println("test example_singleton end")
}
