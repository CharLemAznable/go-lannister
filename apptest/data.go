package apptest

import (
    "errors"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/golog"
    "strings"
)

var (
    GeneratedKeyPair, _ = gokits.GenerateRSAKeyPairDefault()
    PrivateKeyString, _ = GeneratedKeyPair.RSAPrivateKeyEncoded()
    PublicKeyString, _  = GeneratedKeyPair.RSAPublicKeyEncoded()
)

func init() {
    golog.Debugf("Generate Private Key: %s", PrivateKeyString)
    golog.Debugf("Generate Public Key: %s", PublicKeyString)
}

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
    "MerchantName": "1001",
    "MerchantCode": "m1001",
}

var merchant1002 = map[string]string{
    "MerchantId":   "1002",
    "MerchantName": "1002",
    "MerchantCode": "m1002",
}

var merchantById = map[string]map[string]string{
    "1001": merchant1001,
    "1002": merchant1002,
}

var merchantByCode = map[string]map[string]string{
    "m1001": merchant1001,
    "m1002": merchant1002,
}

var present = struct{}{}

var accessorMerchants = map[string]map[string]interface{}{
    "1001": {
        "1001": present,
    },
    "1002": {
        "1002": present,
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

var merchantErrors = map[string]error{
    "1002": errors.New("MockError"),
}

var merchantApiParams = map[string]map[string]map[string]string{}

func getMerchantApiParams(merchantId, apiName string) map[string]string {
    apis, ok := merchantApiParams[merchantId]
    if !ok {
        apis = map[string]map[string]string{}
        merchantApiParams[merchantId] = apis
    }
    params, ok := apis[strings.ToLower(apiName)]
    if !ok {
        params = map[string]string{}
        apis[strings.ToLower(apiName)] = params
    }
    return params
}
