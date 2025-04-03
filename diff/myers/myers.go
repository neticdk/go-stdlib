package myers

import (
	"fmt"
	"strings"

	"github.com/neticdk/go-stdlib/diff"
	"github.com/neticdk/go-stdlib/diff/internal/lcs"
	"github.com/neticdk/go-stdlib/diff/internal/lines"
)

const (
	smallInputThreshold = 100
)

// Diff computes differences between two values using the Myers diff algorithm.
// Currently never returns an error.
func Diff(a, b string, opts ...Option) (string, error) {
	if a == "" && b == "" {
		return "", nil
	}

	aLines := lines.Split(a)
	bLines := lines.Split(b)
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
		script = computeEditScriptLinearSpace(a, b, opts.maxEditDistance, opts.simpleDiffFallbackSize, opts.linearRecursionMaxDepth)
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

	// If no path found within constraints, fall back to a simpler approach
	if !foundPath {
		return simpleDiff(a, b)
	}

	// Backtrack to construct the edit script
	return backtrack(a, b, vs, maxDist)
}

// computeEditScriptLinearSpace implements Myers' algorithm with O(N) space complexity
func computeEditScriptLinearSpace(a, b []string, maxEditDistance int, simpleDiffFallbackSize int, linearRecursionMaxDepth int) []diff.Line {
	// Use standard algorithm for small inputs
	if len(a) < smallInputThreshold || len(b) < smallInputThreshold ||
		(maxEditDistance >= 0 && maxEditDistance < len(a)+len(b)) {
		return computeEditScript(a, b, maxEditDistance)
	}

	if len(a) > simpleDiffFallbackSize || len(b) > simpleDiffFallbackSize {
		return simpleDiff(a, b)
	}
	return linearSpaceMyersRecWithDepth(a, b, 0, len(a), 0, len(b), 0, linearRecursionMaxDepth)
}

// linearSpaceMyersRecWithDepth tracks recursion depth to prevent stack overflow
func linearSpaceMyersRecWithDepth(a, b []string, aStart, aEnd, bStart, bEnd, depth, maxDepth int) []diff.Line { //revive:disable-line:argument-limit
	// Base cases
	if aStart == aEnd {
		result := make([]diff.Line, bEnd-bStart)
		for i := range bEnd - bStart {
			result[i] = diff.Line{Kind: diff.Insert, Text: b[bStart+i]}
		}
		return result
	}

	if bStart == bEnd {
		result := make([]diff.Line, aEnd-aStart)
		for i := range aEnd - aStart {
			result[i] = diff.Line{Kind: diff.Delete, Text: a[aStart+i]}
		}
		return result
	}

	// Fall back to simpler algorithm if recursion gets too deep
	if depth >= maxDepth {
		return simpleDiffSubsequence(a[aStart:aEnd], b[bStart:bEnd])
	}

	// Special case: if sequences are small enough, use direct algorithm
	if aEnd-aStart < 10 && bEnd-bStart < 10 {
		return simpleDiffSubsequence(a[aStart:aEnd], b[bStart:bEnd])
	}

	// Check for common prefix/suffix to reduce problem size
	prefix := 0
	for prefix < aEnd-aStart && prefix < bEnd-bStart && a[aStart+prefix] == b[bStart+prefix] {
		prefix++
	}

	suffix := 0
	for suffix < aEnd-aStart-prefix && suffix < bEnd-bStart-prefix &&
		a[aEnd-1-suffix] == b[bEnd-1-suffix] {
		suffix++
	}

	// If we found prefix/suffix, handle them separately
	if prefix > 0 || suffix > 0 {
		prefixLines := make([]diff.Line, prefix)
		for i := range prefix {
			prefixLines[i] = diff.Line{Kind: diff.Equal, Text: a[aStart+i]}
		}

		// Recursively handle the middle part (with reduced size)
		middleLines := linearSpaceMyersRecWithDepth(
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
	prefixScript := linearSpaceMyersRecWithDepth(
		a, b, aStart, snake.startX, bStart, snake.startY, depth+1, maxDepth,
	)

	// Create the middle snake as equality edits
	middleScript := make([]diff.Line, snake.length)
	for i := range snake.length {
		middleScript[i] = diff.Line{Kind: diff.Equal, Text: a[snake.startX+i]}
	}

	suffixScript := linearSpaceMyersRecWithDepth(
		a, b, snake.endX, aEnd, snake.endY, bEnd, depth+1, maxDepth,
	)

	// Combine all parts
	result := make([]diff.Line, 0, len(prefixScript)+len(middleScript)+len(suffixScript))
	result = append(result, prefixScript...)
	result = append(result, middleScript...)
	result = append(result, suffixScript...)

	return result
}

// simpleDiffSubsequence is a non-recursive fallback for small sequences
func simpleDiffSubsequence(a, b []string) []diff.Line {
	// Similar to the existing simpleDiff function
	edits := []diff.Line{}

	// Use a simple LCS approach for small subsequences
	longest := lcs.LongestCommonSubsequence(a, b)
	aIndex, bIndex := 0, 0

	for _, item := range longest {
		// Add deletions for unmatched items in a
		for aIndex < item.AIndex {
			edits = append(edits, diff.Line{Kind: diff.Delete, Text: a[aIndex]})
			aIndex++
		}

		// Add insertions for unmatched items in b
		for bIndex < item.BIndex {
			edits = append(edits, diff.Line{Kind: diff.Insert, Text: b[bIndex]})
			bIndex++
		}

		// Add the matching item
		edits = append(edits, diff.Line{Kind: diff.Equal, Text: a[aIndex]})
		aIndex++
		bIndex++
	}

	// Handle remaining items
	for aIndex < len(a) {
		edits = append(edits, diff.Line{Kind: diff.Delete, Text: a[aIndex]})
		aIndex++
	}

	for bIndex < len(b) {
		edits = append(edits, diff.Line{Kind: diff.Insert, Text: b[bIndex]})
		bIndex++
	}

	return edits
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
	// Subsequence lengths
	n := aEnd - aStart
	m := bEnd - bStart

	// Base case: If either subsequence is empty, return an empty snake at the start.
	if n <= 0 || m <= 0 {
		return snake{aStart, bStart, aStart, bStart, 0}
	}

	// Optimization: Handle trivial cases where one sequence is very small.
	// This might avoid unnecessary allocation/computation for the vectors below.
	// (Consider adjusting the threshold '1' if needed based on profiling)
	if n == 1 || m == 1 {
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if a[aStart+i] == b[bStart+j] {
					// Found a single matching element, return it as the snake
					return snake{
						startX: aStart + i, startY: bStart + j,
						endX: aStart + i + 1, endY: bStart + j + 1,
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
	// Reverse search starts at (n,m) relative to the subproblem, which is diagonal k = n-m = delta.
	// vr[offset + delta] = n (x=n for k=delta)
	vr[offset+delta] = n

	// Helper to check if we can follow a diagonal (match) safely within bounds.
	// Note: Coordinates x, y are relative to the start of the subproblem (0..n-1, 0..m-1).
	canFollow := func(x, y int) bool {
		// Check bounds relative to the subproblem grid (n x m)
		return x >= 0 && y >= 0 && x < n && y < m &&
			// Check bounds relative to the original full slices a and b
			aStart+x < len(a) && bStart+y < len(b) &&
			// Check for actual character equality
			a[aStart+x] == b[bStart+y]
	}

	// Iterate D (number of edits) from 0 up to maxDiff
	for d := 0; d <= maxDiff; d++ {
		// --- Forward Search ---
		// Explore diagonals k from -d to d, stepping by 2 (essential property of Myers)
		for k := -d; k <= d; k += 2 {
			idx := offset + k // Array index for diagonal k
			if idx < 0 || idx >= vectorSize {
				continue
			} // Bounds check

			// Determine the starting x for this step on diagonal k.
			// We prioritize moving down (from k+1) if it reaches further than moving right (from k-1).
			var x int
			canMoveRight := k > -d && (offset+k-1) >= 0 && vf[offset+k-1] != -1       // Check if k-1 is valid and reached
			canMoveDown := k < d && (offset+k+1) < vectorSize && vf[offset+k+1] != -1 // Check if k+1 is valid and reached

			if !canMoveRight && canMoveDown { //nolint:gocritic  // ifElseChain is fine here
				x = vf[offset+k+1] // Must move down from k+1
			} else if canMoveRight && !canMoveDown {
				x = vf[offset+k-1] + 1 // Must move right from k-1
			} else if canMoveRight && canMoveDown {
				// Choose the path that reached further previously
				if vf[offset+k-1]+1 > vf[offset+k+1] {
					x = vf[offset+k-1] + 1 // Move right from k-1
				} else {
					x = vf[offset+k+1] // Move down from k+1
				}
			} else {
				// Base case for d=0, k=0 (handled by vf[offset]=0 initialization)
				// Or, if neither k-1 nor k+1 was reachable (shouldn't happen for d>0)
				if idx == offset && d == 0 {
					x = vf[idx]
				} else {
					// This path should ideally not be taken if vf is initialized correctly
					// If it happens, it indicates an issue or an unreachable state.
					// For robustness, maybe continue or log? Setting x=0 might be risky.
					// Let's rely on the initial vf[offset]=0 and the logic above.
					// If vf[idx] is still -1, we can't proceed from here this step.
					if vf[idx] == -1 && k != -d && k != d { //nolint:gocritic  // Special case for edges? ifElseChain is fine here
						// Revisit this logic if issues arise. The standard Myers assumes
						// progress is always possible from a valid previous state.
						// Let's assume the logic above correctly finds a valid previous x.
						// If k == -d, MUST come from k+1. If k == d, MUST come from k-1.
						if k == -d {
							x = vf[offset+k+1] // Move down
						} else { // k == d
							x = vf[offset+k-1] + 1 // Move right
						}
					} else if vf[idx] != -1 {
						// If already visited (vf[idx] != -1), potentially update if a new path is better?
						// The standard algorithm updates vf[idx] later unconditionally.
						// Let's stick to the standard: determine x based on prev k's.
						switch k {
						case -d: // Must come from k+1 (move down)
							x = vf[offset+k+1]
						case d: // Must come from k-1 (move right)
							x = vf[offset+k-1] + 1
						default: // Compare k-1 and k+1
							// Choose the path that reached further previously.
							// Prioritize moving right if the previous x-values are equal,
							// as moving right advances x by 1.
							if vf[offset+k-1] >= vf[offset+k+1] {
								x = vf[offset+k-1] + 1 // Move right
							} else {
								x = vf[offset+k+1] // Move down
							}
						}
					} else {
						// Should not happen if d>0 and k in [-d, d]
						continue // Cannot proceed for this k
					}
				}
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

			// Update the furthest reaching x for this diagonal k
			vf[idx] = x

			// --- Overlap Check ---
			// Check if the forward path (vf) has overlapped with the reverse path (vr).
			// Conditions:
			// 1. delta % 2 != 0: Overlap can only happen when D (number of edits) differs
			//    by one between forward and reverse searches. This happens when total
			//    edit path length (N+M) is odd, and we are checking after the forward pass
			//    of step 'd'. The reverse pass is effectively at step 'd-1'.
			// 2. k is within the bounds of diagonals explored by the reverse search
			//    in the *previous* step (d-1). The reverse diagonals range from
			//    delta - (d-1) to delta + (d-1).
			// 3. rvk = offset + k - delta: Calculate the index in 'vr' corresponding to the
			//    diagonal 'k' in the forward search, viewed from the reverse perspective.
			// 4. vr[rvk] >= 0: Ensure the reverse path has actually reached this diagonal.
			// 5. x >= vr[rvk]: The critical overlap condition. The furthest x reached by
			//    the forward search (x) on diagonal k is >= the furthest x reached
			//    by the reverse search *starting from the end* (vr[rvk]) on the
			//    corresponding reverse diagonal. This means the paths have met or crossed.
			if delta%2 != 0 && k >= delta-(d-1) && k <= delta+(d-1) {
				rvk := offset + k - delta
				if rvk >= 0 && rvk < vectorSize && vr[rvk] >= 0 && x >= vr[rvk] {
					// Overlap detected!
					// We only return a snake if we actually moved diagonally (startX < x).
					if startX < x {
						// Return the snake found (coordinates are absolute)
						return snake{
							startX: aStart + startX, startY: bStart + startY,
							endX: aStart + x, endY: bStart + y,
							length: x - startX, // Length of the diagonal run
						}
					}
					// Overlap occurred right at the start point of this step's exploration
					// for diagonal k. We need to find the snake ending here from the *previous* step.
					// The snake effectively has length 0 for *this* D-path step,
					// but the overlap point itself is the 'middle'.
					// We return the point (startX, startY) as a zero-length snake.
					return snake{
						startX: aStart + startX, startY: bStart + startY,
						endX: aStart + startX, endY: bStart + startY,
						length: 0,
					}
				}
			}
		} // End forward k loop

		// --- Reverse Search ---
		// Explore reverse diagonals k' from -d to d (relative to the reverse D-path)
		for k_rev := -d; k_rev <= d; k_rev += 2 {
			// Map the reverse diagonal k_rev to the corresponding forward diagonal k
			k := k_rev + delta
			idx := offset + k // Array index for diagonal k
			if idx < 0 || idx >= vectorSize {
				continue
			} // Bounds check

			// Determine the starting x for this step on diagonal k in the *reverse* direction.
			// We prioritize moving up (from k+1's perspective, which is k+1 in forward terms)
			// if it reaches a smaller x than moving left (from k-1).
			var x int
			// Reverse indices for vr access based on k, not k_rev
			canMoveLeft := k > -d+delta && (offset+k-1) >= 0 && vr[offset+k-1] != -1
			canMoveUp := k < d+delta && (offset+k+1) < vectorSize && vr[offset+k+1] != -1

			if !canMoveLeft && canMoveUp { //nolint:gocritic //  ifElseChain is fine here
				x = vr[offset+k+1] - 1 // Must move up (effectively from k+1, requires x decrements)
			} else if canMoveLeft && !canMoveUp {
				x = vr[offset+k-1] // Must move left (effectively from k-1)
			} else if canMoveLeft && canMoveUp {
				// Choose the path that starts further back (smaller x)
				if vr[offset+k-1] < vr[offset+k+1]-1 { // Compare x values from prev step
					x = vr[offset+k-1] // Move left from k-1
				} else {
					x = vr[offset+k+1] - 1 // Move up from k+1
				}
			} else {
				// Base case: d=0, k_rev=0 => k=delta. Handled by vr[offset+delta]=n init.
				if idx == offset+delta && d == 0 {
					x = vr[idx]
				} else {
					// Similar robustness considerations as forward pass
					if k == -d+delta { //nolint:gocritic // Must come from k+s, ifElseChain is fine here
						x = vr[offset+k+1] - 1 // Move up
					} else if k == d+delta { // Must come from k-1
						x = vr[offset+k-1] // Move left
					} else if vr[idx] != -1 { // If already visited
						// Determine based on comparison
						if vr[offset+k-1] <= vr[offset+k+1]-1 { // Check comparison logic carefully
							x = vr[offset+k-1] // Move left
						} else {
							x = vr[offset+k+1] - 1 // Move up
						}
					} else {
						continue // Cannot proceed
					}
				}
			}

			y := x - k // Calculate y based on x and forward diagonal k

			// Store the ending point of the potential snake for this reverse step
			endX := x
			endY := y

			// Follow the diagonal backwards (matches) as far as possible
			// Need to check bounds carefully when moving backwards (x-1, y-1)
			for x > 0 && y > 0 && // Ensure we don't go below the subproblem start
				aStart+x-1 < len(a) && bStart+y-1 < len(b) && // Check original slice bounds
				a[aStart+x-1] == b[bStart+y-1] { // Check equality
				x--
				y--
			}

			// Update the furthest-back reaching x for this diagonal k in reverse search
			vr[idx] = x

			// --- Overlap Check ---
			// Check if the reverse path (vr) has overlapped with the forward path (vf).
			// Conditions:
			// 1. delta % 2 == 0: Overlap can happen when D is the same for both forward
			//    and reverse searches. This occurs when the total edit path length (N+M)
			//    is even, and we are checking after the reverse pass of step 'd'.
			// 2. k is within the bounds of diagonals explored by the forward search
			//    in the *current* step 'd'. Forward diagonals range from -d to d.
			// 3. fk = offset + k: Calculate the index in 'vf' corresponding to diagonal k.
			// 4. vf[fk] >= 0: Ensure the forward path has actually reached this diagonal.
			// 5. vf[fk] >= x: The critical overlap condition. The furthest x reached by
			//    the forward search (vf[fk]) on diagonal k is >= the furthest x reached
			//    *backwards* by the reverse search (x) on the same diagonal k.
			//    This signifies the paths have met or crossed.
			if delta%2 == 0 && k >= -d && k <= d {
				fk := offset + k
				if fk >= 0 && fk < vectorSize && vf[fk] >= 0 && vf[fk] >= x {
					// Overlap detected!
					// We only return a snake if we actually moved diagonally backwards (x < endX).
					if x < endX {
						// Return the snake found (coords are absolute, start is (x,y), end is (endX,endY))
						return snake{
							startX: aStart + x, startY: bStart + y,
							endX: aStart + endX, endY: bStart + endY,
							length: endX - x, // Length of the diagonal run
						}
					}
					// Overlap occurred right at the end point of this reverse step's exploration.
					// Return the point (x, y) as a zero-length snake.
					return snake{
						startX: aStart + x, startY: bStart + y,
						endX: aStart + x, endY: bStart + y,
						length: 0,
					}
				}
			}
		} // End reverse k loop
	} // End d loop

	// Should not be reached if inputs are valid and algorithm is correct,
	// but return a fallback zero-length snake at the start just in case.
	// This might happen if maxDiff is somehow too small, though calculated correctly.
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

// simpleDiff is a fallback diff algorithm for when Myers becomes too expensive
func simpleDiff(a, b []string) []diff.Line {
	edits := []diff.Line{}

	// Use a simple longest common subsequence approach
	longest := lcs.LongestCommonSubsequence(a, b)
	aIndex, bIndex := 0, 0

	for _, item := range longest {
		// Add deletions for unmatched items in a
		for aIndex < item.AIndex {
			edits = append(edits, diff.Line{Kind: diff.Delete, Text: a[aIndex]})
			aIndex++
		}

		// Add insertions for unmatched items in b
		for bIndex < item.BIndex {
			edits = append(edits, diff.Line{Kind: diff.Insert, Text: b[bIndex]})
			bIndex++
		}

		// Add the matching item
		edits = append(edits, diff.Line{Kind: diff.Equal, Text: a[aIndex]})
		aIndex++
		bIndex++
	}

	// Handle remaining items in a
	for aIndex < len(a) {
		edits = append(edits, diff.Line{Kind: diff.Delete, Text: a[aIndex]})
		aIndex++
	}

	// Handle remaining items in b
	for bIndex < len(b) {
		edits = append(edits, diff.Line{Kind: diff.Insert, Text: b[bIndex]})
		bIndex++
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
