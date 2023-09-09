import { render } from 'preact'
import { useState } from 'preact/hooks';
import { useUrlQueue } from "./hooks/hooks.js"
import { MantineProvider, ColorSchemeProvider, TypographyStylesProvider, Tooltip } from '@mantine/core';
import { Tabs } from '@mantine/core';
import { HiMiniArrowDownTray, HiMiniCog6Tooth, HiMiniQueueList } from "react-icons/hi2"
import { Notifications } from "@mantine/notifications"
import DownloadTab from './tabs/DownloadTab.jsx';
import QueueTab from './tabs/QueueTab.jsx';
import SettingsTab from './tabs/SettingsTab.jsx';


export default function Root() {
    const [ queue, queueHandlers ] = useUrlQueue([])
    const [ colorScheme, setColorScheme ] = useState("dark")


    const mantineTheme = {
        colorScheme: colorScheme,
        defaultRadius: 6,
        activeStyles: {
            transform: "",
        },
        cursorType: "pointer",
        components: {
            Button: {
                defaultProps: {
                    compact: true,
                },
            },
        },
        globalStyles: (theme) => ({
            "#root": {
                backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[8] : theme.colors.gray[1],
                padding: 40,
                minHeight: "100vh"
            }
        })
    }


    const tabsStyle = (theme) => ({
        root: {
            borderRadius: 6,
        },
        tabsList: {
            backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[8] : theme.colors.gray[1],
            border: `1px solid ${theme.colorScheme === "dark" ? theme.colors.dark[6] : theme.colors.gray[3]}`,
            borderRadius: 12,
            padding: 4,
            width: "fit-content",
        },
        panel: {
            marginTop: 20,
        },
        tab: {
            transitionDuration: "100ms",
            "&[aria-selected=true]": {
                backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[6] : theme.colors.gray[3],
                color: theme.colorScheme === "dark" ? theme.colors.dark[0] : theme.black,
            },
            "&[aria-selected=false]:hover": {
                backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[7] : theme.colors.gray[2]
            },
            "&[data-active]:hover": {
                backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[6] : theme.colors.gray[3]
            }
        }
    })


    return (
        <ColorSchemeProvider colorScheme={colorScheme}>
            <MantineProvider theme={mantineTheme} withGlobalStyles withNormalizeCSS>
                <TypographyStylesProvider>
                    <Notifications position="top-center"/>
                    
                    <Tabs defaultValue="download" variant="pills" styles={tabsStyle}>
                        <Tabs.List>
                            <Tooltip label="Download" withArrow position="right">
                                <Tabs.Tab value="download" icon={<HiMiniArrowDownTray size="16"/>}/>
                            </Tooltip>

                            <Tooltip label="Queue" withArrow position="right">
                                <Tabs.Tab value="queue" icon={<HiMiniQueueList size="16"/>}/>
                            </Tooltip>

                            <Tooltip label="Settings" withArrow position="right">
                                <Tabs.Tab value="settings" icon={<HiMiniCog6Tooth size="16"/>}/>
                            </Tooltip>
                        </Tabs.List>

                        <Tabs.Panel value="download">
                            <DownloadTab queue={queue} queueHandlers={queueHandlers}/>
                        </Tabs.Panel>

                        <Tabs.Panel value="queue">
                            <QueueTab queue={queue} queueHandlers={queueHandlers}/>
                        </Tabs.Panel>

                        <Tabs.Panel value="settings">
                            <SettingsTab colorScheme={colorScheme} setColorScheme={setColorScheme}/>
                        </Tabs.Panel>
                    </Tabs>
                </TypographyStylesProvider>
            </MantineProvider>
        </ColorSchemeProvider>
    )
}


render(
    <Root/>, 
    document.getElementById("root")
)