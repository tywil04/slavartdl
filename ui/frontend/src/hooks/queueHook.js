import { useListState } from "@mantine/hooks";
import { DownloadUrl } from "../../wailsjs/go/main/SlavartdlUI.js";


export function useQueue(starting) {
    const [ state, stateHandlers ] = useListState(starting)


    const running = () => state.filter((url) => url.state === "started").length !== 0

    const appendBatch = (urls, settings) => {
        const processed = urls.map((url) => ({
            url: url,
            state: "queued",
            timeStarted: null,
            timeFinished: null,
            settings: settings
        }))
        stateHandlers.append(...processed)
        
        if (!running()) {
            stateHandlers.setItemProp(0, "state", "started")
            stateHandlers.setItemProp(0, "timeStarted", new Date())

            DownloadUrl(
                processed[0].url, 
                processed[0].settings.outputDir, 
                processed[0].settings.quality,
                processed[0].settings.timeout,
                processed[0].settings.cooldown,
                processed[0].settings.skipUnzip,
                processed[0].settings.ignoreCover,
                processed[0].settings.ignoreSubdir,
            ).then((success) => {
                if (success) {
                    stateHandlers.setItemProp(0, "state", "finished")
                    stateHandlers.setItemProp(0, "timeFinished", new Date())
                } else {
                    console.log("failed")
                    stateHandlers.setItemProp(0, "state", "failed")
                }
            })
        }
    }


    return [ state, { ...stateHandlers, appendBatch: appendBatch } ]
}