package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateID() string {
	n := strconv.Itoa(int(time.Now().Unix()))
	r := strconv.Itoa(rand.Int())
	hash := sha256.Sum256([]byte(r))
	readable := hex.EncodeToString(hash[:2])
	id := fmt.Sprintf("%s%s", n, readable)
	return id
}

func GenerateKey() string {
	id := GenerateID()
	hash := sha256.Sum256([]byte(id))
	key := hex.EncodeToString(hash[:16])
	return key
}

func ToJson(v interface{}) string {
	r, err := json.Marshal(&v)
	if err != nil {
		return ""
	}
	return string(r)
}

func Where(clause, column string, values []string) string {
	b := strings.Builder{}
	switch clause {
	case "$in":
		b.WriteString(fmt.Sprintf(`{%s:{"$in":[`, column))
		for i, val := range values {
			b.WriteString(fmt.Sprintf(`"%s"`, val))
			if i != len(values)-1 {
				b.WriteString(",")
			}
		}
		b.WriteString("]}}")
	default:
		b.WriteString(fmt.Sprintf(`{"%s":{"%s":"%s"}}`, column, clause, values[0]))
	}
	return b.String()
}
