package shuffle

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
)

func ShuffleKnuthDurstenfeld[T any](a []T) {
	for i := len(a) - 1; i >= 0; i-- {
		j := rand.Int() % (i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func ShuffleInsideOut[T any](a []T) []T {
	b := slices.Clone(a)
	for i := 0; i < len(a); i++ {
		j := rand.Int() % (i + 1)
		b[i] = b[j]
		b[j] = a[i]
	}
	return b
}

// SamplingReservoir randomly choose m elements from a.
func SamplingReservoir[T any](a []T, m int) ([]T, error) {
	n := len(a)
	if m > n {
		return nil, errors.New(
			fmt.Sprintf("can not choose %d elements from a %d slices", m, n),
		)
	}
	for i := m; i < n; i++ {
		j := rand.Int() % (i + 1)
		if j < m {
			a[i], a[j] = a[j], a[i]
		}
	}
	return a[:m], nil
}
