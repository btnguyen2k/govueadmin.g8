package gvabe

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"hash"
	"io"
	"strconv"
	"time"
)

// RsaMode defines RSA encryption modes
//
// available since template-v0.2.0
type RsaMode int

const (
	RsaModeAuto RsaMode = iota
	RsaModeOAEP
	RsaModePKCS1v15
)

var (
	ErrKeySizeTooSmall = errors.New("key size is too small")
	ErrInvalidRsaMode  = errors.New("invalid RSA mode or mode not supported")
)

// available since template-v0.2.0
type RsaChunkEncryptionFunc func(hash hash.Hash, random io.Reader, pub *rsa.PublicKey, msg []byte) ([]byte, error)

// available since template-v0.2.0
func _rsaEncryptOAEP(hash hash.Hash, random io.Reader, pub *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptOAEP(hash, random, pub, msg, nil)
}

// available since template-v0.2.0
func _rsaEncryptPKCS1v15(_ hash.Hash, random io.Reader, pub *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(random, pub, msg)
}

// available since template-v0.2.0
func rsaEncrypt(rsaMode RsaMode, data []byte, rsaPubKey *rsa.PublicKey) ([]byte, error) {
	if rsaMode == RsaModeAuto {
		rsaMode = RsaModeOAEP
	}
	var hf hash.Hash
	var buffSize int
	var rsaFunc RsaChunkEncryptionFunc
	switch rsaMode {
	case RsaModeOAEP:
		hf = sha256.New()
		buffSize = rsaPubKey.Size() - 2*hf.Size() - 2
		rsaFunc = _rsaEncryptOAEP
	case RsaModePKCS1v15:
		hf = nil
		buffSize = rsaPubKey.Size() - 11
		rsaFunc = _rsaEncryptPKCS1v15
	default:
		return nil, ErrInvalidRsaMode
	}
	if buffSize < 1 {
		return nil, ErrKeySizeTooSmall
	}
	result := make([]byte, 0)
	for i, n := 0, len(data); i < n; i += buffSize {
		end := i + buffSize
		if end > n {
			end = n
		}
		output, err := rsaFunc(hf, rand.Reader, rsaPubKey, data[i:end])
		if err != nil {
			return nil, err
		}
		result = append(result, output...)
	}
	return result, nil
}

// available since template-v0.2.0
func rsaDecrypt(rsaMode RsaMode, encdata []byte, rsaPrivKey *rsa.PrivateKey) ([]byte, error) {
	if len(encdata)%rsaPrivKey.Size() != 0 {
		return nil, rsa.ErrDecryption
	}
	if rsaMode == RsaModeAuto {
		rsaMode = RsaModeOAEP
	}
	var opts crypto.DecrypterOpts = nil
	if rsaMode == RsaModeOAEP {
		opts = &rsa.OAEPOptions{Hash: crypto.SHA256}
	}
	result := make([]byte, 0)
	for i, n, k := 0, len(encdata), rsaPrivKey.Size(); i < n; i += k {
		buff, err := rsaPrivKey.Decrypt(nil, encdata[i:i+k], opts)
		if err != nil {
			return nil, err
		}
		result = append(result, buff...)
	}
	return result, nil
}

// available since template-v0.2.0
func genRsaKey(numBits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, numBits)
}

// available since template-v0.2.0
func parseRsaPublicKeyFromPem(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	switch block.Type {
	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)
	case "PUBLIC KEY":
		if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
			return nil, err
		} else {
			switch pub := pub.(type) {
			case *rsa.PublicKey:
				return pub, nil
			}
		}
	}
	return nil, errors.New("not RSA public key")
}

// padRight adds "0" right right of a string until its length reach a specific value.
func padRight(str string, l int) string {
	for len(str) < l {
		str += "0"
	}
	return str
}

// aesEncrypt encrypts a block of data using AES/CTR mode.
//
// IV is put at the beginning of the cipher data.
func aesEncrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := []byte(padRight(strconv.FormatInt(time.Now().UnixNano(), 16), 16))
	cipherData := make([]byte, 16+len(data))
	copy(cipherData, iv)
	ctr := cipher.NewCTR(block, iv)
	ctr.XORKeyStream(cipherData[16:], data)
	return cipherData, nil
}

// aesDecrypt decrypts a block of encrypted data using AES/CTR mode.
//
// Assuming IV is put at the beginning of the cipher data.
func aesDecrypt(key, encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := encryptedData[0:16]
	data := make([]byte, len(encryptedData)-16)
	ctr := cipher.NewCTR(block, iv)
	ctr.XORKeyStream(data, encryptedData[16:])
	return data, nil
}
