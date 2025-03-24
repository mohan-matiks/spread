// Navigation service to avoid direct window.location usage
// This allows programmatic navigation to be centralized

// Create a navigation object that will be updated by components with navigate function
export const navigationService = {
  navigate: (_: string) => {
    console.warn("Navigation attempted before navigate function was set");
  },
  setNavigate: (navigateFn: (path: string) => void) => {
    navigationService.navigate = navigateFn;
  },
};

// Path constants to avoid string literals
export const ROUTES = {
  LOGIN: "/login",
  DASHBOARD: "/dashboard",
  // Add more routes as needed
};
