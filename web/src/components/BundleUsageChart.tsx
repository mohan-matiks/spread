import React from 'react';
import { Box, Text, Flex } from 'rebass/styled-components';

interface BundleUsageChartProps {
    activeUsers: number;
    totalDownloads: number;
}

const BundleUsageChart: React.FC<BundleUsageChartProps> = ({ activeUsers, totalDownloads }) => {
    // Ensure we don't have a negative inactive count or exceed total downloads
    const sanitizedActiveUsers = Math.min(Math.max(0, activeUsers), totalDownloads);
    const activePercentage = totalDownloads > 0
        ? (sanitizedActiveUsers / totalDownloads) * 100
        : 0;

    return (
        <Box>
            <Flex
                justifyContent="space-between"
                alignItems="center"
                mb={3}
            >
                <Text fontSize="16px" fontWeight="bold">Usage Statistics</Text>
                <Text color="#666" fontSize="14px">
                    {sanitizedActiveUsers.toLocaleString()} active of {totalDownloads.toLocaleString()} downloads
                    <Text as="span" ml={1} color="#34C363" fontWeight="bold">
                        ({activePercentage.toFixed(1)}%)
                    </Text>
                </Text>
            </Flex>

            <Box
                height={24}
                mb={3}
                sx={{
                    borderRadius: '12px',
                    overflow: 'hidden',
                    boxShadow: '0 2px 4px rgba(0,0,0,0.05)',
                    position: 'relative'
                }}
            >
                {/* Background bar (inactive users) */}
                <Box
                    sx={{
                        height: '100%',
                        width: '100%',
                        backgroundColor: '#E0E0E0',
                        position: 'absolute'
                    }}
                />

                {/* Foreground bar (active users) */}
                <Box
                    sx={{
                        height: '100%',
                        width: `${activePercentage}%`,
                        backgroundColor: '#34C363',
                        position: 'absolute',
                        transition: 'width 0.5s ease-in-out'
                    }}
                />

                {/* Percentage label - only show if enough space */}
                {activePercentage > 10 && (
                    <Flex
                        sx={{
                            position: 'absolute',
                            height: '100%',
                            width: '100%',
                            alignItems: 'center',
                            justifyContent: 'center',
                            pointerEvents: 'none'
                        }}
                    >
                        <Text color="white" fontWeight="bold" fontSize="12px">
                            {activePercentage.toFixed(0)}%
                        </Text>
                    </Flex>
                )}
            </Box>

            <Flex justifyContent="space-between" fontSize="12px" color="#666">
                <Flex alignItems="center">
                    <Box
                        width={10}
                        height={10}
                        mr={1}
                        sx={{
                            backgroundColor: '#34C363',
                            borderRadius: '50%'
                        }}
                    />
                    <Text>Active Users</Text>
                </Flex>
                <Flex alignItems="center">
                    <Box
                        width={10}
                        height={10}
                        mr={1}
                        sx={{
                            backgroundColor: '#E0E0E0',
                            borderRadius: '50%'
                        }}
                    />
                    <Text>Inactive Users</Text>
                </Flex>
            </Flex>
        </Box>
    );
};

export default BundleUsageChart; 