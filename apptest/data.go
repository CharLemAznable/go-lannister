package apptest

import (
    "errors"
    . "github.com/CharLemAznable/go-lannister/elf"
    "github.com/kataras/golog"
)

var (
    GeneratedKeyPair, _ = GenerateKeyPairDefault()
    PrivateKeyString, _ = GeneratedKeyPair.PrivateKeyEncoded()
    PublicKeyString, _  = GeneratedKeyPair.PublicKeyEncoded()
)

var accessors = map[string]map[string]string{
    "1001": {
        "AccessorId":     "1001",
        "AccessorName":   "1001",
        "AccessorPubKey": PublicKeyString,
    },
    "1002": {
        "AccessorId":     "1002",
        "AccessorName":   "1002",
        "AccessorPubKey": PublicKeyString,
    },
}

var accessorErrors = map[string]error{
    "1002": errors.New("MockError"),
}

var merchant1001 = map[string]string{
    "MerchantId":   "1001",
    "MerchantName": "",
    "MerchantCode": "m1001",
}

var merchantById = map[string]map[string]string{
    "1001": merchant1001,
}

var merchantByCode = map[string]map[string]string{
    "m1001": merchant1001,
}

var present = struct{}{}

var accessorMerchants = map[string]map[string]interface{}{
    "0": {},
    "1001": {
        "1001": present,
    },
}

func getAccessorMerchants(accessorId string) map[string]interface{} {
    merchants, ok := accessorMerchants[accessorId]
    if !ok {
        merchants = map[string]interface{}{}
        accessorMerchants[accessorId] = merchants
    }
    return merchants
}

func init() {
    golog.Debugf("Generate Private Key: %s", PrivateKeyString)
    golog.Debugf("Generate Public Key: %s", PublicKeyString)
}
