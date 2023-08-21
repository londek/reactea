package render

import "strings"

const floatMultiplier = 1000000.0 // Acceptable precision

// No allocations variant, otherwise could be done
// by spliting each line and comparing len(line)
func WidthOf(str string) int {
	maxWidth := 0
	lineWidth := 0

	for _, char := range str {
		if char == '\n' {
			if maxWidth < lineWidth {
				maxWidth = lineWidth
			}

			lineWidth = 0

			continue
		}

		lineWidth += 1
	}

	if maxWidth < lineWidth {
		maxWidth = lineWidth
	}

	return maxWidth
}

func HeightOf(str string) int {
	return strings.Count(str, "\n") + 1
}

func SizeOf(str string) (int, int) {
	return WidthOf(str), HeightOf(str)
}

func Clip(str string, width, height int) string {
	var b strings.Builder
	b.Grow(width * height)

	lineWidth := 0
	totalLines := 0
	skipLine := false

	for _, char := range str {
		if char == '\n' {
			lineWidth = 0
			totalLines += 1

			if totalLines > height {
				break
			}

			b.WriteRune('\n')

			skipLine = false
			continue
		}

		if !skipLine {
			b.WriteRune(char)

			lineWidth += 1

			if lineWidth > width {
				lineWidth = 0
				skipLine = true
			}
		}
	}

	return b.String()
}

// Slightly more efficient than Clip when its not needed
func ClipWidth(str string, width int) string {
	var b strings.Builder
	b.Grow(width * strings.Count(str, "\n"))

	lineWidth := 0
	skipLine := false

	for _, char := range str {
		if char == '\n' {
			lineWidth = 0

			b.WriteRune('\n')

			skipLine = false
			continue
		}

		if !skipLine {
			b.WriteRune(char)

			lineWidth += 1

			if lineWidth > width {
				lineWidth = 0
				skipLine = true
			}
		}
	}

	return b.String()
}

// Slightly more efficient than Clip when its not needed
func ClipHeight(str string, height int) string {
	var b strings.Builder
	b.Grow(100 * height)

	totalLines := 0

	for _, char := range str {
		if char == '\n' {
			totalLines += 1

			if totalLines > height {
				break
			}

			b.WriteRune('\n')

			continue
		}

		b.WriteRune(char)
	}

	return b.String()
}
