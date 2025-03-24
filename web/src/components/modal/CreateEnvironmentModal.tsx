import React, { useState } from 'react'
import { Box, Text, Flex } from 'rebass/styled-components'
import Button from '../primitives/Button'
import { IoMdClose } from 'react-icons/io'
import Input from '../primitives/Input'
import { z } from 'zod'

const environmentSchema = z.object({
    environmentName: z.string().min(1, 'Environment name is required'),
    appName: z.string().min(1, 'App name is required')
})

export type CreateEnvironmentFormData = z.infer<typeof environmentSchema>

interface CreateEnvironmentModalProps {
    isOpen: boolean
    onClose: () => void
    onAdd: (data: CreateEnvironmentFormData) => void
    appName: string
}

const CreateEnvironmentModal = ({ isOpen, onClose, onAdd, appName }: CreateEnvironmentModalProps) => {
    const [formData, setFormData] = useState<CreateEnvironmentFormData>({
        environmentName: '',
        appName: appName
    })
    const [errors, setErrors] = useState<Record<string, string>>({})

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({ ...formData, environmentName: e.target.value })
        setErrors({})
    }

    const handleSubmit = () => {
        try {
            formData.appName = appName // Ensure app name is set correctly
            environmentSchema.parse(formData)
            onAdd(formData)
            setFormData({ environmentName: '', appName })
        } catch (error) {
            if (error instanceof z.ZodError) {
                const fieldErrors: Record<string, string> = {}
                error.errors.forEach(err => {
                    const field = err.path[0] as string
                    fieldErrors[field] = err.message
                })
                setErrors(fieldErrors)
            }
        }
    }

    if (!isOpen) return null

    return (
        <Box
            sx={{
                position: 'fixed',
                top: 0,
                left: 0,
                right: 0,
                bottom: 0,
                backgroundColor: 'rgba(0, 0, 0, 0.5)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                zIndex: 1000
            }}
        >
            <Box
                sx={{
                    backgroundColor: 'white',
                    borderRadius: '8px',
                    width: '400px',
                    maxWidth: '90%',
                    boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)'
                }}
            >
                <Flex justifyContent="space-between" p={3} sx={{ borderBottom: '1px solid #e0e0e0' }}>
                    <Text fontWeight="bold" fontSize={18}>Create New Environment</Text>
                    <Box onClick={onClose} sx={{ cursor: 'pointer' }}>
                        <IoMdClose size={24} />
                    </Box>
                </Flex>

                <Box p={3}>
                    <Box mb={3}>
                        <Text mb={"3px"} fontSize={"14px"}>Name</Text>
                        <Input
                            type="text"
                            placeholder="e.g. staging, testing"
                            value={formData.environmentName}
                            onChange={handleChange}
                            sx={{
                                width: '100%',
                                padding: '10px',
                                border: errors.environmentName ? '1px solid #ff4d4f' : '1px solid #e0e0e0',
                                borderRadius: '4px',
                                '&:focus': {
                                    outline: 'none',
                                }
                            }}
                        />
                        {errors.environmentName && <Text color="#ff4d4f" fontSize="14px" mt={1}>{errors.environmentName}</Text>}
                    </Box>

                    <Flex justifyContent="flex-end" mt={4}>
                        <Button
                            onClick={onClose}
                            mr={2}
                            sx={{
                                bg: 'transparent',
                                color: '#666',
                                border: '1px solid #e0e0e0',
                                borderRadius: '4px',
                                padding: '8px 16px',
                                cursor: 'pointer',
                                '&:hover': {
                                    bg: '#f5f5f5'
                                }
                            }}
                        >
                            Cancel
                        </Button>
                        <Button onClick={handleSubmit}>Create</Button>
                    </Flex>
                </Box>
            </Box>
        </Box>
    )
}

export default CreateEnvironmentModal 