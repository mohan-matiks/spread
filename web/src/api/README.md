# Network Layer

This module provides a complete network layer for interacting with SpreadServer APIs.

## Features

- Axios-based HTTP client with interceptors
- Authentication handling with token
- Custom hooks for API interactions
- Error handling and response normalization
- Navigation service for React Router integration

## Structure

- `/client.ts` - Base axios client configuration
- `/hooks/` - React hooks for API interactions
  - `useAuth.ts` - Authentication hook (login, logout, session management)
- `/navigation.ts` - Central navigation service to avoid direct window.location usage

## Usage

### Authentication

```tsx
import { useAuth } from '../api';
import { ROUTES } from '../api/navigation';
import { useNavigate } from 'react-router-dom';

const LoginComponent = () => {
  const navigate = useNavigate();
  const { login, loading, error } = useAuth();
  
  const handleLogin = async (credentials) => {
    const response = await login(credentials);
    
    if (response.success) {
      navigate(ROUTES.DASHBOARD);
    }
  };
  
  return (
    // ...
  );
};
```

### API Response Structure

#### Login Response
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Protected Routes

The client automatically attaches the authentication token to requests when available and handles 401 unauthorized responses by redirecting to the login page.

### Navigation Service

The application uses a centralized navigation service that integrates with React Router's `useNavigate` hook:

```typescript
// In App.tsx or other top-level component
import { navigationService } from './api/navigation';
import { useNavigate } from 'react-router-dom';

// Inside component
const navigate = useNavigate();

useEffect(() => {
  // Set the navigate function for use throughout the app
  navigationService.setNavigate(navigate);
}, [navigate]);
```

This allows non-component code (like API interceptors) to navigate without direct references to `window.location`.

## Error Handling

All API requests return a normalized response with the following structure:

```typescript
{
  success: boolean;
  data?: T;  // Generic type for the response data
  error?: string;
}
``` 