package elf

import (
    "github.com/bwmarrin/snowflake"
)

func getNodeId() int64 {
    ipInt64, _ := GetIpInt64()
    return ^(int64(-1) << snowflake.NodeBits) & ipInt64
}

var idGenerator, _ = snowflake.NewNode(getNodeId())

func NextId() string {
    return idGenerator.Generate().String()
}
