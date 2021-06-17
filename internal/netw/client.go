package netw

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net"

	"gitlab.com/4thlabs/perfmon/internal/recording"
)

type Client struct {
	addr net.UDPAddr
	conn *net.UDPConn
}

func NewClient() *Client {
	return &Client{}
}

func (client *Client) Connect(address string, port int) error {
	// addr := net.UDPAddr{Port: port, IP: net.ParseIP(address)}
	// conn, err := net.DialUDP("udp", nil, &addr)

	// if err != nil {
	// 	return err
	// }

	// client.addr = addr
	// client.conn = conn

	return nil
}

func encryt(key []byte, data []byte) []byte {
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

func getHmac(key []byte, data []byte) string {
	mac := hmac.New(sha1.New, key)

	mac.Write(data)

	sha := hex.EncodeToString(mac.Sum(nil))

	return sha
}

func (client *Client) Start(recording *recording.Recording) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	for i := 0; i < 2000; i++ {

		go func(idx int) {
			addr := net.UDPAddr{Port: 1000 + idx, IP: net.ParseIP("10.11.3.16")}
			conn, err := net.DialUDP("udp", nil, &addr)

			if err != nil {
				log.Fatalln(err)
			}

			defer conn.Close()

			for {
				frame := recording.GetInMemoryFrame()

				//for y := 0; y < 100; y++ {
				encryt(key, frame.Data)
				getHmac(key, frame.Data)
				//}

				_, err = conn.Write(frame.Data)

				if err != nil {
					//fmt.Println(err)
				}
			}
		}(i)
	}
}

func (client *Client) Close() {
	//client.conn.Close()
}
