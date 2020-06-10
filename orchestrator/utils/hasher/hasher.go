package hasher

import (
	"crypto/md5"
	"fmt"
	"io"
)

// MD5Hash A function to return an MD5 hash of a given string
func MD5Hash(str string) string {
	hasher := md5.New()
	io.WriteString(hasher, str)

	return fmt.Sprintf("%x", hasher.Sum(nil))
}
