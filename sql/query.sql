-- name: GetInfoByBlock :many
SELECT * FROM missed_blocks WHERE height = ?;

-- name: GetInfoByValidatorID :many
SELECT * FROM missed_blocks WHERE validator_id = ?;

-- name: InsertMissedBlockInfo :exec
INSERT INTO missed_blocks(
   height, validator_id, missed
) VALUES (
   ?, ?, ?
);

-- name: GetLatestHeight :one
SELECT height FROM missed_blocks ORDER BY height DESC LIMIT 1;

-- name: InsertValidator :one
INSERT INTO validators(
    operator_address, pubkey, validator_address, moniker, indentity
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING id;

-- name: GetValidatorByOperatorAddress :one
SELECT * FROM validators WHERE operator_address = ? LIMIT 1;

-- name: GetValidatorByMoniker :one
SELECT * FROM validators WHERE moniker = ? LIMIT 1;

-- name: GetValidatorByValidatorAddress :one
SELECT * FROM validators WHERE validator_address = ? LIMIT 1;
