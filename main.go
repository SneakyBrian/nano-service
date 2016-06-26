package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/SneakyBrian/nano-service/deploy"
	"github.com/SneakyBrian/nano-service/retrieve"
	"github.com/SneakyBrian/nano-service/run"
)

func main() {

	portNumPtr := flag.Uint("port", 8181, "port number for webserver")

	flag.Parse()

	fmt.Printf("Starting nano-service on port %d\n", *portNumPtr)

	http.HandleFunc("/deploy/", deploy.HandleDeploy)
	http.HandleFunc("/run/", run.HandleRun)
	http.HandleFunc("/retrieve/", retrieve.HandleRetrieve)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portNumPtr), nil))

	fmt.Printf("Exiting nano-service\n")

}
