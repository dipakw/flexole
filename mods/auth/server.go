package auth

import (
	"crypto/mlkem"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

func Server(clientConn net.Conn, args *ServerOpts) *Auth {
	res := &Auth{}

	// Step 1: Generate private key.
	privkey, err := mlkem.GenerateKey1024()

	if err != nil {
		return res.re(&Err{
			reason: "failed to generate a private key",
			err:    err,
		})
	}

	// Step 2: Send public key to the client.
	pubkeyb := privkey.EncapsulationKey().Bytes()
	n, err := clientConn.Write(pubkeyb)

	if err != nil || n != mlkem.EncapsulationKeySize1024 {
		return res.re(&Err{
			reason: "failed to send public key to the client",
			err:    err,
		})
	}

	// Step 3: Receive the ciphertext.
	buf, n, err := readonce(clientConn, mlkem.CiphertextSize1024, &readopts{
		timeout: args.Timeout,
		full:    true,
	})

	if err != nil || n != mlkem.CiphertextSize1024 {
		return res.re(&Err{
			reason: "failed to get ciphertext from the client",
			err:    err,
		})
	}

	// Step 4: Decapsulate the chipertext.
	res.Key, err = privkey.Decapsulate(buf)

	if err != nil {
		return res.re(&Err{
			reason: "failed to decapsulate the ciphertext",
			err:    err,
		})
	}

	// Step 5: Send ACK.
	msg := []byte{0, 8, 0, 8}
	msg = append(msg, randbytes(6)...)
	msg, err = res.Encrypt(msg)

	if err != nil {
		return res.re(&Err{
			reason: "failed to encrypt the ack message",
			err:    err,
		})
	}

	n, err = clientConn.Write(msg)

	if err != nil || n != len(msg) {
		return res.re(&Err{
			reason: "failed to send the ack message",
			err:    err,
		})
	}

	// Step 6: Receive the ID and meta data.
	buf, n, err = readonce(clientConn, MAX_ID_META_SIZE, &readopts{
		timeout: args.Timeout,
	})

	if err != nil {
		return res.re(&Err{
			reason: "failed to receive the ID and meta data",
			err:    err,
		})
	}

	idme, err := res.Decrypt(buf[:n])

	if err != nil {
		return res.re(&Err{
			reason: "failed to decrypt the ID and meta data",
			err:    err,
		})
	}

	if len(idme) < 16 {
		return res.re(&Err{
			reason: "ID and meta data message is too short",
			err:    errors.New("ID and meta data message is too short"),
		})
	}

	id, meta, err := decodeIdMeta(idme)

	if err != nil {
		return res.re(&Err{
			reason: "failed to decode the ID and meta data",
			err:    err,
		})
	}

	res.ID = id
	res.Meta = meta

	// Step 7: Send the challenge.
	challenge := randbytes(CHALLENGE_SIZE)
	encryptedChlng, err := res.Encrypt(challenge)

	if err != nil {
		return res.re(&Err{
			reason: "failed to encrypt the challenge message",
			err:    errors.New("failed to encrypt the challenge message"),
		})
	}

	n, err = clientConn.Write(encryptedChlng)

	if err != nil || n != len(encryptedChlng) {
		return res.re(&Err{
			reason: "failed to write the challenge message",
			err:    err,
		})
	}

	// Step 8: Get the encrypted signed message size.
	buf, n, err = readonce(clientConn, 2, &readopts{
		timeout: args.Timeout,
		full:    true,
	})

	if err != nil || n != 2 {
		return res.re(&Err{
			reason: "failed to get the encrypted message size",
			err:    err,
		})
	}

	encsize := binary.BigEndian.Uint16(buf)

	if encsize < args.MinSigSize || encsize > args.MaxSigSize {
		return res.re(&Err{
			reason: "received invalid signature size",
			err:    fmt.Errorf("received: %d, min: %d, max: %d", encsize, args.MinSigSize, args.MaxSigSize),
		})
	}

	// Step 8: Get the signed message.
	buf, n, err = readonce(clientConn, int(encsize), &readopts{
		timeout: args.Timeout,
		full:    true,
	})

	if err != nil {
		return res.re(&Err{
			reason: "failed to get the signature",
			err:    err,
		})
	}

	dsig, err := res.Decrypt(buf[:n])

	if err != nil {
		return res.re(&Err{
			reason: "failed to decrypt the signature",
			err:    err,
		})
	}

	// Step 9: Verify the signature
	if ok, err := args.VerifySig(res, challenge, dsig); !ok {
		if err == nil {
			err = errors.New("signatures didn't match")
		}

		return res.re(&Err{
			reason: "failed to verify the signature",
			err:    err,
		})
	}

	// Step 10: Send the confirmation
	cnfm, err := res.Encrypt(challenge)

	if err != nil {
		return res.re(&Err{
			reason: "failed to encrypt the confirmation message",
			err:    err,
		})
	}

	n, err = clientConn.Write(cnfm)

	if err != nil || n != len(cnfm) {
		return res.re(&Err{
			reason: "failed to send the confirmation message",
			err:    err,
		})
	}

	return res
}
