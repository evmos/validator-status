package config

type Config struct {
	LogFile    string `json:"logfile"`
	HTTPServer struct {
		Address string `json:"address"`
		Port    string `json:"port"`
	} `json:"httpserver"`
	RefreshDuration string `json:"refreshduration"`
	DatabaseFile    string `json:"dbfile"`
	CosmosAPI       string `json:"cosmosapi"`
	CosmosRPC       string `json:"cosmosrpc"`
	PruneOffset     int    `json:"pruneoffset"`
}
