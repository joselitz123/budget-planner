---
# BP-6tuh
title: Reflection Templates System
status: open
type: feature
priority: low
created_at: 2025-12-28T16:51:49Z
updated_at: 2025-12-28T16:51:49Z
---

Implement templated monthly reflections with customizable question sets.

## Background
Current reflection only has basic 3 prompts. Spec includes full template system with multiple question types and mood tracking.

## Acceptance Criteria
- [ ] Reflection template CRUD operations
- [ ] Multiple question types (text, rating, yes/no)
- [ ] Mood tracking per question
- [ ] Template selection when creating reflection
- [ ] Question sequence/ordering
- [ ] Pre-built templates (Monthly Review, Weekly Check-in)

## Technical Details

### Database Tables (backend exists)
- reflection_templates (name, version, is_active)
- template_questions (template_id, question_id, question_text, type, sort_order)
- reflection_questions (reflection_id, question_id, answer, mood, sequence)

### UI Components
- TemplateSelector.svelte (choose template)
- ReflectionQuestion.svelte (single question with mood)
- TemplateEditor.svelte (admin: create/edit templates)

### Question Types
- text: Long-form text answer
- rating: 1-10 scale
- yesno: Yes/No question
- mood: Emoji mood selection

### Pre-built Templates
1. "Monthly Budget Review" (5 questions)
2. "Weekly Financial Check-in" (3 questions)
3. "End of Year Reflection" (10 questions)

### Files to Create
- frontend/src/lib/components/reflection/TemplateSelector.svelte
- frontend/src/lib/components/reflection/ReflectionQuestion.svelte
- frontend/src/lib/stores/reflectionTemplates.ts
- frontend/src/lib/api/reflectionTemplates.ts

### Files to Modify
- frontend/src/routes/+page.svelte (integrate template selector)

## Effort Estimate
2 hours

## Dependencies
- BP-r94p (Backend API Integration) - completed

## Session Date
2025-12-28