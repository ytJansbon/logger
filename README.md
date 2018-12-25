## Logger
An asynchronous log module that supports automatic sharding by date for file size.

## Features
1. Logs are recorded at different levels.
2. Support logs to console and file output.
3. File log sharding supports by date and file size.

## Installation
go get -u github.com/ytJansbon/logger

## Examples
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
    		"\"outputType\":%d" +
    		"}"
    	config = fmt.Sprintf(config, user.HomeDir+"/log", 1024*1024*1, logger.OutputType_console)
    	logger, err := logger.NewLogger(config)
    	if err != nil {
    		fmt.Println("logger new error:", err.Error())
    		return
    	}
    
    	logger.Info("test logger info")
    	time.Sleep(time.Duration(5) * time.Second)
    	fmt.Println("test end")
    }