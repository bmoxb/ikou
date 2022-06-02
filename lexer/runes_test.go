package lexer

import "testing"

func TestRuneIsAlpha(t *testing.T) {
	//var a rune = 'a'
	//var z rune = 'z'
	//t.Fatalf("%d %d", a, z)
	table := map[rune]bool{
		'a': true,
		'g': true,
		'z': true,
		'!': false,
		'5': false,
		0:   false,
		'A': true,
		'J': true,
		'Z': true,
	}

	for r, b := range table {
		isAlpha := runeIsAlpha(r)

		if isAlpha != b {
			t.Errorf("runeIsAlpha('%c') is %t but expected it to be %t", r, isAlpha, b)
		}
	}
}
