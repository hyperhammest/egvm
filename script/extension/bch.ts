import {PrivateKey} from './crypto';

export interface OutPoint {
    HexTxID: string
    Index: number
}

export interface TxIn {
    PreviousOutPoint: OutPoint
    HexPubkey: string
    Sequence: number
    Value: number
}

export interface TxOut {
    Value: number
    HexPubkeyHash: string
    HexDataElements: string[]
}

export interface BchTx {
    HexTxID: string
    Version: number
    TxIn: TxIn[]
    TxOut: TxOut[]
    LockTime: number
}

export declare const ParseTxInHex: (hexStr: string) => BchTx;
export declare const SignTxAndSerialize: (tx: BchTx, ...privateKeys: PrivateKey[]) => string;
export declare const MerkleProofToRootAndMatches: (proof: ArrayBuffer) => Array<ArrayBuffer>;
