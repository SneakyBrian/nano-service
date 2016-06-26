package run

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/SneakyBrian/nano-service/storage"
	"github.com/robertkrimen/otto"
)

// HandleRun handles running the script
func HandleRun(w http.ResponseWriter, r *http.Request) {

	var err error

	query := r.URL.Query()

	//get the name and the hash from the querystring
	name := query.Get("name")
	hash := query.Get("hash")

	if name != "" && hash != "" {
		//remove the name and the hash from the query object to avoid confusion
		query.Del("name")
		query.Del("hash")

		if script, err := storage.Get(name, hash); err == nil {

			//get the configured VM
			vm := getVM(query)

			//run the script
			//if value, err := vm.Run(script); err == nil {
			if value, err := runUnsafe(vm, script); err == nil {

				if responseBody, err := value.ToString(); err == nil {

					w.Write([]byte(responseBody))

					return
				}
			}
		}

		w.Write([]byte(err.Error()))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func getVM(query url.Values) (vm *otto.Otto) {

	//create the runtime
	vm = otto.New()

	//bind the support functions

	//simple synchronous http get with $get
	//Example: var result = $get(url);
	vm.Set("$get", func(call otto.FunctionCall) otto.Value {

		if response, err := http.Get(call.Argument(0).String()); err == nil {
			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			bodyString := string(body)

			value, _ := otto.ToValue(bodyString)

			return value
		}

		return otto.Value{}

	})

	//simple synchronous http post with $post
	//Example: var result = $post(url, "application/json", JSON.stringify(obj))
	vm.Set("$post", func(call otto.FunctionCall) otto.Value {

		buf := bytes.NewBuffer([]byte(call.Argument(2).String()))

		if response, err := http.Post(call.Argument(0).String(), call.Argument(1).String(), buf); err == nil {
			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			bodyString := string(body)

			value, _ := otto.ToValue(bodyString)

			return value
		}

		return otto.Value{}

	})

	//add in the query values as the global $query object
	vm.Set("$query", query)

	return vm
}

var errHalt = errors.New("Script Time Overrun")

const maxScriptTime = 10

func runUnsafe(vm *otto.Otto, unsafe interface{}) (value otto.Value, err error) {

	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if caught := recover(); caught != nil {
			if caught == errHalt {
				fmt.Fprintf(os.Stderr, "Code execution time exceeded! Stopping after %v\n", duration)
				return
			}
			panic(caught) // Something else happened, repanic!
		}
		fmt.Fprintf(os.Stderr, "Ran code successfully in %v\n", duration)
	}()

	vm.Interrupt = make(chan func(), 1) // The buffer prevents blocking

	go func() {
		time.Sleep(maxScriptTime * time.Second) // Stop after maxScriptTime seconds
		vm.Interrupt <- func() {
			panic(errHalt)
		}
	}()

	value, err = vm.Run(unsafe) // Here be dragons (risky code)

	return value, err
}
