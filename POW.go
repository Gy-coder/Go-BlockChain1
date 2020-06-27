package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	//block
	block *Block
	//目标值
	target *big.Int
}

func newProofOfWork(block *Block) *ProofOfWork {
	pow:=ProofOfWork{
		block: block,
	}
	//指定的难度值
	targetStr :=  "0000100000000000000000000000000000000000000000000000000000000000"
	//辅助变量，目的是将上面的难度值转为big.int
	tmpInt := big.Int{}
	//将难度值赋值给big.int 指定为16进制
	tmpInt.SetString(targetStr,16)
	pow.target = &tmpInt
	return &pow
}

func (pow *ProofOfWork) Run() ([]byte,uint64){
	//1.拼装数据(区块数据以及不断变化的随机数)
	//2.做哈希运算
	//3.与POW中的target进行比较，找到的话退出返回，否则继续找，随机数+1
	var nonce uint64
	block:=pow.block
	var hash [32]byte
	//拼接数据
	for {
		//1. 拼装数据（区块的数据，还有不断变化的随机数）
		tmp := [][]byte{
			Uint64ToByte(block.version),
			block.prevHash,
			block.merkleroot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Diffcult),
			Uint64ToByte(nonce),
			//只对区块头做哈希值，区块体通过MerkelRoot产生影响
			//block.Data,
		}

		//将二维的切片数组链接起来，返回一个一维的切片
		blockInfo := bytes.Join(tmp, []byte{})

		//哈希运算
		hash = sha256.Sum256(blockInfo)
		//与pow中的target进行比较
		tmpInt := big.Int{}
		//将得到的hash数组转化为bigInt
		tmpInt.SetBytes(hash[:])
		//比较当前hash与目标hash 如果当前<目标 找到 否则继续找
		// Cmp
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if tmpInt.Cmp(pow.target) == -1{
			// 找到了
			fmt.Printf("挖矿成功！hash:%x,nonce:%d\n",hash,nonce)
			break
		}else{
			// 没找到，继续找 随机数+1
			nonce++
		}
	}
	return hash[:],nonce
}