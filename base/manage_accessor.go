package base

import (
    "github.com/CharLemAznable/gokits"
    "github.com/CharLemAznable/sqlx"
    "sort"
    "strings"
)

type AccessorManage struct {
    BaseResp

    AccessorId      string `json:"accessorId"`
    AccessorName    string `json:"accessorName"`
    AccessorPubKey  string `json:"accessorPubKey"`  // 访问者公钥, 用于平台验证访问者发起的请求
    PayNotifyUrl    string `json:"payNotifyUrl"`    // 支付回调地址
    RefundNotifyUrl string `json:"refundNotifyUrl"` // 退款回调地址
    PubKey          string `json:"pubKey"`          // 平台公钥, 用于访问者验证平台回调的请求
}

type AccessorManageDao interface {
    // 查询访问者信息
    QueryAccessor(accessorId string) (*AccessorManage, error)
    // 更新访问者名称/访问者公钥/支付回调地址/退款回调地址
    UpdateAccessor(accessorId string, manage *AccessorManage) (int64, error)
    // 更新平台分配给访问者的秘钥对
    UpdateKeyPair(accessorId, nonsense, pubKey, prvKey string) error
}

var accessorManageDaoRegistry = NewDaoRegistry("AccessorManageDao")

func RegisterAccessorManageDao(name string, builder func(*sqlx.DB) AccessorManageDao) {
    accessorManageDaoRegistry.Register(name, builder)
}

func GetAccessorManageDao(db *sqlx.DB) AccessorManageDao {
    return accessorManageDaoRegistry.GetDao(db).(AccessorManageDao)
}

/****************************************************************************************************/

const (
    NONSENSE  = "nonsense"
    SIGNATURE = "signature"
)

type AccessorVerify struct {
    AccessorId     string
    AccessorPubKey string
}

func (v *AccessorVerify) Verify(paramMap map[string]string) error {
    signature := paramMap[SIGNATURE]
    delete(paramMap, SIGNATURE)

    names := make([]string, 0)
    for name := range paramMap {
        names = append(names, name)
    }
    sort.Strings(names)

    paramPairs := make([]string, 0)
    for _, name := range names {
        value := paramMap[name]
        if "" == name || "" == value {
            continue
        }
        paramPairs = append(paramPairs, name+"="+value)
    }
    plainText := strings.Join(paramPairs, "&")

    return gokits.SHA1WithRSA.VerifyBase64ByRSAKeyString(
        plainText, signature, v.AccessorPubKey)
}

type AccessorVerifyDao interface {
    // 查询访问者公钥信息
    QueryAccessorById(accessorId string) (*AccessorVerify, error)
}

var accessorVerifyDaoRegistry = NewDaoRegistry("AccessorVerifyDao")

func RegisterAccessorVerifyDao(name string, builder func(*sqlx.DB) AccessorVerifyDao) {
    accessorVerifyDaoRegistry.Register(name, builder)
}

func GetAccessorVerifyDao(db *sqlx.DB) AccessorVerifyDao {
    return accessorVerifyDaoRegistry.GetDao(db).(AccessorVerifyDao)
}
