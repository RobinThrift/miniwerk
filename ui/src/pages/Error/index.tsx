import { Document, DocumentProps } from "../_internal/Document"

export interface ErrorPageProps extends DocumentProps {
    code: number
    title: string
    details?: string
}

export function Error(props: ErrorPageProps) {
    return (
        <Document {...props}>
            <div className="grid h-screen px-4 bg-white place-content-center">
                <div className="text-center">
                    <h1 className="font-black text-9xl text-transparent rainbow bg-clip-text bg-gradient-to-bl">
                        {props.code}
                    </h1>
                    <p className="text-2xl font-bold tracking-tight text-gray-900 sm:text-4xl">
                        {props.title}
                    </p>
                    <p className="mt-4 text-gray-500">{props.details}</p>

                    <a
                        href="/"
                        className="btn btn-primary inline-block mt-6 py-4 md:text-base"
                    >
                        Back Home
                    </a>
                </div>
            </div>
        </Document>
    )
}
