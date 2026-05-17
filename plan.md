# PharmD — Multi-Tenant Pharmacy Management System

## Overview
A multi-tenant pharmacy management system where tenants range from single-location independent pharmacies to multi-location chains.

---

## Architecture Decisions

| Decision | Choice | Status |
|----------|--------|--------|
| Backend    | Go + chi | ✅ |
| Frontend   | React (Vite) | ✅ |
| Database   | MySQL (shared, row-level isolation via tenant_id) | ✅ |
| Auth       | JWT + RBAC | ✅ |
| API style  | REST (JSON) | ✅ |
| API tooling | OpenAPI 3.0 → oapi-codegen (chi-server + types) | ✅ |
| Code gen    | `go generate ./...` via `//go:generate` directives | ✅ |
| Migrations  | olympian (github.com/ichtrojan/olympian) | ✅ |
| Deployment  | Docker | ✅ |

---

## API Design Workflow

Each resource follows this pattern (inspired by `spleak`):

```
backend/api/<resource>/
  <resource>.yaml       — hand-written OpenAPI 3.0.3 spec
  gen.go                — //go:generate directives
  server_gen.go         — generated (chi-server interface + routing)
  server_impl.go        — hand-written business logic
  types_gen.go          — generated (request/response structs)
  server_gen.yaml       — config for oapi-codegen (server)
  types_gen.yaml        — config for oapi-codegen (types)
```

**Process:**
1. Write the OpenAPI spec (`<resource>.yaml`) with paths, schemas, and `$ref` to shared models.
2. Write `gen.go` with `//go:generate oapi-codegen -config=server_gen.yaml <resource>.yaml` and similarly for types.
3. Run `go generate ./...` to produce `server_gen.go` and `types_gen.go`.
4. Hand-write `server_impl.go` implementing the generated `ServerInterface`.

Shared models (User, Address, etc.) live in `common/common.yaml` and are referenced via `$ref`.

---

## Phase 1 — Foundation

### 1. Multi-Tenancy Core
- [ ] Tenant registration & onboarding flow
- [ ] Tenant context middleware (resolve tenant from domain/subdomain/header)
- [ ] Row-level tenant isolation on all data tables
- [ ] Tenant configuration store (settings, preferred currency, locale, tax rules)
- [ ] Tenant status lifecycle (trial → active → suspended → cancelled)

### 2. Authentication & Authorization
- [ ] User registration (invite-only by tenant admin)
- [ ] Role-based access control (RBAC)
  - System Admin (cross-tenant)
  - Tenant Admin
  - Pharmacist-in-Charge
  - Pharmacist
  - Technician
  - Cashier / Front-desk
- [ ] Login / logout (JWT / session)
- [ ] Password reset, MFA (optional)
- [ ] Audit log for auth events

### 3. Organization & Location Management
- [ ] Organization entity (the "tenant" root)
- [ ] Location entity (belongs to organization)
  - Single-location orgs: 1 location created automatically
  - Multi-location orgs: CRUD for locations
- [ ] Address, contact info, license/registration per location
- [ ] Operating hours per location
- [ ] Location-level settings (tax rate, markup, timezone)

### 4. User Management
- [ ] User profile (name, email, phone, role)
- [ ] User-location assignment (which locations a user can operate)
- [ ] User status (active / disabled)
- [ ] Activity log (who did what)

---

## Phase 2 — Core Pharmacy Operations

### 5. Patient / Customer Management
- [ ] Patient registration (name, DOB, gender, contact, address)
- [ ] Patient ID / NHIS / insurance linkage
- [ ] Allergies & medical conditions registry
- [ ] Patient visit history
- [ ] Loyalty / prescription refill reminders
- [ ] Search / deduplication

### 6. Drug & Product Catalog
- [ ] Central product catalog (brand name, generic name, manufacturer)
- [ ] Drug classification (Rx, OTC, Controlled, Narcotic)
- [ ] Drug form (tablet, syrup, injection, cream, etc.)
- [ ] Strength / dosage units
- [ ] Barcode / NDC / RxNorm / ATC codes
- [ ] Drug-drug interaction data interface
- [ ] Storage conditions (temperature, light-sensitive)
- [ ] Alternate products / generic substitution mapping

### 7. Inventory Management
- [ ] Stock per location (not global — each location has its own inventory)
- [ ] Stock batches (batch number, expiry date, quantity, purchase price)
- [ ] Stock movements (received, transferred, dispensed, returned, adjusted)
- [ ] Stock alerts (low stock, about-to-expire, expired)
- [ ] Inventory transfers between locations (multi-tenant chain only)
- [ ] Stock count / reconciliation
- [ ] Inventory valuation (FIFO / average cost)
- [ ] Waste / damage logging (narcotic waste logs for controlled substances)

### 8. Supplier Management
- [ ] Supplier registration (name, contact, payment terms)
- [ ] Supplier price lists
- [ ] Purchase order creation & approval workflow
- [ ] Goods received note (GRN) — match against PO
- [ ] Supplier performance (lead time, fill rate)

### 9. Prescription Management
- [ ] Prescription intake (walk-in, e-prescription, fax)
- [ ] Prescription record (prescriber, patient, date, diagnosis)
- [ ] Line items (drug, quantity, dosage instructions, refills)
- [ ] Prescription validation (DEA number, patient identity check)
- [ ] Refill management & tracking
- [ ] Electronic prescribing integration (future)
- [ ] Prescription history per patient

### 10. Dispensing Workflow
- [ ] Dispensing queue (pending → in-progress → dispensed → collected)
- [ ] Drug interaction check at point of dispensing
- [ ] Allergy check against patient record
- [ ] Substitution workflow (therapeutic / generic)
- [ ] Label printing (patient name, drug, dosage, warnings)
- [ ] Leaflet / patient information sheet generation
- [ ] Partial dispensing support
- [ ] Controlled substance dispensing logs (with witness)

### 11. Point of Sale (POS)
- [ ] Simple POS for over-the-counter (OTC) sales
- [ ] Prescription sales (linked to dispensing)
- [ ] Multiple payment methods (cash, card, mobile money, insurance)
- [ ] Receipt printing (thermal printer support)
- [ ] Invoice & receipt history
- [ ] Daily sales summary / X-report / Z-report
- [ ] Hold / void / refund transactions
- [ ] Price overrides (with approval if needed)

### 12. Pricing & Tax
- [ ] Product base price per location
- [ ] Pricing formulas (cost + markup, fixed price, tiered)
- [ ] Tax calculation per location / jurisdiction
- [ ] Discount management (percentage / fixed, with max limit)
- [ ] Insurance / NHIS price schedules

---

## Phase 3 — Financial & Reporting

### 13. Accounts & Financials
- [ ] Daily cash-up / till management
- [ ] Accounts receivable (credit patients, insurance claims)
- [ ] Accounts payable (supplier invoices)
- [ ] Expense tracking
- [ ] Profit margin reporting per product / category
- [ ] Sales tax reporting

### 14. Reporting & Analytics
- [ ] Sales reports (daily, weekly, monthly, by location, by cashier)
- [ ] Inventory reports (stock value, aging, turnover)
- [ ] Prescription reports (volume, top drugs, by prescriber)
- [ ] Expiry reports
- [ ] Low-stock alerts report
- [ ] Dispensing statistics
- [ ] Dashboard widgets (revenue, alerts, recent activity)
- [ ] Export to CSV / PDF

### 15. Audit & Compliance
- [ ] Audit log for all inventory movements
- [ ] Controlled substances perpetual inventory log
- [ ] Prescription record retention & search
- [ ] HIPAA / GDPR / data privacy compliance
- [ ] Data export / right-to-deletion tooling

---

## Phase 4 — Advanced & Platform

### 16. Multi-Location Features
- [ ] Centralized purchasing (PO from HQ, distributed to locations)
- [ ] Inter-location stock transfer workflow
- [ ] Consolidated reporting across all locations
- [ ] Master product catalog (global) + location-specific overrides
- [ ] Cross-location patient lookup
- [ ] Role hierarchy: regional manager, area supervisor

### 17. Subscription & Billing (for the platform itself)
- [ ] Tenant subscription tiers (single-location vs multi-location)
- [ ] Usage-based metering (active users, transactions)
- [ ] Invoicing & payment collection
- [ ] Upgrade / downgrade / cancellation
- [ ] Trial period management

### 18. Integrations
- [ ] SMS / email notifications (refill reminders, alerts)
- [ ] Insurance claims submission (electronic)
- [ ] National drug registries (NDC / RxNorm lookup)
- [ ] E-prescription network
- [ ] Accounting software export (QuickBooks, etc.)
- [ ] Barcode scanner integration

### 19. System Administration
- [ ] Tenant management dashboard (super admin)
- [ ] Feature flags per tenant
- [ ] Rate limiting & abuse protection
- [ ] Backup & restore
- [ ] Logging & monitoring
- [ ] API rate limits & API key management

---

## Implementation Order (Recommended)

| Phase | Focus | Est. Duration |
|-------|-------|---------------|
| 1 | Foundation — multi-tenancy, auth, org/location, user mgmt | — |
| 2 | Core — patients, catalog, inventory, suppliers, Rx, dispensing, POS, pricing | — |
| 3 | Financial — cash mgmt, AR/AP, reporting, audit | — |
| 4 | Platform — multi-location, billing, integrations, admin | — |

---

## Design Principles
- **Tenant isolation first** — no accidental data leaks between tenants.
- **Location-aware from day one** — all inventory, pricing, and user access is scoped to location.
- **Audit everything** — pharmacy is regulated; every stock movement, dispense, and user action is logged.
- **Offline resilience** — pharmacies cannot afford downtime; design for eventual sync where possible.
- **Regulatory ready** — HIPAA/GDPR, controlled substance tracking, prescription retention laws.
