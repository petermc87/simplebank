-- name: CreateAccount :one
INSERT INTO accounts (
  owner, 
  balance, 
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccount :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :exec
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;
-- We only want to update the balance. The owner and currency stay the same.
-- We return the updated data to the client.

-- name: DeleteAccount :one
DELETE FROM accounts 
WHERE id = $1
RETURNING *;


-- The comment above each is the CRUD opeation for sqlc ORM - GOLang.
-- We want to return all the table, including the id, to the client, after creation.
-- Make sure the dollar numbers match the number of columns.
-- we dont need to add the created at and id as a cloumn here because of auto generation.