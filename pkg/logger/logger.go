package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// init the logger with a fixed format
// func init() {
// 	logrus.SetFormatter(&logrus.JSONFormatter{
// 		TimestampFormat: "2006-01-02 15:05:06",
// 	})
// 	logrus.SetReportCaller(false)
// }

// wirte log with message and file name
func Write(msg string, filename string) {
	setOutPutFile(logrus.InfoLevel, filename)
	logrus.Info(msg)
}

// write the debug log ,args as fields = {"a":"b"}
func Debug(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.DebugLevel, "debug")
	logrus.WithFields(fields).Debug(args...)
}

func Info(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.InfoLevel, "info")
	logrus.WithFields(fields).Info(args...)
}

func Warn(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.WarnLevel, "warn")
	logrus.WithFields(fields).Warn(args...)
}

func Fatal(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.FatalLevel, "fatal")
	logrus.WithFields(fields).Fatal(args...)
}

func Error(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.ErrorLevel, "error")
	logrus.WithFields(fields).Error(args...)
}
func Panic(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.PanicLevel, "panic")
	logrus.WithFields(fields).Panic(args...)
}
func Trace(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.TraceLevel, "trace")
	logrus.WithFields(fields).Trace(args...)
}

// write the log
func setOutPutFile(level logrus.Level, logName string) {
	//check file is exist or not
	if _, err := os.Stat("./runtime/log"); os.IsNotExist(err) {
		err = os.MkdirAll("./runtime/log", 0777)
		if err != nil {
			panic(fmt.Errorf("create log dir '%s' error: %s", "./runtime/log", err))
		}
	}
	//create the time format and file name
	timeStr := time.Now().Format("2006-01-02")
	fileName := path.Join("./runtime/log", logName+"_"+timeStr+".log")

	//if exist , write in, or create and write in
	var err error
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Open log file err", err)
		return
	}
	//wirte in file and set the level
	logrus.SetOutput(file)
	logrus.SetLevel(level)

}

// set log config
func LoggerToFile() gin.LoggerConfig {
	//check file is exist or not
	if _, err := os.Stat("./runtime/log"); os.IsNotExist(err) {
		err = os.MkdirAll("./runtime/log", 0777)
		if err != nil {
			panic(fmt.Errorf("create log dir '%s' error: '%s'", "./runtime/log", err))
		}
	}

	// format the timetamp and combin the file name
	timeStr := time.Now().Format("2006-01-02")
	fileName := path.Join("./runtime/log", "success_"+timeStr+".log")

	//if file is exist to write in, or create and write in
	file, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	//set the print info type
	var conf = gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s-%s \"%s %s %s %d %s \"%s\" %s\"\n",
				param.TimeStamp.Format("2006-01-02 15:05:05"),
				param.ClientIP,
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
		Output: io.MultiWriter(os.Stdout, file),
	}
	return conf
}

// error capture
func Recover(c *gin.Context) {
	defer func() {
		//check file is exist or not
		if err := recover(); err != nil {
			if _, errDir := os.Stat("./runtime/log"); os.IsNotExist(errDir) {
				errDir = os.MkdirAll("./runtime/log", 0777)
				if errDir != nil {
					panic(fmt.Errorf("create log dir '%s' error: %s", "./runtime/log", errDir))
				}
			}
			// format the timetamp and combin the file name
			timeStr := time.Now().Format("2006-01-02")
			fileName := path.Join("./runtime/log", "error_"+timeStr+".log")

			//if exist to write in, or create and write in
			file, errFile := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if errFile != nil {
				fmt.Println(errFile)
			}

			// format the content to write in
			timeFileStr := time.Now().Format("2006-01-02 15:05:05")
			file.WriteString("panic error time:" + timeFileStr + "\n")
			file.WriteString(fmt.Sprintf("%v", err) + "\n")
			file.WriteString("stackTrace from panic:" + string(debug.Stack()) + "\n")
			file.Close()
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  fmt.Sprintf("%v", err),
			})
			c.Abort()
		}
	}()
	c.Next()
}
