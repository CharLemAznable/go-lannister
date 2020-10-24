package elf

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/base64"
    "errors"
)

type KeyPair struct {
    PrivateKey *rsa.PrivateKey
    PublicKey  *rsa.PublicKey
}

func GenerateKeyPairDefault() (*KeyPair, error) {
    return GenerateKeyPair(1024)
}

func GenerateKeyPair(keySize int) (*KeyPair, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
    if err != nil {
        return nil, err
    }
    publicKey := &privateKey.PublicKey
    return &KeyPair{
        PrivateKey: privateKey,
        PublicKey:  publicKey}, nil
}

func (p *KeyPair) PrivateKeyEncoded() (string, error) {
    if nil == p.PrivateKey {
        return "", errors.New("PrivateKeyEmpty")
    }
    bytes, err := x509.MarshalPKCS8PrivateKey(p.PrivateKey)
    if nil != err {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(bytes), nil
}

func (p *KeyPair) PublicKeyEncoded() (string, error) {
    if nil == p.PublicKey {
        return "", errors.New("PublicKeyEmpty")
    }
    bytes, err := x509.MarshalPKIXPublicKey(p.PublicKey)
    if nil != err {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(bytes), nil
}

func PrivateKeyDecoded(privateKeyString string) (*rsa.PrivateKey, error) {
    bytes, err := base64.StdEncoding.DecodeString(privateKeyString)
    if nil != err {
        return nil, err
    }
    privateKey, err := x509.ParsePKCS8PrivateKey(bytes)
    if nil != err {
        return nil, err
    }
    return privateKey.(*rsa.PrivateKey), nil
}

func PublicKeyDecoded(publicKeyString string) (*rsa.PublicKey, error) {
    bytes, err := base64.StdEncoding.DecodeString(publicKeyString)
    if nil != err {
        return nil, err
    }
    publicKey, err := x509.ParsePKIXPublicKey(bytes)
    if nil != err {
        return nil, err
    }
    return publicKey.(*rsa.PublicKey), nil
}

type Signer struct {
    hash crypto.Hash
}

func (s *Signer) SignBase64ByKeyString(plainText, privateKeyString string) (string, error) {
    sign, err := s.SignByKeyString(plainText, privateKeyString)
    if nil != err {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(sign), nil
}

func (s *Signer) SignBase64(plainText string, privateKey *rsa.PrivateKey) (string, error) {
    sign, err := s.Sign(plainText, privateKey)
    if nil != err {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(sign), nil
}

func (s *Signer) SignByKeyString(plainText, privateKeyString string) ([]byte, error) {
    privateKey, err := PrivateKeyDecoded(privateKeyString)
    if nil != err {
        return nil, err
    }
    return s.Sign(plainText, privateKey)
}

func (s *Signer) Sign(plainText string, privateKey *rsa.PrivateKey) ([]byte, error) {
    hash := s.hash.New()
    hash.Write([]byte(plainText))
    return rsa.SignPKCS1v15(rand.Reader, privateKey, s.hash, hash.Sum(nil))
}

func (s *Signer) VerifyBase64ByKeyString(plainText, signText, publicKeyString string) error {
    sign, err := base64.StdEncoding.DecodeString(signText)
    if nil != err {
        return err
    }
    return s.VerifyByKeyString(plainText, sign, publicKeyString)
}

func (s *Signer) VerifyBase64(plainText, signText string, publicKey *rsa.PublicKey) error {
    sign, err := base64.StdEncoding.DecodeString(signText)
    if nil != err {
        return err
    }
    return s.Verify(plainText, sign, publicKey)
}

func (s *Signer) VerifyByKeyString(plainText string, sign []byte, publicKeyString string) error {
    publicKey, err := PublicKeyDecoded(publicKeyString)
    if nil != err {
        return err
    }
    return s.Verify(plainText, sign, publicKey)
}

func (s *Signer) Verify(plainText string, sign []byte, publicKey *rsa.PublicKey) error {
    hash := s.hash.New()
    hash.Write([]byte(plainText))
    return rsa.VerifyPKCS1v15(publicKey, s.hash, hash.Sum(nil), sign)
}

var (
    SHA1WithRSA   = Signer{hash: crypto.SHA1}
    SHA256WithRSA = Signer{hash: crypto.SHA256}
)
