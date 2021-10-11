package main

import (
	"flag"
	"github.com/Alan-Lxc/crypto_contest/src/control"
)

func main() {
	//label := flag.Int("l", 1, "Enter node label")
	counter := flag.Int("c", 1, "Enter number of nodes")
	degree := flag.Int("d", 1, "Enter the polynomial degree")
	metadataPath := flag.String("path", "/mpss/metadata", "Enter the metadata path")
	S0 := flag.String("secret", 1, "Enter the secret")
	//metadataPath := flag.String("path", "/mpss/metadata", "Enter the metadata path")
	flag.Parse()

	clock, _ := control.New(*degree, *counter, *metadataPath, *S0)
	clock.Connect()
	clock.StartHandoff()
}
