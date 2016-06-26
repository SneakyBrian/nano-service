package deploy

import (
	"io/ioutil"
	"net/http"

	"github.com/SneakyBrian/nano-service/storage"
	"github.com/robertkrimen/otto"
)

//HandleDeploy handles the deploy of code
func HandleDeploy(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		var err error

		query := r.URL.Query()
		name := query.Get("name")

		if name != "" {

			defer r.Body.Close()

			body, _ := ioutil.ReadAll(r.Body)

			bodyString := string(body)

			vm := otto.New()

			//check it compiles
			if script, err := vm.Compile(name, bodyString); err == nil {

				if hash, err := storage.Set(name, script); err == nil {
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
