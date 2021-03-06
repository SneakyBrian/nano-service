package run

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SneakyBrian/nano-service/storage"
	"github.com/robertkrimen/otto"
)

//MaxScriptTime is the maximum number of seconds that a script is permitted to run
var MaxScriptTime uint = 10

// HandleRun handles running the script
func HandleRun(w http.ResponseWriter, r *http.Request) {

	var err error

	//get the name and the hash from the path
	urlPart := strings.Split(r.URL.Path, "/")

	name := urlPart[2]
	hash := urlPart[3]

	if name != "" && hash != "" {

		if script, err := storage.Get(name, hash); err == nil {

			//get the configured VM
			vm := getVM(r)

			//run the script
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

func getVM(r *http.Request) (vm *otto.Otto) {

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

	//add in the request values as the global objects

	r.ParseForm()

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	cookies := r.Cookies()
	cookieMap := make(map[string]string)

	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}

	vm.Set("$uri", r.RequestURI)
	vm.Set("$headers", r.Header)
	vm.Set("$cookies", cookieMap)
	vm.Set("$params", r.Form)
	vm.Set("$body", string(body))

	return vm
}

var errHalt = errors.New("Script Time Overrun")

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
		time.Sleep(time.Duration(MaxScriptTime) * time.Second) // Stop after MaxScriptTime seconds
		vm.Interrupt <- func() {
			panic(errHalt)
		}
	}()

	value, err = vm.Run(unsafe) // Here be dragons (risky code)

	return value, err
}
