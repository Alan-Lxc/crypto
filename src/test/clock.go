package main

import (
	"flag"
	"github.com/Alan-Lxc/crypto_contest/src/control"
)

func main() {
	metadataPath := flag.String("path", "/mpss/metadata", "Enter the metadata path")
	flag.Parse()

	clock, _ := control.New(*metadataPath)
	clock.Connect()
	clock.StartHandoff()
}
