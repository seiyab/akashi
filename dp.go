package akashi

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type dpCell struct {
	loss  float64
	entry entry
	fromA int
	fromB int
}

func mixedEntries(
	p diffProcess,
	leftLength, rightLength int,
	getLeft, getRight func(int) reflect.Value,
) ([]entry, error) {
	leading := make([]entry, 0)
	for i := 0; i < leftLength && i < rightLength; i++ {
		t := p.diff(getLeft(i), getRight(i))
		if t.loss() > 0 {
			break
		}
		leading = append(leading, entry{value: p.leftPure(getLeft(i))})
	}
	k := len(leading)

	dp := make([][]dpCell, leftLength-k+1)
	for i := range dp {
		dp[i] = make([]dpCell, rightLength-k+1)
		for j := range dp[i] {
			dp[i][j] = dpCell{loss: math.MaxFloat64}
		}
	}
	dp[0][0] = dpCell{loss: 0}
	for b := 0; k+b < rightLength+1; b++ {
		for a := 0; k+a < leftLength+1; a++ {
			l := dp[a][b].loss
			if k+a < leftLength {
				if l+1 < dp[a+1][b].loss {
					dp[a+1][b] = dpCell{
						loss: l + 1,
						entry: entry{
							leftOnly: true,
							value:    p.leftPure(getLeft(k + a)),
						},
						fromA: a,
						fromB: b,
					}
				}
			}
			if k+b < rightLength {
				if l+1 < dp[a][b+1].loss {
					dp[a][b+1] = dpCell{
						loss: l + 1,
						entry: entry{
							rightOnly: true,
							value:     p.rightPure(getRight(k + b)),
						},
						fromA: a,
						fromB: b,
					}
				}
			}
			if k+a < leftLength && k+b < rightLength {
				t := p.diff(getLeft(k+a), getRight(k+b))
				tl := t.loss()
				switch t.(type) {
				case mixed, cycle, nilNode, format1:
					if l+tl < dp[a+1][b+1].loss {
						dp[a+1][b+1] = dpCell{
							loss:  l + tl,
							entry: entry{value: t},
							fromA: a,
							fromB: b,
						}
					}
				}
			}
		}
	}
	a := len(dp) - 1
	b := len(dp[a]) - 1
	if dp[a][b].loss > 1_000_000 {
		return nil, fmt.Errorf("failed to compute diff")
	}

	trailing := make([]entry, 0, leftLength+rightLength)
	for a, b := leftLength-k, rightLength-k; a > 0 || b > 0; {
		cell := dp[a][b]
		trailing = append(trailing, cell.entry)
		if !(cell.fromA < a || cell.fromB < b) {
			return nil, fmt.Errorf("infinite loop")
		}
		a = cell.fromA
		b = cell.fromB
	}
	reverse(trailing)

	entries := append(leading, trailing...)

	return entries, nil
}

func sliceMixedEntries(v1, v2 reflect.Value, p diffProcess) ([]entry, error) {
	if v1.Kind() != reflect.Slice && v1.Kind() != reflect.Array ||
		v2.Kind() != reflect.Slice && v2.Kind() != reflect.Array {
		return nil, errors.New("unexpected kind")
	}
	es, err := mixedEntries(
		p,
		v1.Len(), v2.Len(),
		func(i int) reflect.Value { return v1.Index(i) },
		func(i int) reflect.Value { return v2.Index(i) },
	)
	if err != nil {
		return nil, err
	}
	return es, nil
}

func multiLineStringEntries(v1, v2 []string, p diffProcess) ([]entry, error) {
	return mixedEntries(
		p,
		len(v1), len(v2),
		func(i int) reflect.Value { return reflect.ValueOf(v1[i]) },
		func(i int) reflect.Value { return reflect.ValueOf(v2[i]) },
	)
}

func reverse(entries []entry) {
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
}
