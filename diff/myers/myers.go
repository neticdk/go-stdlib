package myers

import (
	"fmt"
	"strings"

	"github.com/neticdk/go-stdlib/diff"
	"github.com/neticdk/go-stdlib/diff/internal/diffcore" // Import the new internal package
)

// Diff computes differences between two values using the Myers diff algorithm.
// Currently never returns an error.
func Diff(a, b string, opts ...Option) (string, error) {
	if a == "" && b == "" {
		return "", nil
	}

	aLines := diffcore.SplitLines(a)
	bLines := diffcore.SplitLines(b)
	return DiffStrings(aLines, bLines, opts...)
}

// DiffStrings computes differences between string slices using the Myers diff algorithm.
// Currently never returns an error.
func DiffStrings(a, b []string, opts ...Option) (string, error) {
	return myersDiffStrings(a, b, applyOptions(opts...))
}

func myersDiffStrings(a, b []string, opts options) (string, error) {
	// Compute edit script using appropriate algorithm
	var script []diff.Line
	if opts.linearSpace {
		script = computeEditScriptLinearSpace(a, b, opts)
	} else {
		script = computeEditScript(a, b, opts.maxEditDistance)
	}

	// Format the diff output
	var sb strings.Builder

	// Track line numbers if enabled
	aLineNum := 1
	bLineNum := 1

	// Group edits by type for context-aware output
	chunks := groupEditsByContext(script, opts.contextLines)

	for i, chunk := range chunks {
		// Add separator between chunks
		if i > 0 && opts.contextLines > 0 {
			sb.WriteString("...\n")
		}

		for _, edit := range chunk {
			switch edit.Kind {
			case diff.Equal:
				if opts.showLineNumbers {
					sb.WriteString(fmt.Sprintf("%4d %4d   ", aLineNum, bLineNum))
				} else {
					sb.WriteString("  ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
				aLineNum++
				bLineNum++
			case diff.Delete:
				if opts.showLineNumbers {
					sb.WriteString(fmt.Sprintf("%4d      - ", aLineNum))
				} else {
					sb.WriteString("- ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
				aLineNum++
			case diff.Insert:
				if opts.showLineNumbers {
					sb.WriteString(fmt.Sprintf("     %4d + ", bLineNum))
				} else {
					sb.WriteString("+ ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
				bLineNum++
			}
		}
	}

	return sb.String(), nil
}

// computeEditScript implements Myers' diff algorithm to find the shortest edit script
func computeEditScript(a, b []string, maxEditDistance int) []diff.Line {
	n, m := len(a), len(b)
	maxDist := n + m

	// Handle special cases
	if maxDist == 0 {
		return []diff.Line{}
	}

	if n == 0 {
		result := make([]diff.Line, m)
		for i := range m {
			result[i] = diff.Line{Kind: diff.Insert, Text: b[i]}
		}
		return result
	}

	if m == 0 {
		result := make([]diff.Line, n)
		for i := range n {
			result[i] = diff.Line{Kind: diff.Delete, Text: a[i]}
		}
		return result
	}

	// Limit max edit distance if specified
	if maxEditDistance > 0 && maxEditDistance < maxDist {
		maxDist = maxEditDistance
	}

	// Initialize the edit graph
	v := make([]int, 2*maxDist+1)
	vs := make([][]int, 0, maxDist)

	// Find the shortest edit script path
	var x, y int
	var foundPath bool

	for d := 0; d <= maxDist; d++ {
		// Clone the current v for backtracking
		vClone := make([]int, len(v))
		copy(vClone, v)
		vs = append(vs, vClone)

		for k := -d; k <= d; k += 2 {
			// Determine whether to move down or right
			if k == -d || (k != d && v[maxDist+k-1] < v[maxDist+k+1]) {
				x = v[maxDist+k+1] // Move down
			} else {
				x = v[maxDist+k-1] + 1 // Move right
			}

			y = x - k

			// Follow diagonal (matching) paths as far as possible
			for x < n && y < m && a[x] == b[y] {
				x++
				y++
			}

			v[maxDist+k] = x

			// Check if we've reached the bottom-right corner
			if x >= n && y >= m {
				foundPath = true
				break
			}
		}

		if foundPath {
			break
		}
	}

	// If no path was found within the constraints (maxEditDistance likely reached),
	// fall back to the simpler LCS-based diff approach using the shared implementation.
	if !foundPath {
		// This simple diff doesn't respect maxEditDistance but provides a complete diff.
		return diffcore.ComputeEditsLCS(a, b)
	}

	// Backtrack to construct the edit script
	return backtrack(a, b, vs, maxDist)
}

// computeEditScriptLinearSpace implements Myers' algorithm with O(N) space complexity.
// It selects the diff algorithm based on input size and constraints:
//
//  1. For small inputs (less than opts.smallInputThreshold) or when
//     opts.maxEditDistance is constrained, it uses the standard Myers algorithm
//     (computeEditScript) which may be faster.
//
//  2. For very large inputs (exceeding opts.largeInputThreshold), it falls back
//     to the simpler LCS-based diff algorithm (diffcore.ComputeEditsLCS) to avoid
//     excessive memory usage.
//
//  3. Otherwise, it uses the linear space Myers algorithm (linearSpaceMyersRecWithIndices).
func computeEditScriptLinearSpace(a, b []string, opts options) []diff.Line {
	n := len(a)
	m := len(b)

	// Use standard algorithm (computeEditScript) for small inputs or when a tight maxEditDistance is set.
	// The standard algorithm might be faster for small N or small D.
	if n < opts.smallInputThreshold || m < opts.smallInputThreshold ||
		(opts.maxEditDistance > 0 && opts.maxEditDistance < n+m) { // Use > 0 to allow disabling constraint with -1
		return computeEditScript(a, b, opts.maxEditDistance)
	}

	// If inputs are very large, fall back to the simpler LCS-based diff (diffcore.ComputeEditsLCS)
	// to avoid potential performance issues or excessive memory with Myers variants.
	if n > opts.largeInputThreshold || m > opts.largeInputThreshold {
		return diffcore.ComputeEditsLCS(a, b)
	}

	// Otherwise, proceed with the linear space recursive approach using indices.
	return linearSpaceMyersRecWithIndices(a, b, 0, n, 0, m, 0, opts.linearRecursionMaxDepth)
}

// linearSpaceMyersRecWithIndices is the recursive helper for the linear space algorithm,
// operating on index ranges (aStart..aEnd, bStart..bEnd) within the original slices 'a' and 'b'
// to avoid repeated sub-slice allocations. It tracks recursion depth to prevent stack overflow.
//
// It uses a divide-and-conquer approach based on Hirschberg's algorithm to find a "middle snake"
// and recursively compute the diff in the sub-ranges before and after the snake.
//
// Fallback conditions:
//  1. Recursion depth limit reached (depth >= maxDepth).
//  2. Subsequences become very small (n < 10 && m < 10), where the overhead of Myers/Hirschberg
//     might outweigh the benefits over the simple LCS algorithm.  In these cases,
//     diffcore.ComputeEditsLCS(a[aStart:aEnd], b[bStart:bEnd]) is used.
func linearSpaceMyersRecWithIndices(a, b []string, aStart, aEnd, bStart, bEnd, depth, maxDepth int) []diff.Line { //revive:disable-line:argument-limit
	n := aEnd - aStart
	m := bEnd - bStart

	if n < 0 || m < 0 {
		return []diff.Line{}
	}

	// Base cases: If one subsequence (defined by the index range) is empty,
	// return pure insertions or deletions for the other subsequence's range.
	if n == 0 {
		result := make([]diff.Line, m)
		for i := range m {
			result[i] = diff.Line{Kind: diff.Insert, Text: b[bStart+i]}
		}
		return result
	}

	if m == 0 {
		result := make([]diff.Line, n)
		for i := range n {
			result[i] = diff.Line{Kind: diff.Delete, Text: a[aStart+i]}
		}
		return result
	}

	// Fallback conditions:
	// 1. Recursion depth limit reached.
	// 2. Subsequences become very small (e.g., < 10 elements), where the overhead
	//    of Myers/Hirschberg might outweigh the benefits over simple LCS.
	if depth >= maxDepth || (n < 10 && m < 10) {
		// Use the shared simple diff logic operating on the relevant subsequences.
		// Note: Slicing happens here, but only for the fallback on small/deep parts.
		return diffcore.ComputeEditsLCS(a[aStart:aEnd], b[bStart:bEnd])
	}

	// Check for common prefix/suffix to reduce problem size
	prefix := 0
	for prefix < n && a[aStart+prefix] == b[bStart+prefix] {
		prefix++
	}

	// Early return to avoid recursion issues
	if aStart+prefix > aEnd || bStart+prefix > bEnd {
		return []diff.Line{}
	}

	suffix := 0
	for suffix < n-prefix && a[aEnd-1-suffix] == b[bEnd-1-suffix] {
		suffix++
	}

	// If we found prefix/suffix, handle them separately
	if prefix > 0 || suffix > 0 {
		if aStart+prefix > aEnd || bStart+prefix > bEnd {
			return []diff.Line{}
		}

		prefixLines := make([]diff.Line, prefix)
		for i := range prefix {
			prefixLines[i] = diff.Line{Kind: diff.Equal, Text: a[aStart+i]}
		}

		// Recursively handle the middle part (with reduced size)
		middleLines := linearSpaceMyersRecWithIndices(
			a, b,
			aStart+prefix, aEnd-suffix,
			bStart+prefix, bEnd-suffix,
			depth+1, maxDepth,
		)

		suffixLines := make([]diff.Line, suffix)
		for i := range suffix {
			suffixLines[i] = diff.Line{Kind: diff.Equal, Text: a[aEnd-suffix+i]}
		}

		result := make([]diff.Line, 0, len(prefixLines)+len(middleLines)+len(suffixLines))
		result = append(result, prefixLines...)
		result = append(result, middleLines...)
		result = append(result, suffixLines...)
		return result
	}

	// Find the middle snake (optimized for smaller comparisons)
	snake := findMiddleSnake(a, b, aStart, aEnd, bStart, bEnd)

	// Recursively solve the two subproblems
	prefixScript := linearSpaceMyersRecWithIndices(
		a, b, aStart, snake.startX, bStart, snake.startY, depth+1, maxDepth,
	)

	// Create the middle snake as equality edits
	middleScript := make([]diff.Line, snake.length)
	for i := range snake.length {
		middleScript[i] = diff.Line{Kind: diff.Equal, Text: a[snake.startX+i]}
	}

	suffixScript := linearSpaceMyersRecWithIndices(
		a, b, snake.endX, aEnd, snake.endY, bEnd, depth+1, maxDepth,
	)

	// Combine all parts
	result := make([]diff.Line, 0, len(prefixScript)+len(middleScript)+len(suffixScript))
	result = append(result, prefixScript...)
	result = append(result, middleScript...)
	result = append(result, suffixScript...)

	return result
}

// Snake represents a diagonal run of matches in the edit graph
type snake struct {
	startX, startY int // Start position
	endX, endY     int // End position
	length         int // Length of the snake
}

// findMiddleSnake implements the core of the linear-space Myers algorithm (Hirschberg's algorithm principle).
// It finds a "middle snake" - a sequence of diagonal moves (matches) that lies on *some*
// shortest edit path between a[aStart:aEnd] and b[bStart:bEnd].
// It works by running the standard Myers algorithm forward from the start and backward
// from the end simultaneously, looking for where the paths overlap.
func findMiddleSnake(a, b []string, aStart, aEnd, bStart, bEnd int) snake { //revive:disable-line:argument-limit
	// Subsequence lengths (within the given ranges)
	n := aEnd - aStart
	m := bEnd - bStart

	// Base case: If either subsequence is empty, return an empty snake at the start.
	if n <= 0 || m <= 0 {
		return snake{aStart, bStart, aStart, bStart, 0}
	}

	if aStart < 0 || aEnd > len(a) || bStart < 0 || bEnd > len(b) {
		return snake{aStart, bStart, aStart, bStart, 0}
	}

	// Optimization: Handle trivial cases where one sequence is very small.
	// This might avoid unnecessary allocation/computation for the vectors below.
	if n == 1 || m == 1 {
		for i := range n {
			for j := range m {
				if a[aStart+i] == b[bStart+j] {
					// Found a single matching element, return it as the snake
					return snake{
						startX: aStart + i,
						startY: bStart + j,
						endX:   aStart + i + 1,
						endY:   bStart + j + 1,
						length: 1,
					}
				}
			}
		}
		// No match found in the trivial case.
		return snake{aStart, bStart, aStart, bStart, 0}
	}

	// maxDiff: The maximum number of edit steps (D) we need to explore.
	// We only need to go up to ceil((N+M)/2) because if a shortest path
	// exists, the forward and reverse searches are guaranteed to meet
	// by this point.
	maxDiff := (n + m + 1) / 2
	// delta: Difference in lengths, used for indexing the reverse path.
	delta := n - m

	// Vectors to store the furthest-reaching x-coordinate for each diagonal k.
	// vf: Forward search (from top-left)
	// vr: Reverse search (from bottom-right)
	// Size needs to accommodate k ranging from -maxDiff to +maxDiff.
	vectorSize := 2*maxDiff + 1
	if delta > 0 {
		vectorSize += delta // Ensure space for reverse search diagonals
	}
	// offset: Used to map diagonal k (which can be negative) to non-negative array indices.
	offset := maxDiff

	vf := make([]int, vectorSize)
	vr := make([]int, vectorSize) // Stores relative x-coordinates within the n x m grid

	// Initialize vectors with sentinel values indicating diagonal hasn't been reached.
	// Using -1 works because valid x coordinates are >= 0.
	for i := range vf {
		vf[i] = -1
		vr[i] = -1
	}

	// Initialize starting points:
	// Forward search starts at (0,0) relative to the subproblem, which is diagonal k=0.
	// vf[offset + 0] = 0 (x=0 for k=0)
	vf[offset] = 0

	// Ensure the initial reverse position is within bounds
	reverseIndex := offset + delta
	if reverseIndex >= 0 && reverseIndex < vectorSize {
		vr[reverseIndex] = n
	} else {
		// If we can't properly initialize the reverse search,
		// fall back to simpler diff algorithm
		return snake{aStart, bStart, aStart, bStart, 0}
	}

	// Add bounds checking for vector updates
	safeVectorUpdate := func(v []int, index, value int) bool {
		if index < 0 || index >= len(v) {
			return false
		}
		v[index] = value
		return true
	}

	// Helper to check if we can follow a diagonal (match) safely within bounds.
	// This now uses the aStart, bStart, aEnd, bEnd indices
	canFollow := func(x, y int) bool {
		return x >= 0 && y >= 0 && x < n && y < m &&
			aStart+x < aEnd && bStart+y < bEnd &&
			a[aStart+x] == b[bStart+y]
	}

	// Iterate D (number of edits) from 0 up to maxDiff
	for d := 0; d <= maxDiff; d++ {
		// --- Forward Search ---
		// Explore diagonals k from -d to d, stepping by 2 (essential property of Myers)
		for k := -d; k <= d; k += 2 { //nolint:gocritic  // ifElseChain is fine here
			idx := offset + k // Array index for diagonal k
			if idx < 0 || idx >= vectorSize {
				continue
			}

			// Determine the starting x for this step on diagonal k.
			// We prioritize moving down (from k+1) if it reaches further than moving right (from k-1).
			var x int
			canMoveRight := k > -d && (offset+k-1) >= 0 && vf[offset+k-1] != -1
			canMoveDown := k < d && (offset+k+1) < vectorSize && vf[offset+k+1] != -1

			if !canMoveRight && canMoveDown { //nolint:gocritic  // ifElseChain is fine here
				x = vf[offset+k+1] // Must move down from k+1
			} else if canMoveRight && !canMoveDown {
				x = vf[offset+k-1] + 1 // Must move right from k-1
			} else if canMoveRight && canMoveDown {
				// Choose the path that reached further previously
				x = max(vf[offset+k-1]+1, vf[offset+k+1]) // Choose the better path
			} else {
				// Base case for d=0, k=0 (handled by vf[offset]=0 initialization)
				// Or, if neither k-1 nor k+1 was reachable (shouldn't happen for d>0)
				if idx != offset || d != 0 {
					continue
				}
				x = vf[idx]
			}

			y := x - k // Calculate y based on x and diagonal k

			// Store the starting point of the potential snake for this (d, k)
			startX := x
			startY := y

			// Follow the diagonal (matches) as far as possible
			for canFollow(x, y) {
				x++
				y++
			}

			if !safeVectorUpdate(vf, idx, x) {
				continue
			}

			// Update the furthest reaching x for this diagonal k
			vf[idx] = x

			//  Overlap Check
			// Overlap check
			if delta%2 != 0 && k >= delta-(d-1) && k <= delta+(d-1) {
				reverseIdx := offset + (delta - k)
				if reverseIdx >= 0 && reverseIdx < vectorSize {
					reverseX := vr[reverseIdx]
					if reverseX >= 0 && x >= reverseX {
						return snake{
							startX: aStart + startX,
							startY: bStart + startY,
							endX:   aStart + x,
							endY:   bStart + y,
							length: x - startX,
						}
					}
				}
			}
		} // End forward k loop

		// --- Reverse Search ---
		for k_rev := -d; k_rev <= d; k_rev += 2 {
			k := k_rev + delta
			idx := offset + k
			if idx < 0 || idx >= vectorSize {
				continue
			}

			var x int
			canMoveLeft := k > -d+delta && (offset+k-1) >= 0 && vr[offset+k-1] != -1
			canMoveUp := k < d+delta && (offset+k+1) < vectorSize && vr[offset+k+1] != -1

			switch {
			case !canMoveLeft && canMoveUp:
				x = vr[offset+k+1] - 1
			case canMoveLeft && !canMoveUp:
				x = vr[offset+k-1]
			case canMoveLeft && canMoveUp:
				x = min(vr[offset+k-1], vr[offset+k+1]-1)
			default:
				// Base case, should not normally happen
				continue
			}

			y := x - k

			endX := x
			endY := y

			for x > 0 && y > 0 && canFollow(x-1, y-1) {
				x--
				y--
			}

			if !safeVectorUpdate(vr, idx, x) {
				continue
			}

			vr[idx] = x

			// --- Overlap Check ---
			if delta%2 == 0 && k >= -d && k <= d {
				forwardIdx := offset + k
				if forwardIdx >= 0 && forwardIdx < vectorSize {
					forwardX := vf[forwardIdx]
					if forwardX >= 0 && forwardX >= x {
						return snake{
							startX: aStart + x,
							startY: bStart + y,
							endX:   aStart + endX,
							endY:   bStart + endY,
							length: endX - x,
						}
					}
				}
			}
		} // End reverse k loop
	} // End d loop

	// Should not be reached
	return snake{aStart, bStart, aStart, bStart, 0}
}

// backtrack constructs the edit script by walking backward through the edit graph
// stored in vs. It starts from the end point (len(a), len(b)) and works
// back towards (0, 0).
func backtrack(a, b []string, vs [][]int, midpoint int) []diff.Line {
	// Estimate capacity - could be up to len(a) + len(b) edits
	// This is a heuristic, not a tight bound.
	initialCapacity := len(a) + len(b)
	edits := make([]diff.Line, 0, initialCapacity)

	x, y := len(a), len(b) // Start at the end of the edit graph

	// Iterate backwards through the D-paths recorded in vs
	for d := len(vs) - 1; d >= 0; d-- {
		v := vs[d] // v is the state of furthest-reaching x-values for D-path d
		k := x - y // Current diagonal k = x - y

		// Determine the k-diagonal of the *previous* D-path (d-1) that led to (x, y).
		// If k == -d, the only way to reach it is from k = -d + 1 (move down/insert).
		// If k == d, the only way to reach it is from k = d - 1 (move right/delete).
		// Otherwise, compare the x-values on the two possible previous diagonals (k-1 and k+1)
		// in the *previous* D-path's v state (vs[d-1], accessed via the current v).
		// We choose the diagonal that could reach furthest along the x-axis.
		var prevK int
		if k == -d || (k != d && v[midpoint+k-1] < v[midpoint+k+1]) {
			// Came from k+1: Move was down (insert B's element)
			prevK = k + 1
		} else {
			// Came from k-1: Move was right (delete A's element)
			prevK = k - 1
		}

		// Find the coordinate (prevX, prevY) on the previous D-path
		prevX := v[midpoint+prevK]
		prevY := prevX - prevK

		// Backtrack along the diagonal (equal lines) until we reach (prevX, prevY)
		// These diagonal moves represent matching elements.
		for x > prevX && y > prevY {
			// Append Equal edit (indices are 1-based, so use x-1, y-1)
			// Appending to the end here, will reverse later.
			edits = append(edits, diff.Line{Kind: diff.Equal, Text: a[x-1]})
			x--
			y--
		}

		// If we haven't reached the start (d > 0), add the non-diagonal move
		// (Insert or Delete) that led from (prevX, prevY) to the start of the
		// diagonal we just traversed.
		if d > 0 {
			if prevX < x { // Moved right (x increased): Delete A[x-1]
				edits = append(edits, diff.Line{Kind: diff.Delete, Text: a[x-1]})
				x-- // Decrement x corresponding to the deletion
			} else { // Moved down (y increased relative to x): Insert B[y-1]
				edits = append(edits, diff.Line{Kind: diff.Insert, Text: b[y-1]})
				y-- // Decrement y corresponding to the insertion
			}
		}
	}

	// Reverse the edits slice because we appended them in reverse order
	for i, j := 0, len(edits)-1; i < j; i, j = i+1, j-1 {
		edits[i], edits[j] = edits[j], edits[i]
	}

	return edits
}

// groupEditsByContext groups edits into chunks for context-aware display
func groupEditsByContext(edits []diff.Line, contextLines int) [][]diff.Line {
	// If no context requested, only include non-equal edits
	if contextLines <= 0 {
		var changedEdits []diff.Line
		for _, edit := range edits {
			if edit.Kind != diff.Equal {
				changedEdits = append(changedEdits, edit)
			}
		}
		return [][]diff.Line{changedEdits}
	}

	// If the total edit is small enough, return as one chunk
	if len(edits) <= contextLines*2 {
		return [][]diff.Line{edits}
	}

	// Find indices of non-equal edits (changes)
	var changeIndices []int
	for i, edit := range edits {
		if edit.Kind != diff.Equal {
			changeIndices = append(changeIndices, i)
		}
	}

	// If no changes, return all as one chunk
	if len(changeIndices) == 0 {
		return [][]diff.Line{edits}
	}

	// Pre-allocate chunks - in worst case, each change could be in its own chunk
	// A more conservative estimate would be len(changeIndices)/2 + 1
	chunks := make([][]diff.Line, 0, min(len(changeIndices), len(edits)/2+1))

	// Group changes that are close to each other
	var currentGroup []int
	var groups [][]int

	for i, idx := range changeIndices {
		if i == 0 || idx-changeIndices[i-1] > contextLines*2 {
			// Start a new group
			if len(currentGroup) > 0 {
				groups = append(groups, currentGroup)
			}
			currentGroup = []int{idx}
		} else {
			// Add to current group
			currentGroup = append(currentGroup, idx)
		}
	}

	// Add the last group
	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	// Create chunks based on groups with context
	for _, group := range groups {
		firstChange := group[0]
		lastChange := group[len(group)-1]

		// Calculate chunk boundaries with context
		startIdx := max(0, firstChange-contextLines)
		endIdx := min(len(edits), lastChange+contextLines+1)

		chunk := make([]diff.Line, endIdx-startIdx)
		copy(chunk, edits[startIdx:endIdx])
		chunks = append(chunks, chunk)
	}

	return chunks
}
