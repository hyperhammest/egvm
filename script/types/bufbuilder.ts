export interface BufBuilder {
    readonly Len: () => number;
    readonly Reset: () => void;
    readonly Write: (buf: ArrayBuffer) => number;
    readonly ToBuf: () => ArrayBuffer;
}

export declare const NewBufBuilder: () => BufBuilder;
