package netw

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
)

func Encryt(key []byte, data []byte) []byte {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	return aesgcm.Seal(nil, nonce, data, nil)
}

func Hmac(key []byte, data []byte) string {
	mac := hmac.New(sha1.New, key)

	mac.Write(data)

	sha := hex.EncodeToString(mac.Sum(nil))

	return sha
}

// Create a single random initialised byte array of size.
func GenerateNonce(size int) ([]byte, error) {

	b := make([]byte, size)

	// not checking len here because rand.Read doc reads:
	//             On return, n == len(b) if and only if err == nil.
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func Copy(src []byte, srcI int, dest []byte, destI int, copyLen int) {
	srcI2 := srcI + copyLen
	copy(dest[destI:], src[srcI:srcI2])
}

func EncryptGCM(encKey, input []byte) ([]byte, error) {

	encKeyLen := len(encKey)

	if encKeyLen < 16 {
		return nil, fmt.Errorf("The key must be 16 bytes long")
	}

	encKeySized := encKey

	if encKeyLen > 16 {
		encKeySized = encKey[:16]
	}

	c, err := aes.NewCipher(encKeySized)

	if err != nil {
		return nil, err
	}

	//----------- Create the IV

	// remember that GCM normally takes a 12 byte (96 bit) nounce
	nonceSize := 12
	iv, err := GenerateNonce(nonceSize)
	if err != nil {
		return nil, err
	}

	//----------- Encrypt

	ivLen := len(iv)
	enc, err := cipher.NewGCMWithNonceSize(c, nonceSize)

	if err != nil {
		return nil, err
	}

	cipherText := enc.Seal(nil, iv, input, nil)

	//----------- Pack the message

	// create output tag
	output := make([]byte, 1+ivLen+len(cipherText))

	i := 0
	output[i] = byte(ivLen)
	i++
	Copy(iv, 0, output, i, ivLen)
	i += ivLen

	Copy(cipherText, 0, output, i, len(cipherText))

	return output, nil
}
