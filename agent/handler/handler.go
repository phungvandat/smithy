package handler

import (
	"encoding/json"
	"net/http"

	agentConfig "github.com/dwarvesf/smithy/agent/config"
	handlerCommon "github.com/dwarvesf/smithy/common/handler"
)

// Expose return handler for expose metadata, connection for dashboard
func Expose(cfg *agentConfig.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != cfg.SerectKey {
			handlerCommon.EncodeJSONError(errorMissingAuth{}, w)
			return
		}

		err := encodeAgentConfig(w, cfg)
		if err != nil {
			handlerCommon.EncodeJSONError(err, w)
			return
		}
	}
}

func encodeAgentConfig(w http.ResponseWriter, cfg *agentConfig.Config) error {
	return json.NewEncoder(w).Encode(cfg)
}

type errorMissingAuth struct{}

func (errorMissingAuth) Error() string {
	return "missing auth"
}

// StatusCode implement status code for error missing auth
func (errorMissingAuth) StatusCode() int {
	return http.StatusUnauthorized
}
