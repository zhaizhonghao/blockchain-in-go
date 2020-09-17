package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

//base58 is different from the base64, it removes the following characters:
//0 O 1 I + /

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}

	return decode
}


