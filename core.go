package logger

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	MAX_LOGCHAN_BUFFER_SIZE   = 2048 //channel buffer size
	CHECK_SLICE_FILE_INTERVAL = 180  //check slice file interval time
)

const (
	DEFAULT_CONSOLE_VALUE   = true
	DEFAULT_SLICETYPE_VALUE = 0
	DEFAULT_PREFIX_VALUE    = ""
	DEFAULT_LEVEL_VALUE     = Level_Info
	DEFAULT_FILEDIR_VALUE   = "./log"
	DEFAULT_FILENAME_VALUE  = "logfile.log"
	DEFAULT_FILECOUNT_VALUE = 3
	DEFAULT_FILESIZE_VALUE  = 1024 * 1024 * 3
)

type logHandler struct {
	dataChan chan string
	lock     *sync.Mutex

	console   bool
	sliceType SliceType
	prefix    string
	level     int

	curDate    *time.Time
	fileHandle *os.File

	fileDir   string
	fileName  string
	fileCount int
	fileSize  uint64
	fileIndex int
}

//init logHandler
func (lh *logHandler) init(config string) error {
	//check config params
	err := lh.validate(config)
	if err != nil {
		return errors.New("init logger failed, err: " + err.Error())
	}

	//init channels
	lh.lock = new(sync.Mutex)
	lh.dataChan = make(chan string, MAX_LOGCHAN_BUFFER_SIZE)

	//init arguments
	if !isPathExist(lh.fileDir) {
		os.Mkdir(lh.fileDir, 0755)
	}
	filePath := joinFilePath(lh.fileDir, lh.fileName)
	if lh.sliceType == SliceType_Size {
		for i := 0; i != lh.fileCount; i++ {
			if isPathExist(filePath + "." + strconv.Itoa(i)) {
				break
			}
			lh.fileIndex = i
		}
	} else {
		lh.curDate = getCurrentDate()
	}
	lh.fileHandle, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("init logger failed, err: " + err.Error())
	}
	lh.checkSliceFile()

	//start go workers
	go lh.outputDataWorker()
	go lh.sliceFileWorker()

	return nil
}

//validate config params
func (lh *logHandler) validate(config string) error {
	js, err := simplejson.NewJson([]byte(config))
	if err != nil {
		return errors.New("validate log config failed, err: " + err.Error())
	}
	//validate console
	console, err := js.Get("console").Bool()
	if err == nil {
		lh.console = console
	} else {
		lh.console = DEFAULT_CONSOLE_VALUE
	}
	//validate sliceType
	slicetype, err := js.Get("sliceType").Int()
	if err == nil && slicetype >= int(SliceType_Size) && slicetype <= int(SliceType_Date) {
		lh.sliceType = SliceType(slicetype)
	} else {
		lh.sliceType = DEFAULT_SLICETYPE_VALUE
	}
	//validate prefix
	prefix, err := js.Get("prefix").String()
	if err == nil {
		lh.prefix = prefix
	} else {
		lh.prefix = DEFAULT_PREFIX_VALUE
	}
	//validate level
	level, err := js.Get("level").Int()
	if err == nil && level >= int(Level_All) && level <= int(Level_Fatal) {
		lh.level = level
	} else {
		lh.level = DEFAULT_LEVEL_VALUE
	}
	//validate fileDir
	fileDir, err := js.Get("fileDir").String()
	if err == nil {
		lh.fileDir = fileDir
	} else {
		lh.fileDir = DEFAULT_FILEDIR_VALUE
	}
	//validate fileName
	fileName, err := js.Get("fileName").String()
	if err == nil {
		lh.fileName = fileName
	} else {
		lh.fileName = DEFAULT_FILENAME_VALUE
	}
	//validate fileCount
	fileCount, err := js.Get("fileCount").Int()
	if err == nil && fileCount >= 1 {
		lh.fileCount = fileCount
	} else {
		lh.fileCount = DEFAULT_FILECOUNT_VALUE
	}
	//validate fileSize
	fileSize, err := js.Get("fileSize").Uint64()
	if err == nil {
		lh.fileSize = fileSize
	} else {
		lh.fileSize = DEFAULT_FILESIZE_VALUE
	}
	return nil
}

//consume log data worker
func (lh *logHandler) outputDataWorker() {
	for {
		select {
		case data := <-lh.dataChan:
			lh.outputData(data)
		}
	}
}

//check file slice worker
func (lh *logHandler) sliceFileWorker() {
	timer := time.NewTicker(time.Duration(CHECK_SLICE_FILE_INTERVAL) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			lh.checkSliceFile()
		}
	}
}

//input log data into channels
func (lh *logHandler) inputData(level int, data string) {
	if level < lh.level || level < Level_All || level > Level_Fatal {
		return
	}
	content := time.Now().Format("2006-01-02 15:04:05")
	switch level {
	case Level_All:
		data = content + " [All]:" + lh.prefix + data
	case Level_Debug:
		data = content + " [Debug]:" + lh.prefix + data
	case Level_Info:
		data = content + " [Info]:" + lh.prefix + data
	case Level_Warn:
		data = content + " [Warn]:" + lh.prefix + data
	case Level_Error:
		data = content + " [Error]:" + lh.prefix + data
	case Level_Fatal:
		data = content + " [Fatal]:" + lh.prefix + data
	default:
		return
	}
	lh.dataChan <- data
}

//output log data to console or file
func (lh *logHandler) outputData(data string) {
	if lh.console {
		lh.outputToConsole(data)
	}
	lh.outputToFile(data)
}

//output log data to file
func (lh *logHandler) outputToFile(data string) {
	lh.lock.Lock()
	defer lh.lock.Unlock()
	if lh.fileHandle == nil {
		return
	}
	lh.fileHandle.Write([]byte(data))
}

//output data to console
func (lh *logHandler) outputToConsole(data string) {
	fmt.Println(data)
}

//check slice file
func (lh *logHandler) checkSliceFile() {
	lh.lock.Lock()
	defer lh.lock.Unlock()
	if lh.sliceType == SliceType_Size {
		lh.sliceFileBySize()
	} else {
		lh.sliceFileByDate()
	}
}

//slice file by date
func (lh *logHandler) sliceFileByDate() {
	tm := getCurrentDate()
	if !tm.After(*lh.curDate) { //check if need to slice by date
		return
	}
	filePath := joinFilePath(lh.fileDir, lh.fileName)
	logFileTemp := filePath + "." + lh.curDate.Format("2006-01-02")
	if !isPathExist(logFileTemp) {
		if lh.fileHandle != nil {
			lh.fileHandle.Close()
		}
		os.Rename(filePath, logFileTemp)
		lh.curDate = getCurrentDate()
		lh.fileHandle, _ = os.Create(filePath)
	}
}

//slice file by size
func (lh *logHandler) sliceFileBySize() {
	filePath := joinFilePath(lh.fileDir, lh.fileName)
	if lh.fileCount <= 1 {
		return
	}
	curSize := getFileSize(filePath)
	if curSize < lh.fileSize {
		return
	}
	lh.fileIndex = int(lh.fileIndex%lh.fileCount + 1)
	if lh.fileHandle != nil {
		lh.fileHandle.Close()
	}
	logFileTemp := filePath + "." + strconv.Itoa(lh.fileIndex)
	if isPathExist(logFileTemp) {
		os.Remove(logFileTemp)
	}
	os.Rename(filePath, logFileTemp)
	lh.fileHandle, _ = os.Create(filePath)
}
