package main

import (
	"encoding/hex"
	"fmt"
)

type Bitmap [8]byte

func (bits *Bitmap) IsSet(i int) bool {
	i -= 1
	return bits[i/8]&(1<<uint(7-i%8)) != 0
}
func (bits *Bitmap) Set(i int) { i -= 1; fmt.Println(uint(7 - i%8)); bits[i/8] |= 1 << uint(7-i%8) }
func (bits *Bitmap) Clear(i int) {
	i -= 1
	fmt.Println(uint(7 - i%8))
	bits[i/8] &^= 1 << uint(7-i%8)
}

func (bits *Bitmap) Sets(xs ...int) {
	for _, x := range xs {
		bits.Set(x)
	}
}

func (bits Bitmap) String() string { return hex.EncodeToString(bits[:]) }

// func main() {
// 	var ch Chess = Chess{}

// 	fmt.Println(ch.GetKingMoves(60))
// 	fmt.Println(ch.Count(ch.GetKingMoves(60)))
// }

type Chess struct {
}

func (c *Chess) Count(val uint64) int {
	var cnt int
	for val != 0 {
		if val&1 == 1 {
			cnt++
		}
		val >>= 1
	}
	return cnt
}

func (c *Chess) GetKingMoves(pos int) uint64 {
	var K uint64 = 1 << pos
	var noA uint64 = 0xfefefefefefefefe
	var noH uint64 = 0x7f7f7f7f7f7f7f7f
	var Ka uint64 = K & noA
	var Kh uint64 = K & noH
	var mask uint64 = (Ka << 7) | (K << 8) | (Kh << 9) |
		(Ka >> 1) | (Kh << 1) | (Ka >> 9) |
		(K >> 8) | (Kh >> 7)

	return mask
}
