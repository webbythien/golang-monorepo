package id

import (
	"strconv"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// defaultAlphabet is the alphabet used for ID characters by default.
var defaultAlphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

func NewNanoID(length int) (string, error) {
	return gonanoid.Generate(defaultAlphabet, length)
}

func MustGenerateNanoID(length int) string {
	id, err := gonanoid.Generate(defaultAlphabet, length)
	if err != nil {
		panic(err)
	}
	return id
}

func NewGlobalUID(entity string, source string, localID string) string {
	var sb strings.Builder
	sb.WriteString(entity)
	sb.WriteString("::")
	sb.WriteString(source)
	sb.WriteString("::")
	sb.WriteString(localID)
	return sb.String()
}

type LocalID string

func GetLocalID(globalID string) LocalID {
	return LocalID(globalID[strings.LastIndex(globalID, "::"):])
}

func ToInt(id LocalID) int {
	res, _ := strconv.Atoi(string(id))
	return res
}
