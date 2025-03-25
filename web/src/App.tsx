import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Outlet, useNavigate } from 'react-router-dom';
import { ThemeProvider } from 'styled-components';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import theme from './theme';
import './App.css';
import { navigationService } from './api/navigation';

// Import components
import LoginView from './views/auth/LoginView';
import AppView from './views/dashboard/apps/AppView';
import VersionView from './views/dashboard/versions/VersionView';
import BundleView from './views/dashboard/bundles/BundleView';
import NotFound from './views/404';
import Header from './components/Header';
import AuthInitializer from './components/auth/AuthInitializer';
import ProtectedRoute from './components/auth/ProtectedRoute';

// Navigation provider to initialize navigation service at app level
const NavigationProvider = ({ children }: { children: React.ReactNode }) => {
    const navigate = useNavigate();

    useEffect(() => {
        navigationService.setNavigate(navigate);
    }, [navigate]);

    return <>{children}</>;
};

// Dashboard layout component
const DashboardLayout = () => {
    return (
        <div className="dashboard-layout">
            <Header />
            {/* Navigation and common dashboard elements would go here */}
            <Outlet /> {/* This renders the child routes */}
        </div>
    );
};

const App: React.FC = () => {
    return (
        <ThemeProvider theme={theme}>
            <BrowserRouter basename="/web">
                <NavigationProvider>
                    {/* AuthInitializer checks token on app load */}
                    <AuthInitializer>
                        <Routes>
                            {/* Public routes */}
                            <Route path="/" element={<LoginView />} />
                            <Route path="/login" element={<LoginView />} />

                            {/* Protected dashboard routes */}
                            <Route element={<ProtectedRoute />}>
                                <Route path="/dashboard" element={<DashboardLayout />}>
                                    <Route index element={<AppView />} />
                                    <Route path="version" element={<VersionView />} />
                                    <Route path="version/bundle/:id" element={<BundleView />} />
                                </Route>
                            </Route>

                            {/* Render 404 directly instead of redirecting */}
                            <Route path="*" element={<NotFound />} />
                        </Routes>
                    </AuthInitializer>
                </NavigationProvider>
            </BrowserRouter>
            <ToastContainer
                position="top-right"
                autoClose={5000}
                hideProgressBar={false}
                newestOnTop
                closeOnClick
                rtl={false}
                pauseOnFocusLoss
                draggable
                pauseOnHover
                theme="light"
            />
        </ThemeProvider>
    );
};

export default App;
