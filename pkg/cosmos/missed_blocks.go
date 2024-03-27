package cosmos

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/evmos/validator-status/pkg/database"
	"github.com/evmos/validator-status/pkg/logger"
)

type MissedBlockResponse struct {
	Info []struct {
		Address             string    `json:"address"`
		StartHeight         string    `json:"start_height"`
		IndexOffset         string    `json:"index_offset"`
		JailedUntil         time.Time `json:"jailed_until"`
		Tombstoned          bool      `json:"tombstoned"`
		MissedBlocksCounter string    `json:"missed_blocks_counter"`
	} `json:"info"`
}

func (c *Cosmos) getMissedBlocksByHeight(height string) (MissedBlockResponse, error) {
	// NOTE: assumes that we do not need pagination
	reqGo, err := http.NewRequest("GET", c.apiURL+"/cosmos/slashing/v1beta1/signing_infos", nil)
	reqGo.Header.Set("x-cosmos-block-height", height)

	resp, err := http.DefaultClient.Do(reqGo)
	if resp.StatusCode != http.StatusOK {
		return MissedBlockResponse{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return MissedBlockResponse{}, err
	}

	var v MissedBlockResponse
	err = json.Unmarshal(content, &v)
	if err != nil {
		return MissedBlockResponse{}, err
	}

	return v, nil
}

func (c *Cosmos) addMissingBlocksInfoToDB(height int64, missingBlocks MissedBlockResponse) error {
	for _, v := range missingBlocks.Info {
		val, err := c.db.GetValidatorByValidatorAddress(c.ctx, v.Address)
		if err != nil {
			// Validator is not in the database
			return err
		}
		logger.LogDebug(fmt.Sprintf("adding validator missing blocks to database: %s", v.Address))
		totalMissed, err := strconv.ParseInt(v.MissedBlocksCounter, 10, 64)
		if err != nil {
			return err
		}

		params := database.InsertMissedBlockInfoParams{
			Height:      height,
			ValidatorID: val.ID,
			Missed:      totalMissed,
		}
		if err := c.db.InsertMissedBlockInfo(c.ctx, params); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cosmos) UpdateMissingTable(height int64) error {
	missingBlocks, err := c.getMissedBlocksByHeight(strconv.Itoa(int(height)))
	if err != nil {
		return err
	}
	return c.addMissingBlocksInfoToDB(height, missingBlocks)
}
