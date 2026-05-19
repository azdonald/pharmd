import { apiRequest } from "./client";

export interface Sale {
  id: string;
  location_id: string;
  patient_id: string;
  patient_name: string;
  prescription_id: string;
  sale_type: string;
  status: string;
  subtotal: number;
  tax_total: number;
  discount_total: number;
  grand_total: number;
  paid_amount: number;
  change_amount: number;
  notes: string;
  created_by: string;
  voided_by: string;
  voided_at: string;
  created_at: string;
  updated_at: string;
  items?: SaleItem[];
  payments?: Payment[];
}

export interface SaleItem {
  id: string;
  product_id: string;
  product_name: string;
  quantity: number;
  unit_price: number;
  discount: number;
  line_total: number;
}

export interface Payment {
  id: string;
  sale_id: string;
  method: string;
  amount: number;
  reference: string;
}

export interface SaleListResponse {
  data: Sale[];
  page: number;
  limit: number;
  total: number;
}

export function listSales(page = 1, limit = 20, status = "") {
  const params = new URLSearchParams({ page: String(page), limit: String(limit) });
  if (status) params.set("status", status);
  return apiRequest<SaleListResponse>(`/pos/sales?${params}`);
}

export function getSale(id: string) {
  return apiRequest<Sale>(`/pos/sales/${id}`);
}

export function createSale(data: { location_id: string; patient_id?: string; prescription_id?: string; sale_type?: string; notes?: string; items: { product_id: string; quantity: number; unit_price: number; discount?: number }[] }) {
  return apiRequest<Sale>("/pos/sales", { method: "POST", body: JSON.stringify(data) });
}

export function updateSale(id: string, data: { status: "held" | "voided" | "refunded"; notes?: string }) {
  return apiRequest<Sale>(`/pos/sales/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function recordPayment(saleId: string, payments: { method: string; amount: number; reference?: string }[]) {
  return apiRequest<Sale>("/pos/payments", { method: "POST", body: JSON.stringify({ sale_id: saleId, payments }) });
}

export function getReceipt(id: string) {
  return apiRequest<any>(`/pos/sales/${id}/receipt`);
}

export function getDailySummary(date: string, locationId: string) {
  return apiRequest<any>(`/pos/summary?date=${encodeURIComponent(date)}&location_id=${encodeURIComponent(locationId)}`);
}

export function closeDay(date: string, locationId: string, notes?: string) {
  return apiRequest<any>("/pos/close-day", { method: "POST", body: JSON.stringify({ date, location_id: locationId, notes }) });
}
