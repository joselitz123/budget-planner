---
# BP-pxaa
title: Fix infinite redirect loop on auth routes
status: completed
type: bug
priority: critical
tags:
    - frontend
    - auth
    - critical
    - routing
created_at: 2025-12-29T12:28:11Z
updated_at: 2025-12-29T12:28:11Z
---

Critical bug: sign-in and sign-up pages were causing infinite redirect loops because root +layout.server.ts was applying auth checks to all routes including auth pages. Fixed by adding route check to skip auth for /sign-in and /sign-up paths.