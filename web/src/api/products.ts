import { apiRequest } from "./client";

export interface Product {
  id: string;
  name: string;
  description: string;
  category_id: string;
  classification: string;
  brand_name: string;
  generic_name: string;
  manufacturer: string;
  barcode: string;
  ndc: string;
  unit_of_measure: string;
  strength: string;
  form: string;
  reorder_level: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ProductListResponse {
  data: Product[];
  page: number;
  limit: number;
  total: number;
}

export interface ProductCategory {
  id: string;
  name: string;
  description: string;
  parent_id: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface GenericSubstitution {
  id: string;
  product_id: string;
  substitute_product_id: string;
  substitute_name: string;
  notes: string;
  created_at: string;
}

export function listProducts(page = 1, limit = 20, query = "", categoryId = "") {
  let path = `/products?page=${page}&limit=${limit}`;
  if (query) path += `&query=${encodeURIComponent(query)}`;
  if (categoryId) path += `&category_id=${categoryId}`;
  return apiRequest<ProductListResponse>(path);
}

export function getProduct(id: string) {
  return apiRequest<Product>(`/products/${id}`);
}

export function createProduct(data: Partial<Product>) {
  return apiRequest<Product>("/products", { method: "POST", body: JSON.stringify(data) });
}

export function updateProduct(id: string, data: Partial<Product>) {
  return apiRequest<Product>(`/products/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function deleteProduct(id: string) {
  return apiRequest<void>(`/products/${id}`, { method: "DELETE" });
}

export function barcodeLookup(barcode: string) {
  return apiRequest<Product>("/products/barcode-lookup", { method: "POST", body: JSON.stringify({ barcode }) });
}

export function listCategories() {
  return apiRequest<ProductCategory[]>("/product-categories");
}

export function getCategory(id: string) {
  return apiRequest<ProductCategory>(`/product-categories/${id}`);
}

export function createCategory(data: { name: string; description?: string; parent_id?: string }) {
  return apiRequest<ProductCategory>("/product-categories", { method: "POST", body: JSON.stringify(data) });
}

export function updateCategory(id: string, data: Partial<ProductCategory>) {
  return apiRequest<ProductCategory>(`/product-categories/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function deleteCategory(id: string) {
  return apiRequest<void>(`/product-categories/${id}`, { method: "DELETE" });
}

export function listSubstitutes(id: string) {
  return apiRequest<GenericSubstitution[]>(`/products/${id}/substitutes`);
}

export function addSubstitute(id: string, substituteProductId: string, notes?: string) {
  return apiRequest<GenericSubstitution>(`/products/${id}/substitutes`, {
    method: "POST",
    body: JSON.stringify({ substitute_product_id: substituteProductId, notes }),
  });
}

export function removeSubstitute(productId: string, substituteId: string) {
  return apiRequest<void>(`/products/${productId}/substitutes?substitute_id=${substituteId}`, { method: "DELETE" });
}
