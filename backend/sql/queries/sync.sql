-- name: CreateSyncOperation :one
INSERT INTO sync_operations (user_id, table_name, record_id, operation, local_data, server_data, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSyncOperationsByUser :many
SELECT * FROM sync_operations
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPendingSyncOperations :many
SELECT * FROM sync_operations
WHERE user_id = $1 AND status = 'pending'
ORDER BY created_at ASC;

-- name: GetFailedSyncOperations :many
SELECT * FROM sync_operations
WHERE user_id = $1 AND status = 'failed'
ORDER BY created_at ASC;

-- name: GetSyncOperationByID :one
SELECT * FROM sync_operations
WHERE id = $1
LIMIT 1;

-- name: UpdateSyncOperationStatus :one
UPDATE sync_operations
SET
    status = $1,
    error_message = COALESCE(sqlc.narg('error_message'), error_message),
    attempt_count = attempt_count + 1,
    last_attempt_at = NOW(),
    updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: DeleteSyncOperation :exec
DELETE FROM sync_operations
WHERE id = $1;

-- name: DeleteSyncedOperations :exec
DELETE FROM sync_operations
WHERE user_id = $1 AND status = 'synced' AND created_at < NOW() - INTERVAL '30 days';

-- Sync pull queries - fetch records updated since last sync

-- name: GetBudgetsSince :many
SELECT * FROM budgets
WHERE user_id = $1
  AND deleted = false
  AND ($2 IS NULL OR updated_at > $2)
ORDER BY updated_at ASC;

-- name: GetTransactionsSince :many
SELECT * FROM transactions
WHERE user_id = $1
  AND deleted = false
  AND ($2 IS NULL OR updated_at > $2)
ORDER BY updated_at ASC;

-- name: GetCategoriesSince :many
SELECT * FROM categories
WHERE user_id = $1
  AND deleted = false
  AND ($2 IS NULL OR updated_at > $2)
ORDER BY updated_at ASC;

-- name: CountPendingSyncOperations :one
SELECT COUNT(*) as count
FROM sync_operations
WHERE user_id = $1 AND status = 'pending';

-- name: ResolveSyncOperation :exec
UPDATE sync_operations
SET
    status = $2,
    updated_at = NOW()
WHERE id = $1;
