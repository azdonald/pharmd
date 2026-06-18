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

### 2. Dispensing does not appear to reduce inventory

Dispense records are created, but the dispensing service does not deduct stock batches, enforce available quantity, or select batches by expiry/FEFO. POS sale creation also records sale items but does not appear to reduce inventory.

Recommendation: Make inventory movements the single source of truth for every stock-affecting workflow. Dispensing, POS, receiving, adjustments, voids, refunds, and counts should all create stock movements inside database transactions.

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

## Current Assessment

PharmD is best described today as a broad alpha/MVP foundation for a pharmacy management system. It has many of the right modules and a sensible architecture, but the highest-risk pharmacy behaviors need hardening before it should be treated as production-ready. The most valuable next step is not adding more screens; it is making auth, permissions, inventory, dispensing, POS, audit, and clinical safety dependable end to end.
