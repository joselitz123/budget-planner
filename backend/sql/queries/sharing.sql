-- name: CreateShareInvitation :one
INSERT INTO share_invitations (budget_id, owner_id, recipient_email, permission, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetInvitationByID :one
SELECT * FROM share_invitations
WHERE id = $1
LIMIT 1;

-- name: GetPendingInvitationsByRecipient :many
SELECT si.*, u.name as owner_name, u.email as owner_email, b.name as budget_name, b.month as budget_month
FROM share_invitations si
JOIN users u ON si.owner_id = u.id
JOIN budgets b ON si.budget_id = b.id
WHERE si.recipient_email = $1 AND si.status = 'pending' AND si.expires_at > NOW()
ORDER BY si.created_at DESC;

-- name: GetInvitationsByOwner :many
SELECT si.*, u.name as recipient_name
FROM share_invitations si
LEFT JOIN users u ON si.recipient_email = u.email
WHERE si.owner_id = $1
ORDER BY si.created_at DESC;

-- name: UpdateInvitationStatus :one
UPDATE share_invitations
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: DeleteInvitation :exec
DELETE FROM share_invitations
WHERE id = $1;

-- name: GetShareAccessByBudget :many
SELECT sa.*, u.name as shared_with_name, u.email as shared_with_email
FROM share_access sa
JOIN users u ON sa.shared_with_id = u.id
WHERE sa.budget_id = $1;

-- name: GetShareAccessForUser :many
SELECT sa.*, o.name as owner_name, o.email as owner_email, b.name as budget_name, b.month as budget_month
FROM share_access sa
JOIN users o ON sa.owner_id = o.id
JOIN budgets b ON sa.budget_id = b.id
WHERE sa.shared_with_id = $1
ORDER BY b.month DESC;

-- name: CreateShareAccess :one
INSERT INTO share_access (budget_id, owner_id, shared_with_id, permission)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetShareAccessByID :one
SELECT * FROM share_access
WHERE id = $1
LIMIT 1;

-- name: UpdateShareAccess :one
UPDATE share_access
SET permission = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: DeleteShareAccess :exec
DELETE FROM share_access
WHERE id = $1;

-- name: GetShareAccessForBudgetAndUser :one
SELECT sa.*
FROM share_access sa
WHERE sa.budget_id = $1 AND sa.shared_with_id = $2
LIMIT 1;

-- name: CheckBudgetAccess :one
SELECT 'owner' as permission, true as is_owner
FROM budgets b
WHERE b.id = $1 AND b.user_id = $2 AND b.deleted = false
UNION ALL
SELECT sa.permission, false as is_owner
FROM share_access sa
WHERE sa.budget_id = $1 AND sa.shared_with_id = $2
LIMIT 1;
