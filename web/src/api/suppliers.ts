import { apiRequest } from "./client";

export interface Supplier {
  id: string;
  name: string;
  contact_person: string;
  phone: string;
  email: string;
  address: string;
  city: string;
  state: string;
  country: string;
  payment_terms: string;
  notes: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface SupplierListResponse {
  data: Supplier[];
  page: number;
  limit: number;
  total: number;
}

export interface SupplierProduct {
  id: string;
  supplier_id: string;
  product_id: string;
  product_name: string;
  unit_price: number;
  min_order_qty: number;
  lead_time_days: number;
  notes: string;
  created_at: string;
  updated_at: string;
}

export function listSuppliers(page = 1, limit = 20, query = "") {
  const q = query ? `&query=${encodeURIComponent(query)}` : "";
  return apiRequest<SupplierListResponse>(`/suppliers?page=${page}&limit=${limit}${q}`);
}

export function getSupplier(id: string) {
  return apiRequest<Supplier>(`/suppliers/${id}`);
}

export function createSupplier(data: Partial<Supplier>) {
  return apiRequest<Supplier>("/suppliers", { method: "POST", body: JSON.stringify(data) });
}

export function updateSupplier(id: string, data: Partial<Supplier>) {
  return apiRequest<Supplier>(`/suppliers/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function deleteSupplier(id: string) {
  return apiRequest<void>(`/suppliers/${id}`, { method: "DELETE" });
}

export function listSupplierProducts(id: string) {
  return apiRequest<SupplierProduct[]>(`/suppliers/${id}/products`);
}

export function setSupplierProducts(id: string, products: { product_id: string; unit_price: number; min_order_qty?: number; lead_time_days?: number; notes?: string }[]) {
  return apiRequest<SupplierProduct[]>(`/suppliers/${id}/products`, { method: "PUT", body: JSON.stringify(products) });
}
