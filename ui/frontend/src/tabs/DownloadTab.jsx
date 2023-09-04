import { useState } from "preact/hooks"
import { Textarea, Stack, Checkbox, SimpleGrid, Select, Card, Text, NumberInput, Button } from '@mantine/core';
import { useForm } from "@mantine/form"
import FSInput from "../components/FSInput.jsx";


export default function DownloadTab(props) {
    const [ ignoreErrs, setIgnoreErrs ] = useState(true)
    const [ skipUnzip, setSkipUnzip ] = useState(false)
    const [ ignoreCover, setIgnoreCover ] = useState(false)
    const [ ignoreSubDirs, setIgnoreSubDirs ] = useState(false)
    const [ skipUrlChecking, setSkipUrlChecking ] = useState(false)
    const [ dryRun, setDryRun ] = useState(false)
    const [ outputDir, setOutputDir ] = useState("")
    const [ quality, setQuality ] = useState(0)
    const [ timeout, setTimeout ] = useState(120)
    const [ cooldown, setCooldown ] = useState(0)


    const handleStartJob = () => {
        const data = [
            "SUBMIT DATA:",
            `- ignoreErrs: ${ignoreErrs}`,
            `- skipUnzip: ${skipUnzip}`,
            `- ignoreCover: ${ignoreCover}`,
            `- ignoreSubDirs: ${ignoreSubDirs}`,
            `- skipUrlChecking: ${skipUrlChecking}`,
            `- dryRun: ${dryRun}`,
            `- outputDir: ${outputDir}`,
            `- quality: ${quality}`,
            `- timeout: ${timeout}`,
            `- cooldown: ${cooldown}`,
        ]
        console.log(data.join("\n"))
    }

    const handleCheckbox = (value, setter) => {
        return (event) => {
            event.preventDefault()
            setter(!value)
        }
    }

    const handleNumberInputStepperHold = (time) => {
        return Math.max(1000 / time ** 2, 25)
    }
    

    const addURLsToQueueButtonTheme = (theme) => ({
        position: "absolute",
        right: theme.globalStyles(theme).body.padding,
        marginTop: 6
    })

    const checkboxCardTheme = (theme) => ({
        cursor: "pointer",
        backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[7] : theme.colors.gray[0],
        borderRadius: 12
    })

    const cardTheme = (theme) => ({
        backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[7] : theme.colors.gray[0],
        borderRadius: 12
    })

    const inputTheme = () => ({
        marginTop: -5,
    })


    return (
        <Stack>
            <Textarea 
                minRows={3} 
                label={
                    <>
                        URLs
                        <Button onClick={handleStartJob} sx={addURLsToQueueButtonTheme}>Start Download Job with URLs</Button>
                    </>
                }
                description="Seperate URLs with newlines."
                required
            />

            <Stack spacing={2}>
                <Text size="xs" tt="uppercase" c="dimmed">Bool Options</Text>

                <SimpleGrid cols={3}>
                    <Card withBorder sx={checkboxCardTheme} onClick={handleCheckbox(ignoreErrs, setIgnoreErrs)}>
                        <Checkbox 
                            label="Ignore Errors" 
                            description="If an error occurs while downloading a URL, should the error be ignored."
                            checked={ignoreErrs}
                        />
                    </Card>

                    <Card withBorder sx={checkboxCardTheme} onClick={handleCheckbox(skipUnzip, setSkipUnzip)}>
                        <Checkbox 
                            label="Skip Unzipping" 
                            description="Skip unzipping the downloaded zip file into the output directory."
                            checked={skipUnzip}
                        />
                    </Card>

                    <Card withBorder sx={checkboxCardTheme} onClick={handleCheckbox(ignoreCover, setIgnoreCover)}>
                        <Checkbox 
                            label="Ignore Cover Image" 
                            description="If unzipping, should the cover image be ignored or extracted."
                            checked={ignoreCover}
                        />
                    </Card>

                    <Card withBorder sx={checkboxCardTheme} onClick={handleCheckbox(ignoreSubDirs, setIgnoreSubDirs)}>
                        <Checkbox 
                            label="Ignore Sub Directories" 
                            description="If unzipping, should sub directories be ignored or extracted." 
                            checked={ignoreSubDirs} 
                        />
                    </Card>

                    <Card withBorder sx={checkboxCardTheme} onClick={handleCheckbox(skipUrlChecking, setSkipUrlChecking)}>
                        <Checkbox 
                            label="Skip Checking URLs" 
                            description="Skip checking every URL against known compatible hosts for slavart." 
                            checked={skipUrlChecking} 
                        />
                    </Card>

                    <Card withBorder sx={checkboxCardTheme} onClick={handleCheckbox(dryRun, setDryRun)}>
                        <Checkbox 
                            label="Dry Run" 
                            description="The zip file collected won't be downloaded but all other systems run." 
                            checked={dryRun} 
                        />
                    </Card>
                </SimpleGrid>
            </Stack>

            <Stack spacing={2}>
                <Text size="xs" tt="uppercase" c="dimmed">Options</Text>

                <SimpleGrid cols={2}>
                    <Card withBorder sx={cardTheme}>
                        <FSInput 
                            sx={inputTheme}
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

                    <Card withBorder sx={cardTheme}>
                        <Select
                            sx={inputTheme}
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

                    <Card withBorder sx={cardTheme}>
                        <NumberInput
                            sx={inputTheme}
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

                    <Card withBorder sx={cardTheme}>
                        <NumberInput
                            sx={inputTheme}
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