package jsoner

import (
	"github.com/sunnywalden/sync-data/pkg/types"
	jsoniter "github.com/json-iterator/go"
	"log"
)

func JsonDecode(bodyC []byte) (types.OAToken, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var jsonMap types.OAToken
	err := json.Unmarshal(bodyC, &jsonMap)
	if err != nil {
		log.Fatalf("json decode err!%s\n", err.Error())
		return types.OAToken{}, err
	} else {
		return jsonMap, nil
	}
}
