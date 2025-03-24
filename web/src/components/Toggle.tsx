import React from 'react'
import { Box, Flex, Text } from 'rebass/styled-components'

interface ToggleProps {
    isActive: boolean
    onChange: () => void
    label?: string
    activeColor?: string
    size?: 'small' | 'medium' | 'large'
    disabled?: boolean
}

const Toggle: React.FC<ToggleProps> = ({
    isActive,
    onChange,
    label,
    activeColor = '#34C363',
    size = 'medium',
    disabled = false
}) => {
    // Size configurations
    const sizes = {
        small: {
            width: 32,
            height: 18,
            knobSize: 14,
            knobOffset: 2,
            fontSize: '12px',
        },
        medium: {
            width: 42,
            height: 22,
            knobSize: 18,
            knobOffset: 2,
            fontSize: '14px',
        },
        large: {
            width: 50,
            height: 26,
            knobSize: 22,
            knobOffset: 2,
            fontSize: '16px',
        }
    };

    const { width, height, knobSize, knobOffset, fontSize } = sizes[size];

    return (
        <Flex
            alignItems="center"
            sx={{
                opacity: disabled ? 0.6 : 1,
                pointerEvents: disabled ? 'none' : 'auto',
            }}
        >
            {label && (
                <Text
                    mr={2}
                    fontSize={fontSize}
                    fontWeight="500"
                    color="#666"
                    sx={{
                        lineHeight: 1,
                    }}
                >
                    {label}
                </Text>
            )}

            <Box
                role="switch"
                aria-checked={isActive}
                tabIndex={0}
                onClick={onChange}
                onKeyDown={(e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                        e.preventDefault();
                        onChange();
                    }
                }}
                sx={{
                    position: 'relative',
                    width,
                    height,
                    backgroundColor: isActive ? activeColor : '#E0E0E0',
                    borderRadius: height / 2,
                    transition: 'background-color 0.2s',
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: isActive ? 'flex-end' : 'flex-start',
                    padding: `0 ${knobOffset}px`,
                }}
            >
                <Box
                    sx={{
                        width: knobSize,
                        height: knobSize,
                        backgroundColor: 'white',
                        borderRadius: '50%',
                        transition: 'transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1)',
                        transform: isActive ? 'scale(1.05)' : 'scale(1)',
                        boxShadow: '0 1px 2px rgba(0,0,0,0.1)',
                    }}
                />
            </Box>
        </Flex>
    )
}

export default Toggle 