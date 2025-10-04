package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 3 {
		log.Fatal("not enough args")
	}

	res, err := client.Calculate(os.Args[1], os.Args[2:]...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Result:", res)
}
