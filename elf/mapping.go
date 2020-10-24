package elf

import (
    "net/http"
)

func GetMapping(relativePath, funcName string) (string, string, string) {
    return http.MethodGet, relativePath, funcName
}

func PostMapping(relativePath, funcName string) (string, string, string) {
    return http.MethodPost, relativePath, funcName
}
