

import { BeatLoader } from 'react-spinners'
import { Button, ButtonProps } from 'rebass/styled-components'

interface Props extends ButtonProps {
    loading?: boolean
}

const AppButton = ({ loading = false, children, ...props }: Props) => {
    return (
        <Button
            sx={{
                transition: "0.5s",
                fontWeight: [600, 600, 400],
                "&:hover": {
                    background: "#50D97D",
                    cursor: "pointer"
                },
                "&:disabled": {
                    backgroundColor: '#67c787'
                },
                ...props.style
            }}
            {...props}>
            {loading ? <BeatLoader color='white' size={"10px"} /> :
                children
            }
        </Button>
    )
}

export default AppButton