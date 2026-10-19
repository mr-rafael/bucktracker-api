# Frontend Architecture

The BuckTracker UI is a static multi-page site: HTML pages plus a few shared JavaScript modules. There is no SPA framework, bundler, or `package.json`. Docker copies `frontend/` into an nginx image and serves it on port `3000`. Pages call the REST API at `http://localhost:8080`.

Backend layers and API routes are documented in [ARCHITECTURE.md](ARCHITECTURE.md). HTTP request and response shapes are in [api.md](api.md).

## Tech Stack

| Layer | Technology |
|---|---|
| Markup | HTML5 pages (one file per screen) |
| Styling | [Tailwind CSS v4](https://tailwindcss.com/) via `@tailwindcss/browser` |
| Components | [daisyUI v5](https://daisyui.com/) (`corporate` theme) |
| Scripting | Vanilla JavaScript (ES modules are not used; scripts attach to `window`) |
| HTTP | Browser `fetch` |
| Auth storage | Access JWT in `sessionStorage`; refresh token as an HTTP-only `refresh_token` cookie |
| JSON highlighting | [Highlight.js](https://highlightjs.org/) 11.11.1 on create-success pages |
| Runtime | nginx Alpine (`frontend/Dockerfile`) |

CDN assets are loaded from `cdn.jsdelivr.net` in each page `<head>`. Shared app scripts are served from the same origin (`/authenticated-api.js`, `/formatters.js`, and so on).

## Libraries

| Library | Where it is used | Role |
|---|---|---|
| `@tailwindcss/browser@4` | Every page | Compiles Tailwind utility classes in the browser |
| `daisyui@5` | Every page | Component classes (`btn`, `card`, `input`, `modal`, `table`, `alert`, `menu`) and the `data-theme="corporate"` palette |
| `@highlightjs/cdn-assets@11.11.1` | Create-success pages | Syntax-highlights the JSON body returned after creating a user, loan, savings plan, or payment plan |
| Browser `Intl.NumberFormat` | `formatters.js` | Formats cents as USD |

There is no React, Vue, Alpine, or other UI framework. Page behavior lives in an inline `<script>` at the bottom of each HTML file.

## File Structure

```
frontend/
├── Dockerfile                  # nginx:alpine; copies this directory to /usr/share/nginx/html
├── authenticated-api.js        # Bearer fetch + access-token refresh
├── api-errors.js               # Format HTTP and network errors for the UI
├── formatters.js               # Dollars, dates, percentages
├── plan-display.js             # Render calculator result lists and tables
├── delete-modal.js             # Confirm-and-delete dialog
├── index.html                  # Welcome / health check
├── login.html
├── create-account.html
├── create-account-success.html
├── dashboard.html              # Saved loans and savings lists
├── loan-calculator.html        # Public loan calculate
├── savings-calculator.html     # Public savings calculate
├── create-loan.html
├── create-savings.html
├── create-payment-plan.html
├── loan-success.html
├── savings-success.html
├── payment-plan-success.html
└── detail.html                 # Loan, savings, or payment-plan detail / edit
```

## Shared Scripts

Each script is an IIFE that exposes a frozen object on `window`. Pages include only the files they need.

| Script | Global | Responsibility |
|---|---|---|
| `api-errors.js` | `window.apiErrors` | Read an error body (`error` or `message`) or wrap a network failure as `NETWORK_ERROR: ...` |
| `formatters.js` | `window.formatters` | Convert API cents to USD, ISO timestamps to `YYYY-MM-DD`, and numeric rates to percents |
| `authenticated-api.js` | `window.authenticatedApi` | Attach `Authorization: Bearer ...`, send cookies, refresh an expired access token, redirect home if refresh fails |
| `plan-display.js` | `window.planDisplay` | Build definition lists and zebra tables for calculator results |
| `delete-modal.js` | `window.deleteModal` | daisyUI `<dialog>` for delete confirmation; issues `DELETE` through `authenticatedApi` |

### Auth helper

`authenticatedApi.fetch(url, options)`:

1. Reads `sessionStorage.access_token`.
2. If the JWT `exp` is missing or in the past, calls `POST /app/refresh` with `credentials: "include"`.
3. Sends the request with the Bearer header.
4. On `401`, refreshes once and retries. A second `401` (or a failed refresh) clears the token and replaces the location with `/`.

Login uses a plain `fetch` to `POST /app/login` with `credentials: "include"`, then stores `access_token` and navigates to `dashboard.html`. Logout is a link back to `/` (the token is tab-scoped in `sessionStorage`).

## Pages

| Page | Auth | Role |
|---|---|---|
| `index.html` | No | Welcome screen; pings `GET /api/healthz`; links to login and public calculators |
| `login.html` | No | Email/password login |
| `create-account.html` | No | `POST /app/users/create` |
| `create-account-success.html` | No | Shows the create-user JSON from `sessionStorage` |
| `loan-calculator.html` | No | `POST /app/loans/calculate`; `?from=welcome` or `?from=dashboard` sets the back link |
| `savings-calculator.html` | No | `POST /app/savings/calculate`; same `from` query param |
| `dashboard.html` | Yes | Lists loans (`GET /app/loans/list`) or savings (`GET /app/savings/list`) via `#loans` / `#savings` |
| `create-loan.html` | Yes | `POST /app/loans/save` |
| `create-savings.html` | Yes | `POST /app/savings/save` |
| `create-payment-plan.html` | Yes | Loads the loan, then `POST /app/loans/{loanId}/payment-plans` |
| `loan-success.html` / `savings-success.html` / `payment-plan-success.html` | Yes | Confirmation + highlighted JSON; link to `detail.html` |
| `detail.html` | Yes | `?type=loans\|savings\|payment-plans&id=...` (payment plans also need `loanId`); view, rename, edit fields, delete |

Money inputs on create/calculator forms accept dollars in the UI and convert to cents before calling the API (`Math.round(dollars * 100)`).

## Navigation

```
index.html
  ├── login.html → dashboard.html
  │     └── create-account.html → create-account-success.html → login.html
  ├── loan-calculator.html
  └── savings-calculator.html

dashboard.html
  ├── create-loan.html → loan-success.html → detail.html?type=loans
  ├── create-savings.html → savings-success.html → detail.html?type=savings
  ├── loan/savings calculators (?from=dashboard)
  └── detail.html
        └── create-payment-plan.html?loanId=... → payment-plan-success.html
              └── detail.html?type=payment-plans
```

Create-success pages read a JSON snapshot from `sessionStorage` (`created_loan_response`, `created_savings_response`, `created_payment_plan_response`, `created_user_response`). That payload is written by the create form immediately before navigation.

## How pages talk to the API

```
Browser (http://localhost:3000)
    │
    │  public: fetch()
    │  authenticated: authenticatedApi.fetch()
    ▼
Go API (http://localhost:8080)
```

- Public pages use `fetch` with no `Authorization` header.
- Authenticated pages use `authenticatedApi.fetch` and include cookies so refresh works.
- Compose sets `ALLOWED_ORIGIN=http://localhost:3000` on the backend so CORS and credentialed requests succeed.

## Serving the UI

`frontend/Dockerfile` is a static file server:

1. Start from `nginx:alpine`.
2. Replace the default site with the contents of `frontend/`.
3. Listen on container port `80`, published as host port `3000`.

There is no frontend build step. Edit an HTML or JS file, rebuild the `frontend` Compose service, and nginx serves the new files.
