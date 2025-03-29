// API Response Types
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

// Auth Types
export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
}

export interface AuthState {
  isAuthenticated: boolean;
  token: string | null;
  loading: boolean;
  error: string | null;
}

// App Types - Updated to match API schema
export interface App {
  id: string;
  name: string;
  os: string;
  packageName?: string;
  description?: string;
  isValid?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

// Environment Types
export interface Environment {
  id: string;
  name: string;
  appId: string;
  isValid: boolean;
  createdAt: string;
  updatedAt: string;
}

// Version Types
export interface Version {
  id: string;
  environmentId: string;
  appVersion: string;
  versionNumber: number;
  currentBundleId: string;
  updatedAt: string;
  createdAt: string;
}

// Bundle Types
export interface Bundle {
  id: string;
  environmentId: string;
  versionId: string;
  appId: string;
  sequenceId: number;
  hash: string;
  size: number;
  downloadFile: string;
  isMandatory: boolean;
  failed: number;
  installed: number;
  active: number;
  description: string;
  label: string;
  isValid: boolean;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
  isActive?: boolean; // UI state to track if bundle is active
}
