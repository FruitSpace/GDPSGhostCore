package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

type ThunderAES struct {
	CryptKey []byte
	Blk      cipher.Block
	GCM      cipher.AEAD
}

func (taes *ThunderAES) GenKey(s string) error {
	sha := sha256.New()
	_, err := sha.Write([]byte(s))
	taes.CryptKey = sha.Sum(nil)[:32]
	return err
}

func (taes *ThunderAES) Init() error {
	var err error
	taes.Blk, err = aes.NewCipher(taes.CryptKey)
	if err != nil {
		return err
	}
	taes.GCM, err = cipher.NewGCM(taes.Blk)
	return err
}

func (taes *ThunderAES) Encrypt(text string) (string, error) {
	nonce := make([]byte, taes.GCM.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	raw := taes.GCM.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(raw), err
}

func (taes *ThunderAES) Decrypt(blk string) (string, error) {
	block, err := base64.StdEncoding.DecodeString(blk)
	if err != nil {
		return "", err
	}
	nonce, raw := block[:taes.GCM.NonceSize()], block[taes.GCM.NonceSize():]
	plain, err := taes.GCM.Open(nil, nonce, raw, nil)
	return string(plain), err
}

func (taes *ThunderAES) DecryptLegacy(blk string) (string, error) {
	block, err := base64.StdEncoding.DecodeString(blk)
	if err != nil {
		return "", err
	}
	nSize := taes.GCM.NonceSize()
	nonce := make([]byte, nSize)
	tag := make([]byte, 16)
	copy(nonce, block[:nSize])
	copy(tag, block[nSize:nSize+16])
	raw := append(block[nSize+16:], tag...)
	plain, err := taes.GCM.Open(nil, nonce, raw, nil)
	return string(plain), err
}

func (taes *ThunderAES) EncryptRaw(text string) ([]byte, error) {
	nonce := make([]byte, taes.GCM.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	raw := taes.GCM.Seal(nonce, nonce, []byte(text), nil)
	return raw, err
}

func (taes *ThunderAES) DecryptRaw(block []byte) (string, error) {
	nonce, raw := block[:taes.GCM.NonceSize()], block[taes.GCM.NonceSize():]
	plain, err := taes.GCM.Open(nil, nonce, raw, nil)
	return string(plain), err
}
