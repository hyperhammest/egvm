import {Uint256} from "../types/u256";

export declare const Keccak256: (...args: Array<string | ArrayBuffer | Uint256>) => ArrayBuffer;
export declare const Sha256: (...args: Array<string | ArrayBuffer | Uint256>) => ArrayBuffer;
export declare const Ripemd160: (...args: Array<string | ArrayBuffer | Uint256>) => ArrayBuffer;
export declare const XxHash32: (...args: Array<string | ArrayBuffer | Uint256>) => ArrayBuffer;
export declare const XxHash64: (...args: Array<string | ArrayBuffer | Uint256>) => ArrayBuffer;
export declare const XxHash128: (...args: Array<string | ArrayBuffer | Uint256>) => ArrayBuffer;
export declare const XxHash32Int: (...args: Array<string | ArrayBuffer | Uint256>) => ArrayBuffer;
