import {Uint256} from "./u256";

export interface Sint256 {
    readonly Abs: () => Sint256;
    readonly Neg: () => Sint256;
    readonly Add: (v: Sint256) => Sint256;
    readonly Equal: (v: Sint256) => boolean;
    readonly IsZero: () => boolean;
    readonly Sign: () => number;
    readonly Sub: (v: Sint256) => Sint256;
    readonly Mul: (v: Sint256) => Sint256;
    readonly Div: (v: Sint256) => Sint256;
    readonly Lsh: (v: number) => Sint256;
    readonly Rsh: (v: number) => Sint256;
    readonly Gt: (v: Sint256) => boolean;
    readonly Gte: (v: Sint256) => boolean;
    readonly GtNum: (v: number) => boolean;
    readonly GteNum: (v: number) => boolean;
    readonly Lt: (v: Sint256) => boolean;
    readonly Lte: (v: Sint256) => boolean;
    readonly LtNum: (v: number) => boolean;
    readonly LteNum: (v: number) => boolean;
    readonly IsSafeInteger: () => boolean;
    readonly ToSafeInteger: () => number;
    readonly String: () => string;
    readonly ToU256: () => Uint256;
}

export declare const HexToS256: (hex: string) => Sint256;
export declare const BufToS256: (buf: ArrayBuffer) => Sint256;
export declare const S256: (v: number) => Sint256;
