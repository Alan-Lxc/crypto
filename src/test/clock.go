package main

import (
	"flag"
	"github.com/Alan-Lxc/crypto_contest/src/control"
)

func main() {
	//label := flag.Int("l", 1, "Enter node label")
	metadataPath := flag.String("path", "/mpss/metadata", "Enter the metadata path")
	//metadataPath := flag.String("path", "/mpss/metadata", "Enter the metadata path")
	flag.Parse()

	clock, _ := control.New(*metadataPath)
	clock.Connect()
	clock.StartHandoff()
}
