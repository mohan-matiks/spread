import React from 'react';
import { Box, Text } from 'rebass/styled-components';
import { useSetup } from '../../api/hooks/useSetup';
import SetupView from '../../views/setup/SetupView';

interface SetupWrapperProps {
    children: React.ReactNode;
}

const SetupWrapper: React.FC<SetupWrapperProps> = ({ children }) => {
    const { setupStatus, loading, error, refetch } = useSetup();

    if (loading) {
        return (
            <Box
                display="flex"
                alignItems="center"
                justifyContent="center"
                minHeight="100vh"
                backgroundColor="#f5f5f5"
            >
                <Text fontSize="18px" color="#666">
                    Checking setup status...
                </Text>
            </Box>
        );
    }

    if (error) {
        return (
            <Box
                display="flex"
                alignItems="center"
                justifyContent="center"
                minHeight="100vh"
                backgroundColor="#f5f5f5"
            >
                <Text fontSize="18px" color="#d32f2f">
                    Error: {error}
                </Text>
            </Box>
        );
    }

    // If setup is not completed (no users exist), show setup view
    if (setupStatus === false) {
        return <SetupView onUserCreated={refetch} />;
    }

    // If setup is completed, show the normal app
    return <>{children}</>;
};

export default SetupWrapper; 