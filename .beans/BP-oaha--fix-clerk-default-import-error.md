---
# BP-oaha
title: Fix Clerk default import error
status: completed
type: bug
priority: critical
created_at: 2025-12-31T04:31:55Z
updated_at: 2025-12-31T04:34:31Z
---

Multiple files import Clerk as default export but @clerk/clerk-js only provides named exports. Change `import ClerkJs from` to `import { Clerk } from` in:
- frontend/src/lib/auth/clerkProvider.ts
- frontend/src/routes/(auth)/sign-in/+page.svelte
- frontend/src/routes/(auth)/sign-up/+page.svelte
- frontend/src/routes/+layout.svelte