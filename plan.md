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

```
backend/api/<resource>/
  <resource>.yaml       — hand-written OpenAPI 3.0.3 spec
  gen.go                — //go:generate directives
  server_gen.go         — generated (chi-server interface + routing)
  server_impl.go        — hand-written business logic
  types_gen.go          — generated (request/response structs)
```

1. Write the OpenAPI spec → 2. `go generate ./...` → 3. Hand-write `server_impl.go`

Shared models live in `common/common.yaml`.

---

## Implementation Rule
Every feature must be implemented **backend-first, then frontend** before moving to the next feature. The plan pairs each backend task with its corresponding frontend task.

---

# Features & Iterations

Each feature lists tasks in order:
1. OpenAPI spec & generated code
2. Backend implementation (repo → service → handler)
3. Migration (if applicable)
4. Frontend implementation (API client → pages → components)

---

## Iteration 1 — Auth & Tenant Foundation (P0)

### F1. Organisation Registration & Login

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 1.1 | Register endpoint | `POST /v1/register` | `src/pages/Register.tsx` |
| 1.2 | Login endpoint | `POST /v1/login` | `src/pages/Login.tsx` |
| 1.3 | Token refresh | `POST /v1/refresh` | `src/api/client.ts` (auto-refresh interceptor) |
| 1.4 | Logout | `POST /v1/logout` | Logout button + token clear |
| 1.5 | Change password | `PUT /v1/change-password` | `src/pages/ChangePassword.tsx` |
| 1.6 | JWT auth middleware | `middleware/auth.go` | `src/api/auth.ts` (send token header) |
| 1.7 | Frontend auth context | — | `src/context/AuthContext.tsx` |

**Status:** Backend done ✅ | Frontend ✅

**Backend files:** `api/auth/auth.yaml`, `api/auth/server_impl.go`, `repository/auth_repo_impl.go`, `service/auth_service_impl.go`, `middleware/auth.go`, `utils/token.go`

**Frontend files:** `src/api/auth.ts`, `src/api/client.ts`, `src/context/AuthContext.tsx`, `src/pages/Login.tsx`, `src/pages/Register.tsx`, `src/pages/ChangePassword.tsx`

---

### F2. Roles & Permissions (RBAC)

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 2.1 | Seed permission slugs | `migrations/1767000004_seed_permissions.go` | — |
| 2.2 | Seed default roles (Admin, PIC, Pharmacist, Tech, Cashier) | `migrations/1767000005_seed_default_roles.go` | — |
| 2.3 | List permissions | `GET /permissions` | `src/api/permissions.ts` |
| 2.4 | Role CRUD | `api/roles/` (GET/POST /roles, GET/PUT/DELETE /roles/{id}) | `src/pages/Roles.tsx`, `src/pages/RoleForm.tsx` |
| 2.5 | Assign permissions to role | `PUT /roles/{id}/permissions` | Role edit → permission checkboxes |
| 2.6 | Permission middleware | `middleware/auth.go` (`RequirePermission`) | `src/components/ProtectedRoute.tsx` |

**Status:** Backend done ✅ | Frontend ✅

---

### F3. Location Management

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 3.1 | List locations | `GET /locations` | `src/pages/Locations.tsx` |
| 3.2 | Create location | `POST /locations` | `src/pages/LocationForm.tsx` |
| 3.3 | Get location | `GET /locations/{id}` | — |
| 3.4 | Update location | `PUT /locations/{id}` | Location edit form |
| 3.5 | Delete location | `DELETE /locations/{id}` | Delete with confirmation |
| 3.6 | Auto-create default location on registration | in `service/auth_service_impl.go` | — |

**Status:** Backend done ✅ | Frontend ✅

---

### F4. User Management

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 4.1 | List users | `GET /users` | `src/pages/Users.tsx` |
| 4.2 | Create user (invite with temp password) | `POST /users` | `src/pages/UserForm.tsx` |
| 4.3 | Get user | `GET /users/{id}` | — |
| 4.4 | Update user | `PUT /users/{id}` | Edit user form |
| 4.5 | Delete / deactivate user | `DELETE /users/{id}` | Toggle active |
| 4.6 | Assign role to user | `PUT /users/{id}/roles` | Role dropdown on user form |

**Status:** Backend done ✅ | Frontend ✅

---

## Iteration 2 — Core Pharmacy (P1)

### F5. Patient / Customer Management

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 5.1 | Patient registration | `POST /patients` | `src/pages/PatientForm.tsx` |
| 5.2 | Patient search | `GET /patients?query=` | `src/pages/Patients.tsx` (searchable list) |
| 5.3 | Patient profile | `GET /patients/{id}` | `src/pages/PatientDetail.tsx` |
| 5.4 | Update patient | `PUT /patients/{id}` | Edit form |
| 5.5 | Allergies & conditions | `POST /patients/{id}/allergies` | Condition list editor |
| 5.6 | Insurance linkage | `POST /patients/{id}/insurance` | Insurance form |

**Status:** Backend done ✅ | Frontend ✅

**Backend files:** `api/patients/`, `repository/patient_repo_impl.go`, `service/patient_service_impl.go`, `migrations/1768000000_create_patients_table.go`, `migrations/1768000001_create_patient_allergies_table.go`, `migrations/1768000002_create_patient_conditions_table.go`

**Frontend files:** `src/api/patients.ts`, `src/pages/Patients.tsx`, `src/pages/PatientForm.tsx`, `src/pages/PatientDetail.tsx`

**DB Tables:** `patients`, `patient_allergies`, `patient_conditions`, `patient_insurance`

---

### F6. Drug / Product Catalog

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 6.1 | Product CRUD | `api/products/` | `src/pages/Products.tsx`, `src/pages/ProductForm.tsx` |
| 6.2 | Drug classification | enum field on product | Select dropdown |
| 6.3 | Barcode / NDC tracking | `POST /products/barcode-lookup` | Barcode input |
| 6.4 | Product categories | `api/categories/` | Category tree |
| 6.5 | Generic substitution | `POST /products/{id}/substitutes` | Substitution link editor |

**Frontend:** `src/api/products.ts`, `src/pages/Products.tsx`, `src/pages/ProductForm.tsx`

**DB Tables:** `products`, `product_categories`, `generic_substitutions`

---

### F7. Inventory & Stock Management

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 7.1 | Stock batch entry | `POST /inventory/batches` | Stock receipt form |
| 7.2 | Stock adjustment | `POST /inventory/adjustments` | Adjustment form (waste, damage) |
| 7.3 | Stock list per location | `GET /inventory?location_id=` | Inventory table |
| 7.4 | Low-stock alerts | `GET /inventory/alerts` | Alert badges |
| 7.5 | Expiry tracking | `GET /inventory/expiring?days=30` | Expiry report |
| 7.6 | Stock count / reconciliation | `POST /inventory/counts` | Count sheet |

**Frontend:** `src/api/inventory.ts`, `src/pages/Inventory.tsx`, `src/pages/StockReceipt.tsx`

**DB Tables:** `stock_batches`, `stock_movements`, `inventory_alerts`

---

### F8. Supplier Management

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 8.1 | Supplier CRUD | `api/suppliers/` | `src/pages/Suppliers.tsx` |
| 8.2 | Supplier price list | `PUT /suppliers/{id}/prices` | Price list editor |

**Frontend:** `src/api/suppliers.ts`, `src/pages/Suppliers.tsx`

**DB Tables:** `suppliers`, `supplier_products`

---

### F9. Purchase Orders

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 9.1 | Create PO | `POST /purchases` | PO creation form |
| 9.2 | PO list | `GET /purchases` | PO table |
| 9.3 | PO approval | `PUT /purchases/{id}/approve` | Approve / reject buttons |
| 9.4 | Goods Received Note | `POST /purchases/{id}/receive` | GRN form |
| 9.5 | PO status tracking | status field | Status badges |

**Frontend:** `src/api/purchases.ts`, `src/pages/PurchaseOrders.tsx`, `src/pages/POForm.tsx`

**DB Tables:** `purchase_orders`, `purchase_order_items`, `goods_received_notes`

---

### F10. Pricing & Tax

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 10.1 | Product price per location | `api/pricing/` | Price list editor |
| 10.2 | Pricing formulas | config per location | Formula selector |
| 10.3 | Tax rate per location | on location model | Tax field on location form |
| 10.4 | Discount rules | `api/pricing/discounts` | Discount config |

**Frontend:** `src/api/pricing.ts`, `src/pages/Pricing.tsx`

**DB Tables:** `product_prices`, `tax_rates`, `discount_rules`

---

## Iteration 3 — Prescription & POS (P1)

### F11. Prescription Management

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 11.1 | Prescription intake | `POST /prescriptions` | Rx intake form |
| 11.2 | Prescription list | `GET /prescriptions` | Rx queue |
| 11.3 | Prescription detail | `GET /prescriptions/{id}` | Rx detail view |
| 11.4 | Prescriber management | `api/prescribers/` | Prescriber lookup |
| 11.5 | Refill tracking | `POST /prescriptions/{id}/refill` | Refill button |
| 11.6 | Rx status workflow | status field | Status stepper |

**Frontend:** `src/api/prescriptions.ts`, `src/pages/Prescriptions.tsx`, `src/pages/RxForm.tsx`

**DB Tables:** `prescriptions`, `prescription_items`, `prescribers`, `refills`

---

### F12. Dispensing Workflow

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 12.1 | Dispensing queue | `GET /dispensing?status=pending` | Queue view |
| 12.2 | Dispense medication | `POST /dispensing` | Dispense form |
| 12.3 | Drug interaction check | `GET /drugs/{id}/interactions` | Warning modal |
| 12.4 | Allergy check | auto on dispense | Alert banner |
| 12.5 | Partial dispensing | qty field on dispense | Partial qty input |
| 12.6 | Controlled substance log | additional fields on dispense | Witness signature field |
| 12.7 | Label printing | `GET /dispensing/{id}/label` | Print button |

**Frontend:** `src/api/dispensing.ts`, `src/pages/DispensingQueue.tsx`, `src/pages/DispenseForm.tsx`

**DB Tables:** `dispensing_records`, `dispensing_logs`, `controlled_substance_logs`

---

### F13. Point of Sale

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 13.1 | OTC sale | `POST /pos/sales` | POS interface |
| 13.2 | Prescription sale | link to dispensed record | POS → Rx lookup |
| 13.3 | Multiple payment methods | `POST /pos/payments` | Payment split UI |
| 13.4 | Receipt printing | `GET /pos/sales/{id}/receipt` | Print receipt |
| 13.5 | Hold / void / refund | `PUT /pos/sales/{id}/void` | Void/refund buttons |
| 13.6 | Daily summary (X-report) | `GET /pos/summary` | Summary modal |
| 13.7 | End-of-day closeout (Z-report) | `POST /pos/close-day` | Closeout button |

**Frontend:** `src/api/pos.ts`, `src/pages/POS.tsx`, `src/components/PaymenModal.tsx`

**DB Tables:** `sales`, `sale_items`, `payments`, `refunds`, `daily_summaries`

---

## Iteration 4 — Financial & Reporting (P2)

### F14. Accounts & Financials

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 14.1 | Daily cash-up | `POST /accounts/cashup` | Cash-up form |
| 14.2 | AR (credit patients) | `GET /accounts/receivable` | AR aging table |
| 14.3 | AP (supplier invoices) | `GET /accounts/payable` | AP table |
| 14.4 | Expenses | `POST /accounts/expenses` | Expense form |
| 14.5 | Profit margin report | `GET /reports/profit-margin` | Profit chart |

**DB Tables:** `cashups`, `accounts_receivable`, `accounts_payable`, `expenses`

### F15. Reports & Analytics

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 15.1 | Sales report | `GET /reports/sales` | Filterable report table |
| 15.2 | Inventory report | `GET /reports/inventory` | Report table |
| 15.3 | Prescription report | `GET /reports/prescriptions` | Report table |
| 15.4 | Expiry report | `GET /reports/expiry` | Report table |
| 15.5 | Dashboard widgets | `GET /dashboard` | Dashboard page |

### F16. Audit & Compliance

| # | Task | Backend | Frontend |
|---|------|---------|----------|
| 16.1 | Activity log | `GET /audit` | Audit log table |
| 16.2 | Controlled substance log | `GET /audit/controlled` | Regulatory log |
| 16.3 | Data export | `GET /export` | Export button |

---

## Iteration 5 — Multi-Location & Platform (P2/P3)

### F17. Multi-Location Chain Features

| # | Backend | Frontend |
|---|---------|----------|
| 17.1 | Centralized purchasing (`POST /chain/purchase-orders`) | Chain PO form |
| 17.2 | Inter-location transfers (`POST /chain/transfers`) | Transfer request form |
| 17.3 | Consolidated reports (`GET /chain/reports`) | Chain dashboard |
| 17.4 | Cross-location patient lookup (`GET /chain/patients?query=`) | Global search |

### F18. Platform Administration

| # | Backend | Frontend |
|---|---------|----------|
| 18.1 | Super admin tenant management | Admin panel |
| 18.2 | Feature flags per tenant | Toggle UI |

### F19. Integrations

SMS, email, insurance claims, NDC lookup, accounting export.

---

## Implementation Roadmap

```
Iteration 1: Auth & Foundation     → F1, F2, F3, F4     (P0) ✅ DONE
Iteration 2: Core Pharmacy          → F5, F6, F7, F8, F9, F10 (P1) ← CURRENT
```

---

## Progress

| Feature | Backend | Frontend |
|---------|---------|----------|
| F1. Auth & Registration | ✅ done | ✅ done |
| F2. Roles & Permissions | ✅ done | ✅ done |
| F3. Location Management | ✅ done | ✅ done |
| F4. User Management | ✅ done | ✅ done |
| F5. Patients | ✅ done | ✅ done |
| F6. Products | ❌ | ❌ |
| F7. Inventory | ❌ | ❌ |
| F8. Suppliers | ❌ | ❌ |
| F9. Purchase Orders | ❌ | ❌ |
| F10. Pricing & Tax | ❌ | ❌ |
| F11. Prescriptions | ❌ | ❌ |
| F12. Dispensing | ❌ | ❌ |
| F13. POS | ❌ | ❌ |
| F14. Financials | ❌ | ❌ |
| F15. Reports | ❌ | ❌ |
| F16. Audit | ❌ | ❌ |
| F17. Multi-Location Chain | ❌ | ❌ |
| F18. Platform Admin | ❌ | ❌ |
| F19. Integrations | ❌ | ❌ |

---

## Design Principles
- **Tenant isolation first** — every query includes `organisation_id`; never trust client-side.
- **Location-aware from day one** — inventory, pricing, user access scoped to location.
- **Audit everything** — pharmacy is regulated; every stock movement, dispense, and action logged.
- **Regulatory ready** — HIPAA/GDPR, controlled substance tracking, prescription retention laws.
