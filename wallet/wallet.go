package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	//four bytes
	checksumLength = 4
	version        = byte(0x00)
)

//Package ecdsa implements the Elliptic Curve Digital Signature Algorithm, as defined in FIPS 186-3.
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionHash := append([]byte{version}, pubHash...)

	checksum := Checksum(versionHash)

	fullHash := append(versionHash, checksum...)

	address := Base58Encode(fullHash)

	fmt.Printf("pub key:%x\n", w.PublicKey)
	fmt.Printf("pub hash:%x\n", pubHash)
	fmt.Printf("address:%s\n", address)

	return address
}

// Address  : 1LzkbzcmkkddFxRG46uJSddXxhF1otCp6n
// FullHash : 00248bd9e7a4445245252452624526kv1451vj145h6155
// [Version]: 00
// [Pub Key Hash] 248bd9eu9734715713nv873m51574n5105n
// [CheckSum] 145h6155
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	fmt.Println("pubKeyHash", pubKeyHash)
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	fmt.Println("actualChecksum", actualChecksum)
	version := pubKeyHash[0]
	fmt.Println("version", version)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	fmt.Println("pubKeyHash", pubKeyHash)
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))
	fmt.Println("targetChecksum", targetChecksum)

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	//generate the private key using the curve
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//get the publicKey(this the coordinate of the final point when move the dot D times,D is in privateKey)
	//// PrivateKey represents an ECDSA private key.
	// type PrivateKey struct {
	// 	PublicKey
	// 	D *big.Int
	// }
	//
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)
	//to generate a digest with 160 bit
	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}
	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}
