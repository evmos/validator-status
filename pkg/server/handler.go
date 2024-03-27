package server

import (
	"fmt"
	"net/http"

	"github.com/evmos/validator-status/pkg/logger"
)

const HomePrefix = "/"

func (s *Server) HandlerHome(w http.ResponseWriter, r *http.Request) {
	if SetHandlerCorsForOptions(r, &w) {
		return
	}

	body := []byte("hello hello")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		logger.LogError(fmt.Sprintf("could not write home response: %q", err))
	}

}
