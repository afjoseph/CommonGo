package print

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	LOG_SILENCE = iota
	LOG_INFO    = iota
	LOG_DEBUG   = iota
)

var logLevel uint = LOG_INFO

func SetLevel(level uint) {
	logLevel = level
}

func absToRelFilepath(absFilepath string) (string, bool) {
	if len(absFilepath) == 0 {
		return "", false
	}
	arr := strings.Split(absFilepath, "/")
	if len(arr) < 2 {
		return "", false
	}
	relFilepathArr := arr[len(arr)-2:]
	relFilepath := strings.Join(relFilepathArr, "/")
	return relFilepath, true
}

func SetLogFile(fileNamePrefix string) {
	var fileName string
	if fileNamePrefix == "" {
		fileName = fmt.Sprintf("%d.log", time.Now().Unix())
	} else {
		fileName = fmt.Sprintf("%s_%d.log", fileNamePrefix, time.Now().Unix())
	}
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		// Ignore errors
		return
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(5)
}

// absToSimpleFilePath takes an absolute path and returns only the last two nodes.
// This is used exclusively for printing debug statements
// Example:
//     Input -> /tmp/aaa/bbb/ccc
//     Output -> bbb/ccc
func absToSimpleFilePath(absFilepath string) (string, bool) {
	if len(absFilepath) == 0 {
		return "", false
	}
	arr := strings.Split(absFilepath, "/")
	if len(arr) < 3 {
		return "", false
	}
	relFilepathArr := arr[len(arr)-2:]
	relFilepath := strings.Join(relFilepathArr, "/")
	return relFilepath, true
}

func InfoFunc() {
	if logLevel < LOG_INFO {
		return
	}
	pc, absFilepath, lineNum, ok := runtime.Caller(1)
	if !ok {
		return
	}
	relFilepath, ok := absToSimpleFilePath(absFilepath)
	if !ok {
		return
	}
	details := runtime.FuncForPC(pc)
	fmt.Printf("[+:%s:%d] %s()\n", relFilepath, lineNum, details.Name())
}

func Infoln(text string) {
	if logLevel < LOG_INFO {
		return
	}
	color.Set(color.FgHiBlue)
	fmt.Print("[+] ")
	fmt.Println(text)
	color.Unset()
}

func Infof(format string, v ...interface{}) {
	if logLevel < LOG_INFO {
		return
	}
	color.Set(color.FgHiBlue)
	fmt.Print("[+] ")
	fmt.Printf(format, v...)
	color.Unset()
}

// XXX One caveat here: you can't use error propagation using %w
func Errorf(format string, v ...interface{}) error {
	content := fmt.Sprintf(format, v...)

	_, absFilepath, lineNum, _ := runtime.Caller(1)
	relFilepath, ok := absToRelFilepath(absFilepath)
	var prefix string
	if ok {
		prefix = fmt.Sprintf("[!%s:%d]", relFilepath, lineNum)
	}
	return fmt.Errorf("%s: %s", prefix, content)
}

func Warnln(v ...interface{}) {
	color.Set(color.FgRed)
	fmt.Println(v...)
	color.Unset()
}

func Warnf(format string, v ...interface{}) {
	color.Set(color.FgRed)
	fmt.Printf(format, v...)
	color.Unset()
}

func Debugln(text string) {
	if logLevel < LOG_DEBUG {
		return
	}
	color.Set(color.FgYellow)
	_, absFilepath, lineNum, _ := runtime.Caller(1)
	relFilepath, ok := absToSimpleFilePath(absFilepath)
	if !ok {
		log.Print("[D] ")
	} else {
		fmt.Printf("[D:%s:%d] ", relFilepath, lineNum)
	}
	log.Println(text)
	color.Unset()
}

func Debugf(format string, v ...interface{}) {
	if logLevel < LOG_DEBUG {
		return
	}
	color.Set(color.FgYellow)
	_, absFilepath, lineNum, _ := runtime.Caller(1)
	relFilepath, ok := absToSimpleFilePath(absFilepath)
	if !ok {
		fmt.Print("[D] ")
	} else {
		fmt.Printf("[D:%s:%d] ", relFilepath, lineNum)
	}
	fmt.Printf(format, v...)
	color.Unset()
}

func DebugFunc() {
	if logLevel < LOG_DEBUG {
		return
	}
	color.Set(color.FgYellow)
	pc, absFilepath, lineNum, ok := runtime.Caller(1)
	if !ok {
		return
	}
	relFilepath, ok := absToSimpleFilePath(absFilepath)
	if !ok {
		return
	}
	details := runtime.FuncForPC(pc)
	fmt.Printf("[+%s:%d] %s()\n", relFilepath, lineNum, details.Name())
	color.Unset()
}
