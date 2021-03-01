package common

import (
	"fmt"
	"os"
	"time"
)

var Results = make(chan string)
var Worker = 0
var Start = true
var LogSucTime int64
var LogErr bool
var LogErrTime int64

func LogSuccess(result string) {
	Worker++
	LogSucTime = time.Now().Unix()
	if Start {
		go SaveLog()
		Start = false
	}
	Results <- result
}

func SaveLog() {
	for result := range Results {
		fmt.Println(result)
		if IsSave {
			WriteFile(result, Outputfile)
		}
		Worker--
	}
}

func WriteFile(result string, filename string) {
	var text = []byte(result + "\n")
	fl, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Printf("Open %s error, %v\n", filename, err)
		return
	}
	_, err = fl.Write(text)
	fl.Close()
	if err != nil {
		fmt.Printf("write %s error, %v\n", filename, err)
	}
}

func WaitSave() {
	for {
		if Worker == 0 {
			close(Results)
			return
		}
	}
}

func LogError(errinfo interface{}) {
	if LogErr {
		if (time.Now().Unix()-LogSucTime) > 10 && (time.Now().Unix()-LogErrTime) > 10 {
			fmt.Println(errinfo)
			LogErrTime = time.Now().Unix()
		}
	}
}
