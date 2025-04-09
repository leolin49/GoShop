package shuffle

import (
	"fmt"
	"testing"
)

func TestShuffleKnuthDurstenfeld(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	ShuffleKnuthDurstenfeld(a)
	fmt.Println(a)
}

func TestShuffleInsideOut(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	b := ShuffleInsideOut(a)
	fmt.Println(b)
}

func Test(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	b, _ := SamplingReservoir(a, len(a)/2)
	fmt.Println(b)
}
