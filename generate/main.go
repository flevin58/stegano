package main

import (
	"os"
)

const MAX_ALGOS int = 20

func main() {
	off := [MAX_ALGOS][8]byte{
		{0, 1, 2, 4, 5, 6, 8, 9},
		{9, 8, 6, 5, 4, 2, 1, 0},
		{1, 0, 2, 4, 5, 6, 8, 9},
		{4, 5, 6, 8, 9, 0, 1, 2},
		{8, 9, 0, 1, 2, 4, 5, 6},
		{2, 4, 6, 8, 0, 1, 5, 9},
		{1, 5, 9, 0, 2, 4, 6, 8},
		{6, 8, 4, 2, 0, 9, 5, 1},
		{5, 4, 1, 2, 8, 9, 0, 6},
		{0, 2, 4, 6, 8, 1, 5, 9},
		{9, 1, 8, 2, 6, 4, 5, 0},
		{4, 0, 5, 1, 6, 2, 8, 9},
		{8, 6, 4, 2, 1, 0, 5, 9},
		{1, 2, 4, 5, 0, 6, 8, 9},
		{5, 6, 8, 9, 4, 2, 1, 0},
		{0, 9, 1, 8, 2, 6, 4, 5},
		{2, 1, 0, 4, 5, 6, 8, 9},
		{6, 5, 4, 2, 0, 1, 8, 9},
		{9, 8, 5, 6, 4, 1, 2, 0},
		{4, 2, 1, 0, 8, 9, 6, 5},
	}

	f, err := os.Create("stegano/algo.bin")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buffer := make([]byte, 0, MAX_ALGOS*8)
	for algo := 0; algo < MAX_ALGOS; algo++ {
		for bit := 0; bit < 8; bit++ {
			buffer = append(buffer, off[algo][bit])
		}
	}
	f.Write(buffer)
}
