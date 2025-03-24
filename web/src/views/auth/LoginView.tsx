import React, { useState, useEffect } from 'react'
import { Box, Text, Image } from 'rebass/styled-components'
import { useNavigate } from 'react-router-dom'
import Input from '../../components/primitives/Input'
import Button from '../../components/primitives/Button'
import SpreadLogo from "../../assets/spread-logo.png"
import { loginSchema, LoginFormValues } from '../../schemas'
import { useAuth } from '../../api'
import { ROUTES } from '../../api/navigation'
import useAuthStore from '../../store/authStore'

const LoginView = () => {
    const navigate = useNavigate()
    const { login } = useAuth()
    // Get loading and error from the global auth store
    const { loading, error } = useAuthStore()
    const [formValues, setFormValues] = useState<LoginFormValues>({
        username: '',
        password: '',
    })
    const [formErrors, setFormErrors] = useState<Record<string, string>>({})

    // Redirect to dashboard if already authenticated
    useEffect(() => {
        // This is a backup check - AuthInitializer should handle this in most cases
        const token = localStorage.getItem('token')
        if (token) {
            navigate(ROUTES.DASHBOARD, { replace: true })
        }
    }, [navigate])

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target
        setFormValues({
            ...formValues,
            [name]: value,
        })

        // Clear error when user starts typing
        if (formErrors[name]) {
            setFormErrors({
                ...formErrors,
                [name]: '',
            })
        }
    }

    const validateForm = (): boolean => {
        try {
            loginSchema.parse(formValues)
            setFormErrors({})
            return true
        } catch (error: any) {
            // Handle zod validation errors
            const validationErrors: Record<string, string> = {}
            error.errors.forEach((err: any) => {
                if (err.path[0]) {
                    validationErrors[err.path[0]] = err.message
                }
            })
            setFormErrors(validationErrors)
            return false
        }
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        // Validate form
        if (!validateForm()) return

        // Call login API
        const response = await login(formValues)

        if (response.success) {
            // Redirect to dashboard on successful login
            navigate(ROUTES.DASHBOARD)
        }
    }

    return (
        <>
            <Box backgroundColor={"#fafafa"} height={"100%"} display={"flex"} justifyContent={"center"} alignItems={"center"} flexDirection={"column"}>
                <Image ml={"-40px"} mb={"10px"} src={SpreadLogo} alt="logo" maxWidth={"13%"} />
                <Box
                    as="form"
                    onSubmit={handleSubmit}
                    p={4}
                    sx={{
                        borderRadius: "10px",
                        border: "1px solid #e0e0e0",
                    }}
                    backgroundColor={"#fff"}
                    width={"400px"}
                    height={"auto"}
                >
                    <Box mb={4}>
                        <Text fontSize={"24px"} fontWeight={"bold"}>Login to your account</Text>
                        <Text fontSize={"16px"} color={"#666"}>Welcome back! Please enter your details</Text>
                    </Box>
                    <Input
                        name="username"
                        placeholder="Username"
                        value={formValues.username}
                        onChange={handleChange}
                        error={formErrors.username}
                    />
                    <Box height={"10px"} />
                    <Input
                        name="password"
                        type="password"
                        placeholder="Password"
                        value={formValues.password}
                        onChange={handleChange}
                        error={formErrors.password}
                    />

                    {error && (
                        <Box mt={2}>
                            <Text color="red" fontSize={14}>{error}</Text>
                        </Box>
                    )}

                    <Box height={"10px"} />
                    <Button
                        type="submit"
                        width={"100%"}
                        mt={2}
                        height={"40px"}
                        disabled={loading}
                    >
                        {loading ? 'Logging in...' : 'Login'}
                    </Button>
                </Box>
                <Box mt={2}>
                    <Text fontSize={"14px"} color={"#999999"}>Baked with 💚 by Swish</Text>
                </Box>
            </Box>
        </>
    )
}

export default LoginView