import { h } from "preact"
import { render } from "preact-render-to-string"
import { pages } from "./pages"

const CHUNK_SIZE = 1024
const STDIN = 0
const STDOUT = 1

main()

function main() {
    let input: Input
    try {
        input = readInput()
        // biome-ignore lint/suspicious/noExplicitAny: just typescript things
    } catch (err: any) {
        console.error("error reading input:", err.message)
        return
    }

    let output: Output
    try {
        output = renderPage(input)
        // biome-ignore lint/suspicious/noExplicitAny: just typescript things
    } catch (err: any) {
        console.error("error rendering page:", err.message)
        return
    }

    try {
        writeOutput(output)
        // biome-ignore lint/suspicious/noExplicitAny: just typescript things
    } catch (err: any) {
        console.error("error writing output:", err.message)
        return
    }
}

interface Input {
    page: keyof typeof pages
    // biome-ignore lint/suspicious/noExplicitAny: unknown data
    data: any
}

interface Output {
    html: string
}

function renderPage(input: Input): Output {
    let page = pages[input.page]
    if (!page) {
        throw new Error(`unknown page ${input.page}`)
    }

    return { html: render(h(page, input.data)) }
}

function readInput(): Input {
    // let totalBytes = 0
    // let inputChunks = []

    let finalBuffer = new Uint8Array(readAllInput())

    // while (true) {
    //     let buffer = new Uint8Array(CHUNK_SIZE)
    //     let bytesRead = Javy.IO.readSync(STDIN, buffer)
    //
    //     totalBytes += bytesRead
    //     if (bytesRead === 0) {
    //         break
    //     }
    //
    //     inputChunks.push(buffer.subarray(0, bytesRead))
    // }
    //
    // let { finalBuffer } = inputChunks.reduce(
    //     (context, chunk) => {
    //         context.finalBuffer.set(chunk, context.bufferOffset)
    //         context.bufferOffset += chunk.length
    //         return context
    //     },
    //     { bufferOffset: 0, finalBuffer: new Uint8Array(totalBytes) },
    // )

    return JSON.parse(new TextDecoder().decode(finalBuffer))
}

function writeOutput(output: Output) {
    let encodedOutput = new TextEncoder().encode(output.html)
    let buffer = new Uint8Array(encodedOutput)
    Javy.IO.writeSync(STDOUT, buffer)
}

function* readAllInput() {
    let buffer = new Uint8Array(CHUNK_SIZE)
    while (true) {
        let bytesRead = Javy.IO.readSync(STDIN, buffer)

        if (bytesRead === 0) {
            return
        }

        yield* buffer.subarray(0, bytesRead)
        // buffer.fill(0)
    }
}
