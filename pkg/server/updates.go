package server

import (
	"fmt"

	"github.com/evmos/validator-status/pkg/logger"
)

func (s *Server) RunUpdates() {
	// Keep updating until close
	for {
		select {
		case <-s.doneTicker:
			return
		case <-s.updateTicker.C:
			s.UpdateValidators()
		}
	}
}

func (s *Server) StopUpdates() {
	s.updateTicker.Stop()
	s.doneTicker <- true
}

func (s *Server) UpdateValidators() {
	logger.LogInfo("updating the values...")
	if err := s.cosmos.UpdateValidatorsTable(); err != nil {
		logger.LogError(fmt.Sprintf("problem updating validators table %s", err.Error()))
		return
	}
	// TODO: get the correct height
	if err := s.cosmos.UpdateMissingTable(0); err != nil {
		logger.LogError(fmt.Sprintf("problem updating missng table %s", err.Error()))
		return
	}
}
