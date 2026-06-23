import { apiRequest } from "./client";

export interface User {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  is_active: boolean;
  location_id?: string;
  created_at: string;
  updated_at: string;
}

export interface UserListResponse {
  data: User[];
  page: number;
  limit: number;
  total: number;
}

export function listUsers(page = 1, limit = 20) {
  return apiRequest<UserListResponse>(`/users?page=${page}&limit=${limit}`);
}

export function getUser(id: string) {
  return apiRequest<User>(`/users/${id}`);
}

export function createUser(data: { first_name: string; last_name: string; email: string; role_id?: string; location_id?: string }) {
  return apiRequest<User>("/users", { method: "POST", body: JSON.stringify(data) });
}

export function updateUser(id: string, data: { first_name?: string; last_name?: string; is_active?: boolean; location_id?: string }) {
  return apiRequest<User>(`/users/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function deleteUser(id: string) {
  return apiRequest<void>(`/users/${id}`, { method: "DELETE" });
}

export function assignUserRole(id: string, roleId: string) {
  return apiRequest<void>(`/users/${id}/roles`, { method: "PUT", body: JSON.stringify({ role_id: roleId }) });
}
