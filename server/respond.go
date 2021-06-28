package server

import (
	"encoding/json"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/global"
)

func Respond(w http.ResponseWriter, payload interface{}, statusCode int, args ...interface{}) {
	w.WriteHeader(statusCode)
	switch v := payload.(type) {
	case *global.DefaultRespBody:
		v.Message = args[0]
		v.Error = args[1]
	case *global.LoginRespBody:
		v.Message = args[0]
		v.Error = args[1]
		v.Token = args[2]
	case *global.SecretpageRespBody:
		v.Message = args[0]
		v.Error = args[1]
		v.Data = args[2]
	case *global.ViewCoinsRespBody:
		v.Message = args[0]
		v.Error = args[1]
		v.Coins = args[2]
	case *global.TxnRespBody:
		v.Message = args[0]
		v.Error = args[1]
		v.TxnId = args[2]
	}
	json.NewEncoder(w).Encode(payload)
}
