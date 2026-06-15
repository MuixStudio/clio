# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
pnpm dev        # Start dev server with Turbopack
pnpm build      # Production build
pnpm lint       # ESLint (eslint-config-next)
pnpm format     # Prettier (ts/tsx)
pnpm typecheck  # tsc --noEmit
```

## Architecture

Next.js 16 App Router app with shadcn/ui (radix-nova style, lucide icons), Tailwind CSS v4, React 19, TypeScript strict mode. `@/*` maps to project root.

### Route groups

- `app/(auth)/` — signin/signup/recovery/exchange-code, split-screen layout (`BrandPanel` + form), no app shell
- `app/(doc)/` — static doc/legal pages (e.g. account-agreement), no app shell
- `app/(main)/` — the authenticated app shell
  - `(workbench)/page.tsx` — root `/` dashboard
  - `(zus)/zus/...` — the "Zus" alerting product, all team-scoped (see below): `(overview)`, `(alert)/alerts`, `(alert)/alert/history`, `(connector)/connector[/...][/new][/[connetcor_id]]`, `(team)/members`
  - `(zus)/admin/...` — admin area: teams, users
  - `settings/...`, `table/...` (data-grid demo)

### Layout & root providers

- `app/layout.tsx` — root layout (only Server Component): imports `globals.css`, wraps everything in `Providers` (ThemeProvider from next-themes + TooltipProvider + Toaster)
- `app/(main)/layout.tsx` — the main app shell: `SidebarProvider` → `AppNavbar` + `MenuProvider` → `AppSidebar` + `SidebarInset`
- Every other component is `"use client"`

### Sidebar menu system

Menus are declared statically per route segment, not fetched from an API:

- Each route layout defines a `Menu` object (types in `components/app-sidebar/sidebar.ts`: `Menu` → `NavGroup[]` → `NavItem`, with Lucide icon components imported directly) and registers it with `<SidebarMenuLoader menu={MENU} routeKey="/xxx" />`.
- `MenuProvider`/`useMenu` (`components/app-sidebar/nav-main.tsx`) keep a `Map<routeKey, Menu>` registry; the menu registered under the longest (deepest) `routeKey` wins. The loader unregisters its entry on unmount, so navigating back restores the parent route's menu.
- `AppSidebar` renders `menu.header`/`menu.footer` in dedicated slots and `<NavMain />` for `navGroups`.
- `ApiMenu`/`ApiNavItem`/`ApiNavGroup` types and `lib/sidebar-icons.ts` (icon-name → Lucide mapping) exist for a JSON/API-driven variant but aren't wired into any current route — follow the static `Menu`-object pattern above for new sections.

### Auth & API client

- `service/fetch.ts` `base()` — ky-based HTTP client. Prefixes `NEXT_PUBLIC_BASE_API_PATH`, sends `credentials: "include"`, JSON by default; `download`/`audio`/`zip` content-types are returned as `Blob`. Non-2xx responses toast the error message (unless `silent`) and reject.
- `service/base.ts` `request()` (and `get`/`post`/`put`/`patch`/`del`) wrap `base()`. On a 401: if `reason === "unauthorized"` it clears the session and reloads; otherwise it calls `refreshAccessTokenOrRelogin` (`service/refresh-token.ts` — cross-tab refresh lock via `localStorage`, single retry) and retries the original request once before redirecting to `/signin`.
- The access token is a JWT stored in `localStorage` (`access_token`); `getAccessToken()` in `service/fetch.ts` reads it for the `Authorization`/`X-Session-Token` headers.
- Service modules under `service/` (`zus-team.ts`, `zus-connector.ts`, `zus-alert.ts`, `auth.ts`) are thin typed wrappers returning `{ message: string; data: T }` — match this envelope shape when adding new endpoints.

### Team-scoped pages (`(zus)` routes)

- `useZusTeam()` (`hooks/use-zus-team.ts`) loads the user's teams and tracks the selected team, persisted to `localStorage` (`zus_team_id`) and synced across components/tabs via a custom window event + `storage` event. `<TeamSelect />` (`components/zus/team-select.tsx`) is the picker, usually placed in a menu's `header`/`footer`.
- Data-fetching pages follow the pattern in `connector/page.tsx`:
  - A `useCallback` fetcher (e.g. `loadConnectors`) takes `team`, pagination, and filters as deps and writes both `connectorState` (the fetched payload + the params it was fetched for) and `error`.
  - `useEffect(() => { void (async () => { await loadConnectors() })() }, [loadConnectors])` — the async-IIFE wrapper is required: calling the fetcher directly (`void loadConnectors()`) trips the React Compiler's `react-hooks/set-state-in-effect` lint rule on pages that don't call `useReactTable()` directly (pages that do are exempted via the `incompatible-library` bailout).
  - A `requestIdRef` counter discards stale responses when team/pagination/filters change quickly.
  - Derive `connectors`/`totalCount`/`loading` by checking whether `connectorState`'s params still match the current params (`isCurrentConnectorState`) instead of a separate boolean that can race with data.

### Data table (`components/data-table/data-table.tsx`)

- `<DataTable table={table} emptyState?={node}>{children}</DataTable>` wraps a TanStack `useReactTable` instance built with manual pagination — renders the header/body/`DataTablePagination`. `children` render above the table (typically filters + `<DataTableViewOptions table={table} />`).
- Body falls back to a "No results." row when `table.getRowModel().rows` is empty; pass `emptyState` to render a custom node instead (e.g. `connector-table.tsx` shows a "no connectors" vs "no results for this filter + reset" panel via `ConnectorStatePanel`, while filters/view-options/pagination stay visible).
- `components/data-table/` also ships a fuller advanced-filtering toolkit (`data-table-filter-list/menu/sort-list/advanced-toolbar`, `hooks/use-data-table.ts` with nuqs-based URL state sync, `config/data-table.ts`, `lib/parsers.ts`) — not currently used by any page, available for building richer tables.

### Data grid (`components/data-grid/`, `hooks/use-data-grid.ts`, `types/data-grid.ts`)

A production-grade spreadsheet-like grid built on `@tanstack/react-table` + `@tanstack/react-virtual`. Key design decisions:

- **Custom store** (`DataGridStore`) using `useSyncExternalStore` — avoids re-rendering the entire grid on every cell interaction. State updates are batched via microtasks.
- **TableMeta** (`types/data-grid.ts`) extends TanStack Table's `TableMeta` with all grid operations (cell editing, selection, clipboard, file uploads, search). This is the primary API surface — cells access grid behavior through `tableMeta`.
- **`useDataGrid` hook** returns a stable memoized object containing refs, virtualizer data, and table instance. The `DataGrid` component consumes this as its props.
- **Cell variants**: short-text, long-text, number, select, multi-select, checkbox, date, url, file — defined in `CellOpts` and handled in `data-grid-cell-variants.tsx` and paste logic.
- **Keyboard navigation**: full Excel-like navigation with shift/ctrl modifiers, search (Ctrl+F), clipboard (Ctrl+C/V/X), and row operations (Shift+Enter to add, Ctrl+Backspace to delete rows).

<claude-mem-context>
# Recent Activity

<!-- This section is auto-generated by claude-mem. Edit content outside the tags. -->

*No recent activity*
</claude-mem-context>