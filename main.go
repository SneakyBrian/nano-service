package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/SneakyBrian/nano-service/deploy"
	"github.com/SneakyBrian/nano-service/retrieve"
	"github.com/SneakyBrian/nano-service/run"
	"github.com/SneakyBrian/nano-service/storage"
)

func main() {

	portNumPtr := flag.Uint("port", 8181, "port number for webserver")
	maxScriptTimePtr := flag.Uint("maxScriptTime", 10, "the maximum number of seconds that a script is permitted to run")
	storageDirPtr := flag.String("storageDir", "storage", "the path to store files under")

	flag.Parse()

	run.MaxScriptTime = *maxScriptTimePtr
	storage.RootDir = *storageDirPtr

	fmt.Printf("Starting nano-service on port %d\n", *portNumPtr)
	fmt.Printf("Max Script Run Time Seconds: %d\n", run.MaxScriptTime)
	fmt.Printf("Storage Root Directory: %s\n", storage.RootDir)

	http.HandleFunc("/deploy/", deploy.HandleDeploy)
	http.HandleFunc("/run/", run.HandleRun)
	http.HandleFunc("/retrieve/", retrieve.HandleRetrieve)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portNumPtr), nil))

	fmt.Printf("Exiting nano-service\n")

}
