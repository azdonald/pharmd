import { apiRequest } from "./client";

export interface Location {
  id: string;
  name: string;
  address: string;
  city: string;
  state: string;
  country: string;
  phone: string;
  email: string;
  tax_rate: number;
  timezone: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface LocationListResponse {
  data: Location[];
  page: number;
  limit: number;
  total: number;
}

export function listLocations(page = 1, limit = 20) {
  return apiRequest<LocationListResponse>(`/locations?page=${page}&limit=${limit}`);
}

export function getLocation(id: string) {
  return apiRequest<Location>(`/locations/${id}`);
}

export function createLocation(data: Partial<Location>) {
  return apiRequest<Location>("/locations", { method: "POST", body: JSON.stringify(data) });
}

export function updateLocation(id: string, data: Partial<Location>) {
  return apiRequest<Location>(`/locations/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function deleteLocation(id: string) {
  return apiRequest<void>(`/locations/${id}`, { method: "DELETE" });
}
