import { Stack, Group, Text, UnstyledButton, Collapse, Checkbox, Button, ActionIcon } from "@mantine/core"
import { useDisclosure } from "@mantine/hooks"
import { HiMiniTrash } from "react-icons/hi2"


export default function QueueTab(props) {
    const qualities = [
        "Best Quality Available",
        "128kbps MP3/AAC",
        "320kbps MP3/AAC",
        "16bit 44.1kHz",
        "24bit ≤96kHz",
        "24bit ≤192kHz",
    ]


    const handleCancel = (urlIndex) => {
        return (event) => {
            event.preventDefault()
            props.queueHandlers.remove(urlIndex)
        }
    }


    const buttonStyle = (theme) => ({
        padding: 10,
        paddingLeft: 15,
        paddingRight: 15,
        backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[7] : theme.colors.gray[0],
        border: `1px solid ${theme.colorScheme === "dark" ? theme.colors.dark[6] : theme.colors.gray[2]}`,
        "&:not(:first-of-type)": {
            borderTop: "none"
        },
        "&:first-of-type": {
            borderTopLeftRadius: 6,
            borderTopRightRadius: 6,
        },
        "&:last-of-type": {
            borderBottomLeftRadius: 6,
            borderBottomRightRadius: 6,
        }
    })

    const cancelButtonStyle = (theme) => ({
        pointerEvents: "none",
        opacity: 0,
        borderRadius: 3,
        transitionDuration: "100ms",
        "div>button:has(&):hover &": {
            pointerEvents: "auto",
            opacity: 1,
        },
        "div>button:has(&):hover &:hover": {
            filter: "brightness(0.9)"
        }
    })

    const stackStyle = (theme) => ({ "&>*": { height: 16 }, marginBottom: 5 })
    
    const checkboxStyle = (theme) => ({ pointerEvents: "none" })


    let pendingRows = []
    let startedRows = []
    let finishedRows = []
    for (let i = 0; i < props.queue.length; i++) {
        const url = props.queue[i]

        let stackToUse
        switch (url.state) {
            case "queued": stackToUse = pendingRows; break
            case "started": stackToUse = startedRows; break
            case "finished": stackToUse = finishedRows; break
        }

        if (stackToUse != undefined) {
            const [ buttonOpen, buttonHandlers ] = useDisclosure(false)

            stackToUse.push(
                <UnstyledButton sx={buttonStyle} onClick={buttonHandlers.toggle}>
                    <Group>
                        <Text>{url.url}</Text> 
                        <ActionIcon variant="filled" sx={cancelButtonStyle} ml="auto" color="red" onClick={handleCancel(i)}>
                            <HiMiniTrash size="16"/>
                        </ActionIcon>
                    </Group>
    
                    <Collapse in={buttonOpen} transitionDuration={300} animateOpacity={true}>
                        <Group>
                            <Stack spacing={3} sx={stackStyle}>
                                <Text size="sm" span c="dimmed">Ignore Errors:</Text>
                                <Text size="sm" span c="dimmed">Ignore Cover:</Text>
                                <Text size="sm" span c="dimmed">Ignore Subdirectories:</Text>
                                <Text size="sm" span c="dimmed">Skip Unzipping:</Text>
                                <Text size="sm" span c="dimmed">Skip URL Checking:</Text>
                                <Text size="sm" span c="dimmed">Dry Run:</Text>
                                <Text size="sm" span c="dimmed">Output Directory:</Text>
                                <Text size="sm" span c="dimmed">Quality:</Text>
                                <Text size="sm" span c="dimmed">Timeout:</Text>
                                <Text size="sm" span c="dimmed">Cooldown:</Text>
                            </Stack>
    
                            <Group grow>
                                <Stack spacing={3} sx={stackStyle}>
                                    <Checkbox size="xs" checked={url.settings.ignoreErrs} sx={checkboxStyle}/>
                                    <Checkbox size="xs" checked={url.settings.ignoreCover} sx={checkboxStyle}/>
                                    <Checkbox size="xs" checked={url.settings.ignoreSubDirs} sx={checkboxStyle}/>
                                    <Checkbox size="xs" checked={url.settings.skipUnzip} sx={checkboxStyle}/>
                                    <Checkbox size="xs" checked={url.settings.skipUrlChecking} sx={checkboxStyle}/>
                                    <Checkbox size="xs" checked={url.settings.dryRun} sx={checkboxStyle}/>
                                    <Text size="sm" span>{url.settings.outputDir}</Text>
                                    <Text size="sm" span>{qualities[url.settings.quality]}</Text>
                                    <Text size="sm" span>{url.settings.timeout} Seconds</Text>
                                    <Text size="sm" span>{url.settings.cooldown} Seconds</Text>
                                </Stack>
                            </Group>
                        </Group>
                    </Collapse>
                </UnstyledButton>
            )
        }
    }

    if (startedRows.length === 0) {
        startedRows = <Text>No URLs being downloaded</Text>
    }

    if (pendingRows.length === 0) {
        pendingRows = <Text>No URLs are in the download queue</Text>
    }

    if (finishedRows.length === 0) {
        finishedRows = <Text>No URLs have finished downloaded</Text>
    }


    return (
        <Stack spacing="xl">
            <Stack spacing={2}>
                <Text size="xs" tt="uppercase" c="dimmed">URLs Being Downloaded</Text>

                <Stack spacing={0}>
                    {startedRows}
                </Stack>
            </Stack>

            <Stack spacing={2}>
                <Text size="xs" tt="uppercase" c="dimmed">URLs In Queue</Text>

                <Stack spacing={0}>
                    {pendingRows}
                </Stack>
            </Stack>

            <Stack spacing={2}>
                <Text size="xs" tt="uppercase" c="dimmed">URLs Downloaded Since App Launch</Text>

                <Stack spacing={0}>
                    {finishedRows}
                </Stack>
            </Stack>
        </Stack>
    )
}