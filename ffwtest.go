package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dollarkillerx/urllib"
)

func main() {
	li := 100
	if len(os.Args) == 2 {
		c := os.Args[1]
		lc,err := strconv.Atoi(c)
		if err == nil {
			li = lc
		}
	}
	log.Println(li)
	file, err := ioutil.ReadFile("img.jpg")
	if err != nil {
		log.Fatalln(err)
	}

	outlog := make(chan string, 1000)
	go writeLog(outlog)
	go func() {
		time.Sleep(120 * time.Second)
		close(outlog)
		panic("")
	}()

	limit := make(chan bool,li)
	i := 0
	for {
		i++
		limit<-true
		go func(id int) {
			defer func() {
				<-limit
			}()
			filePath := fmt.Sprintf("http://192.168.88.11:8089/file/file-%d.jpg",id)
			log.Println("In")
			retry, body, err := urllib.Post(filePath).SetTimeout(10).SetJson(file).ByteRetry(3)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("OoO")
			if retry != 200 {
				log.Println(string(body))
				return
			}

			outlog <- fmt.Sprintf("%s\n",filePath)
		}(i)
	}
}

func writeLog(out chan string) {
	logFile, err := os.Create("filelog.log")
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		logFile.Close()
	}()
c:
	for {
		select {
		case path, ex := <-out:
			if !ex {
				break c
			}
			log.Println(path)
			logFile.Write([]byte(path))
		}
	}
}