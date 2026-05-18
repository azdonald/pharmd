import { apiRequest } from "./client";

export interface DispenseRecord {
  id: string;
  prescription_id: string;
  prescription_item_id: string;
  patient_id: string;
  patient_name: string;
  product_id: string;
  product_name: string;
  location_id: string;
  quantity_dispensed: number;
  quantity_prescribed: number;
  pharmacist_id: string;
  pharmacist_name: string;
  technician_id: string;
  status: string;
  notes: string;
  witness_name: string;
  is_controlled: boolean;
  dispensed_at: string;
  created_at: string;
  updated_at: string;
}

export interface DispenseListResponse {
  data: DispenseRecord[];
  page: number;
  limit: number;
  total: number;
}

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

export function listDispensing(page = 1, limit = 20, status = "", prescriptionId = "") {
  const params = new URLSearchParams({ page: String(page), limit: String(limit) });
  if (status) params.set("status", status);
  if (prescriptionId) params.set("prescription_id", prescriptionId);
  return apiRequest<DispenseListResponse>(`/dispensing?${params}`);
}

export function getDispenseRecord(id: string) {
  return apiRequest<DispenseRecord>(`/dispensing/${id}`);
}

export function createDispense(data: {
  prescription_item_id: string;
  quantity_dispensed: number;
  pharmacist_id: string;
  technician_id?: string;
  notes?: string;
  witness_name?: string;
  is_controlled?: boolean;
}) {
  return apiRequest<DispenseRecord>("/dispensing", { method: "POST", body: JSON.stringify(data) });
}

export function updateDispense(id: string, data: { technician_id?: string; witness_name?: string; notes?: string }) {
  return apiRequest<DispenseRecord>(`/dispensing/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function updateDispenseStatus(id: string, status: string) {
  return apiRequest<DispenseRecord>(`/dispensing/${id}/status`, { method: "PUT", body: JSON.stringify({ status }) });
}

export function getLabelData(id: string) {
  return apiRequest<any>(`/dispensing/${id}/label`);
}

export function checkInteractions(productId: string, patientId: string) {
  return apiRequest<{ severity: string; message: string; interacting_product: string }[]>(`/drugs/${productId}/interactions?patient_id=${patientId}`);
}

export function checkAllergies(productId: string, patientId: string) {
  return apiRequest<{ severity: string; allergen: string; reaction: string }[]>(`/drugs/${productId}/allergy-check?patient_id=${patientId}`);
}
