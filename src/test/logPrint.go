package main

import (
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("../metadata/testLog/test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	logger := log.New(file, "", log.LstdFlags|log.Llongfile)
	logger.Println("日志1.")
	logger.Println("日志23")
}
