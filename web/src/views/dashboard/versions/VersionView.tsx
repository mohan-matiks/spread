import { useState, useEffect } from 'react'
import { Box, Text, Flex } from 'rebass/styled-components'
import Button from '../../../components/primitives/Button'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { FaArrowLeft, FaCopy } from 'react-icons/fa'
import CreateEnvironmentModal, { CreateEnvironmentFormData } from '../../../components/modal/CreateEnvironmentModal'
import { apiRequest } from '../../../api'
import { toast } from 'react-toastify'

type Version = {
    id: string
    environmentId: string
    appVersion: string
    versionNumber: number
    currentBundleId: string
    updatedAt: string
    createdAt: string
}

// Custom hook for fetching version by ID
export const useVersionFetcher = () => {
    const fetchVersionById = async (versionId: string): Promise<Version> => {
        try {
            const response = await apiRequest<{ data: Version }>({
                method: 'GET',
                url: `/core/version/${versionId}`
            })

            if (response.success && response.data?.data) {
                return response.data.data
            }
            throw new Error(response.error || 'Failed to fetch version')
        } catch (error) {
            console.error('Failed to fetch version:', error)
            throw error
        }
    }

    return { fetchVersionById }
}

type Environment = {
    id: string
    appId: string
    name: string
    key: string
    updatedAt: string
    createdAt: string
}

type App = {
    id: string
    name: string
    os: string
    createdAt?: string
    updatedAt?: string
}


const VersionView = () => {
    const [searchParams] = useSearchParams()
    const appId = searchParams.get('appId')
    const navigate = useNavigate()

    const [versions, setVersions] = useState<Version[]>([])
    const [environments, setEnvironments] = useState<Environment[]>([])
    const [appDetails, setAppDetails] = useState<App | null>(null)
    const [selectedEnvironment, setSelectedEnvironment] = useState<string>('development')
    const [selectedEnvironmentId, setSelectedEnvironmentId] = useState<string>('')
    const [showCreateEnvModal, setShowCreateEnvModal] = useState(false)
    const [dropdownOpen, setDropdownOpen] = useState(false)
    const [isLoading, setIsLoading] = useState(true)
    const [versionsLoading, setVersionsLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [versionsError, setVersionsError] = useState<string | null>(null)

    // Fetch app details
    useEffect(() => {
        const fetchAppDetails = async () => {
            if (!appId) return;

            try {
                const response = await apiRequest<App>({
                    method: 'GET',
                    url: `/core/app/${appId}`
                });

                if (response.success && response.data) {
                    setAppDetails(response.data);
                }
            } catch (error) {
                console.error('Failed to fetch app details:', error);
            }
        };

        fetchAppDetails();
    }, [appId]);

    // Fetch environments on component mount
    useEffect(() => {
        const fetchEnvironments = async () => {
            if (!appId) return;

            setIsLoading(true)
            try {
                const response = await apiRequest<Environment[]>({
                    method: 'GET',
                    url: `/core/environment/${appId}`
                })

                if (response.success && response.data) {
                    setEnvironments(response.data)
                    // Set first environment as selected by default
                    if (response.data.length > 0) {
                        setSelectedEnvironment(response.data[0].name)
                        setSelectedEnvironmentId(response.data[0].id)
                    }
                }
            } catch (error) {
                console.error('Failed to fetch environments:', error)
                setError('Failed to load environments')
            } finally {
                setIsLoading(false)
            }
        }

        fetchEnvironments()
    }, [appId])

    // Fetch versions for the selected environment
    useEffect(() => {
        const fetchVersions = async () => {
            if (!selectedEnvironmentId) return;

            setVersionsLoading(true)
            setVersionsError(null)

            try {
                const response = await apiRequest<Version[]>({
                    method: 'GET',
                    url: `/core/version?environmentId=${selectedEnvironmentId}`
                })

                if (response.success && response.data) {
                    setVersions(response.data)
                } else {
                    setVersions([])
                    setVersionsError('Failed to load versions')
                }
            } catch (error) {
                console.error('Failed to fetch versions:', error)
                setVersions([])
                setVersionsError('Failed to load versions')
            } finally {
                setVersionsLoading(false)
            }
        }

        fetchVersions()
    }, [selectedEnvironmentId])

    useEffect(() => {
        if (dropdownOpen) {
            const handleClickOutside = (event: MouseEvent) => {
                const target = event.target as Node;
                const dropdown = document.getElementById('env-dropdown');
                if (dropdown && !dropdown.contains(target)) {
                    setDropdownOpen(false);
                }
            };

            document.addEventListener('mousedown', handleClickOutside);
            return () => {
                document.removeEventListener('mousedown', handleClickOutside);
            };
        }
    }, [dropdownOpen]);

    // const handleEnvironmentChange = (event: React.FormEvent<HTMLDivElement>) => {
    //     const target = event.target as HTMLSelectElement;
    //     setSelectedEnvironment(target.value);
    // }

    const handleCreateEnvironment = () => {
        setShowCreateEnvModal(true)
    }

    const handleCloseEnvModal = () => {
        setShowCreateEnvModal(false)
    }

    const handleAddEnvironment = async (data: CreateEnvironmentFormData) => {
        if (!appId) return;

        setIsLoading(true);
        try {
            const response = await apiRequest<Environment>({
                method: 'POST',
                url: `/core/environment`,
                data: {
                    environmentName: data.environmentName,
                    appName: data.appName
                }
            });

            if (response.success && response.data) {
                const newEnvironment: Environment = {
                    id: response.data.id || '',
                    appId: appId,
                    name: response.data.name,
                    key: response.data.key,
                    updatedAt: new Date().toISOString(),
                    createdAt: new Date().toISOString()
                };

                setEnvironments([...environments, newEnvironment]);
                setSelectedEnvironment(response.data.name);
                setSelectedEnvironmentId(response.data.id);
                setShowCreateEnvModal(false);
            } else {
                setError(response.error || 'Failed to create environment');
            }
        } catch (error) {
            console.error('Error creating environment:', error);
            setError('Failed to create environment');
        } finally {
            setIsLoading(false);
        }
    }

    const formatVersionNumber = (versionNumber: number): string => {
        // Convert number like 10001 to string like "1.0.1"
        const str = versionNumber.toString().padStart(5, '0');
        return `${str.charAt(0)}.${str.charAt(1)}.${str.substring(2)}`.replace(/\.0+$/, '');
    };

    const copyToClipboard = async (text: string, event: React.MouseEvent) => {
        event.stopPropagation();
        try {
            await navigator.clipboard.writeText(text);
            toast.success('Copied to clipboard!');
        } catch (err) {
            toast.error('Failed to copy to clipboard');
        }
    };

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
                <Flex justifyContent={"space-between"} flexDirection={"row"}>
                    <Box flexDirection={"column"}>
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
                            <Text fontSize={"28px"} fontWeight={"bold"}>{appDetails?.name || 'Loading...'}</Text>
                        </Flex>
                    </Box>
                    <Flex>
                        {isLoading && environments.length === 0 ? (
                            null
                        ) : error && environments.length === 0 ? (
                            <Box width={240} mr={2} sx={{ display: 'flex', alignItems: 'center' }}>
                                <Text color="red">{error}</Text>
                            </Box>
                        ) : environments.length === 0 ? (
                            <Box width={240} mr={2}></Box>
                        ) : (
                            <Box width={240} mr={2} sx={{ position: 'relative' }}>
                                <Box id="env-dropdown" sx={{ position: 'relative' }}>
                                    <Flex
                                        onClick={() => setDropdownOpen(!dropdownOpen)}
                                        sx={{
                                            alignItems: 'center',
                                            justifyContent: 'space-between',
                                            width: '100%',
                                            height: '40px',
                                            padding: '0 16px',
                                            backgroundColor: '#f8f9fc',
                                            border: '1px solid #e0e0e0',
                                            borderRadius: '6px',
                                            cursor: 'pointer',
                                            fontSize: '14px',
                                            fontWeight: '500',
                                            transition: 'all 0.2s ease',
                                            '&:hover': {
                                                borderColor: '#ccc',
                                                backgroundColor: '#ffffff',
                                                boxShadow: '0 2px 6px rgba(0,0,0,0.05)',
                                            }
                                        }}
                                    >
                                        <Flex alignItems="center">
                                            <Text fontWeight="500" color="#333">
                                                {selectedEnvironment.charAt(0).toUpperCase() + selectedEnvironment.slice(1)}
                                            </Text>
                                            {environments.find(env => env.name === selectedEnvironment)?.key && (
                                                <Box
                                                    onClick={(e) => copyToClipboard(environments.find(env => env.name === selectedEnvironment)?.key || '', e)}
                                                    sx={{
                                                        display: 'flex',
                                                        alignItems: 'center',
                                                        padding: '2px 8px',
                                                        borderRadius: '4px',
                                                        backgroundColor: '#f3f4f6',
                                                        cursor: 'pointer',
                                                        ml: 2,
                                                        '&:hover': {
                                                            backgroundColor: '#e5e7eb'
                                                        }
                                                    }}
                                                >
                                                    <FaCopy size={12} color="#6b7280" />
                                                    <Text ml={1} fontSize="12px" color="#6b7280">
                                                        {environments.find(env => env.name === selectedEnvironment)?.key.slice(0, 8)}...
                                                    </Text>
                                                </Box>
                                            )}
                                        </Flex>
                                        <Box
                                            sx={{
                                                width: '14px',
                                                height: '14px',
                                                display: 'flex',
                                                alignItems: 'center',
                                                justifyContent: 'center',
                                                color: '#9ca3af',
                                                transform: dropdownOpen ? 'rotate(180deg)' : 'rotate(0deg)',
                                                transition: 'transform 0.2s ease',
                                                '&::before': {
                                                    content: '""',
                                                    display: 'block',
                                                    width: '6px',
                                                    height: '6px',
                                                    borderRight: '1.5px solid',
                                                    borderBottom: '1.5px solid',
                                                    transform: 'rotate(45deg) translateY(-1px)',
                                                }
                                            }}
                                        />
                                    </Flex>

                                    {dropdownOpen && (
                                        <Box
                                            sx={{
                                                position: 'absolute',
                                                top: 'calc(100% + 4px)',
                                                left: 0,
                                                width: '100%',
                                                backgroundColor: '#ffffff',
                                                borderRadius: '6px',
                                                boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
                                                border: '1px solid #eaeaea',
                                                zIndex: 10,
                                                maxHeight: '200px',
                                                overflowY: 'auto',
                                                '&::-webkit-scrollbar': {
                                                    width: '6px',
                                                },
                                                '&::-webkit-scrollbar-track': {
                                                    background: '#f1f1f1',
                                                    borderRadius: '3px',
                                                },
                                                '&::-webkit-scrollbar-thumb': {
                                                    background: '#ccc',
                                                    borderRadius: '3px',
                                                }
                                            }}
                                        >
                                            {environments.map(env => (
                                                <Box
                                                    key={env.id}
                                                    onClick={() => {
                                                        setSelectedEnvironment(env.name);
                                                        setSelectedEnvironmentId(env.id);
                                                        setDropdownOpen(false);
                                                    }}
                                                    sx={{
                                                        padding: '10px 16px',
                                                        cursor: 'pointer',
                                                        transition: 'background-color 0.15s ease',
                                                        backgroundColor: selectedEnvironment === env.name ? '#f0f5ff' : 'transparent',
                                                        '&:hover': {
                                                            backgroundColor: '#f8f9fc',
                                                        },
                                                        '&:not(:last-child)': {
                                                            borderBottom: '1px solid #f0f0f0'
                                                        }
                                                    }}
                                                >
                                                    <Flex alignItems="center" justifyContent="space-between">
                                                        <Flex alignItems="center">
                                                            <Box mr={2}>
                                                                {env.name === 'production' ? (
                                                                    <Box as="span" sx={{ color: '#34C363' }}>🚀</Box>
                                                                ) : (
                                                                    <Box as="span" sx={{ color: '#5569ff' }}>🔧</Box>
                                                                )}
                                                            </Box>
                                                            <Text fontWeight={selectedEnvironment === env.name ? 'bold' : 'normal'} color="#333">
                                                                {env.name.charAt(0).toUpperCase() + env.name.slice(1)}
                                                            </Text>
                                                        </Flex>
                                                        <Box
                                                            onClick={(e) => copyToClipboard(env.key, e)}
                                                            sx={{
                                                                display: 'flex',
                                                                alignItems: 'center',
                                                                padding: '4px 8px',
                                                                borderRadius: '4px',
                                                                backgroundColor: '#f3f4f6',
                                                                cursor: 'pointer',
                                                                '&:hover': {
                                                                    backgroundColor: '#e5e7eb'
                                                                }
                                                            }}
                                                        >
                                                            <FaCopy size={12} color="#6b7280" />
                                                            <Text ml={1} fontSize="12px" color="#6b7280">
                                                                {env.key.slice(0, 8)}...
                                                            </Text>
                                                        </Box>
                                                    </Flex>
                                                </Box>
                                            ))}
                                        </Box>
                                    )}
                                </Box>
                            </Box>
                        )}
                        <Box>
                            <Button onClick={handleCreateEnvironment} disabled={isLoading}>
                                <Flex alignItems="center">
                                    <Box as="span" mr={1}>+</Box>
                                    Create Environment
                                </Flex>
                            </Button>
                        </Box>
                    </Flex>
                </Flex>

                {!isLoading && !error && environments.length === 0 ? (
                    <ZeroStateEnvironments />
                ) : versionsLoading ? (
                    <Box mt={4}>
                        <Box sx={{
                            border: '1px solid #e0e0e0',
                            borderRadius: '8px',
                            overflow: 'hidden'
                        }}>
                            <Flex
                                sx={{
                                    borderBottom: '1px solid #e0e0e0',
                                    bg: '#f7f7f7',
                                    p: 3
                                }}
                            >
                                <Box width={1 / 3}>
                                    <Text fontWeight="bold">Version</Text>
                                </Box>
                                <Box width={1 / 3}>
                                    <Text fontWeight="bold">Environment</Text>
                                </Box>
                                <Box width={1 / 3}>
                                    <Text fontWeight="bold">Created At</Text>
                                </Box>
                            </Flex>

                            {[1, 2].map((i) => (
                                <Flex
                                    key={i}
                                    sx={{
                                        p: 3,
                                        borderBottom: '1px solid #e0e0e0',
                                        '&:last-child': {
                                            borderBottom: 'none'
                                        }
                                    }}
                                >
                                    <Box width={1 / 3}>
                                        <Box sx={{
                                            width: '100px',
                                            height: '20px',
                                            bg: '#f0f0f0',
                                            borderRadius: '4px',
                                            animation: 'pulse 1.5s infinite'
                                        }} />
                                    </Box>
                                    <Box width={1 / 3}>
                                        <Box sx={{
                                            width: '120px',
                                            height: '28px',
                                            bg: '#f0f0f0',
                                            borderRadius: '4px',
                                            animation: 'pulse 1.5s infinite'
                                        }} />
                                    </Box>
                                    <Box width={1 / 3}>
                                        <Box sx={{
                                            width: '150px',
                                            height: '20px',
                                            bg: '#f0f0f0',
                                            borderRadius: '4px',
                                            animation: 'pulse 1.5s infinite'
                                        }} />
                                    </Box>
                                </Flex>
                            ))}
                        </Box>
                    </Box>
                ) : versionsError ? (
                    <ZeroStateVersions />
                ) : versions.length === 0 ? (
                    <ZeroStateVersions />
                ) : (
                    <Box mt={4}>
                        <Box sx={{
                            border: '1px solid #e0e0e0',
                            borderRadius: '8px',
                            overflow: 'hidden'
                        }}>
                            <Flex
                                sx={{
                                    borderBottom: '1px solid #e0e0e0',
                                    bg: '#f7f7f7',
                                    p: 3
                                }}
                            >
                                <Box width={1 / 3}>
                                    <Text fontWeight="bold">Version</Text>
                                </Box>
                                <Box width={1 / 3}>
                                    <Text fontWeight="bold">Environment</Text>
                                </Box>
                                <Box width={1 / 3}>
                                    <Text fontWeight="bold">Created At</Text>
                                </Box>
                            </Flex>

                            {versions.map((version) => (
                                <Flex
                                    key={version.id}
                                    sx={{
                                        p: 3,
                                        borderBottom: '1px solid #e0e0e0',
                                        '&:last-child': {
                                            borderBottom: 'none'
                                        },
                                        '&:hover': {
                                            bg: '#f9f9f9',
                                            cursor: 'pointer'
                                        }
                                    }}
                                    onClick={() => {
                                        navigate(`/dashboard/version/bundle/${version.id}`)
                                    }}
                                >
                                    <Box width={1 / 3}>
                                        <Text>{version.appVersion || formatVersionNumber(version.versionNumber)}</Text>
                                    </Box>
                                    <Box width={1 / 3}>
                                        <Text>
                                            <Box
                                                as="span"
                                                sx={{
                                                    display: 'inline-block',
                                                    px: 2,
                                                    py: 1,
                                                    bg: selectedEnvironment === 'production' ? '#f0fff4' : '#f7f7ff',
                                                    color: selectedEnvironment === 'production' ? '#34C363' : '#5569ff',
                                                    borderRadius: '4px',
                                                    fontSize: '14px'
                                                }}
                                            >
                                                {selectedEnvironment.charAt(0).toUpperCase() + selectedEnvironment.slice(1)}
                                            </Box>
                                        </Text>
                                    </Box>
                                    <Box width={1 / 3}>
                                        <Text>{new Date(version.createdAt).toLocaleString('en-US', {
                                            month: 'short',
                                            day: 'numeric',
                                            year: 'numeric',
                                            hour: '2-digit',
                                            minute: '2-digit'
                                        })}</Text>
                                    </Box>
                                </Flex>
                            ))}
                        </Box>
                    </Box>
                )}
            </Box>

            <CreateEnvironmentModal
                isOpen={showCreateEnvModal}
                onClose={handleCloseEnvModal}
                onAdd={handleAddEnvironment}
                appName={appDetails?.name || ''}
            />
        </Box>
    )
}

const ZeroStateEnvironments = () => {
    return (
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
                No environments available
            </Text>
            <Text fontSize={"16px"} color={"#666"} mb={3}>
                Create an environment to start managing your app versions
            </Text>
        </Box>
    )
}

const ZeroStateVersions = () => {
    return (
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
                No versions found
            </Text>
            <Text fontSize={"16px"} color={"#666"} mb={3}>
                Create a version to get started
            </Text>
        </Box>
    )
}

export default VersionView