package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/josnelihurt/go-stuff/data-harvester/cmd/server"
)

//interruptHandler used to handle interruptions
func interruptHandler(errChan chan<- error) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errChan <- fmt.Errorf("%s", <-c)
}

func main() {

	server := server.NewServer()
	go interruptHandler(server.GetErrorChannel())
	go server.Run()

	log.Fatalln(<-server.GetErrorChannel())

}
