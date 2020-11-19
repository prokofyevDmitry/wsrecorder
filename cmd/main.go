package main

import (
	"flag"
	"log"
	"time"
	"wsrecorder"
)

var flagAddr = flag.String("addr", "", "websocket address")
var flagInputFile = flag.String("i", "init_messages.txt", "text file containing the ws init stream")
var flagOutputFile = flag.String("o", "output.txt", "Output file name")
var flagDuration = flag.Duration("d", 5*time.Second, "Duration of the recording")

func main() {

	log.Printf("Started recorder")
	flag.Parse()
	wsrecorder.Record(*flagAddr, *flagInputFile, *flagOutputFile, *flagDuration)
}
