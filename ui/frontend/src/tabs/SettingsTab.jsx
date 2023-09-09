import { Select, Stack, Text } from "@mantine/core";


export default function SettingsTab(props) {
    return (
        <Stack>
            <Stack spacing={2}>
                <Text size="xs" tt="uppercase" c="dimmed">Appearance</Text>
                <Select
                    withinPortal 
                    shadow="lg" 
                    label="Theme" 
                    clearable={false}
                    data={[
                        { value: "dark", label: "Dark" },
                        { value: "light", label: "Light" },
                    ]}
                    value={props.colorScheme}
                    onChange={props.setColorScheme}
                />
            </Stack>
        </Stack>
    )
}