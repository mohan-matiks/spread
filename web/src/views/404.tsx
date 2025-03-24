import { Box, Text } from 'rebass/styled-components'
import { useNavigate } from 'react-router-dom'

const NotFound = () => {
    const navigate = useNavigate()

    return (
        <Box p={4} textAlign="center">
            <Text fontSize="24px" fontWeight="bold" mb={3}>
                404 - Page Not Found
            </Text>
            <Text fontSize="16px" color="#666" mb={4}>
                The page you're looking for doesn't exist or has been moved.
            </Text>
            <Text
                sx={{
                    cursor: 'pointer',
                    color: '#34C363',
                    textDecoration: 'underline',
                    '&:hover': { color: '#2ba352' }
                }}
                onClick={() => navigate('/')}
            >
                Return to Home
            </Text>
        </Box>
    )
}

export default NotFound