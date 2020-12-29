// Hacks to properly type global objects set by wasm 
import outputJSONData from './output/output.json'

// @ts-ignore
const go: any = window.go

export type GoFlag = {
    Name: string,
    Usage: string,
    DefValue: string,
    Type: string,
}

type RunGameReturnTypeWASM = {
    output: string,
    logs: string,
    error: string,
}

export type RunGameReturnType = {
    output: typeof outputJSONData,
    logs: string,
}

type GetFlagsFormatsReturnTypeWASM = {
    output: string,
    error: string,
}

export type GetFlagsFormatsReturnType = GoFlag[]

let loaded = false

let runGameWASM: (() => RunGameReturnTypeWASM) | undefined;
let getFlagsFormatsWASM: (() => GetFlagsFormatsReturnTypeWASM) | undefined;

const load = async () => {
    const { instance } = await WebAssembly.instantiateStreaming(
        fetch(`${process.env.PUBLIC_URL}/SOMAS2020.wasm`),
        go.importObject
    )
    go.run(instance)

    // @ts-ignore
    runGameWASM = window.RunGame
    // @ts-ignore
    getFlagsFormatsWASM = window.GetFlagsFormats
    
    loaded = true
}

export const runGame = async (): Promise<RunGameReturnType> => {
    if (!loaded) {
        await load()
    }
    if (!runGameWASM) {
        throw new Error("Game not loaded properly")
    }

    const result = runGameWASM()
    if (result.error.length > 0) {
        throw new Error(result.error)
    }

    const processedOutput = JSON.parse(result.output) as typeof outputJSONData

    // we need to patch git info
    processedOutput.GitInfo = outputJSONData.GitInfo

    return {
        output: processedOutput,
        logs: result.logs,
    }
}

export const getFlagsFormats = async (): Promise<GetFlagsFormatsReturnType> => {
    if (!loaded) {
        await load()
    }
    if (!getFlagsFormatsWASM) {
        throw new Error("Game not loaded properly")
    }

    const result = getFlagsFormatsWASM()
    if (result.error.length > 0) {
        throw new Error(result.error)
    }

    const processedOutput = JSON.parse(result.output) as GetFlagsFormatsReturnType

    return processedOutput
}