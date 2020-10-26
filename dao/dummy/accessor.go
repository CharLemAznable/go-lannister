package dummy

import (
    "errors"
    "github.com/CharLemAznable/go-lannister/types"
    "github.com/CharLemAznable/sqlx"
)

type AccessorManageDao struct{}

func NewAccessorManageDao(_ *sqlx.DB) types.AccessorManageDao {
    return &AccessorManageDao{}
}

func (d *AccessorManageDao) QueryAccessor(accessorId string) (*types.AccessorManage, error) {
    accessor := accessors[accessorId]
    return &types.AccessorManage{
        AccessorId:     accessor["AccessorId"],
        AccessorName:   accessor["AccessorName"],
        AccessorPubKey: accessor["AccessorPubKey"],
    }, nil
}

func (d *AccessorManageDao) UpdateAccessor(accessorId string, manage *types.AccessorManage) (int64, error) {
    accessor := accessors[accessorId]
    if "" != manage.AccessorName {
        accessor["AccessorName"] = manage.AccessorName
    }
    if "" != manage.AccessorPubKey {
        accessor["AccessorPubKey"] = manage.AccessorPubKey
    }
    if "" != manage.PayNotifyUrl {
        accessor["PayNotifyUrl"] = manage.PayNotifyUrl
    }
    if "" != manage.RefundNotifyUrl {
        accessor["RefundNotifyUrl"] = manage.RefundNotifyUrl
    }
    return 1, nil
}

func (d *AccessorManageDao) UpdateKeyPair(accessorId, nonsense, pubKey, prvKey string) {
    accessor := accessors[accessorId]
    accessor["PubKey"] = pubKey
}

type AccessorVerifyDao struct{}

func NewAccessorVerifyDao(_ *sqlx.DB) types.AccessorVerifyDao {
    return &AccessorVerifyDao{}
}

func (d *AccessorVerifyDao) QueryAccessorById(accessorId string) (*types.AccessorVerify, error) {
    accessor, ok := accessors[accessorId]
    if !ok {
        return &types.AccessorVerify{}, errors.New("AccessorNotExists")
    }
    return &types.AccessorVerify{
        AccessorId:     accessor["AccessorId"],
        AccessorPubKey: accessor["AccessorPubKey"],
    }, nil
}

func init() {
    types.RegisterAccessorManageDaoConstructor("dummy", NewAccessorManageDao)
    types.RegisterAccessorVerifyDaoConstructor("dummy", NewAccessorVerifyDao)
}
