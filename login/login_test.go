package login

import (
	"crypto/sha512"
	"encoding/hex"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordMatch(t *testing.T) {

	password1 := "s3cret222"
	hashFunc := sha512.New384()
	hashFunc.Write([]byte(password1))
	hexStr := hex.EncodeToString(hashFunc.Sum(nil))

	match, err := doPasswordMatch(password1, hexStr, "SHA384", "")
	assert.NoError(t, err)
	assert.True(t, match, "Password should match")

	match, err = doPasswordMatch("anotherPassword123", hexStr, "SHA384", "")
	assert.NoError(t, err)
	assert.False(t, match, "Password should not match")
}
