import React, { useState } from 'react';
import { Box, Text, Image } from 'rebass/styled-components';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import Input from '../../components/primitives/Input';
import Button from '../../components/primitives/Button';
import SpreadLogo from "../../assets/spread-logo.png";
import { setupService } from '../../api/services/setupService';

interface SetupFormData {
    username: string;
    password: string;
    confirmPassword: string;
}

interface SetupViewProps {
    onUserCreated?: () => void;
}

const SetupView: React.FC<SetupViewProps> = ({ onUserCreated }) => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState<SetupFormData>({
        username: '',
        password: '',
        confirmPassword: '',
    });
    const [loading, setLoading] = useState(false);
    const [formErrors, setFormErrors] = useState<Record<string, string>>({});

    const handleInputChange = (field: keyof SetupFormData) => (
        e: React.ChangeEvent<HTMLInputElement>
    ) => {
        const value = e.target.value;
        setFormData(prev => ({
            ...prev,
            [field]: value,
        }));

        // Clear error when user starts typing
        if (formErrors[field]) {
            setFormErrors(prev => ({
                ...prev,
                [field]: '',
            }));
        }
    };

    const validateForm = (): boolean => {
        const errors: Record<string, string> = {};

        if (!formData.username.trim()) {
            errors.username = 'Username is required';
        }

        if (!formData.password) {
            errors.password = 'Password is required';
        } else if (formData.password.length < 6) {
            errors.password = 'Password must be at least 6 characters long';
        }

        if (!formData.confirmPassword) {
            errors.confirmPassword = 'Please confirm your password';
        } else if (formData.password !== formData.confirmPassword) {
            errors.confirmPassword = 'Passwords do not match';
        }

        setFormErrors(errors);
        return Object.keys(errors).length === 0;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        // Validate form
        if (!validateForm()) return;

        setLoading(true);

        try {
            const response = await setupService.initUser({
                username: formData.username.trim(),
                password: formData.password,
            });

            if (response.data) {
                toast.success('Admin user created successfully! Please log in.');

                // Call the callback to refetch setup status
                if (onUserCreated) {
                    onUserCreated();
                }

                navigate('/login');
            }
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : 'Failed to create admin user';
            toast.error(errorMessage);
        } finally {
            setLoading(false);
        }
    };

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
                        <Text fontSize={"24px"} fontWeight={"bold"}>Welcome to Spread</Text>
                        <Text fontSize={"16px"} color={"#666"}>Let's set up your first admin user to get started</Text>
                    </Box>

                    <Input
                        name="username"
                        placeholder="Username"
                        value={formData.username}
                        onChange={handleInputChange('username')}
                        error={formErrors.username}
                    />
                    <Box height={"10px"} />
                    <Input
                        name="password"
                        type="password"
                        placeholder="Password"
                        value={formData.password}
                        onChange={handleInputChange('password')}
                        error={formErrors.password}
                    />
                    <Box height={"10px"} />
                    <Input
                        name="confirmPassword"
                        type="password"
                        placeholder="Confirm Password"
                        value={formData.confirmPassword}
                        onChange={handleInputChange('confirmPassword')}
                        error={formErrors.confirmPassword}
                    />

                    <Box height={"10px"} />
                    <Button
                        type="submit"
                        width={"100%"}
                        mt={2}
                        height={"40px"}
                        disabled={loading}
                    >
                        {loading ? 'Creating Admin User...' : 'Create Admin User'}
                    </Button>
                </Box>
                <Box mt={2}>
                    <Text fontSize={"14px"} color={"#999999"}>Baked with 💚 by Swish</Text>
                </Box>
            </Box>
        </>
    );
};

export default SetupView; 