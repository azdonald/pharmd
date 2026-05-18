import { apiRequest } from "./client";

export interface Permission {
  id: string;
  name: string;
  slug: string;
  description: string;
}

export function listPermissions() {
  return apiRequest<{ data: Permission[] }>("/permissions");
}
