import {OrderedIntMap} from "./orderedintmap";
import {OrderedStrMap} from "./orderedstrmap";
import {OrderedBufMap} from "./orderedbufmap";

export declare const SerializeMaps: (...map: Array<OrderedIntMap | OrderedStrMap | OrderedBufMap>) => ArrayBuffer;
export declare const DeserializeMap: (buf: ArrayBuffer, expectedTag: number) => [OrderedIntMap | OrderedStrMap | OrderedBufMap, ArrayBuffer];

export interface OrderedMapReader {
    readonly Reset: () => void;
    readonly Read: (expectedTag: number) => OrderedIntMap | OrderedStrMap | OrderedBufMap;
}

export declare const NewOrderedMapReader: (buf: ArrayBuffer) => OrderedMapReader;
