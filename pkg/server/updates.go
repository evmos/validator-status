package server

import (
	"context"
	"database/sql"
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
	// Stop the RunUpdates for loop
	s.updateTicker.Stop()
	s.doneTicker <- true

	// Wait for the last UpdateValidators to finish
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()
}

const (
	maxBlocksToProcessAtTheTime = 10
)

func (s *Server) UpdateValidators() {
	// We are using TryLock here to avoid running more than 1 UpdateValidators at the same time
	if !s.updateMutex.TryLock() {
		return
	}
	defer s.updateMutex.Unlock()
	logger.LogDebug("calling update validators")

	// Get the chain latest height
	chainHeight, err := s.cosmos.GetChainHeight()
	if err != nil {
		logger.LogError(fmt.Sprintf("could not get chain height info: %s", err.Error()))
		return
	}

	// Get the latest height in the database
	height, err := s.db.GetLatestHeight(context.Background())
	if err == sql.ErrNoRows {
		height = int64(chainHeight.CurrentHeight) - int64(s.config.PruneOffset)
	} else if err != nil {
		logger.LogError(fmt.Sprintf("could not latest block in db: %s", err.Error()))
		return
	}

	logger.LogDebug(fmt.Sprintf("height %d chainHeight %d", height, chainHeight.CurrentHeight))

	stopHeight := chainHeight.CurrentHeight
	if stopHeight-int(height) > maxBlocksToProcessAtTheTime {
		stopHeight = int(height) + maxBlocksToProcessAtTheTime
	}
	logger.LogDebug(fmt.Sprintf("height %d stopHeight %d", height, stopHeight))

	// Move to the first not indexed block
	height++

	for height <= int64(stopHeight) {
		logger.LogInfo(fmt.Sprintf("updating the values for height: %d", height))

		if err := s.cosmos.UpdateValidatorsTable(height); err != nil {
			logger.LogError(fmt.Sprintf("problem updating validators table %s", err.Error()))
			return
		}

		if err := s.cosmos.UpdateMissingTable(height); err != nil {
			logger.LogError(fmt.Sprintf("problem updating missng table %s", err.Error()))
			return
		}

		height++
	}
}
