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
  "refreshduration": "5s"
}
```

## Database

- Schema: `sql/schema.sql`
- Queries: `sql/queries.sql`

Generate the go files to interact with the database:

```sh
make generate
```

_NOTE: you need `sqlc` to generate the files, run `make install-deps` if you need it_
