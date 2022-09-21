package displayname

import (
	"fmt"
)

// Our prop(s) is a string itself!
type Props = string

// Stateless components?!?!
// Here we are not using width and height, but you can!
// Using lipgloss styles for example
func Renderer(text Props, width, height int) string {
	return fmt.Sprintf("OMG! Hello %s!", text)
}
