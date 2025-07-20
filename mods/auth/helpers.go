package auth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"
)

// Encrypt encrypts the given message with the provided key using AES-GCM.
func encrypt(key, msg []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, msg, nil)
	return ciphertext, nil
}

// Decrypt decrypts the given ciphertext with the provided key using AES-GCM.
func decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < 12 {
		return nil, errors.New("ciphertext too short")
	}

	nonce, encryptedMsg := ciphertext[:12], ciphertext[12:]
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, nonce, encryptedMsg, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func randbytes(size int) []byte {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		return nil
	}
	return key
}

func readonce(conn net.Conn, size int, opts *readopts) ([]byte, int, error) {
	opt := &readopts{
		timeout: 0,
		full:    false,
	}

	if opts != nil {
		opt = opts
	}

	buf := make([]byte, size)

	if opt.timeout > 0 {
		conn.SetReadDeadline(time.Now().Add(opt.timeout))
		defer conn.SetDeadline(time.Time{})
	}

	var n int
	var err error

	if opt.full {
		n, err = io.ReadFull(conn, buf)
	} else {
		n, err = conn.Read(buf)
	}

	if err != nil {
		return nil, 0, err
	}

	return buf, n, nil
}

func uint16bytes(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return b
}

func bytesToUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

// Format:
// [id-len:uint8][id][key-len:uint8][key][value-len:uint16][value]
func encodeIdMeta(id []byte, meta map[string]string) []byte {
	buf := bytes.Buffer{}

	buf.WriteByte(uint8(len(id)))
	buf.Write(id)

	for key, value := range meta {
		buf.WriteByte(uint8(len(key)))
		buf.WriteString(key)
		buf.Write(uint16bytes(uint16(len(value))))
		buf.WriteString(value)
	}

	return buf.Bytes()
}

func decodeIdMeta(b []byte) ([]byte, map[string]string, error) {
	err := errors.New("malformed data")

	if len(b) < 2 {
		return nil, nil, err
	}

	size := len(b)
	meta := map[string]string{}

	idx := 0
	end := 0

	// Get ID length.
	idLen := int(b[idx])
	idx += 1

	// Get the ID.
	end = idx + idLen

	if end > size {
		return nil, nil, err
	}

	id := b[idx:end]
	idx += len(id)

	if idx > size {
		return nil, nil, err
	}

	for {
		// Get the key len.
		keyLen := int(b[idx])
		idx += 1

		// Get the key
		end = idx + keyLen

		if end > size {
			return nil, nil, err
		}

		key := b[idx:end]
		idx += len(key)

		// Get the value len.
		end = idx + 2

		if end > size {
			return nil, nil, err
		}

		valLen := int(bytesToUint16(b[idx:end]))
		idx += 2

		// Get the value.
		end = idx + valLen

		if end > size {
			return nil, nil, err
		}

		val := b[idx:end]
		idx += len(val)

		meta[string(key)] = string(val)

		if idx == size {
			break
		}
	}

	return id, meta, nil
}
