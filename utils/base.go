package utils

import (
    "encoding/base64"
    "fmt"
)

// GetBaseEncode 对传参进行base64解码
func GetBaseEncode(data string) string {
    return base64.StdEncoding.EncodeToString([]byte(data))
}

// GetBaseDecode 对编码结果进行base64解码
func GetBaseDecode(data string) (decode string) {
    decodeBytes, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        fmt.Println(err)
    }
    return string(decodeBytes)
}
