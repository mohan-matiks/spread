import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import useAuthStore from '../../store/authStore';
import { ROUTES } from '../../api/navigation';

interface ProtectedRouteProps {
    children?: React.ReactNode;
}

/**
 * ProtectedRoute component that verifies authentication
 * Redirects to login if not authenticated
 */
const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
    const { isAuthenticated, loading } = useAuthStore();

    // If still loading auth state, you could show a spinner here
    if (loading) {
        return (
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
                Loading...
            </div>
        );
    }

    // If not authenticated, redirect to login
    if (!isAuthenticated) {
        return <Navigate to={ROUTES.LOGIN} replace />;
    }

    // Render children or outlet (for nested routes)
    return children ? <>{children}</> : <Outlet />;
};

export default ProtectedRoute; 