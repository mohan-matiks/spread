import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from "axios";
import { ApiResponse } from "../types/api";
import { navigationService, ROUTES } from "./navigation";

// Create axios instance with default config
const client = axios.create({
  baseURL: "http://localhost:4000",
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor
client.interceptors.request.use(
  (config) => {
    // Get token from localStorage
    const token = localStorage.getItem("token");

    // Add auth header if token exists
    if (token) {
      config.headers["Authorization"] = `Bearer ${token}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
client.interceptors.response.use(
  (response) => {
    return response;
  },
  (error: AxiosError) => {
    // Handle common error scenarios like 401 (unauthorized)
    if (error.response?.status === 401) {
      // Clear local storage and redirect to login
      localStorage.removeItem("token");
      navigationService.navigate(ROUTES.LOGIN);
    }

    return Promise.reject(error);
  }
);

// Generic API request function
export const apiRequest = async <T>(
  config: AxiosRequestConfig
): Promise<ApiResponse<T>> => {
  try {
    const response: AxiosResponse = await client(config);

    console.log(`API Response from ${config.url}:`, response.data);

    // Check if response has the expected structure with nested data
    if (response.data && response.data.success !== undefined) {
      return response.data as ApiResponse<T>;
    }

    // Fallback for APIs that don't follow the standard response structure
    return {
      success: true,
      data: response.data,
    };
  } catch (error) {
    if (axios.isAxiosError(error)) {
      // Handle axios errors
      const errorMessage =
        error.response?.data?.error || error.message || "An error occurred";

      return {
        success: false,
        error: errorMessage,
      };
    }

    // Handle non-axios errors
    return {
      success: false,
      error: "An unexpected error occurred",
    };
  }
};

export default client;
