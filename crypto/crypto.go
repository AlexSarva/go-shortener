package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"log"
)

//var secretkey = []byte("Ag@th@")

var NotValidSingERR = errors.New("sign is not valid")

func Encrypt(uuid uuid.UUID, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write(uuid[:])
	dst := h.Sum(nil)
	var fullCookie []byte
	fullCookie = append(fullCookie, uuid[:]...)
	fullCookie = append(fullCookie, dst...)
	//fmt.Println("ID: ", uuid)
	//fmt.Println("SIGN: ", dst)
	//fmt.Println(fullCookie)
	return hex.EncodeToString(fullCookie)
}

func Decrypt(cookie string, secret []byte) (uuid.UUID, error) {
	var (
		data []byte    // декодированное сообщение с подписью
		id   uuid.UUID // значение идентификатора
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)

	data, err = hex.DecodeString(cookie)
	if err != nil {
		log.Println(err)
		return uuid.New(), NotValidSingERR
	}
	id, idErr := uuid.FromBytes(data[:16])
	if idErr != nil {
		log.Println(idErr)
	}
	h := hmac.New(sha256.New, secret)
	h.Write(data[:16])
	sign = h.Sum(nil)

	if hmac.Equal(sign, data[16:]) {
		return id, nil
	} else {
		return uuid.New(), NotValidSingERR
	}
}
