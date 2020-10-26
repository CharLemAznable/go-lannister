package dummy

import (
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
}

func init() {
    golog.Debugf("Generate Private Key: %s", PrivateKeyString)
    golog.Debugf("Generate Public Key: %s", PublicKeyString)
}
