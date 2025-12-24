-- name: ListUserReflections :many
SELECT * FROM reflections
WHERE user_id = $1 AND deleted = false
ORDER BY created_at DESC;

-- name: GetReflectionByID :one
SELECT * FROM reflections
WHERE id = $1 AND deleted = false
LIMIT 1;

-- name: GetReflectionByBudget :one
SELECT * FROM reflections
WHERE budget_id = $1 AND deleted = false
LIMIT 1;

-- name: CreateReflection :one
INSERT INTO reflections (user_id, budget_id, overall_rating, is_private)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateReflection :one
UPDATE reflections
SET
    overall_rating = COALESCE(sqlc.narg('overall_rating'), overall_rating),
    is_private = COALESCE(sqlc.narg('is_private'), is_private),
    updated_at = NOW()
WHERE id = $1 AND deleted = false
RETURNING *;

-- name: DeleteReflection :exec
UPDATE reflections
SET deleted = true, updated_at = NOW()
WHERE id = $1;

-- name: GetReflectionQuestions :many
SELECT * FROM reflection_questions
WHERE reflection_id = $1
ORDER BY sequence ASC;

-- name: CreateReflectionQuestion :one
INSERT INTO reflection_questions (reflection_id, sequence, question_id, question_text, answer, mood)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateReflectionQuestion :one
UPDATE reflection_questions
SET
    question_text = COALESCE(sqlc.narg('question_text'), question_text),
    answer = COALESCE(sqlc.narg('answer'), answer),
    mood = COALESCE(sqlc.narg('mood'), mood),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListReflectionTemplates :many
SELECT * FROM reflection_templates
WHERE is_active = true
ORDER BY name ASC;

-- name: GetTemplateByID :one
SELECT * FROM reflection_templates
WHERE id = $1
LIMIT 1;

-- name: GetTemplateQuestions :many
SELECT * FROM template_questions
WHERE template_id = $1
ORDER BY sort_order ASC;

-- name: CreateReflectionTemplate :one
INSERT INTO reflection_templates (name, is_active, version)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateReflectionTemplate :one
UPDATE reflection_templates
SET
    name = COALESCE(sqlc.narg('name'), name),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    version = COALESCE(sqlc.narg('version'), version),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteReflectionTemplate :exec
DELETE FROM reflection_templates
WHERE id = $1;
