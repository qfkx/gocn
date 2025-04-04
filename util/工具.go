package util

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

// UUIDString 生成32位UUID
func UUIDString() string {
	u := uuid.New()
	var buf [32]byte
	hex.Encode(buf[:], u[:])
	return string(buf[:])
}

// GetNanoTimeLabel 获取唯一的时间标识(21位)
func GetNanoTimeLabel() (cur string) {
	now := time.Now()
	cur = fmt.Sprintf("%s%09d", now.Format("060102150405"), now.UnixNano()%1e9)
	return
}

// Base64X
func Base64X(fileName string) string {
	ff, _ := os.Open(fileName)
	defer ff.Close()
	buf := make([]byte, 600000)
	n, _ := ff.Read(buf)
	res := base64.StdEncoding.EncodeToString(buf[:n])
	return res
}
