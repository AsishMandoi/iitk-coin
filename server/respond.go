package server

import (
	"encoding/json"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/global"
)

func Respond(w http.ResponseWriter, payload interface{}, status_code int, args ...string) {
	w.WriteHeader(status_code)
	switch v := payload.(type) {
	case *global.SignupRespBodyFormat:
		v.Message = args[0]
		v.Error = args[1]
	case *global.LoginRespBodyFormat:
		v.Message = args[0]
		v.Error = args[1]
		v.Token = args[2]
	case *global.SecretpageRespBodyFormat:
		v.Message = args[0]
		v.Error = args[1]
		v.Data = args[2]
	}
	json.NewEncoder(w).Encode(payload)
}
