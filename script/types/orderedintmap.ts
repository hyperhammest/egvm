export interface OrderedIntMapIter {
    readonly Close: () => void;
    readonly Next: () => [string, number];
    readonly Prev: () => [string, number];
}

export interface OrderedIntMap {
    readonly Clear: () => void;
    readonly Delete: (k: string) => void;
    readonly Get: (k: string) => [number, boolean];
    readonly Len: () => number;
    readonly Set: (k: string, v: number) => void;
    readonly Seek: (k: string) => [OrderedIntMapIter, boolean];
    readonly SeekFirst: () => OrderedIntMapIter;
    readonly SeekLast: () => OrderedIntMapIter;
}

export declare const NewOrderedIntMap: () => OrderedIntMap;
