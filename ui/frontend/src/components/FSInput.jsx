import { useRef } from "preact/hooks"
import { TextInput, Button } from "@mantine/core"
import { OpenFileDialog, SaveFileDialog, OpenDirectoryDialog } from "../../wailsjs/go/main/SlavartdlUI.js";


export default function FSInput(props) {
    const inputRef = useRef()


    const handleBrowseButton = async () => {
        inputRef.current.focus()

        let func 
        switch (props.func) {
            case "openDirectory": func = OpenDirectoryDialog; break
            case "openFile": func = OpenFileDialog; break
            case "saveFile": func = SaveFileDialog; break
        }

        if (func != undefined) {
            props.onChange?.(await func(props.funcData))
        }
    }

    const handleOnChange = (event) => {
        props.onChange?.(event.target.value)
    }
    

    const browseButtonTheme = (theme) => ({
        color: theme.colorScheme === "dark" ? theme.colors.dark[1] : theme.colors.dark[3],
        paddingLeft: 12, 
        paddingRight: 12, 
        borderRadius: 0, 
        borderTopRightRadius: 6,
        borderBottomRightRadius: 6,
        borderLeft: `1px solid ${theme.colorScheme === "dark" ? theme.colors.dark[4] : theme.colors.gray[4]}`,
        height: "calc(100% - 2px)",
    })

    const browseButton = (
        <Button compact={false} variant="subtle" size="xs" color="gray" sx={browseButtonTheme} onClick={handleBrowseButton}>
            Browse
        </Button>
    )


    let inputProps = {...props} 
    delete inputProps.func 
    delete inputProps.funcData

    return (
        <TextInput {...inputProps} onChange={handleOnChange} ref={inputRef} rightSection={browseButton} rightSectionWidth={66}/>
    )
}