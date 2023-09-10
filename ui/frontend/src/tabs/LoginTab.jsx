import { PasswordInput, Stack, Text, TextInput, Button } from "@mantine/core";
import { useInputState } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { Login } from "../../wailsjs/go/main/SlavartdlUI";
import { HiMiniCheck, HiMiniExclamationTriangle } from "react-icons/hi2";


export default function LoginTab() {
    const [ email, setEmail ] = useInputState("")
    const [ password, setPassword ] = useInputState("")


    const handleLogin = () => {
        Login(email, password).then((success) => {
            if (success) {
                notifications.show({
                    title: "Success",
                    message: "Successfully logged into your divolt account.",
                    color: "blue",
                    icon: <HiMiniCheck size="16"/>,
                })
            } else {
                notifications.show({
                    title: "Error",
                    message: "Failed to login to your divolt account.",
                    color: "red",
                    icon: <HiMiniExclamationTriangle size="16"/>,
                })
            }
        })
    }

    
    const buttonStyle = (theme) => ({ 
        width: "fit-content",
        marginTop: 10
    })


    return (
        <Stack>
            <Text>To download using this app, you must login to a divolt account. If the account hasn't joined the Slavart server, it will join it automatically.</Text>
        
            <Stack spacing={2}>
                <Text size="xs" tt="uppercase" c="dimmed">Login</Text>

                <Stack>
                    <TextInput
                        required
                        label="Email" 
                        description="The email for your divolt account."
                        placeholder="e.g. your@email.com"
                        value={email}
                        onChange={setEmail}
                    />

                    <PasswordInput
                        required
                        label="Password" 
                        description="The password for your divolt account."
                        placeholder="e.g. your-super-secure-password"
                        value={password}
                        onChange={setPassword}
                    />

                    <Button compact={false} sx={buttonStyle} onClick={handleLogin}>Login</Button>
                </Stack>
            </Stack>
        </Stack>
    )
}