package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"flag"
	"io"
	"os"
)

func main() {
	// Parse flags
	decryptVar := flag.Bool("d", false, "Decrypt mode")
	flag.Parse()
	decrypt := *decryptVar

	// Setup AES
	key := []byte("pecman123")

	for len(key) < 32 {
		key = append(key, byte(0))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Create cipher writer
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	if decrypt {
		reader := base64.NewDecoder(base64.StdEncoding, os.Stdin)
		transcoder := &cipher.StreamReader{S: stream, R: reader}

		// Pump
		if _, err := io.Copy(os.Stdout, transcoder); err != nil {
			panic(err)
		}
	} else {
		writer := base64.NewEncoder(base64.StdEncoding, os.Stdout)
		transcoder := &cipher.StreamWriter{S: stream, W: writer}
		defer transcoder.Close()

		// Pump
		if _, err := io.Copy(transcoder, os.Stdin); err != nil {
			panic(err)
		}
	}
}
