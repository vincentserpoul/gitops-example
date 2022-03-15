-- name: SaveHappycatFact :exec
INSERT INTO happycat_facts (id, fact)
VALUES ($1, $2);

-- name: ListHappycatFacts :many
SELECT id, fact, created_at
FROM happycat_facts
ORDER BY id;

-- name: GetHappycatFact :one
SELECT id, fact, created_at
FROM happycat_facts
WHERE id = $1
ORDER BY id;