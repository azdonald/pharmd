import { apiRequest } from "./client";

export interface ProductPrice {
  id: string;
  product_id: string;
  product_name: string;
  location_id: string;
  location_name: string;
  selling_price: number;
  cost_price: number;
  min_price: number;
  max_discount: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface PriceListResponse {
  data: ProductPrice[];
  page: number;
  limit: number;
  total: number;
}

export interface DiscountRule {
  id: string;
  name: string;
  type: "percentage" | "fixed";
  value: number;
  min_order_value: number;
  max_discount_amount: number;
  applies_to: "all" | "category" | "product";
  applies_to_id: string;
  is_active: boolean;
  start_date: string;
  end_date: string;
  created_at: string;
  updated_at: string;
}

export interface DiscountRuleListResponse {
  data: DiscountRule[];
  page: number;
  limit: number;
  total: number;
}

export function listPrices(productId = "", locationId = "", page = 1, limit = 20) {
  const params = new URLSearchParams({ page: String(page), limit: String(limit) });
  if (productId) params.set("product_id", productId);
  if (locationId) params.set("location_id", locationId);
  return apiRequest<PriceListResponse>(`/pricing?${params}`);
}

export function upsertPrice(data: { product_id: string; location_id: string; selling_price: number; cost_price?: number; min_price?: number; max_discount?: number }) {
  return apiRequest<ProductPrice>("/pricing", { method: "POST", body: JSON.stringify(data) });
}

export function deletePrice(id: string) {
  return apiRequest<void>(`/pricing/${id}`, { method: "DELETE" });
}

export function listDiscountRules(page = 1, limit = 20) {
  return apiRequest<DiscountRuleListResponse>(`/pricing/rules?page=${page}&limit=${limit}`);
}

export function getDiscountRule(id: string) {
  return apiRequest<DiscountRule>(`/pricing/rules/${id}`);
}

export function createDiscountRule(data: Partial<DiscountRule>) {
  return apiRequest<DiscountRule>("/pricing/rules", { method: "POST", body: JSON.stringify(data) });
}

export function updateDiscountRule(id: string, data: Partial<DiscountRule>) {
  return apiRequest<DiscountRule>(`/pricing/rules/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function deleteDiscountRule(id: string) {
  return apiRequest<void>(`/pricing/rules/${id}`, { method: "DELETE" });
}
