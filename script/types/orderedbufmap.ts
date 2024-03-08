export interface OrderedBufMapIter {
    readonly Close: () => void;
    readonly Next: () => [string, ArrayBuffer];
    readonly Prev: () => [string, ArrayBuffer];

}

export interface OrderedBufMap {
    readonly Clear: () => void;
    readonly Delete: (key: string) => void;
    readonly Get: (key: string) => [ArrayBuffer, boolean];
    readonly Len: () => number;
    readonly Set: (k: string, v: ArrayBuffer) => void;
    readonly Seek: (k: string) => [OrderedBufMapIter, boolean]
    readonly SeekFirst: () => OrderedBufMapIter;
    readonly SeekLast: () => OrderedBufMapIter;
}

export declare const NewOrderedBufMap: () => OrderedBufMap;
