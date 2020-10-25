package types

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestRegisterAccessorDao(t *testing.T) {
    a := assert.New(t)
    a.Nil(GetAccessorManageDao(nil))
    a.Nil(GetAccessorVerifyDao(nil))
}
