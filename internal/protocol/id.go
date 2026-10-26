package protocol

import (
    "crypto/rand"
    "encoding/hex"
)

// newID returns a 128-bit random ID as 32 hex characters.
func NewID() string {
    b := make([]byte, 16)
    if _, err := rand.Read(b); err != nil {
        panic("crypto/rand failed: " + err.Error()) // unreachable in practice
    }
    return hex.EncodeToString(b)
}
