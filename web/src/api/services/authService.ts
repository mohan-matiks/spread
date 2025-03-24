import { apiRequest } from "../client";
import useAuthStore from "../../store/authStore";

// Type for user data from API
interface UserResponse {
  id: string;
  username: string;
  password?: string;
  roles: string[];
  isValid: boolean;
  createdAt: string;
  updatedAt: string;
}

/**
 * Validates the current token by making a request to the user endpoint
 * Updates auth store based on the result
 */
export const validateToken = async (): Promise<boolean> => {
  const token = localStorage.getItem("token");

  // If no token exists, user is not authenticated
  if (!token) {
    useAuthStore.getState().reset();
    return false;
  }

  // Set token in store
  useAuthStore.getState().setToken(token);
  useAuthStore.getState().setLoading(true);

  try {
    const response = await apiRequest<UserResponse>({
      method: "GET",
      url: "/core/user",
    });

    if (response.success && response.data) {
      // User is authenticated and we have user data
      useAuthStore.getState().setUser(response.data);
      useAuthStore.getState().setAuthenticated(true);
      useAuthStore.getState().setError(null);
      return true;
    } else {
      // Invalid response - clear auth state
      localStorage.removeItem("token");
      useAuthStore.getState().reset();
      return false;
    }
  } catch (error) {
    // Request failed - clear auth state
    localStorage.removeItem("token");
    useAuthStore.getState().reset();
    return false;
  } finally {
    useAuthStore.getState().setLoading(false);
  }
};

/**
 * Logs the user out by clearing token and auth state
 */
export const logout = (): void => {
  localStorage.removeItem("token");
  useAuthStore.getState().reset();
};
