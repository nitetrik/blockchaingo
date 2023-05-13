package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

type SmartContract struct {
	Code string
}

type Transaction struct {
	Sender   string
	Receiver string
	Amount   float64
}

type Block struct {
	Index     int
	Data      []*Transaction
	Contracts []*SmartContract
	Hash      string
	PrevHash  string
	Timestamp time.Time
	Nonce     int
}

type Blockchain struct {
	Blocks     []*Block
	Difficulty int
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks:     []*Block{NewGenesisBlock()},
		Difficulty: 3,
	}
}

func NewGenesisBlock() *Block {
	return &Block{
		Index:     0,
		Data:      []*Transaction{},
		Hash:      "",
		PrevHash:  "",
		Timestamp: time.Now(),
		Nonce:     0,
	}
}

func NewBlock(data []*Transaction, contracts []*SmartContract, prevHash string) *Block {
	return &Block{
		Index:     0,
		Data:      data,
		Contracts: contracts,
		Hash:      "",
		PrevHash:  prevHash,
		Timestamp: time.Now(),
		Nonce:     0,
	}
}

func (b *Block) MineBlock(difficulty int) {
	for {
		b.Nonce = rand.Intn(1000000000)
		b.Hash = b.CalculateHash()
		if b.IsValid(difficulty) {
			break
		}
	}
}

func (b *Block) CalculateHash() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(b.String())))
}

func (b *Block) IsValid(difficulty int) bool {
	return b.Hash[:difficulty] == strings.Repeat("0", difficulty)
}

func (bc *Blockchain) AddBlock(b *Block) {
	bc.Blocks = append(bc.Blocks, b)
}

func (bc *Blockchain) GetLatestBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) Print() {
	for _, b := range bc.Blocks {
		fmt.Println(b)
	}
}

type Validator struct {
	Blockchain *Blockchain
}

func NewValidator(blockchain *Blockchain) *Validator {
	return &Validator{
		Blockchain: blockchain,
	}
}

func (v *Validator) Validate(block *Block) error {
	if !block.IsValid(v.Blockchain.Difficulty) {
		return fmt.Errorf("Invalid block hash")
	}

	if block.PrevHash != v.Blockchain.GetLatestBlock().Hash {
		return fmt.Errorf("Invalid previous hash")
	}

	if len(block.Data) == 0 {
		return fmt.Errorf("Empty block data")
	}

	if block.Nonce == 0 {
		return fmt.Errorf("Invalid nonce")
	}

	return nil
}

func TestCreateBlockchain(t *testing.T) {
	bc := NewBlockchain()

	// Assert that the blockchain is not nil
	if bc == nil {
		t.Errorf("Failed to create a new blockchain")
	}

	// Assert that the blockchain contains the genesis block
	if len(bc.Blocks) != 1 {
		t.Errorf("Expected the blockchain to contain the genesis block only")
	}
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockchain()

	// Create transactions
	transactions := []*Transaction{
		{Sender: "Alice", Receiver: "Bob", Amount: 5.0},
		{Sender: "Bob", Receiver: "Charlie", Amount: 2.5},
	}

	// Create smart contracts
	contracts := []*SmartContract{
		{Code: "Smart Contract 1"},
		{Code: "Smart Contract 2"},
	}

	b := NewBlock(transactions, contracts, bc.GetLatestBlock().Hash)
	b.MineBlock(bc.Difficulty)
	bc.AddBlock(b)

	// Assert that the blockchain contains two blocks (genesis block + new block)
	if len(bc.Blocks) != 2 {
		t.Errorf("Expected the blockchain to contain two blocks")
	}

	// Assert that the new block is the latest block
	if bc.GetLatestBlock() != b {
		t.Errorf("Expected the latest block to be the new block")
	}
}

func TestMineBlock(t *testing.T) {
	bc := NewBlockchain()

	// Create transactions
	transactions := []*Transaction{
		{Sender: "Alice", Receiver: "Bob", Amount: 5.0},
	}

	b := NewBlock(transactions, bc.GetLatestBlock().Hash)
	b.MineBlock(bc.Difficulty)

	// Assert that the block's hash meets the difficulty requirement
	if !b.IsValid(bc.Difficulty) {
		t.Errorf("Failed to mine the block")
	}
}

func TestValidateBlock(t *testing.T) {
	bc := NewBlockchain()

	// Create transactions
	transactions := []*Transaction{
		{Sender: "Alice", Receiver: "Bob", Amount: 5.0},
	}

	b := NewBlock(transactions, bc.GetLatestBlock().Hash)
	b.MineBlock(bc.Difficulty)

	validator := NewValidator(bc)
	err := validator.Validate(b)

	// Assert that there is no error
	if err != nil {
		t.Errorf("Error validating block: %v", err)
	}
}

// Run all the tests
func TestMain(m *testing.M) {
	// Run the tests
	t.Run("CreateBlockchain", TestCreateBlockchain)
	t.Run("AddBlock", TestAddBlock)
	t.Run("MineBlock", TestMineBlock)
	t.Run("ValidateBlock", TestValidateBlock)

	// Run other tests
	m.Run()
}

func (bc *Blockchain) PrintJSON() {
	bcJSON, err := json.MarshalIndent(bc.Blocks, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling blockchain to JSON:", err)
		return
	}

	fmt.Println("Blockchain (JSON):")
	fmt.Println(string(bcJSON))
}

func main() {
	// Create a new blockchain
	blockchain := NewBlockchain()

	// Create transactions
	transactions := []*Transaction{
		{Sender: "Alice", Receiver: "Bob", Amount: 5.0},
		{Sender: "Bob", Receiver: "Charlie", Amount: 2.5},
	}

	// Create smart contracts
	contracts := []*SmartContract{
		{Code: "Smart Contract 1"},
		{Code: "Smart Contract 2"},
	}

	// Create a new block
	block := NewBlock(transactions, contracts, blockchain.GetLatestBlock().Hash)
	block.MineBlock(blockchain.Difficulty)

	// Add the new block to the blockchain
	blockchain.AddBlock(block)

	// Print the blockchain
	fmt.Println("Blockchain:")
	blockchain.Print()

	// Validate the new block
	validator := NewValidator(blockchain)
	err := validator.Validate(block)
	if err != nil {
		fmt.Println("Block validation failed:", err)
	} else {
		fmt.Println("Block validation successful.")
	}

	// Wait for 3 seconds to demonstrate the timestamp
	time.Sleep(3 * time.Second)

	// Create more transactions and add them to a new block
	transactions = []*Transaction{
		{Sender: "Charlie", Receiver: "David", Amount: 1.5},
		{Sender: "David", Receiver: "Eve", Amount: 3.0},
	}
	block = NewBlock(transactions, contracts, blockchain.GetLatestBlock().Hash)
	block.MineBlock(blockchain.Difficulty)
	blockchain.AddBlock(block)

	// Print the updated blockchain
	fmt.Println("\nUpdated Blockchain:")
	blockchain.Print()

	// Print the blockchain in JSON format
	blockchain.PrintJSON()
}
