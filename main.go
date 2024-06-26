package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"warships/gameclient"
	"warships/httpclient"
)

func main() {

	httpc := &httpclient.HttpClient{
		Client: &http.Client{Timeout: time.Second * 20},
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go gameclient.MainMenu(httpc)

	<-sigs
	fmt.Println("\nAbandoning game")
	httpc.Abandon()
	fmt.Println("Game abandoned")
}
