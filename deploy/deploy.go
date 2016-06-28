package deploy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/SneakyBrian/nano-service/storage"
	"github.com/robertkrimen/otto"
)

//HandleDeploy handles the deploy of code
func HandleDeploy(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		var err error

		urlPart := strings.Split(r.URL.Path, "/")

		name := urlPart[2]

		if name != "" {

			defer r.Body.Close()

			body, _ := ioutil.ReadAll(r.Body)

			bodyString := string(body)

			vm := otto.New()

			//check it compiles
			if script, err := vm.Compile(name, bodyString); err == nil {

				if hash, err := storage.Set(name, script); err == nil {
					fmt.Printf("Deployed Script %s (%s)\n", name, hash)
					w.Write([]byte(hash))
					return
				}
			}

			w.Write([]byte(err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

}
