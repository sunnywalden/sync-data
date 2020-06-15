package helper

import (
	"github.com/sunnywalden/sync-data/pkg/types"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)


// ResponseWithJson, encode response data
func ResponseWithJson(w http.ResponseWriter, code int, payload types.Response) {

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	response, err := json.Marshal(payload)
	if err != nil {
		code = http.StatusInternalServerError
		response = []byte{}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

