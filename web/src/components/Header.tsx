import { useState } from 'react'
import { FaUser, FaRightFromBracket, FaKey } from 'react-icons/fa6'
import { Flex, Box, Image, Text } from 'rebass/styled-components'
import { useNavigate } from 'react-router-dom'
import SpreadIcon from '../assets/spread-logo.png'
import useAuthStore from '../store/authStore'
import { logout } from '../api/services/authService'

const Header = () => {
    const [showDropdown, setShowDropdown] = useState(false);
    const { user } = useAuthStore();
    const navigate = useNavigate();

    const toggleDropdown = () => {
        setShowDropdown(!showDropdown);
    };

    const handleLogout = () => {
        logout();
        setShowDropdown(false);
    };

    const handleAuthKeys = () => {
        navigate('/dashboard/auth-keys');
        setShowDropdown(false);
    };

    return (
        <Flex sx={{
            borderBottom: "1px solid #e0e0e0",
        }}>
            <Flex flexDirection={"row"} justifyContent={"space-between"} width={"100%"} px={"10px"} py={1}>
                <Box>
                    <Image src={SpreadIcon} alt="app-icon" maxWidth={"140px"} />
                </Box>
                <Box display={"flex"} justifyContent={"center"} alignItems={"center"} sx={{ position: "relative" }}>
                    <Text
                        onClick={toggleDropdown}
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            "&:hover": {
                                cursor: "pointer",
                                color: "#666"
                            }
                        }}
                    >
                        <Box as="span" mb={"-5px"} mr={2}><FaUser size={18} /></Box>
                        Account
                    </Text>

                    {showDropdown && (
                        <Box
                            sx={{
                                position: 'absolute',
                                top: '100%',
                                right: 0,
                                width: '200px',
                                backgroundColor: 'white',
                                boxShadow: '0 4px 8px rgba(0,0,0,0.1)',
                                borderRadius: '4px',
                                padding: '12px',
                                zIndex: 10,
                                marginTop: '8px'
                            }}
                        >
                            <Text fontWeight="bold" mb={2}>
                                {user?.username || 'Account'}
                            </Text>
                            <Text
                                sx={{
                                    color: '#666',
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: '6px',
                                    "&:hover": {
                                        cursor: "pointer",
                                        opacity: 0.7
                                    }
                                }}
                                onClick={handleAuthKeys}
                                mb={1}
                            >
                                <FaKey size={14} /> Auth Keys
                            </Text>
                            <Text
                                sx={{
                                    color: '#666',
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: '6px',
                                    "&:hover": {
                                        cursor: "pointer",
                                        opacity: 0.7
                                    }
                                }}
                                onClick={handleLogout}
                            >
                                <FaRightFromBracket size={14} /> Logout
                            </Text>
                        </Box>
                    )}
                </Box>
            </Flex>
        </Flex>
    )
}

export default Header