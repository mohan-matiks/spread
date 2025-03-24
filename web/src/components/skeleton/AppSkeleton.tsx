import React from 'react'
import { Box, Flex } from 'rebass/styled-components'

type AppSkeletonProps = {
    count?: number
}

const AppSkeleton: React.FC<AppSkeletonProps> = ({ count = 4 }) => {
    const skeletons = Array(count).fill(0)

    return (
        <Box mt={4}>
            <Flex flexDirection={"row"} flexWrap={"wrap"}>
                {skeletons.map((_, index) => (
                    <Box
                        key={index}
                        p={3}
                        m={2}
                        sx={{
                            border: '1px solid #e8e8e8',
                            borderRadius: '8px',
                            width: ['100%', '45%', '30%'],
                            position: 'relative',
                            background: '#fafafa',
                            animation: 'pulse 1.8s infinite ease-in-out',
                            '@keyframes pulse': {
                                '0%': { opacity: 0.5 },
                                '50%': { opacity: 0.7 },
                                '100%': { opacity: 0.5 },
                            }
                        }}
                    >
                        <Flex alignItems="center">
                            <Box
                                mr={3}
                                p={2}
                                sx={{
                                    borderRadius: '50%',
                                    bg: '#efefef',
                                    border: '1px solid #e0e0e0',
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'center',
                                    width: '36px',
                                    height: '36px'
                                }}
                            />
                            <Box
                                sx={{
                                    bg: '#efefef',
                                    height: '20px',
                                    width: '150px',
                                    borderRadius: '4px'
                                }}
                            />
                        </Flex>
                    </Box>
                ))}
            </Flex>
        </Box>
    )
}

export default AppSkeleton 