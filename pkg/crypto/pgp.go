package crypto

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/openpgp"
)

// EncryptPGP шифрует строку открытым ключом (entity) и возвращает armored текст
func EncryptPGP(plaintext string, entity *openpgp.Entity) (string, error) {
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		return "", err
	}
	_, err = w.Write([]byte(plaintext))
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}
	if buf.Len() == 0 {
		return "", errors.New("encryption resulted in empty data")
	}
	return buf.String(), nil
}
