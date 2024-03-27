package cosmos

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ChainHeight struct {
	CurrentHeight  int
	EarliestHeight int
}

type StatusResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		NodeInfo struct {
			ProtocolVersion struct {
				P2P   string `json:"p2p"`
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"protocol_version"`
			ID         string `json:"id"`
			ListenAddr string `json:"listen_addr"`
			Network    string `json:"network"`
			Version    string `json:"version"`
			Channels   string `json:"channels"`
			Moniker    string `json:"moniker"`
			Other      struct {
				TxIndex    string `json:"tx_index"`
				RPCAddress string `json:"rpc_address"`
			} `json:"other"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHash     string    `json:"latest_block_hash"`
			LatestAppHash       string    `json:"latest_app_hash"`
			LatestBlockHeight   string    `json:"latest_block_height"`
			LatestBlockTime     time.Time `json:"latest_block_time"`
			EarliestBlockHash   string    `json:"earliest_block_hash"`
			EarliestAppHash     string    `json:"earliest_app_hash"`
			EarliestBlockHeight string    `json:"earliest_block_height"`
			EarliestBlockTime   time.Time `json:"earliest_block_time"`
			CatchingUp          bool      `json:"catching_up"`
		} `json:"sync_info"`
		ValidatorInfo struct {
			Address string `json:"address"`
			PubKey  struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"pub_key"`
			VotingPower string `json:"voting_power"`
		} `json:"validator_info"`
	} `json:"result"`
}

func (c *Cosmos) GetChainHeight() (ChainHeight, error) {
	// NOTE: assumes that we do not need pagination
	resp, err := http.Get(c.rpcURL + "/status")
	if err != nil {
		return ChainHeight{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ChainHeight{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChainHeight{}, err
	}

	var v StatusResponse
	err = json.Unmarshal(content, &v)
	if err != nil {
		return ChainHeight{}, err
	}

	currentHeight, err := strconv.ParseInt(v.Result.SyncInfo.LatestBlockHeight, 10, 64)
	if err != nil {
		return ChainHeight{}, err
	}

	earliestHeight, err := strconv.ParseInt(v.Result.SyncInfo.EarliestBlockHeight, 10, 64)
	if err != nil {
		return ChainHeight{}, err
	}

	return ChainHeight{
		CurrentHeight:  int(currentHeight),
		EarliestHeight: int(earliestHeight),
	}, nil
}
