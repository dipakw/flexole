package auth

import (
	"bytes"
	"crypto/mlkem"
	"errors"
	"fmt"
	"net"
)

func Client(serverConn net.Conn, args *ClientOpts) *Auth {
	res := &Auth{}

	// Step 1: Get the public key
	buf, n, err := readonce(serverConn, mlkem.EncapsulationKeySize1024, &readopts{
		timeout: args.Timeout,
		full:    true,
	})

	if err != nil || n != mlkem.EncapsulationKeySize1024 {
		return res.re(&Err{
			reason: "failed to the the public key",
			err:    err,
		})
	}

	pubkey, err := mlkem.NewEncapsulationKey1024(buf)

	if err != nil {
		return res.re(&Err{
			reason: "failed to load the public key",
			err:    err,
		})
	}

	// Step 2: Send the ciphertext
	enckey, ct := pubkey.Encapsulate()

	n, err = serverConn.Write(ct)

	if err != nil || n != len(ct) {
		return res.re(&Err{
			reason: "failed to send the ciphertext",
			err:    err,
		})
	}

	res.Key = enckey

	// Step 3: Receive ACK
	buf, n, err = readonce(serverConn, 256, &readopts{
		timeout: args.Timeout,
	})

	if err != nil {
		return res.re(&Err{
			reason: "failed to receive the ACK",
			err:    err,
		})
	}

	ackm, err := res.Decrypt(buf[:n])

	if err != nil {
		return res.re(&Err{
			reason: "failed to decrypt the ACK",
			err:    err,
		})
	}

	if !bytes.Equal(ackm[0:4], []byte{0, 8, 0, 8}) {
		return res.re(&Err{
			reason: "invalid ACK",
			err:    errors.New("invalid ACK"),
		})
	}

	// Step 4: Send ID and Meta
	idm := encodeIdMeta(args.ID, args.Meta)
	idme, err := res.Encrypt(idm)

	if err != nil {
		return res.re(&Err{
			reason: "failed to encrypt the ID and meta data",
			err:    err,
		})
	}

	n, err = serverConn.Write(idme)

	if err != nil || n != len(idme) {
		return res.re(&Err{
			reason: "failed to send ID and meta data",
			err:    err,
		})
	}

	// Step 5: Get the challenge
	buf, n, err = readonce(serverConn, 128, &readopts{
		timeout: args.Timeout,
	})

	if err != nil {
		return res.re(&Err{
			reason: "failed to get the challenge",
			err:    err,
		})
	}

	chnm, err := res.Decrypt(buf[:n])

	if err != nil {
		return res.re(&Err{
			reason: "failed to decrypt the challenge",
			err:    err,
		})
	}

	if len(chnm) != CHALLENGE_SIZE {
		return res.re(&Err{
			reason: "challenge message is too short",
			err:    fmt.Errorf("received %d/%d", len(chnm), CHALLENGE_SIZE),
		})
	}

	// Step 6: Sign the challenge
	sig, err := args.SignMsg(chnm)

	if err != nil {
		return res.re(&Err{
			reason: "failed to sign the challenge",
			err:    err,
		})
	}

	encsig, err := res.Encrypt(sig)

	if err != nil {
		return res.re(&Err{
			reason: "failed to encrypt the signature",
			err:    err,
		})
	}

	// Step 7: Send the signature
	sized := make([]byte, 2+len(encsig))

	copy(sized, uint16bytes(uint16(len(encsig))))
	copy(sized[2:], encsig)

	n, err = serverConn.Write(sized)

	if err != nil || n != len(sized) {
		return res.re(&Err{
			reason: "failed to send the signature",
			err:    err,
		})
	}

	// Step 8: Get the confirmation
	buf, n, err = readonce(serverConn, 128, &readopts{
		timeout: args.Timeout,
	})

	if err != nil {
		return res.re(&Err{
			reason: "failed to receive the confirmation",
			err:    err,
		})
	}

	dcnf, err := res.Decrypt(buf[:n])

	if err != nil {
		return res.re(&Err{
			reason: "failed to decrypt the confirmation",
			err:    err,
		})
	}

	if len(dcnf) != CHALLENGE_SIZE {
		return res.re(&Err{
			reason: "confirmation message is too short",
			err:    fmt.Errorf("received %d/%d", len(chnm), CHALLENGE_SIZE),
		})
	}

	// Step 9: Verify the confirmation
	if !bytes.Equal(dcnf, chnm) {
		return res.re(&Err{
			reason: "invalid confirmation",
			err:    errors.New("challenges do not match"),
		})
	}

	return res
}
