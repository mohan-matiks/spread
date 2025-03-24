import { Box } from 'rebass/styled-components'

interface RadioProps {
    checked: boolean
    onChange: () => void
}

const Radio = ({ checked, onChange }: RadioProps) => {
    const color = checked ? "#34C363" : "#b5b5b5"
    return (
        <Box onClick={onChange}>
            <Box
                display="grid"
                alignContent="center"
                justifyContent="center"
                width={20}
                height={20}
                backgroundColor="#fff"
                sx={{
                    borderRadius: "100%",
                    border: `2px solid ${color}`
                }}
            >
                <Box
                    width={16}
                    height={16}
                    backgroundColor={color}
                    sx={{
                        borderRadius: "100%",
                    }}
                />
            </Box>
        </Box>
    )
}

export default Radio