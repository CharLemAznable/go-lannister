package elf

import (
    "encoding/binary"
    "net"
)

func GetIpInt64() (int64, error) {
    ip, err := FirstNonLoopbackAddress()
    if nil != err {
        return 0, err
    }
    return int64(binary.BigEndian.Uint32(ip)), nil
}

func FirstNonLoopbackAddress() (net.IP, error) {
    addrs, err := net.InterfaceAddrs()
    if nil != err {
        return nil, err
    }
    for _, value := range addrs {
        if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.To4(), nil
            }
        }
    }
    return net.IPv4(127, 0, 0, 1).To4(), nil
}
