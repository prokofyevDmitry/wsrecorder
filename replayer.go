package wsrecorder

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
	"time"
	"wsrecorder/internal/errhdl"

	"github.com/gorilla/websocket"
)

func Record(addr string, inputFilePath string, outputFilePath string, dration time.Duration) {
	log.Printf("Started recorder")

	flag.Parse()

	// websocket connection
	u, err := url.Parse(addr)
	errhdl.PanicIf(err)
	log.Printf("Connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	errhdl.PanicIf(err)
	defer c.Close()
	log.Printf("Connected to %s", u.String())

	// Open file read and write
	f, err := os.Open(*flagInputFile)
	panicIf(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	f2, err := os.Create(*flagOutputFile)
	panicIf(err)
	defer f2.Close()
	w := bufio.NewWriter(f2)

	// Send init messages
	for s.Scan() {
		// bMsg, _ := json.Marshal(s.Text())
		err = c.WriteMessage(websocket.TextMessage, []byte(s.Text()))
		log.Printf("send %s", s.Text())
		panicIf(err)
	}

	// launch the listenning of new mgs
	new_mgs := make(chan []byte)
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			new_mgs <- msg
		}
	}()

	// store the new messages and wait for the end of the recording
	timer := time.NewTimer(*flagDuration)
	defer timer.Stop()
	msgCount := 0
	end := false
	for !end {
		select {
		case <-timer.C:
			log.Printf("close conn")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			logIf(err)
			log.Printf("record stopped after %v", flagDuration)
			end = true

		case msg := <-new_mgs:
			for _, b := range msg {
				logIf(w.WriteByte(b))
			}
			logIf(w.WriteByte('\n'))
			msgCount++
		}
	}

	log.Printf("wrote %d lines in %s", msgCount, *flagOutputFile)
}
