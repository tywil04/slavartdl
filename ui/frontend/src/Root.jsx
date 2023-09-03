import { render } from 'preact'
import { MantineProvider, SimpleGrid, Button, ColorSchemeProvider } from '@mantine/core';
import { Tabs } from '@mantine/core';
import DownloadTab from './tabs/DownloadTab.jsx';
import { useState } from 'preact/hooks';


export default function Root() {
    const [ colorScheme, setColorScheme ] = useState("dark")


    const handleTheme = (event) => {
        setColorScheme(colorScheme === "dark" ? "light" : "dark")
    }


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
            body: {
                backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[8] : theme.colors.gray[1],
                padding: 60
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
            marginTop: 30,
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

    const themeButtonTheme = (theme) => ({
        border: `1px solid ${theme.colorScheme === "dark" ? theme.colors.dark[6] : theme.colors.gray[3]}`,
        color: theme.colorScheme === "dark" ? theme.colors.dark[0] : theme.black,
        fontWeight: "normal",
        marginTop: "auto",
        marginBottom: "auto",
        marginLeft: "auto",
        "&:hover": {
            backgroundColor: theme.colorScheme === "dark" ? theme.colors.dark[7] : theme.colors.gray[2]
        }
    })


    return (
        <ColorSchemeProvider colorScheme={colorScheme} toggleColorScheme={handleTheme}>
            <MantineProvider theme={mantineTheme} withGlobalStyles withNormalizeCSS>
                <Tabs defaultValue="download" variant="pills" styles={tabsStyle}>
                    <SimpleGrid cols={2}>
                        <Tabs.List>
                            <Tabs.Tab value="download">Download Tab</Tabs.Tab>
                            <Tabs.Tab value="queue">Queue Tab</Tabs.Tab>
                        </Tabs.List>
                        
                        <Button 
                            variant="outline" 
                            compact={false} 
                            sx={themeButtonTheme} 
                            tt="capitalize"
                            onClick={handleTheme}
                        >
                            Use {colorScheme === "dark" ? "light" : "dark"} Mode
                        </Button>
                    </SimpleGrid>

                    <Tabs.Panel value="download">
                        <DownloadTab/>
                    </Tabs.Panel>

                    <Tabs.Panel value="queue">
                        <p>This is queue</p>
                    </Tabs.Panel>
                </Tabs>
            </MantineProvider>
        </ColorSchemeProvider>
    )
}


render(
    <Root/>, 
    document.getElementById("root")
)