package cosmos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/evmos/validator-status/pkg/database"
	"github.com/evmos/validator-status/pkg/logger"
)

type Validator struct {
	OperatorAddress string `json:"operator_address"`
	Status          string `json:"status"`
	Jailed          bool   `json:"jailed"`
	ConsensusPubKey struct {
		TypeURL string `json:"@type"`
		Value   string `json:"key"`
	} `json:"consensus_pubkey"`
	Description struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
	} `json:"description"`
}

type ValidatorsResponse struct {
	Validators []Validator
}

func (c *Cosmos) getValidators(height string) (ValidatorsResponse, error) {
	reqGo, err := http.NewRequest("GET", c.apiURL+"/cosmos/staking/v1beta1/validators", nil)
	if err != nil {
		return ValidatorsResponse{}, fmt.Errorf("error creating the request %s", err.Error())
	}

	reqGo.Header.Set("x-cosmos-block-height", height)

	resp, err := http.DefaultClient.Do(reqGo)
	if err != nil {
		return ValidatorsResponse{}, fmt.Errorf("error making the request %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return ValidatorsResponse{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return ValidatorsResponse{}, err
	}

	var v ValidatorsResponse
	err = json.Unmarshal(content, &v)
	if err != nil {
		return ValidatorsResponse{}, err
	}

	return v, nil
}

func (c *Cosmos) addValidatorsToDatabase(vals []Validator) error {
	for _, v := range vals {
		_, err := c.db.GetValidatorByOperatorAddress(c.ctx, v.OperatorAddress)

		// Validator is not in the database
		if err == sql.ErrNoRows {
			logger.LogInfo(fmt.Sprintf("adding validator to database: %s", v.Description.Moniker))
			validatorAddress, err := pubkeyToValCons(v.ConsensusPubKey.Value)
			if err != nil {
				return err
			}
			params := database.InsertValidatorParams{
				OperatorAddress:  v.OperatorAddress,
				Pubkey:           v.ConsensusPubKey.Value,
				ValidatorAddress: validatorAddress,
				Moniker:          v.Description.Moniker,
				Indentity:        v.Description.Identity,
			}
			if _, err := c.db.InsertValidator(c.ctx, params); err != nil {
				return err
			}
			continue
		}

		// Database error
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cosmos) UpdateValidatorsTable(height int64) error {
	validators, err := c.getValidators(strconv.Itoa(int(height)))
	if err != nil {
		return err
	}
	return c.addValidatorsToDatabase(validators.Validators)
}
