import type { ComponentChildren } from "preact"
import { Head, HeadProps } from "./Head"

export interface DocumentProps extends HeadProps {
    htmlClassName: string
    children: ComponentChildren
}

export function Document(props: DocumentProps) {
    return (
        <html lang="en" class={props.htmlClassName}>
            <Head title={props.title} />
            <body>{props.children}</body>
        </html>
    )
}
