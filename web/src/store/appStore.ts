import { create } from "zustand";
import { App, Environment } from "../types/api";
import { apiRequest } from "../api";

interface AppState {
  // State
  selectedApp: App | null;
  selectedEnvironment: Environment | null;
  environments: Environment[];
  apps: App[];
  loadingApps: boolean;
  loadingEnvironments: boolean;
  error: string | null;

  // Actions
  setSelectedApp: (app: App | null) => void;
  setSelectedEnvironment: (env: Environment | null) => void;
  setEnvironments: (environments: Environment[]) => void;
  setApps: (apps: App[]) => void;
  setLoadingApps: (loading: boolean) => void;
  setLoadingEnvironments: (loading: boolean) => void;
  setError: (error: string | null) => void;

  // Async actions
  fetchApps: () => Promise<void>;
  fetchEnvironments: (appId: string) => Promise<void>;

  reset: () => void;
}

const initialState = {
  selectedApp: null,
  selectedEnvironment: null,
  environments: [],
  apps: [],
  loadingApps: false,
  loadingEnvironments: false,
  error: null,
};

const useAppStore = create<AppState>((set) => ({
  ...initialState,

  // Actions
  setSelectedApp: (app) => set({ selectedApp: app }),
  setSelectedEnvironment: (env) => set({ selectedEnvironment: env }),
  setEnvironments: (environments) => set({ environments }),
  setApps: (apps) => set({ apps }),
  setLoadingApps: (loading) => set({ loadingApps: loading }),
  setLoadingEnvironments: (loading) => set({ loadingEnvironments: loading }),
  setError: (error) => set({ error }),

  // Async Actions
  fetchApps: async () => {
    try {
      set({ loadingApps: true, error: null });
      const token = localStorage.getItem("token");

      const response = await apiRequest<App[]>({
        url: "/core/app",
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.success && response.data) {
        set({ apps: response.data });
      } else {
        set({ error: "Failed to fetch apps" });
      }
    } catch (err) {
      console.error("Error fetching apps:", err);
      set({ error: "Error loading apps" });
    } finally {
      set({ loadingApps: false });
    }
  },

  fetchEnvironments: async (appId: string) => {
    try {
      set({ loadingEnvironments: true, error: null });
      const token = localStorage.getItem("token");

      const response = await apiRequest<Environment[]>({
        url: `/core/environment/${appId}`,
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.success && response.data) {
        set({ environments: response.data });
      } else {
        set({ error: "Failed to fetch environments" });
      }
    } catch (err) {
      console.error("Error fetching environments:", err);
      set({ error: "Error loading environments" });
    } finally {
      set({ loadingEnvironments: false });
    }
  },

  // Reset state to initial values
  reset: () => set(initialState),
}));

export default useAppStore;
