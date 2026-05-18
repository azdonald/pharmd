import { apiRequest } from "./client";

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    organisation_name: string;
  };
}

export function register(
  email: string,
  password: string,
  organisationName: string,
  firstName: string,
  lastName: string
) {
  return apiRequest<AuthResponse>("/v1/register", {
    method: "POST",
    body: JSON.stringify({
      email,
      password,
      organisation_name: organisationName,
      first_name: firstName,
      last_name: lastName,
    }),
  });
}

export function login(email: string, password: string) {
  return apiRequest<AuthResponse>("/v1/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export function logout() {
  return apiRequest<void>("/v1/logout", { method: "POST" });
}

export function changePassword(oldPassword: string, newPassword: string) {
  return apiRequest<void>("/v1/change-password", {
    method: "PUT",
    body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
  });
}
