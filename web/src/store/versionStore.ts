import { create } from "zustand";
import { Version, Bundle } from "../types/api";
import { apiRequest } from "../api";

interface VersionState {
  // State
  selectedVersion: Version | null;
  versions: Version[];
  bundles: Bundle[];
  loadingVersions: boolean;
  loadingBundles: boolean;
  error: string | null;

  // Actions
  setSelectedVersion: (version: Version | null) => void;
  setVersions: (versions: Version[]) => void;
  setBundles: (bundles: Bundle[]) => void;
  setLoadingVersions: (loading: boolean) => void;
  setLoadingBundles: (loading: boolean) => void;
  setError: (error: string | null) => void;

  // Async actions
  fetchVersions: (environmentId: string) => Promise<void>;
  fetchVersion: (versionId: string) => Promise<Version | null>;
  fetchBundles: (versionId: string) => Promise<Bundle[]>;

  reset: () => void;
}

const initialState = {
  selectedVersion: null,
  versions: [],
  bundles: [],
  loadingVersions: false,
  loadingBundles: false,
  error: null,
};

const useVersionStore = create<VersionState>((set, get) => ({
  ...initialState,

  // Actions
  setSelectedVersion: (version) => set({ selectedVersion: version }),
  setVersions: (versions) => set({ versions }),
  setBundles: (bundles) => set({ bundles }),
  setLoadingVersions: (loading) => set({ loadingVersions: loading }),
  setLoadingBundles: (loading) => set({ loadingBundles: loading }),
  setError: (error) => set({ error }),

  // Async Actions
  fetchVersions: async (environmentId: string) => {
    try {
      set({ loadingVersions: true, error: null });
      const token = localStorage.getItem("token");

      const response = await apiRequest<Version[]>({
        url: `/core/environment/${environmentId}/version`,
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.success && response.data) {
        set({ versions: response.data });
      } else {
        set({ error: "Failed to fetch versions" });
      }
    } catch (err) {
      console.error("Error fetching versions:", err);
      set({ error: "Error loading versions" });
    } finally {
      set({ loadingVersions: false });
    }
  },

  fetchVersion: async (versionId: string) => {
    try {
      set({ loadingVersions: true, error: null });
      const token = localStorage.getItem("token");

      const response = await apiRequest<Version>({
        url: `/core/version/${versionId}`,
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.success && response.data) {
        set({ selectedVersion: response.data });
        return response.data;
      } else {
        set({ error: "Failed to fetch version details" });
        return null;
      }
    } catch (err) {
      console.error("Error fetching version:", err);
      set({ error: "Error loading version details" });
      return null;
    } finally {
      set({ loadingVersions: false });
    }
  },

  fetchBundles: async (versionId: string) => {
    try {
      set({ loadingBundles: true, error: null });
      const token = localStorage.getItem("token");

      // First get the version to know the current bundle ID
      const versionResponse = await apiRequest<Version>({
        url: `/core/version/${versionId}`,
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const version =
        versionResponse.success && versionResponse.data
          ? versionResponse.data
          : null;

      // Then fetch all bundles for this version
      const response = await apiRequest<Bundle[]>({
        url: `/core/version/bundle/${versionId}`,
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.success && response.data) {
        const bundlesWithActiveState = response.data.map((bundle: Bundle) => ({
          ...bundle,
          isActive: bundle.id === version?.currentBundleId,
        }));

        set({ bundles: bundlesWithActiveState });
        return bundlesWithActiveState;
      } else {
        set({ error: "Failed to fetch bundles" });
        return [];
      }
    } catch (err) {
      console.error("Error fetching bundles:", err);
      set({ error: "Error loading bundles" });
      return [];
    } finally {
      set({ loadingBundles: false });
    }
  },

  // Reset state to initial values
  reset: () => set(initialState),
}));

export default useVersionStore;
