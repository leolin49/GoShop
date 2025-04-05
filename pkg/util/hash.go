package util

import "github.com/spaolacci/murmur3"

func HashMurmur32(s string) uint32 {
	return murmur3.Sum32([]byte(s))
}

func HashMurmur64(s string) uint64 {
	return murmur3.Sum64([]byte(s))
}
