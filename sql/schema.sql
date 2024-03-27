CREATE TABLE IF NOT EXISTS missed_blocks(
    id INTEGER NOT NULL PRIMARY KEY,
    height INTEGER NOT NULL,
    validator_id INTEGER NOT NULL,
    missed INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS heigthindex on missed_blocks (height);
CREATE INDEX IF NOT EXISTS validatorindex on missed_blocks (validator_id);

CREATE TABLE IF NOT EXISTS validators(
    id INTEGER NOT NULL PRIMARY KEY,
    operator_address TEXT NOT NULL,
    pubkey TEXT NOT NULL,
    validator_address TEXT NOT NULL,
    moniker TEXT NOT NULL,
    indentity TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS operator_address_index on validators (operator_address);
CREATE INDEX IF NOT EXISTS validator_address_index on validators (operator_address);
CREATE INDEX IF NOT EXISTS moniker_index on validators (moniker);

