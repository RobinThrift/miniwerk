import preact from "preact"
import { Index } from "../pages/Index"
import { Error } from "../pages/Error"

export const pages: Record<string, preact.FunctionComponent> = {
    Index: Index,
    Error: Error,
}
