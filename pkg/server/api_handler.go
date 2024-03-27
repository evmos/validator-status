package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/evmos/validator-status/pkg/database"
	"github.com/evmos/validator-status/pkg/logger"
)

const APIPrefix = "/api"

type AllValidatorsResponse struct {
	Values []database.GetInfoBetweenBlocksRow `json:"values"`
}

type SimpleResponse struct {
	Values []database.GetValidatorInfoBetweenBlocksRow `json:"values"`
}

func writeError(w http.ResponseWriter, error string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	errorValue := `{"value":"` + error + `"}`
	if _, err := w.Write([]byte(errorValue)); err != nil {
		logger.LogError(fmt.Sprintf("could not write home response: %q", err))
	}
}

func (s *Server) HandlerAPI(w http.ResponseWriter, r *http.Request) {
	if SetHandlerCorsForOptions(r, &w) {
		return
	}
	// Params
	start := r.URL.Query().Get("start")
	if start == "" {
		logger.LogError("start is missing")
		return
	}

	startNumber, err := strconv.Atoi(start)
	if err != nil {
		// Invalid number
		errMsg := fmt.Sprintf("start is an invalid number %q", err)
		writeError(w, errMsg)
		return
	}

	endNumber := startNumber
	end := r.URL.Query().Get("end")
	if end != "" {
		endNumber, err = strconv.Atoi(start)
		if err != nil {
			errMsg := fmt.Sprintf("end is an invalid number %q", err)
			writeError(w, errMsg)
			return
		}
	}

	if endNumber > startNumber+s.config.APIMaxLimit {
		endNumber = startNumber + s.config.APIMaxLimit
	}

	validator := r.URL.Query().Get("validator")
	if validator != "" {
		params := database.GetValidatorInfoBetweenBlocksParams{
			Height:          int64(startNumber),
			Height_2:        int64(endNumber),
			OperatorAddress: validator,
		}
		res, err := s.db.GetValidatorInfoBetweenBlocks(context.Background(), params)
		if err != nil {
			errMsg := fmt.Sprintf("error getting validator info between blocks %q", err)
			writeError(w, errMsg)
			return
		}

		bytes, err := json.Marshal(SimpleResponse{Values: res})
		if err != nil {
			errMsg := fmt.Sprintf("error converting validator info between blocks %q", err)
			writeError(w, errMsg)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(bytes); err != nil {
			logger.LogError(fmt.Sprintf("could not write home response: %q", err))
		}
		return
	}

	params := database.GetInfoBetweenBlocksParams{
		FromHeight: int64(startNumber),
		ToHeight:   int64(endNumber),
	}
	res, err := s.db.GetInfoBetweenBlocks(context.Background(), params)
	if err != nil {
		errMsg := fmt.Sprintf("error getting info between blocks %q", err)
		writeError(w, errMsg)
		return
	}

	bytes, err := json.Marshal(AllValidatorsResponse{Values: res})
	if err != nil {
		errMsg := fmt.Sprintf("error converting info between blocks %q", err)
		writeError(w, errMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		logger.LogError(fmt.Sprintf("could not write home response: %q", err))
	}
}
