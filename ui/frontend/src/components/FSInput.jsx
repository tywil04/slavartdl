import { useRef } from "preact/hooks"
import { useUncontrolled } from "@mantine/hooks";
import { TextInput, Button } from "@mantine/core"
import { OpenFileDialog, SaveFileDialog, OpenDirectoryDialog } from "../../wailsjs/go/main/SlavartdlUI.js";


export default function FSInput(props) {
    const [value, onChange] = useUncontrolled({
        value: props.value,
        onChange: props.onChange,
    });
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
            onChange(await func(props.funcData))
        }
    }

    const handleOnChange = (event) => {
        onChange(event.target.value)
    }
    

    const browseButtonStyle = (theme) => ({
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
        <Button 
            compact={false} 
            variant="subtle" 
            size="xs" 
            color="gray" 
            sx={browseButtonStyle} 
            onClick={handleBrowseButton}
        >
            Browse
        </Button>
    )


    let inputProps = {...props} 
    delete inputProps.func 
    delete inputProps.funcData

    
    return (
        <TextInput {...inputProps} value={value} onChange={handleOnChange} ref={inputRef} rightSection={browseButton} rightSectionWidth={66}/>
    )
}