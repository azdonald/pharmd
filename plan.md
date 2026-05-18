# PharmD — Multi-Tenant Pharmacy Management System

## Overview
A multi-tenant pharmacy management system where tenants range from single-location independent pharmacies to multi-location chains.

---

## Architecture Decisions

| Decision | Choice | Status |
|----------|--------|--------|
| Backend | Go + chi | ✅ |
| Frontend | React (Vite) | ✅ |
| Database | MySQL (shared, row-level isolation via tenant_id) | ✅ |
| Auth | JWT + RBAC | ✅ |
| API style | REST (JSON) | ✅ |
| API tooling | OpenAPI 3.0 → oapi-codegen (chi-server + types) | ✅ |
| Code gen | `go generate ./...` via `//go:generate` directives | ✅ |
| Migrations | olympian (github.com/ichtrojan/olympian) | ✅ |
| Deployment | Docker | ✅ |

---

## API Design Workflow

Each resource follows this pattern:

```
backend/api/<resource>/
  <resource>.yaml       — hand-written OpenAPI 3.0.3 spec
  gen.go                — //go:generate directives
  server_gen.go         — generated (chi-server interface + routing)
  server_impl.go        — hand-written business logic
  types_gen.go          — generated (request/response structs)
```

**Process:**
1. Write the OpenAPI spec (`<resource>.yaml`) with paths, schemas, and `$ref` to shared models.
2. Write `gen.go` with `//go:generate oapi-codegen -generate=chi-server,types ...`
3. Run `go generate ./...` → produces `server_gen.go` + `types_gen.go`.
4. Hand-write `server_impl.go` implementing the generated `ServerInterface`.

Shared models live in `common/common.yaml` and are referenced via `$ref`.

---

## Project Structure

```
pharmd/
├── backend/
│   ├── main.go
│   ├── cmd/migrate/main.go
│   ├── server/http.go            — router, middleware wiring, startup
│   ├── db/db.go                  — MySQL connection
│   ├── models/                   — domain structs
│   ├── repository/               — data access (SQL)
│   ├── service/                  — business logic
│   ├── middleware/               — auth, tenant, audit
│   ├── utils/                    — response, token, permissions
│   ├── api/
│   │   ├── common/               — shared schemas + types
│   │   ├── auth/                 — authentication
│   │   ├── locations/            — location CRUD
│   │   ├── users/                — user CRUD
│   │   ├── patients/             — patient management
│   │   ├── products/             — drug/product catalog
│   │   ├── inventory/            — stock management
│   │   ├── suppliers/            — supplier management
│   │   ├── purchases/            — purchase orders + GRN
│   │   ├── prescriptions/        — prescription management
│   │   ├── dispensing/           — dispensing workflow
│   │   ├── pos/                  — point of sale
│   │   ├── pricing/              — pricing & tax rules
│   │   ├── reports/              — reporting
│   │   └── ...                   — more as needed
│   ├── migrations/               — olympian migration files
│   ├── Dockerfile
│   └── .env
├── web/                          — React + Vite frontend
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── api/                  — API client
│   │   ├── hooks/
│   │   └── ...
│   └── Dockerfile
├── docker-compose.yml
└── plan.md
```

---

# Feature List & Implementation Plan

## Priority Levels
- **P0** — must-have for any release (blocking)
- **P1** — core pharmacy functionality
- **P2** — important but not blocking
- **P3** — nice-to-have / future

---

## Iteration 1 — Auth & Tenant Foundation

### F1. Organisation Registration & Login
**Priority:** P0 | **API:** `auth/` | **Frontend:** Login/Register pages

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 1.1 | Register endpoint (org + admin user created together) | `POST /v1/register` | Register page |
| 1.2 | Login endpoint (email + password → JWT) | `POST /v1/login` | Login page |
| 1.3 | Token refresh endpoint | `POST /v1/refresh` | Auto-refresh logic |
| 1.4 | Logout endpoint | `POST /v1/logout` | Logout button |
| 1.5 | Change password endpoint | `PUT /v1/change-password` | Change password form |
| 1.6 | JWT auth middleware | `middleware/auth.go` | — |
| 1.7 | Tenant context middleware (resolve org_id from JWT) | same as above | — |

**Files to create:**
- `api/auth/auth.yaml` + `gen.go` → generate → `server_impl.go` ✅ *done*
- `repository/auth_repo_impl.go` ✅ *done*
- `service/auth_service_impl.go` ✅ *done*
- `middleware/auth.go` ✅ *done*
- `web/src/pages/Login.tsx`, `Register.tsx`
- `web/src/api/auth.ts`

**DB Migrations:**
- `1765136578_create_organisations_table.go` ✅ *done*
- `1765310022_create_user_table.go` ✅ *done*

---

### F2. RBAC (Roles & Permissions)
**Priority:** P0 | **API:** `permissions/`, `roles/` | **Frontend:** Role management UI

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 2.1 | Seeded permission slugs (org:create, users:*, locations:*, etc.) | migration + `utils/permissions.go` | — |
| 2.2 | Role CRUD per organisation | `api/roles/` | Roles page |
| 2.3 | Assign permissions to role | `api/roles/{id}/permissions` | Role edit page |
| 2.4 | Assign role to user | `api/users/{id}/roles` | User edit page |
| 2.5 | Permission-checking middleware | `middleware/auth.go` `RequirePermission()` | Route guards |
| 2.6 | Seed default roles (Admin, PIC, Pharmacist, Tech, Cashier) | migration | — |

**DB Migrations:**
- `1767000000_create_permissions_table.go` ✅ *done*
- `1767000001_create_roles_table.go` ✅ *done*
- `1767000002_create_role_permissions_table.go` ✅ *done*
- `1767000003_create_user_roles_table.go` ✅ *done*
- `1767000004_seed_permissions.go`
- `1767000005_seed_default_roles.go`

---

### F3. Location Management
**Priority:** P0 | **API:** `locations/` | **Frontend:** Location settings page

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 3.1 | List/create locations per org | `api/locations/` | Location list page |
| 3.2 | Get/update/delete location | `api/locations/{id}` | Location edit form |
| 3.3 | Auto-create default location for single-location orgs | on org registration | — |
| 3.4 | Location-scoped middleware (inject location_id from header or user context) | middleware | — |

**DB Migrations:**
- `1766000000_create_locations_table.go` ✅ *done*

---

### F4. User Management
**Priority:** P0 | **API:** `users/` | **Frontend:** User management page

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 4.1 | List users (paginated, scoped to org) | `api/users/` | User list |
| 4.2 | Create user (invite with temp password) | `POST /users` | Create user form |
| 4.3 | Update user profile | `PUT /users/{id}` | Edit user form |
| 4.4 | Toggle user active/inactive | `PATCH /users/{id}/status` | Toggle switch |
| 4.5 | Assign user to locations | `POST /users/{id}/locations` | Location assignment |
| 4.6 | Assign user role | handled in F2.4 | handled in F2.4 |

---

## Iteration 2 — Core Pharmacy

### F5. Patient / Customer Management
**Priority:** P1 | **API:** `patients/` | **Frontend:** Patient module

| # | Task |
|---|------|
| 5.1 | Patient registration (name, DOB, gender, phone, email, address) |
| 5.2 | Patient search (by name, phone, ID) |
| 5.3 | Patient profile (visit history, prescriptions, allergies) |
| 5.4 | Allergies & medical conditions registry |
| 5.5 | Insurance / NHIS linkage |
| 5.6 | Duplicate detection on registration |

**DB Tables:** `patients`, `patient_allergies`, `patient_conditions`, `patient_insurance`

---

### F6. Drug / Product Catalog
**Priority:** P1 | **API:** `products/` | **Frontend:** Product management

| # | Task |
|---|------|
| 6.1 | Product CRUD (brand name, generic name, manufacturer, form, strength) |
| 6.2 | Drug classification (Rx, OTC, Controlled, Narcotic) |
| 6.3 | Barcode / NDC code tracking |
| 6.4 | Product categories |
| 6.5 | Generic substitution mapping |
| 6.6 | Storage conditions flag |

**DB Tables:** `products`, `product_categories`, `generic_substitutions`

---

### F7. Inventory & Stock Management
**Priority:** P1 | **API:** `inventory/` | **Frontend:** Inventory dashboard

| # | Task |
|---|------|
| 7.1 | Stock batch tracking (batch#, expiry, qty, cost price) per location |
| 7.2 | Stock receipt (add stock to location) |
| 7.3 | Stock adjustment (waste, damage, theft, found) |
| 7.4 | Low-stock alerts (configurable threshold per product) |
| 7.5 | Expiry tracking & alerts |
| 7.6 | Stock count / reconciliation |
| 7.7 | Inventory valuation (FIFO) |

**DB Tables:** `stock_batches`, `stock_movements`, `inventory_alerts`

---

### F8. Supplier Management
**Priority:** P1 | **API:** `suppliers/` | **Frontend:** Supplier list

| # | Task |
|---|------|
| 8.1 | Supplier CRUD (name, contact, payment terms) |
| 8.2 | Supplier price list per product |

**DB Tables:** `suppliers`, `supplier_products`

---

### F9. Purchase Orders
**Priority:** P1 | **API:** `purchases/` | **Frontend:** PO creation & tracking

| # | Task |
|---|------|
| 9.1 | Create purchase order (select supplier, line items, expected date) |
| 9.2 | PO approval workflow |
| 9.3 | Goods Received Note (GRN) — match against PO, update stock |
| 9.4 | PO status tracking (draft → sent → approved → received → cancelled) |

**DB Tables:** `purchase_orders`, `purchase_order_items`, `goods_received_notes`

---

### F10. Pricing & Tax
**Priority:** P1 | **API:** `pricing/` | **Frontend:** Pricing rules

| # | Task |
|---|------|
| 10.1 | Base selling price per product per location |
| 10.2 | Pricing formula (cost + markup %, fixed price) |
| 10.3 | Tax rate per location |
| 10.4 | Discount rules (percentage cap, max amount, approval threshold) |
| 10.5 | Insurance price schedules |

**DB Tables:** `product_prices`, `tax_rates`, `discount_rules`

---

## Iteration 3 — Prescription & Dispensing

### F11. Prescription Management
**Priority:** P1 | **API:** `prescriptions/` | **Frontend:** Prescription intake

| # | Task |
|---|------|
| 11.1 | Prescription intake (walk-in, external) |
| 11.2 | Prescriber details (name, license#, phone) |
| 11.3 | Line items (drug, qty, dosage, frequency, duration, refills) |
| 11.4 | Prescription validation (patient exists, drug exists, refill check) |
| 11.5 | Prescription status workflow (active → dispensed → exhausted → expired) |
| 11.6 | Refill tracking |

**DB Tables:** `prescriptions`, `prescription_items`, `prescribers`, `refills`

---

### F12. Dispensing Workflow
**Priority:** P1 | **API:** `dispensing/` | **Frontend:** Dispensing queue

| # | Task |
|---|------|
| 12.1 | Dispensing queue per location (pending → in-progress → dispensed → collected) |
| 12.2 | Drug interaction check at point of dispensing |
| 12.3 | Allergy check against patient record |
| 12.4 | Generic / therapeutic substitution with approval |
| 12.5 | Partial dispensing (dispense less than prescribed qty) |
| 12.6 | Controlled substance logs (with witness name, ID check) |
| 12.7 | Label printing |
| 12.8 | Patient information leaflet |

**DB Tables:** `dispensing_records`, `dispensing_logs`, `controlled_substance_logs`

---

### F13. Point of Sale
**Priority:** P1 | **API:** `pos/` | **Frontend:** POS interface

| # | Task |
|---|------|
| 13.1 | OTC sales (add products, calculate total, apply discount) |
| 13.2 | Prescription sales (link to dispensed record) |
| 13.3 | Payment methods (cash, card, transfer, mobile money, insurance) |
| 13.4 | Receipt printing |
| 13.5 | Hold / void / refund transactions |
| 13.6 | Daily sales summary (X-report) |
| 13.7 | End-of-day closeout (Z-report) |

**DB Tables:** `sales`, `sale_items`, `payments`, `refunds`, `daily_summaries`

---

## Iteration 4 — Financial & Reporting

### F14. Accounts & Financials
**Priority:** P2 | **API:** `accounts/` | **Frontend:** Financial dashboard

| # | Task |
|---|------|
| 14.1 | Daily cash-up / till reconciliation |
| 14.2 | Accounts receivable (credit patients, insurance claims aging) |
| 14.3 | Accounts payable (supplier invoice tracking) |
| 14.4 | Expense tracking (utilities, rent, salaries) |
| 14.5 | Profit margin per product / category report |

**DB Tables:** `cashups`, `accounts_receivable`, `accounts_payable`, `expenses`

### F15. Reports & Analytics
**Priority:** P2 | **API:** `reports/` | **Frontend:** Reports dashboard

| # | Task |
|---|------|
| 15.1 | Sales report (by date range, location, cashier, payment method) |
| 15.2 | Inventory report (stock value, aging, turnover) |
| 15.3 | Prescription report (volume, top drugs, top prescribers) |
| 15.4 | Expiry report (products expiring in next N days) |
| 15.5 | Low-stock report |
| 15.6 | Dispensing statistics |
| 15.7 | Dashboard widgets (today's revenue, pending Rx count, alerts) |
| 15.8 | Export to CSV / PDF |

**DB Tables:** `report_cache` (if needed)

### F16. Audit & Compliance
**Priority:** P2 | **API:** `audit/` | **Frontend:** Audit log viewer

| # | Task |
|---|------|
| 16.1 | Activity log for all mutations (who did what, when, to what) |
| 16.2 | Controlled substance perpetual inventory log |
| 16.3 | Prescription record retention & search |
| 16.4 | Data export tooling (GDPR right-to-access) |
| 16.5 | Data deletion tooling (GDPR right-to-erasure) |

**DB Tables:** `activity_logs`

---

## Iteration 5 — Multi-Location & Platform

### F17. Multi-Location Chain Features
**Priority:** P2 | **API:** `chain/` | **Frontend:** Chain dashboard

| # | Task |
|---|------|
| 17.1 | Centralized purchasing (HQ creates POs, distributes to locations) |
| 17.2 | Inter-location stock transfer request → approval → fulfillment |
| 17.3 | Consolidated reporting across all locations |
| 17.4 | Cross-location patient lookup |
| 17.5 | Regional manager / area supervisor role with multi-location scope |
| 17.6 | Location-level override on master catalog (different prices per store) |

### F18. Platform Administration
**Priority:** P3 | **API:** `admin/` | **Frontend:** Admin panel

| # | Task |
|---|------|
| 18.1 | Super admin tenant management (view all orgs, suspend, activate) |
| 18.2 | Feature flags per tenant |
| 18.3 | Rate limiting |
| 18.4 | Backup / restore CLI |
| 18.5 | Health monitoring dashboard |

### F19. Integrations
**Priority:** P3

| # | Task |
|---|------|
| 19.1 | SMS notifications (refill reminders, pickup alerts) |
| 19.2 | Email notifications |
| 19.3 | Insurance claims submission |
| 19.4 | NDC / RxNorm drug code lookup |
| 19.5 | Accounting export (CSV for QuickBooks) |

---

## Implementation Roadmap

```
Iteration 1: Auth & Foundation     → F1, F2, F3, F4    → P0 features
Iteration 2: Core Pharmacy          → F5, F6, F7, F8, F9, F10 → P1 features
Iteration 3: Prescription & POS     → F11, F12, F13     → P1 features
Iteration 4: Financial & Reporting  → F14, F15, F16     → P2 features
Iteration 5: Multi-Location & Ops   → F17, F18, F19     → P2/P3 features
```

---

## Design Principles
- **Tenant isolation first** — every query includes `organisation_id` filter; never trust client-side.
- **Location-aware from day one** — inventory, pricing, and user access scoped to location.
- **Audit everything** — pharmacy is regulated; every stock movement, dispense, and user action is logged.
- **Regulatory ready** — HIPAA/GDPR, controlled substance tracking, prescription retention laws.
