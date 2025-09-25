import React, { useState } from 'react';
import { Box, Text, Flex } from 'rebass/styled-components';
import { toast } from 'react-toastify';
import { FaPlus, FaArrowLeft, FaEye, FaEyeSlash, FaCopy } from 'react-icons/fa6';
import { useNavigate } from 'react-router-dom';
import { useAuthKeys } from '../../../api/hooks/useAuthKeys';
import Button from '../../../components/primitives/Button';
import Input from '../../../components/primitives/Input';

interface CreateAuthKeyModalProps {
    isOpen: boolean;
    onClose: () => void;
    onCreate: (name: string) => Promise<void>;
}

const CreateAuthKeyModal: React.FC<CreateAuthKeyModalProps> = ({ isOpen, onClose, onCreate }) => {
    const [name, setName] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!name.trim()) {
            toast.error('Name is required');
            return;
        }

        setLoading(true);
        try {
            await onCreate(name.trim());
            setName('');
            onClose();
        } catch (error) {
            // Error is handled by the parent component
        } finally {
            setLoading(false);
        }
    };

    if (!isOpen) return null;

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
                zIndex: 1000,
            }}
            onClick={onClose}
        >
            <Box
                sx={{
                    backgroundColor: 'white',
                    borderRadius: '8px',
                    padding: '24px',
                    width: '400px',
                    maxWidth: '90vw',
                }}
                onClick={(e) => e.stopPropagation()}
            >
                <Text fontSize="20px" fontWeight="bold" mb={3}>
                    Create New Auth Key
                </Text>

                <form onSubmit={handleSubmit}>
                    <Box mb={3}>
                        <Text fontSize="14px" fontWeight="bold" mb={1}>
                            Name
                        </Text>
                        <Input
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="Enter auth key name"
                        />
                    </Box>

                    <Flex justifyContent="flex-end">
                        <Box mr={2}>
                            <Button
                                type="button"
                                onClick={onClose}
                                sx={{
                                    backgroundColor: '#f5f5f5',
                                    color: '#333',
                                    '&:hover': {
                                        backgroundColor: '#e0e0e0',
                                    },
                                }}
                            >
                                Cancel
                            </Button>
                        </Box>
                        <Button
                            type="submit"
                            disabled={loading}
                        >
                            {loading ? 'Creating...' : 'Create'}
                        </Button>
                    </Flex>
                </form>
            </Box>
        </Box>
    );
};

const AuthKeysView: React.FC = () => {
    const { authKeys, loading, createAuthKey } = useAuthKeys();
    const [showModal, setShowModal] = useState(false);
    const [showKeys, setShowKeys] = useState<Record<string, boolean>>({});
    const navigate = useNavigate();

    const handleCreateAuthKey = async (name: string) => {
        const newKey = await createAuthKey(name);
        if (newKey) {
            toast.success('Auth key created successfully!');
            // Show the key for the newly created auth key
            setShowKeys(prev => ({ ...prev, [newKey]: true }));
        }
    };

    const toggleKeyVisibility = (key: string) => {
        setShowKeys(prev => ({ ...prev, [key]: !prev[key] }));
    };

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
        toast.success('Auth key copied to clipboard!');
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    if (loading) {
        return (
            <Box>
                <Box
                    backgroundColor="#fff"
                    minHeight="90vh"
                    sx={{
                        border: "1px solid #e0e0e0",
                        borderRadius: "10px",
                        padding: "20px",
                        margin: "10px",
                    }}
                >
                    {/* Header Skeleton */}
                    <Flex justifyContent="space-between" alignItems="center" mb={4}>
                        <Flex alignItems="center">
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
                                    animation: 'pulse 1.5s ease-in-out infinite',
                                }}
                            />
                        </Flex>
                        <Box
                            sx={{
                                width: '140px',
                                height: '40px',
                                bg: '#f0f0f0',
                                borderRadius: '4px',
                                animation: 'pulse 1.5s ease-in-out infinite',
                            }}
                        />
                    </Flex>

                    {/* Auth Keys Skeleton */}
                    <Box>
                        {[1, 2, 3].map((i) => (
                            <Box
                                key={i}
                                p={3}
                                mb={2}
                                sx={{
                                    border: '1px solid #e0e0e0',
                                    borderRadius: '8px',
                                    backgroundColor: '#fafafa',
                                    animation: 'pulse 1.5s ease-in-out infinite',
                                    animationDelay: `${i * 0.1}s`,
                                }}
                            >
                                <Flex justifyContent="space-between" alignItems="center" mb={2}>
                                    <Box>
                                        <Box
                                            sx={{
                                                width: '180px',
                                                height: '20px',
                                                bg: '#f0f0f0',
                                                borderRadius: '4px',
                                                mb: 1,
                                                animation: 'pulse 1.5s ease-in-out infinite',
                                            }}
                                        />
                                        <Box
                                            sx={{
                                                width: '220px',
                                                height: '16px',
                                                bg: '#f0f0f0',
                                                borderRadius: '4px',
                                                animation: 'pulse 1.5s ease-in-out infinite',
                                            }}
                                        />
                                    </Box>
                                </Flex>

                                <Flex alignItems="center" mt={2}>
                                    <Box
                                        sx={{
                                            width: '8px',
                                            height: '8px',
                                            borderRadius: '50%',
                                            bg: '#f0f0f0',
                                            mr: 2,
                                            animation: 'pulse 1.5s ease-in-out infinite',
                                        }}
                                    />
                                    <Box
                                        sx={{
                                            width: '60px',
                                            height: '16px',
                                            bg: '#f0f0f0',
                                            borderRadius: '4px',
                                            animation: 'pulse 1.5s ease-in-out infinite',
                                        }}
                                    />
                                </Flex>
                            </Box>
                        ))}
                    </Box>
                </Box>
            </Box>
        );
    }

    return (
        <Box>
            <Box
                backgroundColor="#fff"
                minHeight="90vh"
                sx={{
                    border: "1px solid #e0e0e0",
                    borderRadius: "10px",
                    padding: "20px",
                    margin: "10px",
                }}
            >
                <Flex justifyContent="space-between" alignItems="center" mb={4}>
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
                            onClick={() => navigate('/dashboard')}
                        >
                            <FaArrowLeft size={16} />
                        </Box>
                        <Text fontSize="28px" fontWeight="bold">
                            Auth Keys
                        </Text>
                    </Flex>
                    <Button
                        onClick={() => setShowModal(true)}
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                        }}
                    >
                        <Box mr={1}>
                            <FaPlus size={14} />
                        </Box>
                        Create Auth Key
                    </Button>
                </Flex>

                {authKeys.length === 0 ? (
                    <Box
                        p={4}
                        sx={{
                            textAlign: 'center',
                            backgroundColor: '#fafafa',
                            borderRadius: '8px',
                            border: '1px dashed #e0e0e0',
                        }}
                    >
                        <Text fontSize="18px" fontWeight="bold" mb={2}>
                            No auth keys found
                        </Text>
                        <Text fontSize="16px" color="#666" mb={3}>
                            Create your first auth key to get started
                        </Text>
                        <Button onClick={() => setShowModal(true)}>
                            <Flex alignItems="center">
                                <Box mr={1}>
                                    <FaPlus size={14} />
                                </Box>
                                Create Auth Key
                            </Flex>
                        </Button>
                    </Box>
                ) : (
                    <Box>
                        {authKeys.map((authKey) => (
                            <Box
                                key={authKey.id}
                                p={3}
                                mb={2}
                                sx={{
                                    border: '1px solid #e0e0e0',
                                    borderRadius: '8px',
                                    backgroundColor: '#fafafa',
                                }}
                            >
                                <Flex justifyContent="space-between" alignItems="center" mb={2}>
                                    <Box>
                                        <Text fontWeight="bold" fontSize="16px">
                                            {authKey.name}
                                        </Text>
                                        <Text fontSize="14px" color="#666">
                                            Created by {authKey.createdBy} on {formatDate(authKey.createdAt)}
                                        </Text>
                                    </Box>
                                </Flex>

                                <Flex alignItems="center" mt={2}>
                                    <Box
                                        sx={{
                                            width: '8px',
                                            height: '8px',
                                            borderRadius: '50%',
                                            backgroundColor: authKey.isValid ? '#4caf50' : '#f44336',
                                        }}
                                    />
                                    <Text fontSize="12px" color="#666" ml={2}>
                                        {authKey.isValid ? 'Active' : 'Inactive'}
                                    </Text>
                                </Flex>

                                {/* Auth Key Display with Masking */}
                                <Box mt={3} p={2} sx={{ backgroundColor: '#f8f9fa', borderRadius: '4px' }}>
                                    <Flex alignItems="center" justifyContent="space-between">
                                        <Box flex={1}>
                                            <Text fontSize="12px" color="#666" mb={1}>
                                                Auth Key
                                            </Text>
                                            <Text
                                                fontSize="14px"
                                                fontFamily="monospace"
                                                sx={{
                                                    wordBreak: 'break-all',
                                                    userSelect: showKeys[authKey.key] ? 'text' : 'none'
                                                }}
                                            >
                                                {showKeys[authKey.key] ? authKey.key : '••••••••••••••••••••••••••••••••'}
                                            </Text>
                                        </Box>
                                        <Flex ml={2}>
                                            <Box
                                                as="button"
                                                mr={1}
                                                sx={{
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    justifyContent: 'center',
                                                    height: '32px',
                                                    width: '32px',
                                                    border: '1px solid #e0e0e0',
                                                    borderRadius: '4px',
                                                    bg: 'transparent',
                                                    cursor: 'pointer',
                                                    '&:hover': {
                                                        bg: '#f0f0f0'
                                                    }
                                                }}
                                                onClick={() => toggleKeyVisibility(authKey.key)}
                                                title={showKeys[authKey.key] ? 'Hide key' : 'Show key'}
                                            >
                                                {showKeys[authKey.key] ? <FaEyeSlash size={14} /> : <FaEye size={14} />}
                                            </Box>
                                            {showKeys[authKey.key] && (
                                                <Box
                                                    as="button"
                                                    sx={{
                                                        display: 'flex',
                                                        alignItems: 'center',
                                                        justifyContent: 'center',
                                                        height: '32px',
                                                        width: '32px',
                                                        border: '1px solid #e0e0e0',
                                                        borderRadius: '4px',
                                                        bg: 'transparent',
                                                        cursor: 'pointer',
                                                        '&:hover': {
                                                            bg: '#f0f0f0'
                                                        }
                                                    }}
                                                    onClick={() => copyToClipboard(authKey.key)}
                                                    title="Copy to clipboard"
                                                >
                                                    <FaCopy size={14} />
                                                </Box>
                                            )}
                                        </Flex>
                                    </Flex>
                                </Box>
                            </Box>
                        ))}
                    </Box>
                )}
            </Box>

            <CreateAuthKeyModal
                isOpen={showModal}
                onClose={() => setShowModal(false)}
                onCreate={handleCreateAuthKey}
            />
        </Box>
    );
};

export default AuthKeysView; 