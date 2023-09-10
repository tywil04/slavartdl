import { Textarea, Stack, Checkbox, SimpleGrid, Select, Card, Text, NumberInput, Button } from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { useInputState } from "@mantine/hooks"
import FSInput from "../components/FSInput.jsx";
import { GetAllowedHosts } from '../../wailsjs/go/main/SlavartdlUI.js';
import { useRef } from 'preact/hooks';
import { HiMiniCheck, HiMiniExclamationTriangle } from 'react-icons/hi2';


export default function DownloadTab(props) {
    const [ skipUnzip, setSkipUnzip ] = useInputState(false)
    const [ ignoreCover, setIgnoreCover ] = useInputState(false)
    const [ ignoreSubDirs, setIgnoreSubDirs ] = useInputState(false)
    const [ outputDir, setOutputDir ] = useInputState("")
    const [ quality, setQuality ] = useInputState(0)
    const [ timeout, setTimeout ] = useInputState(120)
    const [ cooldown, setCooldown ] = useInputState(0)
    const [ urls, setUrls ] = useInputState("")
    const textareaRef = useRef()


    const handleStartJob = async () => {
        let errorWhileFiltering = false
        const allowed = await GetAllowedHosts()
        const rawUrls = urls.split("\n")
        const trimmedUrls = rawUrls.map((url) => url.trim())
        const filteredUrls = trimmedUrls.filter((url) => {
            if (url.length === 0) {
                // remove this "url" because its actually an empty line, no need to throw error though
                return false
            }

            try {
                const parsedUrl = new URL(url)
                if (!allowed.includes(parsedUrl.hostname)) {
                    errorWhileFiltering = true
                    notifications.show({
                        title: "Error",
                        message: "A URL with an unknown host has been found, aborting.",
                        color: "red",
                        icon: <HiMiniExclamationTriangle size="16"/>,
                    })
                    return false
                }
            } catch {
                errorWhileFiltering = true
                notifications.show({
                    title: "Error",
                    message: "A non-URL was found, aborting.",
                    color: "red",
                    icon: <HiMiniExclamationTriangle size="16"/>,
                })
                return false
            }

            // no errors were encountered, so include url
            return true
        })

        if (errorWhileFiltering) {
            // if there was a significant error while filtering, the message would have been shown we just need to finish this function call
            return
        }

        if (filteredUrls.length === 0) {
            setUrls("")
            textareaRef.current.focus()
            notifications.show({
                title: "Error",
                message: "No URLs found.",
                color: "red",
                icon: <HiMiniExclamationTriangle size="16"/>,
            })
            return
        }

        setOutputDir(outputDir.trim())
        if (outputDir === "") {
            notifications.show({
                title: "Error",
                message: "No output directory found.",
                color: "red",
                icon: <HiMiniExclamationTriangle size="16"/>,
            })
            return
        }

        props.queueHandlers.appendBatch(filteredUrls, {
            skipUnzip,
            ignoreCover,
            ignoreSubDirs,
            outputDir,
            quality,
            timeout,
            cooldown,
        })

        notifications.show({
            title: "Success",
            message: "Successfully added URLs to download queue.",
            color: "blue",
            icon: <HiMiniCheck size="16"/>,
        })

        setUrls("")
    }

    const handleCheckbox = (value, setter) => {
        return (event) => {
            event.preventDefault()
            setter(!value)
        }
    }

    const handleNumberInputStepperHold = (time) => Math.max(1000 / time ** 2, 25)
    

    const addURLsToQueueButtonStyle = (theme) => ({
        position: "absolute",
        right: theme.globalStyles(theme)["#root"].padding,
        marginTop: 6
    })

    const checkboxCardStyle = (theme) => ({
        cursor: "pointer",
        backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[7] : theme.colors.gray[0],
        border: `1px solid ${theme.colorScheme === "dark" ? theme.colors.dark[6] : theme.colors.gray[2]}`,
        borderRadius: 12
    })

    const cardStyle = (theme) => ({
        backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[7] : theme.colors.gray[0],
        border: `1px solid ${theme.colorScheme === "dark" ? theme.colors.dark[6] : theme.colors.gray[2]}`,
        borderRadius: 12
    })

    const inputStyle = (theme) => ({ marginTop: -5 })


    return (
        <Stack spacing="xl">
            <Textarea 
                ref={textareaRef}
                required
                minRows={3} 
                label={
                    <>
                        URLs
                        <Button onClick={handleStartJob} sx={addURLsToQueueButtonStyle}>Add URLs to Queue</Button>
                    </>
                }
                description="Seperate URLs with newlines."
                value={urls}
                onChange={setUrls}
            />

            <Stack spacing={2}>
                <Text size="xs" tt="uppercase" c="dimmed">Bool Options</Text>

                <SimpleGrid cols={3} spacing="lg">
                    <Card sx={checkboxCardStyle} onClick={handleCheckbox(skipUnzip, setSkipUnzip)}>
                        <Checkbox 
                            label="Skip Unzipping" 
                            description="Skip unzipping the downloaded zip file into the output directory."
                            checked={skipUnzip}
                        />
                    </Card>

                    <Card sx={checkboxCardStyle} onClick={handleCheckbox(ignoreCover, setIgnoreCover)}>
                        <Checkbox 
                            label="Ignore Cover Image" 
                            description="If unzipping, should the cover image be ignored or extracted."
                            checked={ignoreCover}
                        />
                    </Card>

                    <Card sx={checkboxCardStyle} onClick={handleCheckbox(ignoreSubDirs, setIgnoreSubDirs)}>
                        <Checkbox 
                            label="Ignore Sub Directories" 
                            description="If unzipping, should sub directories be ignored or extracted." 
                            checked={ignoreSubDirs} 
                        />
                    </Card>
                </SimpleGrid>
            </Stack>

            <Stack spacing={2}>
                <Text size="xs" tt="uppercase" c="dimmed">Options</Text>

                <SimpleGrid cols={2} spacing="lg">
                    <Card sx={cardStyle}>
                        <FSInput 
                            sx={inputStyle}
                            required
                            func="openDirectory"
                            funcData="Output Directory"
                            label="Output Directory" 
                            description="The directory to output the resulting file structure/zip archive."
                            placeholder="e.g. /path/to/directory"
                            value={outputDir}
                            onChange={setOutputDir}
                        />
                    </Card>

                    <Card sx={cardStyle}>
                        <Select
                            sx={inputStyle}
                            required
                            withinPortal 
                            shadow="lg" 
                            label="Quality" 
                            description="The quality of the music that should be downloaded."
                            clearable={false}
                            data={[
                                { value: 0, label: "Best Quality Available" },
                                { value: 1, label: "128kbps MP3/AAC" },
                                { value: 2, label: "320kbps MP3/AAC" },
                                { value: 3, label: "16bit 44.1kHz" },
                                { value: 4, label: "24bit ≤96kHz" },
                                { value: 5, label: "24bit ≤192kHz" },
                            ]}
                            value={quality}
                            onChange={setQuality}
                        />
                    </Card>

                    <Card sx={cardStyle}>
                        <NumberInput
                            sx={inputStyle}
                            required
                            label="Timeout" 
                            description="The number of seconds before a download times out."
                            stepHoldDelay={500}
                            stepHoldInterval={handleNumberInputStepperHold}
                            min={0}
                            value={timeout}
                            onChange={setTimeout}
                        />
                    </Card>

                    <Card sx={cardStyle}>
                        <NumberInput
                            sx={inputStyle}
                            required
                            label="Cooldown" 
                            description="The number of seconds to wait before starting the next download."
                            stepHoldDelay={500}
                            stepHoldInterval={handleNumberInputStepperHold}
                            min={0}
                            value={cooldown}
                            onChange={setCooldown}
                        />
                    </Card>
                </SimpleGrid>
            </Stack>
        </Stack>
    )
}