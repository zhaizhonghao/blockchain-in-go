package blockchain

type Blockchain struct {
	Blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

//To create the gensis block
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	//mining
	nonce, hash := pow.Run()
	//mining successfully
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}
