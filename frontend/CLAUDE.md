# Budget Planner Frontend - Development Guide

This file provides comprehensive guidance for Claude Code when working with the SvelteKit frontend of the Budget Planner project.

## Frontend Overview

**Tech Stack:**
- SvelteKit 2.0 - Full-stack web framework
- Svelte 5.0 - Reactive UI framework
- TypeScript 5.0 - Type safety
- Vite 5.0 - Build tool and dev server
- Shadcn-Svelte - Component library (Radix UI-based)
- Tailwind CSS 3.4 - Utility-first CSS framework
- IndexedDB (via `idb` library) - Offline storage
- Clerk - Authentication

**Current Status:** ~75% complete

**Task Tracking:** All tasks tracked in Beans system (`.beans/` directory)
- Use `beans list --tag frontend` to view frontend tasks
- Use `beans list --priority critical` to view high-priority tasks

---

## Quick Commands

```bash
cd /workspace/budget-planner/frontend

# Install dependencies
npm install

# Run development server (with host access for devcontainer)
npm run dev

# Type checking
npm run check
npm run check:watch

# Build for production
npm run build

# Preview production build
npm run preview

# Linting and formatting
npm run lint
npm run format

# Generate PWA icons
npm run pwa:generate-icons
```

---

## Tech Stack

**Core Framework:**
- `svelte@^5.0.0` - Reactive UI framework with runes
- `@sveltejs/kit@^2.0.0` - Full-stack web framework
- `vite@^5.0.0` - Build tool and dev server
- `typescript@^5.0.0` - Type safety

**UI & Styling:**
- `@tailwindcss/forms@^0.5.9` - Form styling plugins
- `tailwindcss@^3.4.10` - Utility-first CSS
- `class-variance-authority@^0.7.0` - Component variants
- `lucide-svelte@^0.469.0` - Icon library

**Offline Storage:**
- `idb@^8.0.0` - IndexedDB wrapper
- `workbox-window@^7.1.0` - Service worker management

**Authentication:**
- `@clerk/sveltekit@^2.0.0` - Clerk authentication

**Development Tools:**
- `@sveltejs/vite-plugin-svelte@^3.0.0` - Svelte plugin for Vite
- `svelte-check@^3.0.0` - Type checking
- `vite-plugin-pwa@^0.20.0` - PWA support
- `prettier@^3.0.0` - Code formatting

**See `package.json` for complete dependency list.**

---

## Implementation Status

**Completed (Iteration 1):**
- ✅ SvelteKit 2.0 project scaffolding
- ✅ Vite 5.0 configuration with PWA plugin
- ✅ TypeScript 5.0 strict mode
- ✅ Tailwind CSS 3.4 with custom notebook theme
- ✅ PWA manifest and service worker configuration
- ✅ IndexedDB client with full CRUD operations
- ✅ Svelte stores for state management
- ✅ API client with error handling
- ✅ Root layout with navigation and theme toggle
- ✅ Budget Overview page with monthly reflection
- ✅ Expense Tracker page with transaction list
- ✅ Bill Payment page with status tracking
- ✅ Transaction Modal with form validation
- ✅ Budget Sharing feature with invitations
- ✅ Type checking: 0 errors, 0 warnings
- ✅ Build successful

**To Implement:**
- ⚠️ Backend API integration
- ⚠️ Month navigation
- ⚠️ Mark bill paid functionality
- ⚠️ Settings page
- ⚠️ Full offline sync testing

**See Beans (`.beans/` directory) for detailed remaining tasks and next steps.**

---

## Architecture Patterns

### SvelteKit File-Based Routing

SvelteKit uses file-based routing. The filesystem determines the routes:

```
routes/                    → /
routes/about              → /about
routes/blog/             → /blog
routes/blog/[slug]       → /blog/:slug (dynamic)
routes/[lang]/about      → /:lang/about (dynamic layout)
```

**Special Files:**
- `+page.svelte` - Page component
- `+layout.svelte` - Layout wrapper for child routes
- `+load.ts` - Server-side data loading
- `+page.server.ts` - Server actions
- `+error.svelte` - Error boundary

### Offline-First Architecture with IndexedDB

The frontend uses IndexedDB for offline data persistence:

**Pattern:**
1. **IndexedDB** stores all data locally for offline access
2. **Sync queue** captures changes when offline
3. **Sync API** (`/api/sync/push`, `/api/sync/pull`) handles bidirectional sync
4. **Conflict resolution** uses timestamp-based and owner-priority strategies

**Database Structure:**
```typescript
// Expected IndexedDB structure (using idb library)
interface BudgetDB {
  budgets: {
    key: string;
    value: Budget;
    indexes: { 'by-month': string };
  };
  categories: { key: string; value: Category };
  transactions: { key: string; value: Transaction };
  reflections: { key: string; value: Reflection };
  syncQueue: { key: string; value: SyncOperation };
}
```

**Data Flow:**
1. **Load data** from IndexedDB on app init
2. **Sync with backend** via API on change
3. **Update IndexedDB** when data changes
4. **React to changes** using Svelte stores

### PWA Configuration

**Service Worker** (`vite-plugin-pwa`):
- Registers on app load
- Caches static assets for offline access
- Provides offline fallback
- Syncs data when connection restored

**Web App Manifest** (`static/manifest.json`):
- App name and short name
- Icons (generated via `npm run pwa:generate-icons`)
- Theme and background colors
- Display mode (standalone)

### State Management with Svelte Stores

**Writable Stores** (for client state):
```typescript
// src/lib/stores/budgets.ts
import { writable } from 'svelte/store';

export const budgets = writable<Budget[]>([]);
export const currentBudget = writable<Budget | null>(null);
```

**Derived Stores** (computed values):
```typescript
import { derived } from 'svelte/store';

const totalSpent = derived(budgets, ($budgets) =>
  $budgets.reduce((sum, b) => sum + b.spent, 0)
);
```

**Using Stores in Components:**
```svelte
<script>
  import { budgets } from '$lib/stores/budgets';
</script>

{#each $budgets as budget}
  <BudgetCard {budget} />
{/each}
```

### Clerk Authentication Integration

**Pattern:**
```typescript
// Using @clerk/sveltekit
import { clerkClient } from '@clerk/sveltekit/sdk';

// In load functions
export const load = async (event) => {
  const { userId } = await event.locals.auth();
  if (!userId) throw redirect(302, '/sign-in');
  return { userId };
};
```

**Environment Variables:**
```env
PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_...
PUBLIC_CLERK_SIGN_IN_URL=/sign-in
PUBLIC_CLERK_SIGN_UP_URL=/sign-up
PUBLIC_CLERK_AFTER_SIGN_IN_URL=/dashboard
PUBLIC_CLERK_AFTER_SIGN_UP_URL=/dashboard
```

### Budget Sharing

**Pattern:**
The budget sharing feature allows users to invite collaborators to their budgets with different permission levels.

**API Client** (`src/lib/api/shares.ts`):
```typescript
import { sharesApi } from '$lib/api/shares';

// Create invitation
await sharesApi.createInvitation({
  budgetId: 'budget-uuid',
  recipientEmail: 'friend@example.com',
  permission: 'view' // or 'edit'
});

// Get pending invitations
const invitations = await sharesApi.getMyInvitations();

// Respond to invitation
await sharesApi.respondToInvitation(invitationId, { status: 'accepted' });

// Get shared budgets
const shared = await sharesApi.getSharedBudgets();

// Get who has access to a budget
const accessList = await sharesApi.getBudgetSharing(budgetId);

// Remove access
await sharesApi.removeAccess(accessId);
```

**State Management** (`src/lib/stores/shares.ts`):
```typescript
import {
  invitations,
  pendingInvitations,
  sharedBudgets,
  loadInvitations,
  loadSharedBudgets,
  createInvitation,
  respondToInvitation
} from '$lib/stores/shares';
```

**Components:**
- `ShareBudgetDialog.svelte` - Modal for inviting users
- `InvitationList.svelte` - List of pending invitations
- `SharedBudgetCard.svelte` - Card displaying shared budget

**Permission Levels:**
- `view` - Read-only access
- `edit` - Can modify transactions and categories

**Manual Notification Flow:**
- Invitations are stored in database but emails are NOT sent
- Owner must manually notify recipient to check the app
- Recipient visits `/shared` route to see pending invitations

---

## Component Organization

### Shadcn-Svelte Components

**Usage Pattern:**
1. Components are in `src/lib/components/ui/`
2. Manually created (CLI unavailable in devcontainer)
3. Import and use in Svelte files

**Example:**
```svelte
<script>
  import { Button } from '$lib/components/ui/button';
  import { Card } from '$lib/components/ui/card';
</script>

<Card>
  <Button>Click me</Button>
</Card>
```

**Common Components:**
- `button` - Button with variants (default, outline, ghost, destructive)
- `card` - Card container with header, content, footer
- `input` - Form input with labels
- `select` - Dropdown selection
- `dialog` - Modal/dialog component (CustomModal for notebook aesthetic)
- `badge` - Status badges

**Note:** bits-ui Dialog had type definition issues, so CustomModal was created for the notebook aesthetic.

### Layout Components

**Root Layout** (`src/routes/+layout.svelte`):
- Navigation bar
- Theme provider
- Offline status indicator
- Mobile bottom navigation

**Feature Layouts:**
- Dashboard layout with sidebar
- Auth layout (minimal, centered)
- Settings layout with tabs

---

## Styling

### Tailwind CSS Configuration

**Custom Notebook Theme:**
```javascript
// tailwind.config.js
theme: {
  extend: {
    colors: {
      primary: '#333333',
      'paper-light': '#fdfbf7',
      'line-light': '#e5e7eb',
    },
    fontFamily: {
      'display': ['Playfair Display', 'serif'],
      'body': ['Inter', 'sans-serif'],
      'handwriting': ['Caveat', 'cursive'],
    }
  }
}
```

**CSS Classes:**
- `.notebook-lines` - Ruled paper background pattern
- `.binding-holes` - Spiral binding holes
- `.binding-coil` - Gold gradient coil
- `.custom-scrollbar` - Styled scrollbars

### Component Styling Patterns

**Class Variants** (using `class-variance-authority`):
```typescript
import { cva } from 'class-variance-authority';

const buttonVariants = cva(
  'inline-flex items-center justify-center rounded-md text-sm',
  {
    variants: {
      variant: {
        default: 'bg-primary text-primary-foreground',
        destructive: 'bg-destructive text-destructive-foreground',
        outline: 'border border-input bg-background',
      }
    }
  }
);
```

**Tailwind Merge** (combine class names):
```typescript
import { cn } from '$lib/utils';

<div class={cn('base-class', condition && 'conditional-class')} />
```

---

## API Integration

### Fetch API with TypeScript

**Typed API Client** (`src/lib/api/client.ts`):
```typescript
interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: { message: string; code: string };
}

export async function apiGet<T>(path: string): Promise<T> {
  const token = await clerkClient.session?.getToken();
  const response = await fetch(`${import.meta.env.VITE_PUBLIC_API_URL}${path}`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.json();
}
```

**Usage:**
```typescript
import { apiGet } from '$lib/api/client';

const budgets = await apiGet<Budget[]>('/budgets');
```

### Loading States

**SvelteKit Load Function:**
```typescript
// src/routes/budgets/+page.ts
export const load = async ({ fetch }) => {
  const response = await fetch('/api/budgets');
  const budgets = await response.json();
  return { budgets };
};
```

**In Component:**
```svelte
<script>
  export let data;
  $: ({ budgets } = data);
</script>

{#await budgets}
  <LoadingSpinner />
{:then items}
  <BudgetList {items} />
{:catch error}
  <ErrorDisplay {error} />
{/await}
```

---

## Common Patterns

### Route Protection

**Server Load Function:**
```typescript
// +page.server.ts
import { redirect } from '@sveltejs/kit';
import { clerkClient } from '@clerk/sveltekit/sdk';

export const load = async (event) => {
  const { userId } = await event.locals.auth();
  if (!userId) throw redirect(302, '/sign-in');
  return { userId };
};
```

### Form Validation

**Svelte Actions:**
```svelte
<script>
  import { enhance } from '$app/forms';

  let { form } = $props();
  let { data } = $form;
</script>

<form method="POST" use:enhance>
  <input name="name" value={data?.name || ''} />
  {#if form?.errors}
    <div class="error">{form.errors.name}</div>
  {/if}
  <button type="submit">Submit</button>
</form>
```

### Date Formatting

**Using Intl API:**
```typescript
const formatDate = (date: Date) => {
  return new Intl.DateTimeFormat('en-US', {
    month: 'long',
    year: 'numeric'
  }).format(date);
};
```

---

## Environment Configuration

### Environment Variables

**File:** `.env` (create from `.env.example` if provided)

```env
# Clerk Authentication
PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_...
PUBLIC_CLERK_SIGN_IN_URL=/sign-in
PUBLIC_CLERK_SIGN_UP_URL=/sign-up
PUBLIC_CLERK_AFTER_SIGN_IN_URL=/dashboard
PUBLIC_CLERK_AFTER_SIGN_UP_URL=/dashboard

# API Configuration
PUBLIC_API_URL=http://localhost:8080/api

# PWA Configuration
PUBLIC_APP_NAME="Budget Planner"
PUBLIC_APP_SHORT_NAME="Budget"
PUBLIC_APP_DESCRIPTION="Offline-first budget planning application"
PUBLIC_APP_THEME_COLOR="#3b82f6"
PUBLIC_APP_BACKGROUND_COLOR="#ffffff"

# Offline Sync Settings
PUBLIC_SYNC_INTERVAL=30000      # 30 seconds
PUBLIC_OFFLINE_RETRY_DELAY=5000 # 5 seconds
PUBLIC_MAX_OFFLINE_OPERATIONS=100
```

**Access in SvelteKit:**
```svelte
<script>
  const apiUrl = import.meta.env.VITE_PUBLIC_API_URL;
</script>
```

---

## Important Files

- **`.beans/`** - **⭐ Beans task management system - current task tracking - START HERE**
- `.beans.yml` - Beans configuration file (root directory)
- `starting-point.md` (root) - Full project specification
- `CLAUDE.md` (root) - Project-level documentation
- `backend/CLAUDE.md` - Backend development guide

---

## Project Structure

```
frontend/
├── src/
│   ├── lib/
│   │   ├── components/        # UI components (Button, Input, CustomModal, etc.)
│   │   │   ├── ui/             # Shadcn-Svelte base components
│   │   │   └── sharing/        # Budget sharing components
│   │   ├── stores/            # Svelte stores for state management
│   │   ├── db/                # IndexedDB client setup
│   │   ├── api/               # API client functions
│   │   └── utils/             # Helper functions
│   ├── routes/                # SvelteKit file-based routing
│   │   ├── +layout.svelte      # Root layout with navigation
│   │   ├── +page.svelte        # Budget Overview page
│   │   ├── shared/            # Shared budgets page
│   │   ├── transactions/       # Expense Tracker + AddExpenseModal
│   │   └── bills/              # Bill Payment page
│   └── app.html               # HTML template
├── static/
│   ├── icons/                 # PWA icons
│   └── manifest.json          # PWA manifest
├── svelte.config.js           # SvelteKit configuration
├── vite.config.ts             # Vite configuration
├── tailwind.config.js         # Tailwind CSS configuration
├── tsconfig.json              # TypeScript configuration
├── package.json               # Dependencies
└── .env                       # Environment variables
```

---

## Development Workflow

### Hot Reload with Vite

- Vite dev server provides instant hot module replacement (HMR)
- Changes to `.svelte` files update in browser without refresh
- Changes to CSS update immediately

### Type Checking

```bash
npm run check           # One-time type check
npm run check:watch     # Watch mode for development
```

**TypeScript Config** (`tsconfig.json`):
- Strict mode enabled
- Path aliases: `$lib` → `src/lib`
- SvelteKit type checking enabled

### Linting and Formatting

```bash
npm run lint    # Check for issues
npm run format  # Format code with Prettier
```

### Building for Production

```bash
npm run build
```

**Output:**
- `build/` directory with optimized assets
- Server-side rendering (SSR) build
- Client-side bundle
- PWA assets and service worker
