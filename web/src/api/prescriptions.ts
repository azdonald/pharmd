import { apiRequest } from "./client";

export interface PrescriptionItem {
  id: string;
  product_id: string;
  product_name: string;
  dosage: string;
  frequency: string;
  duration: string;
  quantity: number;
  refills_authorized: number;
  refills_used: number;
  notes: string;
}

export interface Prescription {
  id: string;
  patient_id: string;
  patient_name: string;
  prescriber_id: string;
  prescriber_name: string;
  location_id: string;
  status: string;
  diagnosis: string;
  notes: string;
  issued_date: string;
  expiry_date: string;
  created_by: string;
  created_at: string;
  updated_at: string;
  items?: PrescriptionItem[];
}

export interface PrescriptionListResponse {
  data: Prescription[];
  page: number;
  limit: number;
  total: number;
}

export function listPrescriptions(page = 1, limit = 20, status = "", patientId = "") {
  const params = new URLSearchParams({ page: String(page), limit: String(limit) });
  if (status) params.set("status", status);
  if (patientId) params.set("patient_id", patientId);
  return apiRequest<PrescriptionListResponse>(`/prescriptions?${params}`);
}

export function getPrescription(id: string) {
  return apiRequest<Prescription>(`/prescriptions/${id}`);
}

export function createPrescription(data: {
  patient_id: string;
  prescriber_id: string;
  location_id: string;
  diagnosis?: string;
  notes?: string;
  issued_date?: string;
  expiry_date?: string;
  items: { product_id: string; dosage: string; frequency: string; duration?: string; quantity: number; refills_authorized?: number; notes?: string }[];
}) {
  return apiRequest<Prescription>("/prescriptions", { method: "POST", body: JSON.stringify(data) });
}

export function updatePrescription(id: string, data: { diagnosis?: string; notes?: string; status?: string; expiry_date?: string }) {
  return apiRequest<Prescription>(`/prescriptions/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function deletePrescription(id: string) {
  return apiRequest<void>(`/prescriptions/${id}`, { method: "DELETE" });
}

export function recordRefill(id: string, itemId: string) {
  return apiRequest<Prescription>(`/prescriptions/${id}/refill`, { method: "POST", body: JSON.stringify({ item_id: itemId }) });
}
