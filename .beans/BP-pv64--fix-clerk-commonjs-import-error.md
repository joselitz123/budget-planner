---
# BP-pv64
title: Fix Clerk CommonJS import error
status: completed
type: bug
priority: critical
tags:
    - frontend
    - auth
    - critical
    - clerk
    - typescript
created_at: 2025-12-29T12:28:11Z
updated_at: 2025-12-29T12:28:11Z
---

Critical bug: @clerk/clerk-js is a CommonJS module that doesn't support named exports. Code was importing { Clerk } which caused build errors. Fixed by changing to default import (import ClerkJs from '@clerk/clerk-js') in sign-in, sign-up, +layout.svelte, and clerkProvider.ts.