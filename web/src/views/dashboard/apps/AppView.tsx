import { useState, useEffect } from 'react'
import { Box, Text, Flex } from 'rebass/styled-components'
import Button from '../../../components/primitives/Button'
import AddAppModal, { AddAppFormData } from '../../../components/modal/AddAppModal'
import { IoLogoAndroid, IoLogoApple } from 'react-icons/io5'
import { useNavigate } from 'react-router-dom'
import AppSkeleton from '../../../components/skeleton/AppSkeleton'
import { useApps } from '../../../api/hooks'
import { App, CreateAppRequest } from '../../../api/schemas/appSchemas'
import useAppStore from '../../../store/appStore'

const AppView = () => {
    const navigate = useNavigate()
    const { apps, loading, error, validationErrors, createApp } = useApps()
    const [showModal, setShowModal] = useState(false)
    const [newApp, setNewApp] = useState<AddAppFormData>({ name: '', os: 'iOS' })
    const [isCreating, setIsCreating] = useState(false)
    const [createError, setCreateError] = useState<string | null>(null)

    // Use global app store
    const {
        setApps: setGlobalApps,
        setSelectedApp,
        fetchEnvironments
    } = useAppStore()

    // Sync local apps state with global store
    useEffect(() => {
        if (apps && apps.length > 0) {
            setGlobalApps(apps);
        }
    }, [apps, setGlobalApps]);

    const handleOpenModal = () => {
        setShowModal(true)
        setCreateError(null)
    }

    const handleCloseModal = () => {
        setShowModal(false)
        setNewApp({ name: '', os: 'iOS' })
        setCreateError(null)
    }

    const handleAddApp = async () => {
        if (newApp.name.trim()) {
            setIsCreating(true)
            setCreateError(null)

            try {
                // Get the lowercase OS value and validate it as a proper OS type
                const osValue = newApp.os.toLowerCase();

                // Create request with proper typing
                const appRequest: CreateAppRequest = {
                    appName: newApp.name.trim(),
                    os: osValue === 'ios' ? 'ios' : 'android'
                }

                const success = await createApp(appRequest)

                if (success) {
                    handleCloseModal()
                }
            } catch (err) {
                const errorMessage = err instanceof Error ? err.message : 'Failed to create app'
                setCreateError(errorMessage)
            } finally {
                setIsCreating(false)
            }
        }
    }

    const handleAppClick = async (app: App) => {
        // Save the selected app to global store
        setSelectedApp(app);

        // Fetch environments for this app and save to global store
        await fetchEnvironments(app.id);

        // Navigate to versions view
        navigate(`/dashboard/version?appId=${app.id}`)
    }

    const getOsIcon = (os: string) => {
        const normalizedOs = os.toLowerCase()
        if (normalizedOs === 'ios') {
            return <Box as={IoLogoApple} color="primary" size={20} />
        } else if (normalizedOs === 'android') {
            return <Box as={IoLogoAndroid} color="primary" size={20} />
        }
        return null
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
                <Flex justifyContent={"space-between"} flexDirection={"row"}>
                    <Box flexDirection={"column"}>
                        <Text fontSize={"28px"} fontWeight={"bold"}>Hello, Saran✨</Text>
                    </Box>
                    <Box>
                        <Button onClick={handleOpenModal} disabled={loading}>
                            <Flex alignItems="center">
                                <Box as="span" mr={1}>+</Box>
                                Add App
                            </Flex>
                        </Button>
                    </Box>
                </Flex>

                {error && (
                    <Box p={3} bg="#ffe0e0" color="#d32f2f" mt={3} sx={{ borderRadius: 2 }}>
                        <Text>{error}</Text>
                    </Box>
                )}

                {loading ? (
                    <AppSkeleton />
                ) : apps.length === 0 ? (
                    <ZeroStateApps onAddClick={handleOpenModal} />
                ) : (
                    <Box>
                        <Flex flexDirection={"row"} flexWrap={"wrap"} mt={4}>
                            {apps.map((app: App) => (
                                <Box
                                    key={app.id}
                                    p={3}
                                    m={2}
                                    onClick={() => handleAppClick(app)}
                                    sx={{
                                        border: '1px solid #e0e0e0',
                                        borderRadius: '8px',
                                        width: ['100%', '45%', '30%'],
                                        transition: 'border-color 0.8s ease-in-out',
                                        position: 'relative',
                                        '&:hover': {
                                            borderColor: 'primary',
                                            backgroundColor: '#fafffb',
                                            cursor: 'pointer'
                                        }
                                    }}
                                >
                                    <Flex alignItems="center">
                                        <Box
                                            mr={3}
                                            p={2}
                                            sx={{
                                                borderRadius: '50%',
                                                bg: app.os.toLowerCase() === 'ios' ? '#f8f9fa' : '#f0fff4',
                                                border: '1px solid',
                                                borderColor: app.os.toLowerCase() === 'ios' ? '#dee2e6' : '#d1fae5',
                                                display: 'flex',
                                                alignItems: 'center',
                                                justifyContent: 'center',
                                                cursor: 'pointer'
                                            }}
                                        >
                                            {getOsIcon(app.os)}
                                        </Box>
                                        <Text fontWeight="bold">{app.name}</Text>
                                    </Flex>
                                </Box>
                            ))}
                        </Flex>
                    </Box>
                )}
            </Box>

            <AddAppModal
                isOpen={showModal}
                appData={newApp}
                onChangeAppData={setNewApp}
                onAdd={handleAddApp}
                onClose={handleCloseModal}
                isLoading={isCreating}
                error={createError}
                validationErrors={validationErrors}
            />
        </Box>
    )
}

interface ZeroStateProps {
    onAddClick: () => void
}

const ZeroStateApps = ({ onAddClick }: ZeroStateProps) => {
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
                No apps found
            </Text>
            <Text fontSize={"16px"} color={"#666"} mb={3}>
                Create your first app to get started
            </Text>
            <Button onClick={onAddClick}>
                <Flex alignItems="center">
                    <Box as="span" mr={1}>+</Box>
                    Create App
                </Flex>
            </Button>
        </Box>
    )
}

export default AppView