import { apiRequest } from "./client";

export interface StockBatch {
  id: string;
  product_id: string;
  product_name: string;
  location_id: string;
  batch_number: string;
  quantity: number;
  remaining_qty: number;
  unit_cost: number;
  selling_price: number;
  manufacturing_date: string;
  expiry_date: string;
  received_date: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface BatchSummary {
  id: string;
  batch_number: string;
  quantity: number;
  unit_cost: number;
  selling_price: number;
  expiry_date: string;
  is_active: boolean;
}

export interface StockItem {
  product_id: string;
  product_name: string;
  brand_name: string;
  generic_name: string;
  classification: string;
  total_quantity: number;
  reorder_level: number;
  batches: BatchSummary[];
}

export interface StockListResponse {
  data: StockItem[];
  page: number;
  limit: number;
  total: number;
}

export interface StockAlert {
  product_id: string;
  product_name: string;
  brand_name: string;
  total_quantity: number;
  reorder_level: number;
  location_id: string;
}

export interface ExpiringBatch {
  id: string;
  product_id: string;
  product_name: string;
  batch_number: string;
  remaining_qty: number;
  expiry_date: string;
  days_until_expiry: number;
}

export interface StockMovement {
  id: string;
  product_id: string;
  product_name: string;
  batch_id: string;
  movement_type: string;
  quantity: number;
  notes: string;
  created_by: string;
  created_at: string;
}

export function listStock(locationId: string, page = 1, limit = 20, query = "") {
  let path = `/inventory?location_id=${locationId}&page=${page}&limit=${limit}`;
  if (query) path += `&query=${encodeURIComponent(query)}`;
  return apiRequest<StockListResponse>(path);
}

export function createBatch(data: {
  product_id: string;
  location_id: string;
  batch_number?: string;
  quantity: number;
  unit_cost?: number;
  selling_price?: number;
  manufacturing_date?: string;
  expiry_date?: string;
}) {
  return apiRequest<StockBatch>("/inventory/batches", { method: "POST", body: JSON.stringify(data) });
}

export function createAdjustment(data: {
  product_id: string;
  location_id: string;
  batch_id?: string;
  quantity: number;
  movement_type: string;
  notes?: string;
}) {
  return apiRequest<StockMovement>("/inventory/adjustments", { method: "POST", body: JSON.stringify(data) });
}

export function listAlerts(locationId = "") {
  let path = "/inventory/alerts";
  if (locationId) path += `?location_id=${locationId}`;
  return apiRequest<StockAlert[]>(path);
}

export function listExpiring(locationId = "", days = 30) {
  let path = `/inventory/expiring?days=${days}`;
  if (locationId) path += `&location_id=${locationId}`;
  return apiRequest<ExpiringBatch[]>(path);
}

export function stockCount(items: {
  product_id: string;
  location_id: string;
  batch_id?: string;
  counted_qty: number;
  notes?: string;
}[]) {
  return apiRequest<{ adjustments: number }>("/inventory/counts", { method: "POST", body: JSON.stringify(items) });
}
