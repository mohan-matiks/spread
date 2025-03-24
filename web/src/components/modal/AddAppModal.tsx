import React from 'react'
import { Box, Text, Button, Flex, Heading } from 'rebass/styled-components'
import { Label, Radio } from '@rebass/forms'
import { IoLogoApple, IoLogoAndroid, IoClose } from 'react-icons/io5'
import Input from '../primitives/Input'

export type AppOS = 'iOS' | 'Android'

export interface AddAppFormData {
    name: string
    os: AppOS
}

interface AddAppModalProps {
    isOpen: boolean
    appData: AddAppFormData
    onChangeAppData: (data: AddAppFormData) => void
    onAdd: () => void
    onClose: () => void
    isLoading?: boolean
    error?: string | null
    validationErrors?: Record<string, string> | null
}

const AddAppModal: React.FC<AddAppModalProps> = ({
    isOpen,
    appData,
    onChangeAppData,
    onAdd,
    onClose,
    isLoading = false,
    error = null,
    validationErrors = null
}) => {
    if (!isOpen) return null

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        onChangeAppData({
            ...appData,
            name: e.target.value
        })
    }

    const handleOsChange = (os: AppOS) => {
        onChangeAppData({
            ...appData,
            os
        })
    }

    // Check for field-specific validation errors
    const getFieldError = (field: string): string | null => {
        if (!validationErrors) return null;
        return validationErrors[field] || null;
    };

    const nameError = getFieldError('appName');
    const osError = getFieldError('os');
    const formError = error || getFieldError('_form');

    return (
        <Box
            sx={{
                position: 'fixed',
                top: 0,
                left: 0,
                right: 0,
                bottom: 0,
                bg: 'rgba(0, 0, 0, 0.5)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                zIndex: 1000
            }}
            onClick={onClose}
        >
            <Box
                sx={{
                    bg: 'white',
                    borderRadius: '8px',
                    width: '100%',
                    maxWidth: '400px',
                    p: 4,
                    position: 'relative'
                }}
                onClick={(e: React.MouseEvent) => e.stopPropagation()}
            >
                <Flex justifyContent="space-between" alignItems="center" mb={3}>
                    <Heading fontSize={3}>Add New App</Heading>
                    <Box
                        as={IoClose}
                        size={24}
                        onClick={onClose}
                        sx={{
                            cursor: 'pointer',
                            color: 'gray',
                        }}
                    />
                </Flex>

                <Box mb={3}>
                    <Label htmlFor="appName" mb={"3px"} fontSize={"14px"}>App Name</Label>
                    <Input
                        id="appName"
                        value={appData.name}
                        onChange={handleInputChange}
                        placeholder="Enter app name"
                        disabled={isLoading}
                        sx={{
                            borderColor: nameError ? '#d32f2f' : 'lightgray',
                            borderRadius: '4px',
                            p: 2,
                            '::placeholder': { color: '#aaa' },
                            "&:focus": {
                                outline: "none",
                            },
                        }}
                    />
                    {nameError && (
                        <Text color="#d32f2f" fontSize="12px" mt={1}>
                            {nameError}
                        </Text>
                    )}
                </Box>

                {formError && (
                    <Box
                        mb={3}
                        p={2}
                        sx={{
                            backgroundColor: '#ffebee',
                            color: '#c62828',
                            borderRadius: '4px',
                            fontSize: '14px'
                        }}
                    >
                        {formError}
                    </Box>
                )}

                <Box mb={4}>
                    <Label mb={"3px"} fontSize={"14px"}>OS</Label>
                    <Flex>
                        <Label
                            sx={{
                                display: 'flex',
                                alignItems: 'center',
                                cursor: isLoading ? 'not-allowed' : 'pointer',
                                mr: 4,
                                opacity: isLoading ? 0.7 : 1
                            }}
                        >
                            <Radio
                                name="os"
                                checked={appData.os === 'iOS'}
                                onChange={() => !isLoading && handleOsChange('iOS')}
                                mr={2}
                                disabled={isLoading}
                            />
                            <Box as={IoLogoApple} color={"primary"} size={20} mr={1} />
                            iOS
                        </Label>

                        <Label
                            sx={{
                                display: 'flex',
                                alignItems: 'center',
                                cursor: isLoading ? 'not-allowed' : 'pointer',
                                opacity: isLoading ? 0.7 : 1
                            }}
                        >
                            <Radio
                                name="os"
                                checked={appData.os === 'Android'}
                                onChange={() => !isLoading && handleOsChange('Android')}
                                mr={2}
                                disabled={isLoading}
                            />
                            <Box as={IoLogoAndroid} color={"primary"} size={18} mr={1} />
                            Android
                        </Label>
                    </Flex>
                    {osError && (
                        <Text color="#d32f2f" fontSize="12px" mt={1}>
                            {osError}
                        </Text>
                    )}
                </Box>

                <Flex justifyContent="flex-end">
                    <Button
                        variant="outline"
                        onClick={onClose}
                        mr={2}
                        disabled={isLoading}
                        sx={{
                            px: 3,
                            py: 2,
                            borderRadius: '4px',
                            opacity: isLoading ? 0.7 : 1
                        }}
                    >
                        Cancel
                    </Button>
                    <Button
                        onClick={onAdd}
                        disabled={!appData.name.trim() || isLoading}
                        sx={{
                            bg: 'primary',
                            color: 'white',
                            px: 3,
                            py: 2,
                            borderRadius: '4px',
                            cursor: (appData.name.trim() && !isLoading) ? 'pointer' : 'not-allowed',
                            opacity: (appData.name.trim() && !isLoading) ? 1 : 0.7,
                            position: 'relative'
                        }}
                    >
                        {isLoading ? 'Creating...' : 'Add App'}
                    </Button>
                </Flex>
            </Box>
        </Box>
    )
}

export default AddAppModal 