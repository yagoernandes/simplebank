-- name: CreateAccount :one
INSERT INTO account (
  owner, currency, balance
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM account
WHERE id = $1
LIMIT 1;

-- name: GetAccountByOwner :one
SELECT * FROM account
WHERE owner = $1
LIMIT 1;

-- name: GetAccounts :many
SELECT * FROM account
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE account
SET balance = $1
WHERE id = $2
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1;
