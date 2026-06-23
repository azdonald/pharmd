# PharmD Guide

## My Understanding

PharmD is intended to be a multi-tenant pharmacy management system for independent pharmacies and multi-location pharmacy groups. The product should let a pharmacy register an organisation, manage users and roles, maintain locations, keep patient and product records, track inventory, manage suppliers and purchase orders, process prescriptions and dispensing, and complete point-of-sale transactions.

The app is currently a Go + chi backend with MySQL persistence and a React/Vite frontend. The backend is organized around OpenAPI resource modules, repositories, services, migrations, and route handlers. The frontend is a protected admin workspace under `/app` with pages for dashboard, users, roles, permissions, locations, patients, products, categories, inventory, suppliers, purchases, pricing, discounts, prescribers, prescriptions, dispensing, POS, and sales history.

At a product level, the app appears to be aiming for a full pharmacy operating system:

- Tenant setup: organisation registration, login, password change, onboarding state.
- Administration: users, roles, permissions, and locations.
- Pharmacy master data: patients, allergies, conditions, prescribers, products, categories, substitutes, prices, discounts, suppliers.
- Supply chain: stock batches, adjustments, counts, low-stock alerts, expiring stock, purchase orders, approvals, receiving.
- Clinical workflow: prescription intake, refills, dispensing queue, dispense record, labels, controlled-drug witness fields.
- Retail workflow: POS sale creation, payments, receipts, sales history, void/refund/hold status, daily summary and closeout.

## What Exists Now

The app has a broad working skeleton and many implemented backend/frontend pairs:

- Auth and tenant foundation exist through `/v1/register`, `/v1/login`, `/v1/refresh`, `/v1/logout`, and `/v1/change-password`.
- JWT middleware adds `organisation_id`, `user_id`, and user permissions to request context.
- Most resource routes enforce backend permissions with `RequirePermission`.
- Data models and migrations cover organisations, users, roles, permissions, locations, patients, products, inventory, suppliers, purchase orders, pricing, prescriptions, dispensing, POS, payments, and daily summaries.
- Repositories generally scope queries by `organisation_id`, which is the right tenant-isolation direction.
- The frontend has navigable pages and API clients for the major workflows.
- Docker Compose can run MySQL, backend, and web services.

## Important Gaps

### 1. Clinical safety checks are placeholders

Dispensing exposes drug-interaction and allergy-check endpoints, but `CheckInteractions` and `CheckAllergies` currently return empty arrays. This means the UI/API can imply clinical safety checks exist even though no real checks are happening.

Recommendation: Treat real clinical checks as P0 before any production use. Start with a clear rules table for allergies and contraindications, then later integrate a drug knowledge source. The UI should clearly distinguish "not checked" from "checked and no issues found."

### 2. Dispensing and POS now deduct inventory (fixed June 2026)

Dispensing (`backend/repository/dispensing_repo_impl.go:110`) and POS payment completion (`backend/repository/pos_repo_impl.go:227`) now use a shared FEFO deduction helper (`backend/repository/inventory_helper.go:11`) that selects batches with remaining stock ordered by earliest expiry, deducts quantity across batches with `SELECT ... FOR UPDATE` row locking, and creates stock movements — all inside the same database transaction. Inventory is also restored on void/refund of completed sales (`RestoreSaleInventory`, `pos_repo_impl.go:262`).

Remaining gaps:
- Dispensing void/refund does not yet restore inventory.
- The POS frontend does not display available stock at the point of adding items to the cart.
- The stock count FEFO helper lives in the `repository` package rather than `service`; consider elevating to a service-layer concern once cross-repository transactions are supported.

### 3. Audit/compliance is not yet system-wide

There are dispense logs and stock movements, but there is no general audit log for user actions, access to patient data, role/permission changes, voids/refunds, prescription changes, or controlled-substance events.

Recommendation: Add an append-only `audit_events` table and log who did what, when, from where, against which entity, including before/after snapshots for high-risk actions. Pharmacy workflows need this from the foundation, not as a later reporting feature.

### 4. Dashboard appears mostly static

The dashboard uses hard-coded metrics, charts, top sellers, and transactions. It is a useful visual target, but it does not look connected to backend summary/report endpoints.

Recommendation: Replace static dashboard values with API-backed widgets. Start with daily revenue, orders, low stock count, expiring stock, pending prescriptions, and open purchase orders.

### 5. Frontend route protection is login-only

The backend enforces many permissions, but the frontend sidebar and routes do not appear to hide or block pages based on user permissions. Users may see navigation they cannot use.

Recommendation: Include permissions in the auth response or add a `/me` endpoint, then drive sidebar visibility and protected route checks from permission slugs.

### 6. Refresh token handling is incomplete

The frontend refreshes in memory and stores tokens in `localStorage`, but refreshed tokens are not written back to `localStorage`. Logout is stateless on the backend, so refresh tokens are not revoked server-side.

Recommendation: Persist refreshed tokens consistently and add server-side refresh token rotation/revocation. For production, consider httpOnly secure cookies or another strategy that reduces XSS token exposure.

### 7. Error responses are inconsistent

Some backend handlers use `http.Error`, while the frontend expects JSON fields like `message` or `error`. `apiRequest` always attempts `res.json()` for non-204 responses, which can break on plain-text errors.

Recommendation: Standardize API error envelopes across all handlers and make the frontend resilient to non-JSON responses.

### 8. OpenAPI specs and registered routes are not perfectly aligned

The auth OpenAPI file describes `/auth/...` paths, while the implemented backend/frontend use `/v1/...` paths. This can confuse generated clients, documentation, and future contributors.

Recommendation: Make the OpenAPI specs match the real route surface, then regenerate code. Going forward, treat OpenAPI as the source of truth.

### 9. No automated tests were found

I did not find Go tests or frontend test files. Given the app handles auth, tenant isolation, stock, dispensing, payments, and regulated workflows, this is the largest engineering risk.

Recommendation: Add focused tests before expanding features. Prioritize tenant isolation, permission enforcement, inventory movement accounting, prescription/dispensing workflows, POS payments, and token refresh behavior.

### 10. Production readiness is still early

Current CORS allows all origins, Docker Compose uses a development token secret, database credentials are simple defaults, and there is no visible CI, seed/demo data strategy, observability, backup/restore plan, or deployment runbook.

Recommendation: Create separate development/staging/production configuration, lock down CORS, require strong secrets, add structured logging, health checks that include DB readiness, and document backup/restore expectations.

## Missing Product Features

The current app has a strong foundation, but these features are still missing or incomplete for a complete pharmacy platform:

- Insurance claims, adjudication, coverage, copays, reversals, and payer integrations.
- Patient insurance records, despite being mentioned in the plan.
- E-prescription intake/import and prescription image/document attachment.
- Pharmacist verification workflow before dispense completion.
- Drug utilization review, interaction checking, duplicate therapy, allergy warnings, and clinical override documentation.
- Controlled-substance register with immutable logs, witness verification, and report/export support.
- Label and receipt printing as real printable templates, not just data endpoints.
- Inventory transfers between locations.
- Returns, recalls, quarantine, expired-stock disposal, and destruction workflows.
- Supplier invoices, accounts payable, expenses, AR/credit accounts, cash-up, and accounting export.
- Reporting module for sales, inventory, expiry, prescriptions, margins, audit, and regulatory exports.
- Platform administration for tenant management and feature flags.
- Notification system for low stock, expiring stock, pending prescriptions, approvals, and failed payments.
- Import/export tools for product catalogs, patients, inventory, and reports.

## Recommended Roadmap

### P0: Make the Existing Workflows Trustworthy

1. Add automated tests for auth, permissions, tenant isolation, inventory, dispensing, and POS.
2. Standardize API errors and align OpenAPI specs with implemented routes.
3. Fix token refresh persistence and add refresh-token revocation.
4. Connect dashboard metrics to real backend data.
5. Add frontend permission-aware navigation and route protection.

### P1: Make Stock and Dispensing Pharmacy-Safe

1. Route every stock-changing action through transactional stock movements.
2. Deduct inventory on dispense and POS sale; restore or reverse correctly on void/refund.
3. Implement FEFO batch selection, available-stock validation, and negative-stock prevention.
4. Implement real allergy and interaction checks, including clinical override documentation.
5. Add pharmacist verification steps and controlled-substance logging.

### P2: Add Compliance and Reporting

1. Add append-only audit events for sensitive actions.
2. Build printable labels, receipts, and controlled-drug reports.
3. Add sales, inventory, expiry, prescription, and margin reports.
4. Add cash-up, daily close reconciliation, and accounting/export workflows.

### P3: Expand Platform Capabilities

1. Add inter-location transfers and chain-level reporting.
2. Add insurance/e-prescribing integrations.
3. Add tenant administration, feature flags, and operational monitoring.
4. Add import/export tooling and migration-safe demo data.

## Development Principles Going Forward

- Keep tenant isolation server-side and test it for every repository.
- Treat OpenAPI as the contract and keep generated code in sync.
- Use transactions for any workflow that writes multiple business records.
- Never mark clinical or compliance behavior as complete while it is stubbed.
- Prefer small, vertical slices: backend contract, migration, repository/service, frontend API, UI, and tests together.
- Make audit and stock movements foundational records, not side effects.
- Keep dashboards and reports derived from real data.
- Avoid expanding surface area until current workflows are correct, testable, and reversible.

## Workflows

### OTC Medicine Purchase (Customer Walk-in)

1. Customer approaches the POS counter or consultation window.
2. Pharmacist or staff confirms the OTC product is appropriate (no prescription needed, no clinical red flags).
3. Staff selects or creates the customer as a patient record (name, date of birth, phone number at minimum; may link to existing profile or create a walk-in guest entry).
4. Staff adds the OTC product(s) to the POS sale — the system checks available stock in real time:
   - If insufficient stock, the staff is notified and can offer an alternative or rain check.
   - If stock is available, the system reserves the quantity (or deducts on completion).
5. Staff applies any relevant discounts (loyalty, promotion, senior, or manual price override with manager approval if needed).
6. Staff completes payment — the system supports cash, card, or other tender types and records the transaction.
7. System generates a receipt (printed or digital) and updates inventory with a deduction from the selected stock batch (FEFO-based if available).
8. Stock movement record is created with type `sale`, linking to the sale item.
9. System records the sale in POS history with status `completed`.
10. Customer receives the product and leaves. If a consultation was provided (e.g., recommendation for a cough medicine), a brief note can optionally be attached to the patient record.

### Products vs Inventory

Products and inventory are separate but related concepts in PharmD. A **product** is a record in the master catalog — drug name, brand, generic name, strength, form, barcode, NDC, manufacturer, and classification (OTC, prescription, controlled). It represents *what* can be stocked, but holds no quantity, location, cost, or expiry information.

**Inventory** is the physical stock-on-hand of a product at a specific location, tracked per batch. Each batch record sits in `stock_batches` with a `product_id` foreign key to `products`, plus `location_id`, `batch_number`, `quantity`, `remaining_qty`, `unit_cost`, `selling_price`, `manufacturing_date`, and `expiry_date`. A single product can have many stock batches across multiple locations — e.g., an organisation has aspirin in the catalog once, but may have three batches at one location and two at another.

The `stock_movements` table provides the audit trail for every quantity change on any batch, recording the type (`receipt`, `sale`, `adjustment`, `count_correction`, `dispense`, `transfer`), the reference document, who performed it, and a timestamp.

**Why cost isn't on the product**: A product's cost varies per purchase — the same ibuprofen can be bought from Supplier A at $0.05/tablet and Supplier B at $0.06/tablet, on different purchase orders, at different times. Cost is a property of each received batch (`stock_batches.unit_cost`), not of the catalog record. Selling price can also vary by batch and by location — there is a separate `product_prices` table keyed by `(product_id, location_id)` for location-specific standard pricing, but the batch-level `selling_price` allows per-lot pricing (e.g., old stock marked down before expiry). The product catalog stays pure: it defines *what* the drug is (name, strength, form, classification, barcode), not *what it costs* or *where it lives*.

**How selling OTC drugs now reduces inventory**: When payments are recorded for a sale (`POSService.RecordPayments`), the service calls `POSRepoImpl.CompleteSale` which records payments, deducts from stock batches (FEFO — earliest expiry first), creates `sale`-type stock movements, and marks the sale `completed` — all in a single transaction. The same pattern applies to dispensing (`DispensingRepoImpl.Create` now deducts within a transaction with `dispense`-type movements). On void or refund of a completed sale, `RestoreSaleInventory` reverses the deductions. If stock is insufficient at payment time, the transaction rolls back and the error propagates to the caller.

**Key distinction**: Product management (CRUD on the catalog) is organisation-wide — create a product once, use it everywhere. Inventory management is per-location and per-batch — you receive, adjust, dispense, and sell physical stock against specific batches. The product catalog API lives at `/products` with `products.*` permissions; the stock API lives at `/inventory` with `inventory.*` permissions. The two repositories and services are completely separate — `ProductRepository` works only with `products` and `generic_substitutions` tables, while `InventoryRepository` works with `stock_batches`, `stock_movements`, and joins to `products` for display.

### Discounts

The Discounts module manages discount rules — configurable policies that define what discounts are available during POS transactions. Each rule has a type (`percentage` or `fixed`), a value, an optional scope (`all`, `category`, or `product`), optional validity dates, a minimum order value, and a maximum discount cap for percentage rules. Rules are reference data — the administrator configures them, but the POS operator applies discounts manually per line item at checkout. The discount rules table uses hard deletes. API at `/pricing/rules` with `discounts:*` permissions.

### Prescribers

A Prescriber is a healthcare professional (doctor, dentist, nurse practitioner) authorised to write prescriptions. The Prescribers module is a directory of these providers — each record holds the prescriber's name, license number, DEA number (for controlled substances), NPI number, phone, email, specialty, and address. Prescribers are a required foreign key in the Prescriptions workflow: every prescription must reference a prescriber, and the prescription list/detail views JOIN to the prescribers table to display the prescriber's name. Uses soft deletes. API at `/prescribers` with `prescribers:*` permissions.
