package elf

import (
    "crypto/rand"
    "crypto/rsa"
    "github.com/stretchr/testify/assert"
    "testing"
)

var _ = func() bool {
    testing.Init()
    return true
}()

func TestKeyPair(t *testing.T) {
    a := assert.New(t)

    plainText := "{ mac=\"MAC Address\", appId=\"16位字符串\", signature=SHA1(\"appId=xxx&mac=yyy\") }"
    keyPair, _ := GenerateKeyPairDefault()
    privateKeyString, _ := keyPair.PrivateKeyEncoded()
    publicKeyString, _ := keyPair.PublicKeyEncoded()

    privateKey, _ := PrivateKeyDecoded(privateKeyString)
    publicKey, _ := PublicKeyDecoded(publicKeyString)

    cipherBytes, _ := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plainText))
    plainBytes, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherBytes)
    a.Equal(plainText, string(plainBytes))

    pair := KeyPair{}
    _, errPrv := pair.PrivateKeyEncoded()
    a.NotNil(errPrv)
    _, errPub := pair.PublicKeyEncoded()
    a.NotNil(errPub)
}

func TestSigner(t *testing.T) {
    a := assert.New(t)

    plainText := "{ mac=\"MAC Address\", appId=\"16位字符串\", signature=SHA1(\"appId=xxx&mac=yyy\") }"
    keyPair, _ := GenerateKeyPairDefault()
    privateKeyString, _ := keyPair.PrivateKeyEncoded()
    publicKeyString, _ := keyPair.PublicKeyEncoded()
    privateKey := keyPair.PrivateKey
    publicKey := keyPair.PublicKey

    sign1, _ := SHA1WithRSA.SignBase64ByKeyString(plainText, privateKeyString)
    a.Nil(SHA1WithRSA.VerifyBase64ByKeyString(plainText, sign1, publicKeyString))

    a.NotNil(SHA1WithRSA.VerifyBase64ByKeyString(plainText, sign1, publicKeyString[1:]))
    a.NotNil(SHA1WithRSA.VerifyBase64ByKeyString(plainText, sign1[1:], publicKeyString))
    _, err1 := SHA1WithRSA.SignBase64ByKeyString(plainText, privateKeyString[1:])
    a.NotNil(err1)

    sign256, _ := SHA256WithRSA.SignBase64(plainText, privateKey)
    a.Nil(SHA256WithRSA.VerifyBase64(plainText, sign256, publicKey))

    a.NotNil(SHA256WithRSA.VerifyBase64(plainText, sign256[1:], publicKey))
}
