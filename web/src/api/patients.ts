import { apiRequest } from "./client";

export interface Patient {
  id: string;
  first_name: string;
  last_name: string;
  date_of_birth: string;
  gender: string;
  phone: string;
  email: string;
  address: string;
  city: string;
  state: string;
  country: string;
  blood_group: string;
  genotype: string;
  notes: string;
  emergency_contact_name: string;
  emergency_contact_phone: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface PatientListResponse {
  data: Patient[];
  page: number;
  limit: number;
  total: number;
}

export interface PatientAllergy {
  id: string;
  patient_id: string;
  allergy: string;
  severity: string;
  notes: string;
  created_at: string;
}

export interface PatientCondition {
  id: string;
  patient_id: string;
  condition: string;
  notes: string;
  created_at: string;
}

export function listPatients(page = 1, limit = 20, query = "") {
  const q = query ? `&query=${encodeURIComponent(query)}` : "";
  return apiRequest<PatientListResponse>(`/patients?page=${page}&limit=${limit}${q}`);
}

export function getPatient(id: string) {
  return apiRequest<Patient>(`/patients/${id}`);
}

export function createPatient(data: Partial<Patient>) {
  return apiRequest<Patient>("/patients", { method: "POST", body: JSON.stringify(data) });
}

export function updatePatient(id: string, data: Partial<Patient>) {
  return apiRequest<Patient>(`/patients/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function deletePatient(id: string) {
  return apiRequest<void>(`/patients/${id}`, { method: "DELETE" });
}

export function listPatientAllergies(id: string) {
  return apiRequest<PatientAllergy[]>(`/patients/${id}/allergies`);
}

export function addPatientAllergy(id: string, data: { allergy: string; severity?: string; notes?: string }) {
  return apiRequest<PatientAllergy>(`/patients/${id}/allergies`, { method: "POST", body: JSON.stringify(data) });
}

export function removePatientAllergy(patientId: string, allergyId: string) {
  return apiRequest<void>(`/patients/${patientId}/allergies?allergy_id=${allergyId}`, { method: "DELETE" });
}

export function listPatientConditions(id: string) {
  return apiRequest<PatientCondition[]>(`/patients/${id}/conditions`);
}

export function addPatientCondition(id: string, data: { condition: string; notes?: string }) {
  return apiRequest<PatientCondition>(`/patients/${id}/conditions`, { method: "POST", body: JSON.stringify(data) });
}
