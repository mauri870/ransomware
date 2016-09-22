package crypto

import (
	"strconv"
	"testing"
)

var (
	stubs = []string{"Hello World", "The Quick Brown Fox", "The difficult we do immediately; the impossible takes a little longer."}
	key   = []byte("hDmPPK2b76ROBqd4uWQcu0ruUyz0tVXd")
)

func TestEncryptDecrypt(t *testing.T) {
	for _, text := range stubs {
		ciphertext, err := Encrypt(key, []byte(text))
		if err != nil {
			t.Error(err)
		}

		decryptedText, err := Decrypt(key, ciphertext)
		if err != nil {
			t.Error(err)
		}

		if string(decryptedText) != text {
			t.Errorf("Expect %s but got %s on decryption", text, decryptedText)
		}
	}
}

func BenchmarkEncryptDecrypt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ciphertext, _ := Encrypt(key, []byte("bench"+strconv.Itoa(n)))
		Decrypt(key, ciphertext)
	}
}
