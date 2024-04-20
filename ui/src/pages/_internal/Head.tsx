export interface HeadProps {
    title: string
    csrfToken?: string
}

export function Head(props: HeadProps) {
    return (
        <head>
            <meta charset="utf-8" />
            <meta name="viewport" content="width=device-width" />
            <meta
                name="viewport"
                content="width=device-width, initial-scale=1, minimum-scale=1"
            />
            <meta name="mobile-web-app-capable" content="yes" />
            {props.csrfToken && (
                <meta name="csrf-token" content={props.csrfToken} />
            )}
            <link rel="stylesheet" href="/static/styles.css" />
            <script defer src="/static/index.min.js" type="module"></script>

            <title>MiniWerk - {props.title}</title>
        </head>
    )
}
