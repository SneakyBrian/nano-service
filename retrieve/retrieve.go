package retrieve

import (
	"net/http"

	"github.com/SneakyBrian/nano-service/storage"
)

// HandleRetrieve handles getting the script back out
func HandleRetrieve(w http.ResponseWriter, r *http.Request) {
	var err error

	query := r.URL.Query()

	//get the name and the hash from the querystring
	name := query.Get("name")
	hash := query.Get("hash")

	if name != "" && hash != "" {

		if script, err := storage.Get(name, hash); err == nil {

			src := script.String()

			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte(src))

			return
		}

		w.Write([]byte(err.Error()))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
