package main

import (
	"log"
	"web2kafka/producer"
)

func main() {

	err := producer.KafkaSend("", "Hello!")
	if err != nil {
		log.Printf("sending logs failed: %v", err)
	}

}
