package main

import (
	"fmt"
	"log"

	"github.com/devwarrior777/atomicswap/libs"
)

func main() {
	fmt.Println("qtest")
	xzcAddr, err := libs.NewAddress("Xzc", true, "localhost", "dev", "dev")
	if err != nil {
		log.Fatalf("cannot get address %v\n", err)
	}
	fmt.Printf("New address: %s\n", xzcAddr)

	ltcAddr, err := libs.NewAddress("Ltc", true, "localhost", "dev", "dev")
	if err != nil {
		log.Fatalf("cannot get address %v\n", err)
	}
	fmt.Printf("New address: %s\n", ltcAddr)

	badAddr, err := libs.NewAddress("bad", true, "localhost", "dev", "dev")
	if err != nil {
		log.Fatalf("cannot get address %v\n", err)
	}
	fmt.Printf("New address: %s\n", badAddr)
}
