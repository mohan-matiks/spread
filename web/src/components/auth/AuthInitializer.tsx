import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { validateToken } from '../../api/services/authService';
import useAuthStore from '../../store/authStore';
import { ROUTES } from '../../api/navigation';

interface AuthInitializerProps {
    children: React.ReactNode;
}

/**
 * Component that checks authentication status on app initialization
 * Validates token and redirects if necessary
 */
const AuthInitializer: React.FC<AuthInitializerProps> = ({ children }) => {
    const navigate = useNavigate();
    const [initializing, setInitializing] = useState(true);
    const { setLoading } = useAuthStore();

    // Function to handle logout
    const handleLogout = () => {
        // Remove token from localStorage
        localStorage.removeItem('token'); // Adjust key name if different
        // Redirect to login page
        navigate(ROUTES.LOGIN, { replace: true });
    };

    useEffect(() => {
        const initializeAuth = async () => {
            setLoading(true);
            try {
                // Check if token exists and is valid
                const isAuthenticated = await validateToken();

                // If authenticated and on login page, redirect to dashboard
                const currentPath = window.location.pathname;
                if (isAuthenticated && (currentPath === '/' || currentPath === ROUTES.LOGIN)) {
                    navigate(ROUTES.DASHBOARD, { replace: true });
                }
            } catch (error: any) {
                console.error('Auth initialization error:', error);

                // Check if error is unauthorized (401)
                // This depends on how your API errors are structured
                if (error.response?.status === 401 ||
                    error.message?.includes('unauthorized') ||
                    error.message?.includes('invalid token')) {
                    handleLogout();
                }
            } finally {
                setLoading(false);
                setInitializing(false);
            }
        };

        initializeAuth();
    }, [navigate]);

    // Display loading indicator during initial token validation
    if (initializing) {
        return (
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh', flexDirection: 'column', background: '#f8f9fa' }}>
                <div style={{
                    width: '60px',
                    height: '60px',
                    border: '5px solid #e0e0e0',
                    borderTopColor: '#34C363',
                    borderRadius: '50%',
                    animation: 'spin 1.5s linear infinite',
                    marginBottom: '20px'
                }}>
                </div>
                <div style={{
                    fontSize: '18px',
                    color: '#4a5568',
                    fontWeight: 500,
                    letterSpacing: '0.5px'
                }}>
                    Putting things together...
                </div>
                <style>
                    {`
                    @keyframes spin {
                        0% { transform: rotate(0deg); }
                        100% { transform: rotate(360deg); }
                    }
                    `}
                </style>
            </div>
        );
    }

    return <>{children}</>;
};

export default AuthInitializer; 