package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"time"
	//"./bolt"
)
//0.定义结构
type Block struct {
	//版本号
	version uint64
	//merkel根(就是一个hash值)
	merkleroot []byte
	//时间戳
	TimeStamp uint64
	//难度值
	Diffcult uint64
	//随机数，挖矿要找的数据
	Nonce uint64
	prevHash []byte
	Hash []byte
	Data []byte
}
//1.创建一个区块
func NewBlock(data string,prevBlockHash []byte) *Block{
	block := Block{
		version: 00,
		merkleroot: []byte{},
		TimeStamp: uint64(time.Now().Unix()),
		Diffcult: 0,	//无效值
		Nonce: 0,	//同无效
		prevHash: prevBlockHash,
		Hash: []byte{},
		Data:[]byte(data),
	}
	//block.setHash()
	pow:=newProofOfWork(&block)
	hash,nonce :=pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}
//2.生成hash

func (block *Block) setHash(){
	var blockInfo []byte
	//1.拼接数据
	//blockInfo=append(blockInfo,block.Data...)
	//blockInfo=append(blockInfo,block.merkleroot...)
	//blockInfo=append(blockInfo,UintToByte(block.version)...)
	//blockInfo=append(blockInfo,UintToByte(block.TimeStamp)...)
	//blockInfo=append(blockInfo,UintToByte(block.Diffcult)...)
	//blockInfo=append(blockInfo,UintToByte(block.Nonce)...)

	//2.sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
//3.引入区块链
type BlockChain struct {
	//定义一个区块链数组
	blockchain []*Block
	//使用数据库改写区块链
	//db *bolt.DB
}
//4.定义一个区块链
func newBlockChain() *BlockChain{
	//创建一个创世块，并作为第一个区块
	genesisBlock:=GenesisBlock()
	return &BlockChain{
		blockchain:[]*Block{genesisBlock},
	}
}
//创世块
func GenesisBlock() *Block{
	return NewBlock("创世块",[]byte{})
}
//5.添加区块
func (blockchain *BlockChain) addBlock(data string){
	//获取前区块hash
	lastBlock := blockchain.blockchain[len(blockchain.blockchain)-1]
	prevHash := lastBlock.Hash
	//创建一个区块
	block := NewBlock(data,prevHash)
	//添加到区块链数组中
	blockchain.blockchain = append(blockchain.blockchain,block)
}
//一个辅助函数 功能是吧uint转换为[]byte
func Uint64ToByte(num uint64) []byte  {
	var buffer bytes.Buffer
	err:=binary.Write(&buffer,binary.BigEndian,num)
	if err != nil{
		log.Panic(err)
	}
	return buffer.Bytes()
}
func main(){
	bc:=newBlockChain()
	bc.addBlock("第一笔交易")
	bc.addBlock("第二笔交易")
	for i,block := range bc.blockchain{
		fmt.Printf("=========当前块高度:%d==========\n",i)
		fmt.Printf("区块pre哈希:%x\n",block.prevHash)
		fmt.Printf("区块哈希:%x\n",block.Hash)
		fmt.Printf("区块数据:%s\n",block.Data)
		fmt.Printf("区块时间:%d\n",block.TimeStamp)
	}
}