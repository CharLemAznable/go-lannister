package dummy

import (
    "errors"
    . "github.com/CharLemAznable/go-lannister/elf"
    "github.com/CharLemAznable/go-lannister/types"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/golog"
)

type AccessorManageDao struct {
    data map[string]*types.AccessorManage
}

var (
    TestKeyPair, _          = GenerateKeyPairDefault()
    TestPrivateKeyString, _ = TestKeyPair.PrivateKeyEncoded()
    TestPublicKeyString, _  = TestKeyPair.PublicKeyEncoded()
)

func NewAccessorManageDao(_ *sqlx.DB) types.AccessorManageDao {
    return &AccessorManageDao{data: map[string]*types.AccessorManage{
        "1001": {
            AccessorId:     "1001",
            AccessorName:   "1001",
            AccessorPubKey: TestPublicKeyString,
        },
    }}
}

func (d *AccessorManageDao) QueryAccessorById(accessorId string) (*types.AccessorManage, error) {
    return d.data[accessorId], nil
}

func (d *AccessorManageDao) UpdateAccessorById(accessorId string, manage *types.AccessorManage) (int64, error) {
    origin := d.data[accessorId]
    if "" != manage.AccessorName {
        origin.AccessorName = manage.AccessorName
    }
    if "" != manage.AccessorPubKey {
        origin.AccessorPubKey = manage.AccessorPubKey
    }
    if "" != manage.PayNotifyUrl {
        origin.PayNotifyUrl = manage.PayNotifyUrl
    }
    if "" != manage.RefundNotifyUrl {
        origin.RefundNotifyUrl = manage.RefundNotifyUrl
    }
    return 1, nil
}

func (d *AccessorManageDao) UpdateKeyPairById(accessorId, nonsense, pubKey, prvKey string) {
    origin := d.data[accessorId]
    origin.PubKey = pubKey
}

type AccessorVerifyDao struct {
    data map[string]*types.AccessorVerify
}

func NewAccessorVerifyDao(_ *sqlx.DB) types.AccessorVerifyDao {
    return &AccessorVerifyDao{data: map[string]*types.AccessorVerify{
        "1001": {
            AccessorId:     "1001",
            AccessorPubKey: TestPublicKeyString,
        },
    }}
}

func (d *AccessorVerifyDao) QueryAccessorById(accessorId string) (*types.AccessorVerify, error) {
    verify, ok := d.data[accessorId]
    if !ok {
        return new(types.AccessorVerify), errors.New("AccessorNotExists")
    }
    return verify, nil
}

func init() {
    golog.Debugf("Generate Private Key: %s", TestPrivateKeyString)
    golog.Debugf("Generate Public Key: %s", TestPublicKeyString)

    types.RegisterAccessorManageDaoConstructor("dummy", NewAccessorManageDao)
    types.RegisterAccessorVerifyDaoConstructor("dummy", NewAccessorVerifyDao)
}
