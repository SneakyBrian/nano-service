package main

//go:generate esc -o static.go static

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
	useLocalPtr := flag.Bool("useLocal", false, "use local filesystem for serving static resources for debugging purposes")
	maxScriptTimePtr := flag.Uint("maxScriptTime", 10, "the maximum number of seconds that a script is permitted to run")
	storageDirPtr := flag.String("storageDir", "storage", "the path to store files under")

	flag.Parse()

	run.MaxScriptTime = *maxScriptTimePtr
	storage.RootDir = *storageDirPtr

	fmt.Printf("Starting nano-service on port %d\n", *portNumPtr)
	fmt.Printf("Max Script Run Time Seconds: %d\n", run.MaxScriptTime)
	fmt.Printf("Storage Root Directory: %s\n", storage.RootDir)
	fmt.Printf("Use Local FileSystem for Static Resources: %t\n", *useLocalPtr)

	http.HandleFunc("/deploy/", deploy.HandleDeploy)
	http.HandleFunc("/run/", run.HandleRun)
	http.HandleFunc("/retrieve/", retrieve.HandleRetrieve)

	http.Handle("/static/", http.FileServer(FS(*useLocalPtr)))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portNumPtr), nil))

	fmt.Printf("Exiting nano-service\n")

}
