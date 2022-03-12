package hash

import (
	"crypto/sha256"
	"fmt"
)

func Hash(word string) []byte {
	h := sha256.New()
	h.Write([]byte(word))
	bs := h.Sum(nil)

	fmt.Printf("hashed %s\n", word)

	return bs
}
