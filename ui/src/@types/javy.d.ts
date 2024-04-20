declare global {
    var Javy: Javy
}

interface Javy {
    IO: {
        readSync(fd: 0 | 1, buffer: Uint8Array): number
        writeSync(fd: 0 | 1, buffer: Uint8Array): number
    }
}

export type {}
