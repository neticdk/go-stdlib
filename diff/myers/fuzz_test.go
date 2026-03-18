package myers

import (
	"strings"
	"testing"

	"github.com/neticdk/go-stdlib/diff"
	"github.com/neticdk/go-stdlib/diff/internal/diffcore"
)

func FuzzMyersLinearSpace(f *testing.F) {
	f.Add("hello", "world")
	f.Add("abc", "def")
	f.Add("A B C D E", "X Y A B C D E Z")
	f.Add("X X X X X", "X X X X")
	f.Add("", "")
	f.Add("A", "")
	f.Add("", "B")
	f.Add("A B C D E F G H I J K", "X Y A B C D E F G H I J K Z")

	f.Fuzz(func(t *testing.T, aStr, bStr string) {
		a := []string{}
		if aStr != "" {
			a = strings.Split(aStr, "")
		}
		b := []string{}
		if bStr != "" {
			b = strings.Split(bStr, "")
		}

		opts := options{
			linearSpace:             true,
			linearRecursionMaxDepth: 100, // standard
		}

		// Test that it doesn't panic
		script := computeEditScriptLinearSpace(a, b, opts)

		// Also compute with LCS for comparison of edit distance length
		lcsScript := diffcore.ComputeEditsLCS(a, b)

		// Verify both scripts correctly transform a to b
		verifyScript(t, a, b, script)
		verifyScript(t, a, b, lcsScript)

		// Count edits to ensure both found optimal paths (same number of edits)
		editsMyers := 0
		for _, l := range script {
			if l.Kind != diff.Equal {
				editsMyers++
			}
		}

		editsLCS := 0
		for _, l := range lcsScript {
			if l.Kind != diff.Equal {
				editsLCS++
			}
		}

		if editsMyers != editsLCS {
			t.Errorf("Myers gave %d edits, LCS gave %d edits. A=%q, B=%q", editsMyers, editsLCS, aStr, bStr)
		}
	})
}

func verifyScript(t *testing.T, a, b []string, script []diff.Line) {
	aIdx := 0
	bIdx := 0

	for _, l := range script {
		switch l.Kind {
		case diff.Equal:
			if aIdx >= len(a) || bIdx >= len(b) {
				t.Fatalf("Equal operation out of bounds: aIdx=%d, bIdx=%d, aLen=%d, bLen=%d", aIdx, bIdx, len(a), len(b))
			}
			if a[aIdx] != l.Text || b[bIdx] != l.Text {
				t.Fatalf("Equal text mismatch: text=%q, a[%d]=%q, b[%d]=%q", l.Text, aIdx, a[aIdx], bIdx, b[bIdx])
			}
			aIdx++
			bIdx++
		case diff.Insert:
			if bIdx >= len(b) {
				t.Fatalf("Insert operation out of bounds: bIdx=%d, bLen=%d", bIdx, len(b))
			}
			if b[bIdx] != l.Text {
				t.Fatalf("Insert text mismatch: text=%q, b[%d]=%q", l.Text, bIdx, b[bIdx])
			}
			bIdx++
		case diff.Delete:
			if aIdx >= len(a) {
				t.Fatalf("Delete operation out of bounds: aIdx=%d, aLen=%d", aIdx, len(a))
			}
			if a[aIdx] != l.Text {
				t.Fatalf("Delete text mismatch: text=%q, a[%d]=%q", l.Text, aIdx, a[aIdx])
			}
			aIdx++
		}
	}

	if aIdx != len(a) {
		t.Fatalf("Did not consume all of a. aIdx=%d, aLen=%d", aIdx, len(a))
	}
	if bIdx != len(b) {
		t.Fatalf("Did not consume all of b. bIdx=%d, bLen=%d", bIdx, len(b))
	}
}
