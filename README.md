# Validator status

## Config

The program requieres a `config.json` file to be executed, by default it will look at the path `./config.json` but it can be modified by the flag `--config`

Config Example:

```json
{
  "logfile": "",
  "dbfile": "status.db",
  "httpserver": {
    "address": "0.0.0.0",
    "port": "42069"
  },
  "cosmosapi": "localhost:1317",
  "cosmosrpc": "localhost:26656",
  "pruneoffset": 100,
  "refreshduration": "5s"
}
```

- `Pruneoffset` is the value of your node pruned settings. The program will use the information from the latest height - pruneoffset to index the chain
- `Refreshduration` is the time the program will wait until pulling new blocks values

## Database

- Schema: `sql/schema.sql`
- Queries: `sql/queries.sql`

Generate the go files to interact with the database:

```sh
make generate
```

_NOTE: you need `sqlc` to generate the files, run `make install-deps` if you need it_
