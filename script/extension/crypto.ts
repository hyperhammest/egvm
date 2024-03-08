export interface PrivateKey {
    readonly GetPublicKey: () => PublicKey;
    readonly ECDH: (pubkey: PublicKey) => ArrayBuffer;
    readonly Encapsulate: (pubkey: PublicKey) => ArrayBuffer;
    readonly Sign: (msg: ArrayBuffer) => ArrayBuffer;
    readonly Equal: (other: PrivateKey) => boolean;
    readonly Hex: () => string;
    readonly Serialize: () => ArrayBuffer;
    readonly Decrypt: (buf: ArrayBuffer) => [ArrayBuffer, boolean];
    readonly VrfProve: (alpha: ArrayBuffer) => [ArrayBuffer, ArrayBuffer];
}

export declare const BufToPrivateKey: (buf: ArrayBuffer) => PrivateKey;

export interface PublicKey {
    readonly Decapsulate: (privKey: PrivateKey) => ArrayBuffer;
    readonly Equal: (other: PublicKey) => boolean;
    readonly Hex: (compressed: boolean) => string;
    readonly SerializeCompressed: () => ArrayBuffer;
    readonly SerializeUncompressed: () => ArrayBuffer;
    readonly Encrypt: (buf, entropy: ArrayBuffer) => ArrayBuffer;
    readonly ToEvmAddress: () => ArrayBuffer;
    readonly ToCashAddress: () => ArrayBuffer;
    readonly VrfVerify: (alpha, pi: ArrayBuffer) => ArrayBuffer;
}

export declare const BufToPublicKey: (buf: ArrayBuffer) => PublicKey;
export declare const GetEthSignedMessage: (msg: ArrayBuffer) => ArrayBuffer;
export declare const VerifySignature: (pubkey, digestHash, signature: ArrayBuffer) => boolean;
export declare const Ecrecover: (hash, sig: ArrayBuffer) => PublicKey;
export declare const AesGcmEncrypt: (secret, nonce, msg: ArrayBuffer) => ArrayBuffer;
export declare const AesGcmDecrypt: (secret, nonce, ciphertext: ArrayBuffer) => [ArrayBuffer, boolean];
