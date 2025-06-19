import { useState, useEffect } from "react";
import { setupService } from "../services/setupService";

export const useSetup = () => {
  const [setupStatus, setSetupStatus] = useState<boolean | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const checkSetupStatus = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await setupService.getSetupStatus();
      setSetupStatus(response.data.completed);
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to check setup status";
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkSetupStatus();
  }, []);

  return {
    setupStatus,
    loading,
    error,
    refetch: checkSetupStatus,
  };
};
