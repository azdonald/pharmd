import { apiRequest } from "./client";

export interface Role {
  id: string;
  name: string;
  slug: string;
  description: string;
  is_system: boolean;
  created_at: string;
  updated_at: string;
}

export interface RoleListResponse {
  data: Role[];
  page: number;
  limit: number;
  total: number;
}

export function listRoles(page = 1, limit = 20) {
  return apiRequest<RoleListResponse>(`/roles?page=${page}&limit=${limit}`);
}

export function getRole(id: string) {
  return apiRequest<Role>(`/roles/${id}`);
}

export function createRole(data: { name: string; description?: string; permission_ids?: string[] }) {
  return apiRequest<Role>("/roles", { method: "POST", body: JSON.stringify(data) });
}

export function updateRole(id: string, data: { name?: string; description?: string }) {
  return apiRequest<Role>(`/roles/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export function deleteRole(id: string) {
  return apiRequest<void>(`/roles/${id}`, { method: "DELETE" });
}

export function getRolePermissions(id: string) {
  return apiRequest<{ permission_ids: string[] }>(`/roles/${id}/permissions`);
}

export function setRolePermissions(id: string, permissionIds: string[]) {
  return apiRequest<void>(`/roles/${id}/permissions`, {
    method: "PUT",
    body: JSON.stringify({ permission_ids: permissionIds }),
  });
}
