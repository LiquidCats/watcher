-- name: GetByKey :one
select key, value, created_at, updated_at from state where key = $1;

-- name: CreateState :exec
INSERT INTO state (
    key, value, created_at, updated_at
) VALUES (
    $1, $2, $3, $4
 );

-- name: UpdateState :exec
UPDATE state SET value = $2, created_at = $3, updated_at = $4
WHERE key = $1;