import { useState, useEffect } from "react";
import { apiRequest } from "../client";
import { App, validateCreateAppRequest } from "../schemas/appSchemas";

interface ApiResponse<T> {
  data: T;
  success: boolean;
  error?: string;
}

const useAppsHook = () => {
  const [apps, setApps] = useState<App[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [validationErrors, setValidationErrors] = useState<Record<
    string,
    string
  > | null>(null);

  const fetchApps = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await apiRequest<any>({
        method: "GET",
        url: "/core/app",
      });

      // Handle different response structure formats
      if (response.success && response.data && Array.isArray(response.data)) {
        // Direct array in response.data
        setApps(response.data);
      } else if (
        response.success &&
        response.data &&
        Array.isArray(response.data.data)
      ) {
        // Nested array in response.data.data
        setApps(response.data.data);
      } else if (
        response.data &&
        response.data.success &&
        Array.isArray(response.data.data)
      ) {
        // API returned { data: { success: true, data: [app1, app2, ...] } }
        setApps(response.data.data);
      } else {
        setApps([]);
        setError("No apps found");
      }
    } catch (err) {
      setError("Failed to load apps");
      setApps([]);
    } finally {
      setLoading(false);
    }
  };

  const createApp = async (appData: {
    appName: string;
    os: string;
  }): Promise<boolean> => {
    // Reset states
    setLoading(true);
    setError(null);
    setValidationErrors(null);

    // Validate the request data
    const validation = validateCreateAppRequest(appData);

    if (!validation.success) {
      setValidationErrors(validation.error || { _form: "Invalid data" });
      setLoading(false);
      return false;
    }

    try {
      const response = await apiRequest<ApiResponse<App>>({
        method: "POST",
        url: "/core/app",
        data: validation.data,
      });

      if (response.success) {
        // Refresh the app list after successful creation
        await fetchApps();
        return true;
      } else {
        setError(response.error || "Failed to create app");
        return false;
      }
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to create app";
      setError(errorMessage);
      return false;
    } finally {
      setLoading(false);
    }
  };

  // Fetch apps on mount
  useEffect(() => {
    fetchApps();
  }, []);

  return {
    apps,
    loading,
    error,
    validationErrors,
    fetchApps,
    createApp,
  };
};

export const useApps = useAppsHook;
export default useAppsHook;
