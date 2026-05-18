import { apiRequest } from "./client";

export interface Prescriber {
  id: string;
  name: string;
  license_number: string;
  phone: string;
  email: string;
  specialty: string;
  dea_number: string;
  npi_number: string;
  address: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface PrescriberListResponse {
  data: Prescriber[];
  page: number;
  limit: number;
  total: number;
}

export function listPrescribers(page = 1, limit = 20, query = "") {
  const q = query ? `&query=${encodeURIComponent(query)}` : "";
  return apiRequest<PrescriberListResponse>(`/prescribers?page=${page}&limit=${limit}${q}`);
}

export function getPrescriber(id: string) {
  return apiRequest<Prescriber>(`/prescribers/${id}`);
}

export function createPrescriber(data: Partial<Prescriber>) {
  return apiRequest<Prescriber>("/prescribers", { method: "POST", body: JSON.stringify(data) });
}

export function updatePrescriber(id: string, data: Partial<Prescriber>) {
  return apiRequest<Prescriber>(`/prescribers/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function deletePrescriber(id: string) {
  return apiRequest<void>(`/prescribers/${id}`, { method: "DELETE" });
}
