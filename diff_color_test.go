package akashi_test

import (
	"strings"
	"testing"

	"github.com/seiyab/akashi"
)

func TestDiff_WithColor_DefaultProfile(t *testing.T) {
	red := "\x1b[31m"
	green := "\x1b[32m"
	reset := "\x1b[0m"

	runTest(t, "hello", "world", strings.Join([]string{
		red + `- "hello"` + reset,
		green + `+ "world"` + reset,
	}, "\n"), akashi.WithColor())
}

func TestDiff_WithColor_CustomProfile(t *testing.T) {
	profile := akashi.ColorProfile{
		Left:  "[L]",
		Right: "[R]",
		Reset: "[X]",
	}

	runTest(t, 1, 2, strings.Join([]string{
		`[L]- 1[X]`,
		`[R]+ 2[X]`,
	}, "\n"), akashi.WithColor(profile))
}
