package uuidx

import (
	"crypto/rand"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/oklog/ulid"
	"strings"
)

func UID() string {
	uid, _ := uuid.NewV4()
	return strings.ReplaceAll(uid.String(), "-", "")
}

func ULID() string {
	// 生成ULID
	id, err := ulid.New(ulid.Now(), ulid.Monotonic(rand.Reader, 0))
	if err != nil {
		fmt.Println("Error generating ULID:", err)
		return ""
	}
	return id.String()
}
