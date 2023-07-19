-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
-- We are filtering by account_id. We only want entries from that account.

-- -- name: UpdateEntry :one
-- UPDATE entries
-- SET amount = $2
-- WHERE id = $1
-- RETURNING *;

-- -- name: DeleteEntry :one
-- DELETE FROM entries
-- WHERE id = $1;
