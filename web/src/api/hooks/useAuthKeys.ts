import { useState, useEffect } from "react";
import { authKeyService, AuthKey } from "../services/authKeyService";

export const useAuthKeys = () => {
  const [authKeys, setAuthKeys] = useState<AuthKey[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchAuthKeys = async () => {
    try {
      setLoading(true);
      setError(null);
      const keys = await authKeyService.getAllAuthKeys();
      setAuthKeys(keys);
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to fetch auth keys";
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const createAuthKey = async (name: string): Promise<string | null> => {
    try {
      setError(null);
      const newKey = await authKeyService.createAuthKey(name);
      // Refresh the list after creating a new key
      await fetchAuthKeys();
      return newKey;
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to create auth key";
      setError(errorMessage);
      return null;
    }
  };

  useEffect(() => {
    fetchAuthKeys();
  }, []);

  return {
    authKeys,
    loading,
    error,
    createAuthKey,
    refetch: fetchAuthKeys,
  };
};
