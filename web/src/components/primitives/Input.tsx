import { Input as RebassInput, InputProps } from '@rebass/forms'
import { Box, Text } from 'rebass/styled-components'

interface Props extends InputProps {
    error?: string;
}

const AppInput = ({ error, ...props }: Props) => {
    return (
        <Box>
            <RebassInput
                sx={{
                    borderRadius: 4,
                    height: 50,
                    borderColor: error ? "#e53e3e" : "#ccc",
                    ...props.style,
                    "&:focus": {
                        outline: "none",
                    },
                }}
                {...props} />
            {error && (
                <Text color="#e53e3e" fontSize="14px" mt={1}>
                    {error}
                </Text>
            )}
        </Box>
    )
}

export default AppInput