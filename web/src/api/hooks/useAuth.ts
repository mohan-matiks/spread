import { useState } from "react";
import { apiRequest } from "../client";
import { LoginRequest, LoginResponse, ApiResponse } from "../../types/api";
import { navigationService, ROUTES } from "../navigation";
import useAuthStore from "../../store/authStore";

// User type from API response
interface UserResponse {
  id: string;
  username: string;
  password?: string;
  roles: string[];
  isValid: boolean;
  createdAt: string;
  updatedAt: string;
}

const useAuth = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Get auth state from zustand store
  const { isAuthenticated, user, setAuthenticated, setUser, setToken } =
    useAuthStore();

  // Login function
  const login = async (
    credentials: LoginRequest
  ): Promise<ApiResponse<LoginResponse>> => {
    setLoading(true);
    setError(null);

    try {
      const response = await apiRequest<LoginResponse>({
        method: "POST",
        url: "/login",
        data: credentials,
      });

      if (response.success && response.data) {
        // Store token in localStorage
        const token = response.data.access_token;
        localStorage.setItem("token", token);

        // Update auth store
        setToken(token);

        // Fetch user info after successful login
        const userResponse = await apiRequest<UserResponse>({
          method: "GET",
          url: "/core/user",
        });

        if (userResponse.success && userResponse.data) {
          setUser(userResponse.data);
          setAuthenticated(true);
        }
      } else if (response.error) {
        setError(response.error);
      }

      return response;
    } catch (err) {
      const errorMessage = "Login failed. Please try again.";
      setError(errorMessage);
      return { success: false, error: errorMessage };
    } finally {
      setLoading(false);
    }
  };

  // Logout function
  const logout = (): void => {
    localStorage.removeItem("token");
    // Navigate to login page using the navigation service
    navigationService.navigate(ROUTES.LOGIN);
  };

  return {
    login,
    logout,
    isAuthenticated,
    user,
    loading,
    error,
  };
};

export default useAuth;
