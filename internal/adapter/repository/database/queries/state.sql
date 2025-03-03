-- name: GetByKey :one
select key, value, created_at, updated_at from state where key = $1;

-- name: CreateState :exec
INSERT INTO state (key, value) VALUES ($1, $2);

-- name: UpdateState :exec
UPDATE state SET value = $2, updated_at = current_timestamp WHERE key = $1;