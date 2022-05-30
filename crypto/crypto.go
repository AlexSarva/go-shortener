package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
)

//var secretkey = []byte("Ag@th@")

func toByteArray(i int) (arr [4]byte) {
	binary.BigEndian.PutUint32(arr[0:4], uint32(i))
	return
}

func Encrypt(id int, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	workID := toByteArray(id)
	h.Write(workID[:])
	dst := h.Sum(nil)
	var fullCookie []byte
	fullCookie = append(fullCookie, workID[:]...)
	fullCookie = append(fullCookie, dst...)
	//fmt.Println("ID: ", workID)
	//fmt.Println("SIGN: ", dst)
	//fmt.Println(fullCookie)
	return hex.EncodeToString(fullCookie)
}

func Decrypt(cookie string, secret []byte) (int, error) {
	var (
		data []byte // декодированное сообщение с подписью
		id   uint32 // значение идентификатора
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)

	data, err = hex.DecodeString(cookie)
	if err != nil {
		panic(err)
	}
	id = binary.BigEndian.Uint32(data[:4])
	h := hmac.New(sha256.New, secret)
	h.Write(data[:4])
	sign = h.Sum(nil)

	if hmac.Equal(sign, data[4:]) {
		return int(id), nil
	} else {
		return 0, errors.New("подпись не верна")
	}
}
