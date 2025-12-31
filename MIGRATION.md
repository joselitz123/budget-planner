# Migration to Beans Task Management System

**Date:** 2025-12-28
**Migrated from:** todo.md-based tracking
**Migrated to:** Beans task management system

## Overview

This document describes the migration from session-based todo.md files to the Beans task management system for the Budget Planner project. Beans provides better queryability, agent integration, and project memory preservation.

## Beans Created

### Backend (1 bean)

- **BP-0g0d**: Backend Complete - All 10 Handlers with 100% Test Coverage (completed)
  - Type: feature
  - Priority: critical
  - Status: completed
  - Tags: backend, testing, completed
  - Effort: 40 hours
  - All 10 handlers implemented with 48/48 tests passing

### Frontend (8 beans)

#### Priority 1 (Critical) - 2 beans

- **BP-r94p**: Backend API Integration (todo)
  - Type: feature
  - Priority: critical
  - Tags: frontend, api, sync, backend
  - Effort: 3 hours
  - Wire up actual API calls to Go backend for all CRUD operations

- **BP-mmy2**: Mark Bill Paid Functionality (todo)
  - Type: feature
  - Priority: critical
  - Tags: frontend, ui
  - Effort: 30 minutes
  - Implement mark bill paid button in bills page

#### Priority 2 (High) - 2 beans

- **BP-mfu9**: Month Navigation (todo)
  - Type: feature
  - Priority: high
  - Tags: frontend, ui
  - Effort: 30 minutes
  - Implement prev/next month buttons

- **BP-7v89**: Settings Page (todo)
  - Type: feature
  - Priority: high
  - Tags: frontend, ui
  - Effort: 1 hour
  - Create settings page with theme toggle, currency options, data export/import

#### Priority 3 (Medium) - 2 beans

- **BP-8fwz**: Loading States (todo)
  - Type: task (enhancement)
  - Priority: normal
  - Tags: frontend, ui
  - Effort: 1 hour
  - Add spinners/skeletons during data loading

- **BP-0mr6**: Error Handling (todo)
  - Type: task (enhancement)
  - Priority: normal
  - Tags: frontend
  - Effort: 1 hour
  - Toast notifications for errors. Error boundaries.

#### Priority 4 (Low) - 2 beans

- **BP-7k0g**: Analytics/Insights Page (todo)
  - Type: feature
  - Priority: low
  - Tags: frontend, analytics
  - Effort: 2 hours
  - Spending trends, category comparison charts, monthly summaries

- **BP-kocv**: PWA Icons Generation (todo)
  - Type: task (chore)
  - Priority: low
  - Tags: frontend, pwa
  - Effort: 30 minutes
  - Create or download app logo. Generate icon assets

### Technical Debts (3 beans)

- **BP-577b**: Shadcn-Svelte CLI Unavailable in Devcontainer (todo)
  - Type: bug (tech-debt)
  - Priority: normal
  - Tags: frontend, tech-debt
  - Effort: 2 hours
  - CLI requires TTY unavailable in devcontainer. All UI components created manually.

- **BP-nbok**: bits-ui Dialog Type Issues - CustomModal Created Instead (completed)
  - Type: bug (tech-debt)
  - Priority: low
  - Tags: frontend, tech-debt
  - Status: completed
  - Effort: 1 hour
  - CustomModal created instead. Works well for notebook aesthetic.

- **BP-8qmq**: TypeScript Type Definition Simplification (todo)
  - Type: bug (tech-debt)
  - Priority: low
  - Tags: frontend, tech-debt
  - Effort: 2 hours
  - Components use simplified types (less type-safe but functional).

**Total Beans Created:** 12 (10 open, 2 completed)

## Priority Mapping

| Frontend Priority | Beans Priority | Beans Count |
|-------------------|----------------|-------------|
| Priority 1        | critical       | 2           |
| Priority 2        | high           | 2           |
| Priority 3        | normal         | 2           |
| Priority 4        | low            | 2           |

## Type Mapping

| todo.md Section          | Beans Type | Bean Count |
|--------------------------|------------|------------|
| Implementation Progress  | feature    | 6          |
| Technical Debt           | bug        | 3          |
| Enhancement              | task       | 2          |
| Chore                    | task       | 1          |

## Legacy Files

Legacy todo.md files were removed on 2025-12-31 after Beans migration was complete. All task tracking is now managed in the `.beans/` directory.

## WRAP-UP Workflow Integration

The WRAP-UP workflow has been enhanced to integrate with Beans:
- Beans auto-created/updated during WRAP-UP
- Git commits include bean metadata
- Technical debts tracked automatically
- Session history preserved in bean metadata

### New WRAP-UP Process

1. Run `beans prime` to review context
2. Create/update beans for completed work
3. Check .gitignore for security files
4. Commit with bean metadata in commit message
5. Push changes

### Commit Message Format

```
{type}: {title}

- Bean: {bean_id}
- Session: {session_date}
- Files: {files_modified}

Co-authored-by: Claude Sonnet <noreply@anthropic.com>
```

## Claude Integration

### Hooks Configuration

File: `.claude/hooks.json`

```json
{
  "hooks": {
    "SessionStart": [
      {
        "type": "command",
        "command": "beans prime"
      }
    ],
    "PreCompact": [
      {
        "type": "command",
        "command": "beans prime"
      }
    ]
  }
}
```

### Permissions

Beans permissions added to `.claude/settings.local.json`:
- `Bash(beans create:*)`
- `Bash(beans list:*)`
- `Bash(beans update:*)`
- `Bash(beans show:*)`
- `Bash(beans prime:*)`
- `Bash(beans check:*)`
- `Bash(beans archive:*)`
- `Bash(beans graphql:*)`

## Beans Configuration

File: `.beans.yml`

**Project Settings:**
- Name: Budget Planner PWA
- Bean prefix: BP-
- ID length: 4
- Default status: todo
- Default type: feature

**Custom Types Defined:**
- feature, bug, enhancement, tech-debt, documentation, testing, chore
- Note: Beans CLI only supports: milestone, epic, bug, feature, task
- Mapping: enhancement→task, tech-debt→bug, documentation→task, testing→task, chore→task

**Custom Priorities:**
- critical (Level 1) - Frontend Priority 1
- high (Level 2) - Frontend Priority 2
- normal (Level 3) - Frontend Priority 3
- low (Level 4) - Frontend Priority 4

**Custom Labels:**
- Frontend/Backend: frontend, backend
- Feature areas: auth, database, api, testing, ui, pwa, sync, indexeddb
- Status indicators: blocked, in-progress, review-needed

**Effort Tracking:**
- Enabled: true
- Units: hours, days
- Default: hours

**Archive Configuration:**
- Auto-archive: true
- Archive after: 30 days
- Keep metadata: true

## Documentation Updates

### CLAUDE.md (Root)

- WRAP-UP workflow section updated with Beans integration
- Added beans prime step
- Updated commit message format to include bean metadata
- Added Beans commands reference

### backend/CLAUDE.md

- Added "Task Tracking" section referencing Beans
- Updated "Important Files" to reference `.beans/` and `.beans.yml`
- Noted that `todo.md.legacy` is preserved for reference

### frontend/CLAUDE.md

- Added "Task Tracking" section referencing Beans
- Updated references from `todo.md` to `.beans/` directory
- Updated "Important Files" to prioritize Beans over todo.md.legacy

## Useful Beans Commands

```bash
# Show all open beans
beans list

# Show beans by status
beans list --status todo
beans list --status completed

# Show beans by priority
beans list --priority critical
beans list --priority high

# Show beans by tag
beans list --tag frontend
beans list --tag backend
beans list --tag tech-debt

# Show bean details
beans show BP-r94p

# Update bean status
beans update BP-r94p --status in-progress
beans update BP-r94p --status completed

# Archive completed beans
beans archive

# Validate configuration
beans check

# Run GraphQL query
beans graphql '{ beans { id title status } }'
```

## Migration Benefits

1. **Better Queryability:** Filter beans by status, priority, tags, type
2. **Agent Integration:** Automatic bean creation via hooks
3. **Project Memory:** Archived beans preserve history
4. **Git Integration:** Commit messages include bean metadata
5. **Structured Data:** Consistent metadata vs. free-form markdown
6. **Dependency Tracking:** Link related beans
7. **Technical Debt Tracking:** Dedicated bug type for tech-debts

## Next Steps

1. **Monitor Bean Usage**
   - Track bean creation rate during development
   - Monitor technical debt accumulation
   - Review priority distribution weekly

2. **Refine Configuration**
   - Add custom labels if needed
   - Adjust priority definitions
   - Tune auto-archive settings

3. **Optimize Workflow**
   - Adjust WRAP-UP based on usage
   - Fine-tune Claude hooks
   - Add custom GraphQL queries

4. **Generate Reports** (Optional)
   - Sprint velocity from completed beans
   - Technical debt dashboard
   - Priority distribution charts

## References

- Beans documentation: https://github.com/hmans/beans
- Project configuration: `.beans.yml`
- Beans data: `.beans/data/*.md`
- Root CLAUDE.md: WRAP-UP workflow section

## Rollback Plan

Legacy todo.md files were removed on 2025-12-31. To restore from git history:

```bash
# 1. Restore todo.md files from git history
git checkout <commit-before-cleanup>~1 -- backend/todo.md.legacy frontend/todo.md.legacy

# 2. Remove Beans directory
rm -rf .beans/
rm .beans.yml

# 3. Revert documentation changes
git checkout CLAUDE.md
git checkout backend/CLAUDE.md
git checkout frontend/CLAUDE.md

# 4. Remove Claude hooks
rm .claude/hooks.json

# 5. Verify restoration
git status
```

## Validation Checklist

- [x] All 12 beans created successfully
- [x] Priority mapping correct (critical/high/normal/low)
- [x] Tags applied correctly ([frontend|backend], feature areas)
- [x] Effort estimates preserved
- [x] WRAP-UP workflow updated in CLAUDE.md
- [x] Claude hooks configured
- [x] Beans permissions added to settings.local.json
- [x] Legacy todo.md files renamed to .legacy
- [x] Migration documentation created (MIGRATION.md)
- [x] Legacy todo.md files removed (2025-12-31)
- [ ] Git commit successful with all changes

---

**Migration Status:** ✅ Complete
**Beans System:** Active
**Legacy Files:** Preserved as reference
