package blockchain

// BlockChain is a basic structure for a blockchain
type BlockChain struct {
	Blocks []*Block
}

// InitBlockChain creates a new blockchain
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

// AddBlock adds a block to the blockchain
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}

