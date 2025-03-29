import { useState, useEffect } from 'react'
import { Box, Text, Flex } from 'rebass/styled-components'
import { useParams, useNavigate } from 'react-router-dom'
import { FaDownload, FaUsers, FaUndo, FaCloudDownloadAlt, FaCheckCircle, FaArrowLeft } from 'react-icons/fa'
import Toggle from '../../../components/Toggle'
import BundleUsageChart from '../../../components/BundleUsageChart'
import { apiRequest } from '../../../api'
import useAppStore from '../../../store/appStore'
import useVersionStore from '../../../store/versionStore'
import { Bundle } from '../../../types/api'

interface ApiResponse {
    success: boolean
    error?: string
}

const BundleView = () => {
    const params = useParams<{ id: string }>()
    const id = params.id
    const navigate = useNavigate()

    // Use the global app store for app and environment data
    const {
        selectedApp,
        selectedEnvironment,
        environments,
        fetchEnvironments
    } = useAppStore()

    // Use the global version store for version and bundle data
    const {
        selectedVersion,
        bundles,
        loadingVersions,
        loadingBundles,
        error: versionError,
        fetchVersion,
        fetchBundles,
        setBundles
    } = useVersionStore()

    // Local loading and error states
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<string | null>(null)
    const [loadingRequiredToggles, setLoadingRequiredToggles] = useState<{ [key: string]: boolean }>({})
    const [loadingActivation, setLoadingActivation] = useState<{ [key: string]: boolean }>({})

    // Load version and bundles
    useEffect(() => {
        if (!id) return;

        const loadData = async () => {
            setLoading(true);
            try {
                // First fetch the version
                const versionData = await fetchVersion(id);

                // Then fetch bundles for this version
                await fetchBundles(id);

                // If we have version data but no environments, fetch environments for this version's app
                if (versionData && environments.length === 0 && versionData.environmentId) {
                    // First get full environment info from API to get appId
                    const token = localStorage.getItem('token');
                    const envResponse = await apiRequest<{ appId: string }>({
                        url: `/core/environment/${versionData.environmentId}`,
                        method: 'GET',
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });

                    // If we successfully got the environment with appId, fetch all environments for this app
                    if (envResponse.success && envResponse.data && envResponse.data.appId) {
                        await fetchEnvironments(envResponse.data.appId);
                    }
                }
            } catch (err) {
                console.error('Error loading version data:', err);
                setError('Failed to load version data');
            } finally {
                setLoading(false);
            }
        };

        loadData();
    }, [id, fetchVersion, fetchBundles, fetchEnvironments, environments.length]);

    useEffect(() => {
        // Set error from store if it exists
        if (versionError) {
            setError(versionError);
        }
    }, [versionError]);

    // Find environment by ID from global store
    const getEnvironmentName = (environmentId: string) => {
        // First try to find the environment in the global store
        const environment = environments.find(env => env.id === environmentId)
        if (environment) {
            // If we found it in the store, use its name
            return environment.name
        }

        // If not found in the store, apply formatting rules based on ID
        switch (environmentId?.toLowerCase()) {
            case 'prod':
            case 'production':
                return 'Production';
            case 'dev':
            case 'development':
                return 'Development';
            case 'qa':
            case 'test':
                return 'Testing';
            case 'staging':
                return 'Staging';
            default:
                // Capitalize first letter for other environment IDs
                return environmentId ? environmentId.charAt(0).toUpperCase() + environmentId.slice(1) : 'Unknown';
        }
    }

    // Add rollback function
    const handleRollback = async (bundleId: string) => {
        try {
            const token = localStorage.getItem('token')
            const bundle = bundles.find(b => b.id === bundleId)

            if (!bundle) return

            const response = await apiRequest<ApiResponse>({
                url: `/core/rollback`,
                method: 'POST',
                data: {
                    appId: bundle.appId,
                    environmentId: bundle.environmentId,
                    versionId: bundle.versionId
                },
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            })

            if (response.success && id) {
                // Refresh bundles from the store
                await fetchBundles(id)
            }
        } catch (err) {
            console.error('Error rolling back bundle:', err)
        }
    }

    // Updates bundle active status via API
    const toggleBundleStatus = async (bundleId: string) => {
        const currentBundle = bundles.find(b => b.id === bundleId)
        if (!currentBundle) {
            console.error('Bundle not found')
            return
        }

        const previousIsValid = currentBundle.isValid
        try {
            setLoadingActivation(prev => ({ ...prev, [bundleId]: true }))
            const token = localStorage.getItem('token')

            // Call API to toggle bundle active status
            const response = await apiRequest<ApiResponse>({
                url: `/core/version/bundle/${bundleId}/active`,
                method: 'PUT',
                headers: {
                    Authorization: `Bearer ${token}`
                }
            })

            if (response.success) {
                // Update local state immediately for better UX
                const updatedBundles = bundles.map(b =>
                    b.id === bundleId
                        ? { ...b, isValid: !previousIsValid }
                        : b
                );
                setBundles(updatedBundles);

                // Refresh bundles to get the latest state from server
                if (id) {
                    await fetchBundles(id);
                }
            } else {
                console.error('Failed to toggle bundle activation:', response.error)
                // Revert the local state if the API call failed
                const revertedBundles = bundles.map(b =>
                    b.id === bundleId
                        ? { ...b, isValid: previousIsValid }
                        : b
                );
                setBundles(revertedBundles);
            }
        } catch (err) {
            console.error('Error toggling bundle status:', err)
            // Revert the local state if there was an error
            const revertedBundles = bundles.map(b =>
                b.id === bundleId
                    ? { ...b, isValid: previousIsValid }
                    : b
            );
            setBundles(revertedBundles);
        } finally {
            setLoadingActivation(prev => ({ ...prev, [bundleId]: false }))
        }
    }

    // Updates bundle mandatory status via API
    const toggleBundleRequired = async (bundleId: string) => {
        try {
            setLoadingRequiredToggles(prev => ({ ...prev, [bundleId]: true }))
            const token = localStorage.getItem('token')
            const bundle = bundles.find(b => b.id === bundleId)

            if (!bundle) return

            // Call API to toggle mandatory status
            const response = await apiRequest<ApiResponse>({
                url: `/core/version/bundle/${bundleId}/mandatory`,
                method: 'PUT',
                data: { isMandatory: !bundle.isMandatory },
                headers: {
                    Authorization: `Bearer ${token}`
                }
            })

            if (response.success) {
                // Update local state
                const updatedBundles = bundles.map(b =>
                    b.id === bundleId
                        ? { ...b, isMandatory: !b.isMandatory }
                        : b
                );
                setBundles(updatedBundles);
            }
        } catch (err) {
            console.error('Error toggling bundle required status:', err)
        } finally {
            setLoadingRequiredToggles(prev => ({ ...prev, [bundleId]: false }))
        }
    }

    const handleDownload = async (bundleId: string) => {
        try {
            const bundle = bundles.find(b => b.id === bundleId)

            if (!bundle) return

            // Use downloadFile field directly from the bundle
            if (bundle.downloadFile) {
                window.open(bundle.downloadFile, '_blank')
            }
        } catch (err) {
            console.error('Error downloading bundle:', err)
        }
    }

    // Format file size to readable format
    const formatFileSize = (bytes: number): string => {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
    }

    // Get sorted bundles (descending by sequenceId)
    const sortedBundles = [...bundles].sort((a, b) => b.sequenceId - a.sequenceId);

    // Split bundles into active and history
    const activeBundle = sortedBundles.find(bundle => bundle.isActive);
    const historyBundles = sortedBundles.filter(bundle => !bundle.isActive);

    // Show loading state
    if (loading || loadingVersions || loadingBundles) {
        return (
            <Box p={4}>
                {/* Header Skeleton */}
                <Flex alignItems="center" mb={5}>
                    <Box
                        sx={{
                            width: '36px',
                            height: '36px',
                            bg: '#f0f0f0',
                            borderRadius: '4px',
                            mr: 3,
                            animation: 'pulse 1.5s ease-in-out infinite',
                            '@keyframes pulse': {
                                '0%': { opacity: 0.6 },
                                '50%': { opacity: 1 },
                                '100%': { opacity: 0.6 }
                            }
                        }}
                    />
                    <Box
                        sx={{
                            width: '120px',
                            height: '32px',
                            bg: '#f0f0f0',
                            borderRadius: '4px',
                            mr: 3,
                            animation: 'pulse 1.5s ease-in-out infinite',
                        }}
                    />
                    <Box
                        sx={{
                            width: '80px',
                            height: '24px',
                            bg: '#f0f0f0',
                            borderRadius: '16px',
                            animation: 'pulse 1.5s ease-in-out infinite',
                        }}
                    />
                </Flex>

                {/* Active Bundle Skeleton */}
                <Box mb={4}>
                    <Box
                        sx={{
                            width: '120px',
                            height: '24px',
                            bg: '#f0f0f0',
                            borderRadius: '4px',
                            mb: 3,
                            animation: 'pulse 1.5s ease-in-out infinite',
                        }}
                    />
                    <Box
                        sx={{
                            border: '1px solid #e0e0e0',
                            borderRadius: '8px',
                            padding: '24px',
                            borderLeftColor: '#f0f0f0',
                            borderLeftWidth: '4px',
                            bg: '#fff',
                        }}
                    >
                        <Flex justifyContent="space-between">
                            {/* Left section skeleton */}
                            <Box width="60%">
                                <Box
                                    sx={{
                                        width: '280px',
                                        height: '24px',
                                        bg: '#f0f0f0',
                                        borderRadius: '4px',
                                        mb: 4,
                                        animation: 'pulse 1.5s ease-in-out infinite',
                                    }}
                                />
                                <Box
                                    sx={{
                                        width: '400px',
                                        height: '16px',
                                        bg: '#f0f0f0',
                                        borderRadius: '4px',
                                        mb: 2,
                                        animation: 'pulse 1.5s ease-in-out infinite',
                                    }}
                                />
                                <Box
                                    sx={{
                                        width: '320px',
                                        height: '16px',
                                        bg: '#f0f0f0',
                                        borderRadius: '4px',
                                        mb: 4,
                                        animation: 'pulse 1.5s ease-in-out infinite',
                                    }}
                                />
                                <Flex flexWrap="wrap">
                                    {[1, 2, 3, 4].map((i) => (
                                        <Box
                                            key={i}
                                            sx={{
                                                width: '120px',
                                                height: '32px',
                                                bg: '#f0f0f0',
                                                borderRadius: '4px',
                                                mr: 3,
                                                mb: 2,
                                                animation: 'pulse 1.5s ease-in-out infinite',
                                                animationDelay: `${i * 0.1}s`,
                                            }}
                                        />
                                    ))}
                                </Flex>
                            </Box>

                            {/* Right section skeleton */}
                            <Box width="30%" sx={{ textAlign: 'right' }}>
                                {[1, 2, 3].map((i) => (
                                    <Box
                                        key={i}
                                        sx={{
                                            width: '140px',
                                            height: '32px',
                                            bg: '#f0f0f0',
                                            borderRadius: '4px',
                                            mb: 3,
                                            ml: 'auto',
                                            animation: 'pulse 1.5s ease-in-out infinite',
                                            animationDelay: `${i * 0.1}s`,
                                        }}
                                    />
                                ))}
                            </Box>
                        </Flex>
                    </Box>
                </Box>

                {/* History Skeleton */}
                <Box>
                    <Box
                        sx={{
                            width: '120px',
                            height: '24px',
                            bg: '#f0f0f0',
                            borderRadius: '4px',
                            mb: 3,
                            animation: 'pulse 1.5s ease-in-out infinite',
                        }}
                    />
                    {[1, 2].map((i) => (
                        <Box
                            key={i}
                            sx={{
                                border: '1px solid #e0e0e0',
                                borderRadius: '8px',
                                padding: '24px',
                                mb: 3,
                                opacity: 0.85,
                                bg: '#fff',
                            }}
                        >
                            <Flex justifyContent="space-between">
                                <Box width="60%">
                                    <Box
                                        sx={{
                                            width: '240px',
                                            height: '24px',
                                            bg: '#f0f0f0',
                                            borderRadius: '4px',
                                            mb: 4,
                                            animation: 'pulse 1.5s ease-in-out infinite',
                                        }}
                                    />
                                    <Box
                                        sx={{
                                            width: '360px',
                                            height: '16px',
                                            bg: '#f0f0f0',
                                            borderRadius: '4px',
                                            mb: 2,
                                            animation: 'pulse 1.5s ease-in-out infinite',
                                        }}
                                    />
                                    <Box
                                        sx={{
                                            width: '280px',
                                            height: '16px',
                                            bg: '#f0f0f0',
                                            borderRadius: '4px',
                                            mb: 4,
                                            animation: 'pulse 1.5s ease-in-out infinite',
                                        }}
                                    />
                                    <Flex flexWrap="wrap">
                                        {[1, 2, 3, 4].map((j) => (
                                            <Box
                                                key={j}
                                                sx={{
                                                    width: '110px',
                                                    height: '32px',
                                                    bg: '#f0f0f0',
                                                    borderRadius: '4px',
                                                    mr: 3,
                                                    mb: 2,
                                                    animation: 'pulse 1.5s ease-in-out infinite',
                                                    animationDelay: `${j * 0.1}s`,
                                                }}
                                            />
                                        ))}
                                    </Flex>
                                </Box>
                                <Box width="30%" sx={{ textAlign: 'right' }}>
                                    {[1, 2, 3].map((j) => (
                                        <Box
                                            key={j}
                                            sx={{
                                                width: '140px',
                                                height: '32px',
                                                bg: '#f0f0f0',
                                                borderRadius: '4px',
                                                mb: 3,
                                                ml: 'auto',
                                                animation: 'pulse 1.5s ease-in-out infinite',
                                                animationDelay: `${j * 0.1}s`,
                                            }}
                                        />
                                    ))}
                                </Box>
                            </Flex>
                        </Box>
                    ))}
                </Box>
            </Box>
        )
    }

    // Show error state
    if (error) {
        return (
            <Box p={4} textAlign="center" color="red">
                <Text fontSize="18px">{error}</Text>
            </Box>
        )
    }

    return (
        <Box>
            <Box
                backgroundColor={"#fff"}
                minHeight={"90vh"}
                sx={{
                    border: "1px solid #e0e0e0",
                    borderRadius: "10px",
                    padding: "20px",
                    margin: "10px",
                }}>
                <Flex flexDirection={"column"}>
                    <Flex alignItems="center">
                        <Box
                            as="button"
                            mr={3}
                            sx={{
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                height: '36px',
                                width: '36px',
                                border: '1px solid #e0e0e0',
                                borderRadius: '4px',
                                bg: 'transparent',
                                cursor: 'pointer',
                                '&:hover': {
                                    bg: '#f7f7f7'
                                }
                            }}
                            onClick={() => navigate(-1)}
                        >
                            <FaArrowLeft size={16} />
                        </Box>
                        <Text fontSize={"28px"} fontWeight={"bold"} mr={3}>Releases</Text>
                        {selectedVersion && (
                            <>
                                <Box
                                    sx={{
                                        display: 'inline-block',
                                        bg: '#f0f0f0',
                                        color: '#666',
                                        px: 2,
                                        py: 1,
                                        borderRadius: '16px',
                                        fontSize: '14px',
                                        fontWeight: '500',
                                        mr: 2
                                    }}
                                >
                                    {selectedVersion.appVersion}
                                </Box>
                                <Box
                                    sx={{
                                        display: 'inline-block',
                                        bg: '#e8f5e9',
                                        color: '#2e7d32',
                                        px: 2,
                                        py: 1,
                                        borderRadius: '16px',
                                        fontSize: '14px',
                                        fontWeight: '500'
                                    }}
                                >
                                    {getEnvironmentName(selectedVersion.environmentId)}
                                </Box>
                            </>
                        )}
                    </Flex>
                </Flex>

                <Box mt={5}></Box>

                {/* Active Bundle Section */}
                {activeBundle ? (
                    <Box mb={4}>
                        <Text fontSize="18px" fontWeight="600" mb={3} color="#333">
                            Active Bundle
                        </Text>
                        <Box
                            sx={{
                                border: '1px solid #e0e0e0',
                                borderRadius: '8px',
                                padding: '16px',
                                borderLeftColor: '#34C363',
                                borderLeftWidth: '4px',
                            }}
                        >
                            <Flex justifyContent="space-between" alignItems="center">
                                {/* Left section - Bundle information */}
                                <Box width="40%">
                                    <Flex alignItems="center" mb={2}>
                                        <Text fontWeight="600" fontSize="16px" mr={2}>
                                            #{activeBundle.sequenceId}
                                        </Text>
                                    </Flex>

                                    <Text
                                        fontSize="15px"
                                        mb={3}
                                    >
                                        {activeBundle.description}
                                    </Text>

                                    <Flex flexWrap="wrap">
                                        <Flex
                                            alignItems="center"
                                            mr={4}
                                            mb={2}
                                            sx={{
                                                px: 2,
                                                py: 1,
                                                borderRadius: '4px',
                                                bg: '#f7f7f7',
                                            }}
                                        >
                                            <Box mr={2} color="#34C363">
                                                <FaUsers size={14} />
                                            </Box>
                                            <Text fontSize="14px" fontWeight="500">
                                                {activeBundle.active.toLocaleString()} active
                                            </Text>
                                        </Flex>

                                        <Flex
                                            alignItems="center"
                                            mr={4}
                                            mb={2}
                                            sx={{
                                                px: 2,
                                                py: 1,
                                                borderRadius: '4px',
                                                bg: '#f7f7f7',
                                            }}
                                        >
                                            <Box mr={2} color="#f59e0b">
                                                <FaUndo size={14} />
                                            </Box>
                                            <Text fontSize="14px" fontWeight="500">
                                                {activeBundle.failed} failed
                                            </Text>
                                        </Flex>

                                        <Flex
                                            alignItems="center"
                                            mr={4}
                                            mb={2}
                                            sx={{
                                                px: 2,
                                                py: 1,
                                                borderRadius: '4px',
                                                bg: '#f7f7f7',
                                            }}
                                        >
                                            <Box mr={2} color="#34C363">
                                                <FaCloudDownloadAlt size={14} />
                                            </Box>
                                            <Text fontSize="14px" fontWeight="500">
                                                {activeBundle.installed.toLocaleString()} installed
                                            </Text>
                                        </Flex>

                                        <Flex
                                            alignItems="center"
                                            mb={2}
                                            sx={{
                                                px: 2,
                                                py: 1,
                                                borderRadius: '4px',
                                                bg: '#f7f7f7',
                                            }}
                                        >
                                            <Box mr={2} color="#34C363">
                                                <FaDownload size={14} />
                                            </Box>
                                            <Text fontSize="14px" fontWeight="500">
                                                {formatFileSize(activeBundle.size)}
                                            </Text>
                                        </Flex>
                                    </Flex>

                                    <Text fontSize="14px" color="#666" mt={2}>
                                        Created {new Date(activeBundle.createdAt).toLocaleDateString('en-US', {
                                            month: 'short',
                                            day: 'numeric',
                                            year: '2-digit'
                                        })}{activeBundle.createdBy ? ` by ${activeBundle.createdBy}` : ''}
                                    </Text>
                                </Box>

                                {/* Middle section - Donut chart */}
                                <Box width="30%" sx={{ textAlign: 'center' }}>
                                    <BundleUsageChart
                                        activeUsers={activeBundle.active}
                                        totalDownloads={activeBundle.installed}
                                    />
                                </Box>

                                {/* Right section - Actions */}
                                <Flex width="30%" justifyContent="flex-end" alignItems="center">
                                    <Flex flexDirection="column" alignItems="flex-end">
                                        {activeBundle.isMandatory && (
                                            <Box
                                                mb={3}
                                                sx={{
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    color: '#34C363',
                                                }}
                                            >
                                                <FaCheckCircle size={15} />
                                                <Text ml={2} fontSize="14px" fontWeight="500">Required</Text>
                                            </Box>
                                        )}
                                        <Box
                                            as="button"
                                            mb={3}
                                            sx={{
                                                display: 'flex',
                                                alignItems: 'center',
                                                justifyContent: 'center',
                                                height: '28px',
                                                px: 3,
                                                border: '1px solid #f59e0b',
                                                color: '#f59e0b',
                                                borderRadius: '4px',
                                                bg: 'transparent',
                                                fontSize: '14px',
                                                fontWeight: '500',
                                                cursor: 'pointer',
                                                '&:hover': {
                                                    bg: '#fff8e6'
                                                }
                                            }}
                                            onClick={() => handleRollback(activeBundle.id)}
                                        >
                                            <Box as="span" mr={1}><FaUndo size={12} /></Box>
                                            Rollback
                                        </Box>
                                        <Box mb={3}>
                                            <Toggle
                                                isActive={activeBundle.isValid}
                                                onChange={() => toggleBundleStatus(activeBundle.id)}
                                                label="Enabled"
                                                size="small"
                                                disabled={loadingActivation[activeBundle.id]}
                                            />
                                        </Box>
                                        <Box mb={3}>
                                            <Toggle
                                                isActive={activeBundle.isMandatory}
                                                onChange={() => toggleBundleRequired(activeBundle.id)}
                                                label="Required"
                                                size="small"
                                                disabled={loadingRequiredToggles[activeBundle.id]}
                                            />
                                        </Box>
                                        {activeBundle.downloadFile && (
                                            <Box
                                                as="button"
                                                sx={{
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    justifyContent: 'center',
                                                    height: '28px',
                                                    px: 3,
                                                    border: '1px solid #666',
                                                    color: '#666',
                                                    borderRadius: '4px',
                                                    bg: 'transparent',
                                                    fontSize: '14px',
                                                    fontWeight: '500',
                                                    cursor: 'pointer',
                                                    '&:hover': {
                                                        bg: '#f7f7f7'
                                                    }
                                                }}
                                                onClick={() => handleDownload(activeBundle.id)}
                                            >
                                                <Box as="span" mr={1}><FaDownload size={12} /></Box>
                                                Download
                                            </Box>
                                        )}
                                    </Flex>
                                </Flex>
                            </Flex>
                        </Box>
                    </Box>
                ) : (
                    <Box
                        mb={4}
                        p={4}
                        sx={{
                            borderRadius: "8px",
                            backgroundColor: "#fafafa",
                            textAlign: "center",
                            border: "1px dashed #e0e0e0"
                        }}
                    >
                        <Text fontSize={"18px"} fontWeight={"bold"} color="#333">
                            No Active Bundle
                        </Text>
                        <Text fontSize={"16px"} color={"#666"} mt={1}>
                            Activate a bundle from the history section below
                        </Text>
                    </Box>
                )}

                {/* History Section */}
                <Box>
                    <Text fontSize="18px" fontWeight="600" mb={3} color="#333">
                        History
                    </Text>

                    {historyBundles.length > 0 ? (
                        historyBundles.map((bundle) => (
                            <Box
                                key={bundle.id}
                                sx={{
                                    border: '1px solid #e0e0e0',
                                    borderRadius: '8px',
                                    padding: '16px',
                                    marginBottom: '12px',
                                    opacity: 0.85,
                                }}
                            >
                                <Flex justifyContent="space-between" alignItems="center">
                                    {/* Left section - Bundle information */}
                                    <Box width="40%">
                                        <Flex alignItems="center" mb={2}>
                                            <Text fontWeight="600" fontSize="16px" mr={2}>
                                                #{bundle.sequenceId}
                                            </Text>
                                        </Flex>

                                        <Text
                                            fontSize="15px"
                                            mb={3}
                                        >
                                            {bundle.description}
                                        </Text>

                                        <Flex flexWrap="wrap">
                                            <Flex
                                                alignItems="center"
                                                mr={4}
                                                mb={2}
                                                sx={{
                                                    px: 2,
                                                    py: 1,
                                                    borderRadius: '4px',
                                                    bg: '#f7f7f7',
                                                }}
                                            >
                                                <Box mr={2} color="#666">
                                                    <FaUsers size={14} />
                                                </Box>
                                                <Text fontSize="14px" fontWeight="500" color="#666">
                                                    {bundle.active.toLocaleString()} active
                                                </Text>
                                            </Flex>

                                            <Flex
                                                alignItems="center"
                                                mr={4}
                                                mb={2}
                                                sx={{
                                                    px: 2,
                                                    py: 1,
                                                    borderRadius: '4px',
                                                    bg: '#f7f7f7',
                                                }}
                                            >
                                                <Box mr={2} color="#f59e0b">
                                                    <FaUndo size={14} />
                                                </Box>
                                                <Text fontSize="14px" fontWeight="500">
                                                    {bundle.failed} failed
                                                </Text>
                                            </Flex>

                                            <Flex
                                                alignItems="center"
                                                mr={4}
                                                mb={2}
                                                sx={{
                                                    px: 2,
                                                    py: 1,
                                                    borderRadius: '4px',
                                                    bg: '#f7f7f7',
                                                }}
                                            >
                                                <Box mr={2} color="#666">
                                                    <FaCloudDownloadAlt size={14} />
                                                </Box>
                                                <Text fontSize="14px" fontWeight="500" color="#666">
                                                    {bundle.installed.toLocaleString()} installed
                                                </Text>
                                            </Flex>

                                            <Flex
                                                alignItems="center"
                                                mb={2}
                                                sx={{
                                                    px: 2,
                                                    py: 1,
                                                    borderRadius: '4px',
                                                    bg: '#f7f7f7',
                                                }}
                                            >
                                                <Box mr={2} color="#666">
                                                    <FaDownload size={14} />
                                                </Box>
                                                <Text fontSize="14px" fontWeight="500" color="#666">
                                                    {formatFileSize(bundle.size)}
                                                </Text>
                                            </Flex>
                                        </Flex>

                                        <Text fontSize="14px" color="#666" mt={2}>
                                            Created {new Date(bundle.createdAt).toLocaleDateString('en-US', {
                                                month: 'short',
                                                day: 'numeric',
                                                year: '2-digit'
                                            })}{bundle.createdBy ? ` by ${bundle.createdBy}` : ''}
                                        </Text>
                                    </Box>

                                    {/* Middle section - Donut chart */}
                                    <Box width="30%" sx={{ textAlign: 'center' }}>
                                        <BundleUsageChart
                                            activeUsers={bundle.active}
                                            totalDownloads={bundle.installed}
                                        />
                                    </Box>

                                    {/* Right section - Actions */}
                                    <Flex width="30%" justifyContent="flex-end" alignItems="center">
                                        <Flex flexDirection="column" alignItems="flex-end">
                                            {bundle.isMandatory && (
                                                <Box
                                                    mb={3}
                                                    sx={{
                                                        display: 'flex',
                                                        alignItems: 'center',
                                                        color: '#666',
                                                    }}
                                                >
                                                    <FaCheckCircle size={15} />
                                                    <Text ml={2} fontSize="14px" fontWeight="500">Required</Text>
                                                </Box>
                                            )}
                                            <Box mb={3}>
                                                <Toggle
                                                    isActive={bundle.isValid}
                                                    onChange={() => toggleBundleStatus(bundle.id)}
                                                    label="Enabled"
                                                    size="small"
                                                    disabled={loadingActivation[bundle.id]}
                                                />
                                            </Box>
                                            <Box mb={3}>
                                                <Toggle
                                                    isActive={bundle.isMandatory}
                                                    onChange={() => toggleBundleRequired(bundle.id)}
                                                    label="Required"
                                                    size="small"
                                                    disabled={loadingRequiredToggles[bundle.id]}
                                                />
                                            </Box>
                                            {bundle.downloadFile && (
                                                <Box
                                                    as="button"
                                                    sx={{
                                                        display: 'flex',
                                                        alignItems: 'center',
                                                        justifyContent: 'center',
                                                        height: '28px',
                                                        px: 3,
                                                        border: '1px solid #666',
                                                        color: '#666',
                                                        borderRadius: '4px',
                                                        bg: 'transparent',
                                                        fontSize: '14px',
                                                        fontWeight: '500',
                                                        cursor: 'pointer',
                                                        '&:hover': {
                                                            bg: '#f7f7f7'
                                                        }
                                                    }}
                                                    onClick={() => handleDownload(bundle.id)}
                                                >
                                                    <Box as="span" mr={1}><FaDownload size={12} /></Box>
                                                    Download
                                                </Box>
                                            )}
                                        </Flex>
                                    </Flex>
                                </Flex>
                            </Box>
                        ))
                    ) : (
                        <Box
                            p={4}
                            sx={{
                                borderRadius: "8px",
                                backgroundColor: "#fafafa",
                                textAlign: "center",
                                border: "1px dashed #e0e0e0"
                            }}
                        >
                            <Text fontSize={"16px"} color={"#666"}>
                                No inactive bundles in history
                            </Text>
                        </Box>
                    )}
                </Box>

                {bundles.length === 0 && (
                    <Box
                        mt={4}
                        p={4}
                        sx={{
                            borderRadius: "8px",
                            backgroundColor: "#fafafa",
                            textAlign: "center",
                            border: "1px dashed #e0e0e0"
                        }}
                    >
                        <Text fontSize={"18px"} fontWeight={"bold"} mb={2}>
                            No bundles found
                        </Text>
                        <Text fontSize={"16px"} color={"#666"} mb={3}>
                            Create a bundle to get started
                        </Text>
                    </Box>
                )}
            </Box>
        </Box>
    )
}

export default BundleView