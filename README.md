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
  "refreshduration": "5s",
  "apimaxlimit": 100
}
```

- `Pruneoffset` is the value of your node pruned settings. The program will use the information from the latest height - pruneoffset to index the chain
- `Refreshduration` is the time the program will wait until pulling new blocks values
- `APIMaxLimit` is the max range allowed by the api, if the range if bigger than the one set, the query will overwrite the `end` value

## Run

- Clone the repo
- Create the `config.json` file in the root of the repo
- Run `make run`

## Database

- Schema: `sql/schema.sql`
- Queries: `sql/queries.sql`

Generate the go files to interact with the database:

```sh
make generate
```

_NOTE: you need `sqlc` to generate the files, run `make install-deps` if you need it_

## Queries:

### Request all the validators info by block height

- Request

```sh
curl "http://localhost:42069/api?start=7313145"
```

_NOTE: optional parameter `end` to set a rangeof blocks_

- Response

```json
{
  "values": [
    {
      "ID": 3337,
      "Height": 7313145,
      "ValidatorID": 20,
      "Missed": 0,
      "ID_2": 20,
      "OperatorAddress": "ethmvaloper1kdtjxywfvwq94jsst2uyshwkel6dwdv5vlf4l2",
      "Pubkey": "TP8Ncm5KUUT3bQgN5SsQADVtPEwaW7O0MY5yX+ztJuE=",
      "ValidatorAddress": "ethmvalcons1qwcmnlr2kwuunpmh0s4pshrgz5zqd9m7ljs8lu",
      "Moniker": "evmOS",
      "Indentity": ""
    },
    {
      "ID": 3338,
      "Height": 7313145,
      "ValidatorID": 13,
      "Missed": 2,
      "ID_2": 13,
      "OperatorAddress": "ethmvaloper1weh6nan3p8cpg7hsfke8teksjnf5pdl8ve97ny",
      "Pubkey": "LxQKOI9n5enbjflu802ZL77lyYLHuejIu9o3FXQ3EOc=",
      "ValidatorAddress": "ethmvalcons1qkxa6u5tdr820pv9kjc5rlf8p0r7qqr25fenyw",
      "Moniker": "peersyst-node-3",
      "Indentity": "14329789E9E20C43"
    },
    ...
  ]
}
```

### Request all one validators info by block height

- Request

```sh
curl "http://localhost:42069/api?start=7313145&validator=ethmvaloper1weh6nan3p8cpg7hsfke8teksjnf5pdl8ve97ny"
```

_NOTE: optional parameter `end` to set a rangeof blocks_

- Response

```json
{
  "values": [
    {
      "ID": 3338,
      "Height": 7313145,
      "ValidatorID": 13,
      "Missed": 2,
      "ID_2": 13,
      "OperatorAddress": "ethmvaloper1weh6nan3p8cpg7hsfke8teksjnf5pdl8ve97ny",
      "Pubkey": "LxQKOI9n5enbjflu802ZL77lyYLHuejIu9o3FXQ3EOc=",
      "ValidatorAddress": "ethmvalcons1qkxa6u5tdr820pv9kjc5rlf8p0r7qqr25fenyw",
      "Moniker": "peersyst-node-3",
      "Indentity": "14329789E9E20C43"
    }
  ]
}
```
