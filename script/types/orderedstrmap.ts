export interface OrderedStrMapIter {
    readonly Close: () => void;
    readonly Next: () => [string, string];
    readonly Prev: () => [string, string];
}

export interface OrderedStrMap {
    readonly Clear: () => void;
    readonly Delete: (k: string) => void;
    readonly Get: (k: string) => [string, boolean];
    readonly Len: () => number;
    readonly Set: (k: string, v: string) => void;
    readonly Seek: (k: string) => [OrderedStrMapIter, boolean];
    readonly SeekFirst: () => OrderedStrMapIter;
    readonly SeekLast: () => OrderedStrMapIter;
}

export declare const NewOrderedStrMap: () => OrderedStrMap;
