import {Sint256} from './S256';
export interface Uint256 {
    readonly Incr: () => Uint256;
    readonly Add: (v: Uint256) => Uint256;
    readonly UnsafeAdd: (v: Uint256) => Uint256;
    readonly And: (v: Uint256) => Uint256;
    readonly Div: (v: Uint256) => Uint256;
    readonly Mod: (v: Uint256) => Uint256;
    readonly DivMod: (v: Uint256) => [Uint256, Uint256];
    readonly Exp: (v: Uint256) => Uint256;
    readonly Gt: (v: Uint256) => boolean;
    readonly Gte: (v: Uint256) => boolean;
    readonly GtNum: (v: number) => boolean;
    readonly GteNum: (v: number) => boolean;
    readonly IsZero: () => boolean;
    readonly Equal: (v: Uint256) => boolean;
    readonly Lt: (v: Uint256) => boolean;
    readonly Lte: (v: Uint256) => boolean;
    readonly LtNum: (v: number) => boolean;
    readonly LteNum: (v: number) => boolean;
    readonly Mul: (v: Uint256) => Uint256;
    readonly UnsafeMul: (v: Uint256) => Uint256;
    readonly Not: () => Uint256;
    readonly Or: (v: Uint256) => Uint256;
    readonly Lsh: (v: number) => Uint256;
    readonly Rsh: (v: number) => Uint256;
    readonly Sqrt: () => Uint256;
    readonly Sub: (v: Uint256) => Uint256;
    readonly UnsafeSub: (v: Uint256) => Uint256;
    readonly IsSafeInteger: () => boolean;
    readonly ToSafeInteger: () => number;
    readonly String: () => string;
    readonly ToS256:() => Sint256;
}

export declare const HexToU256: (hex: string) => Uint256;
export declare const BufToU256: (buf: ArrayBuffer) => Uint256;
export declare const U256: (v: number) => Uint256;
