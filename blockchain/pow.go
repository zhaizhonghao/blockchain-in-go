package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

//Take the data from the block

//create a counter(nonce) which starts at 0

//create a hash of the data plus the counter

//check the hash to see if it meets a set of requirement

//Requirements:
//The First few bytes must contain 0s

//The difficulty cannot be modifed once the blockchain has already stored in the database
//If the protocol want to change the difficulty, it must replay the existed blockchain
const Difficulty = 18

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	//func (z *Int) Lsh(x *Int, n uint) *Int
	//Lsh sets z = x << n and returns z.
	//bigger difficulty produce the smaller target
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	//MaxUint64 = 1<<64 - 1
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\rmining:%x", hash)
		//SetBytes interprets buf as the bytes of a big-endian unsigned integer,
		//sets z to that value, and returns z.
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			//mining successfully
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	//A Buffer is a variable-sized buffer of bytes with Read and Write methods.
	//The zero value for Buffer is an empty buffer ready to use.
	buff := new(bytes.Buffer)
	//func Write(w io.Writer, order ByteOrder, data interface{}) error
	//Write writes the binary representation of data into w.
	//Data must be a fixed-size value or a slice of fixed-size values, or a pointer to such data.
	//Boolean values encode as one byte: 1 for true, and 0 for false.
	//Bytes written to w are encoded using the specified byte order and read from successive fields of the data.
	//When writing structs, zero values are written for fields with blank (_) field names.
	//大端模式Big-Endian就是高位字节排放在内存的低地址端，低位字节排放在内存的高地址端。
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
