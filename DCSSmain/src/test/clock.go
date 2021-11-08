package main

import (
	"github.com/Alan-Lxc/crypto_contest/src/control"
)

func main() {
	//label := flag.Int("l", 1, "Enter node label")
	//metadataPath := flag.String("path", "/mpss/metadata", "Enter the metadata path")
	////metadataPath := flag.String("path", "/mpss/metadata", "Enter the metadata path")
	//flag.Parse()
	metadataPath := "/home/alan/Desktop/crypto_contest/src/metadata"
	clock, _ := control.New(metadataPath)
	clock.Connect()
	clock.StartHandoff()
}
