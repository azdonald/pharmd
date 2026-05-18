import { apiRequest } from "./client";

export interface PurchaseOrder {
  id: string;
  po_number: string;
  supplier_id: string;
  supplier_name?: string;
  location_id: string;
  status: string;
  order_date: string;
  expected_date: string;
  notes: string;
  subtotal: number;
  tax_total: number;
  grand_total: number;
  created_by: string;
  approved_by?: string;
  approved_at?: string;
  created_at: string;
  updated_at: string;
  items?: PurchaseOrderItem[];
}

export interface PurchaseOrderItem {
  id: string;
  product_id: string;
  product_name?: string;
  quantity_ordered: number;
  quantity_received: number;
  unit_cost: number;
  line_total: number;
}

export interface POListResponse {
  data: PurchaseOrder[];
  page: number;
  limit: number;
  total: number;
}

export interface POItemInput {
  product_id: string;
  quantity_ordered: number;
  unit_cost: number;
}

export function listPurchaseOrders(page = 1, limit = 20, status = "") {
  const s = status ? `&status=${encodeURIComponent(status)}` : "";
  return apiRequest<POListResponse>(`/purchases?page=${page}&limit=${limit}${s}`);
}

export function getPurchaseOrder(id: string) {
  return apiRequest<PurchaseOrder>(`/purchases/${id}`);
}

export function createPurchaseOrder(data: { supplier_id: string; location_id: string; expected_date?: string; notes?: string; items: POItemInput[] }) {
  return apiRequest<PurchaseOrder>("/purchases", { method: "POST", body: JSON.stringify(data) });
}

export function approvePurchaseOrder(id: string) {
  return apiRequest<PurchaseOrder>(`/purchases/${id}/approve`, { method: "PUT" });
}

export function rejectPurchaseOrder(id: string) {
  return apiRequest<PurchaseOrder>(`/purchases/${id}/reject`, { method: "PUT" });
}

export function receiveGoods(id: string, data: { items: { item_id: string; quantity_received: number }[]; notes?: string }) {
  return apiRequest<PurchaseOrder>(`/purchases/${id}/receive`, { method: "POST", body: JSON.stringify(data) });
}
