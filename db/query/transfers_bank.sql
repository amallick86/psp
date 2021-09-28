-- name: CreateTransferBank :one
INSERT INTO transfers_bank (
  account_id,
  bank_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransferBank :one
SELECT * FROM transfers_bank
WHERE id = $1 LIMIT 1;

-- name: ListTransfersBank :many
SELECT * FROM transfers_bank
WHERE 
    account_id = $1 OR
    bank_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;