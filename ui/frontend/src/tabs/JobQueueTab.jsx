import { useId } from "preact/hooks"

export default function JobQueueTab(props) {
    const job = {
        urls: [],
        id: useId()
    }
    console.log(job)


    return (
        <p>{props.jobQueue.join(", ")}</p>
    )
}