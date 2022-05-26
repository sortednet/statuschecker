-- name: GetService :one
SELECT * FROM services
WHERE name = $1 LIMIT 1;

-- name: GetServices :many
SELECT * FROM services
ORDER BY name;

-- name: RegisterService :one
INSERT INTO services (
    name, url
) VALUES (
    $1, $2
) RETURNING *;

-- name: UnregisterService :exec
DELETE FROM services
WHERE name = $1;
