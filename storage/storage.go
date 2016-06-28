package storage

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto"
)

//RootDir is the path to store files under
var RootDir = "storage"

var cache = make(map[string]*otto.Script)

//Set stores the script in storage
func Set(name string, script *otto.Script) (hash string, err error) {

	src := script.String()

	data := []byte(src)

	hasher := sha256.New()

	hasher.Write(data)

	hash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	hash = strings.Replace(hash, "=", "", -1)

	pathName := getJSPath(name, hash)

	//store in the cache
	cache[pathName] = script

	//write to disk
	pathOnly := filepath.Dir(pathName)
	os.MkdirAll(pathOnly, 0644)
	err = ioutil.WriteFile(pathName, data, 0644)

	return hash, err
}

//Get retrieves the script
func Get(name string, hash string) (script *otto.Script, err error) {

	pathName := getJSPath(name, hash)

	script = cache[pathName]

	if script == nil {
		if data, err := ioutil.ReadFile(pathName); err == nil {

			src := string(data)

			vm := otto.New()

			script, err = vm.Compile(hash, src)
		}
	}

	return script, err
}

func getJSPath(name string, hash string) (path string) {
	return fmt.Sprintf("%s/%s/%s.js", RootDir, name, hash)
}
