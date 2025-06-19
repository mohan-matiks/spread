import { apiRequest } from "../client";

export interface AuthKey {
  id: string;
  name: string;
  key: string;
  isValid: boolean;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateAuthKeyRequest {
  name: string;
}

export interface CreateAuthKeyResponse {
  data: string; // The generated auth key string
}

export const authKeyService = {
  async getAllAuthKeys(): Promise<AuthKey[]> {
    const response = await apiRequest<AuthKey[]>({
      method: "GET",
      url: "/core/auth-keys",
    });

    if (response.success && response.data) {
      return response.data;
    }
    throw new Error(response.error || "Failed to fetch auth keys");
  },

  async createAuthKey(name: string): Promise<string> {
    const response = await apiRequest<string>({
      method: "POST",
      url: "/core/auth-key/create",
      data: { name },
    });

    if (response.success && response.data) {
      return response.data;
    }
    throw new Error(response.error || "Failed to create auth key");
  },
};
