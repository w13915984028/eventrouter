package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/crewjam/rfc5424"
)

func handler(_ http.ResponseWriter, r *http.Request) {
	log.Printf("request method=%s from=%s", r.Method, r.RemoteAddr)
	if r.Body == nil {
		return
	}
	defer r.Body.Close()

	m := new(rfc5424.Message)
	discardBuf := make([]byte, 1)
	for {
		_, err := m.ReadFrom(r.Body)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Parsing rfc5424 message failed: %+v", err)
		}
		log.Printf("%s", m.Message)

		// read the extraneous \n at the end of the message and discard
		_, _ = io.ReadFull(r.Body, discardBuf)
	}
}

func main() {
	log.Println("starting httpsink server")
	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
