import {PrivateKey} from "./crypto";

export interface Bip32Key {
    readonly B58Serialize: () => string;
    readonly NewChildKey:(childIdx: number) => Bip32Key;
    readonly PublicKey:() => Bip32Key;
    readonly Serialize:() => ArrayBuffer;
    readonly IsPrivate:() => boolean;
    readonly ToPrivateKey:() => PrivateKey;
    readonly DeriveWithBytes32:() =>Bip32Key;
    readonly Derive:(n0,n1,n2,n3,n4 : number) => Bip32Key;
}

export declare const B58ToBip32Key: (data: string) => Bip32Key;
export declare const BufToBip32Key: (buf: ArrayBuffer) => Bip32Key;
export declare const GenerateRandomBip32Key: () => Bip32Key;

