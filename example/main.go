package main

import (
	"fmt"

	color "github.com/craftions/go-term-color"
)

func main() {
	fmt.Println("=====================================")
	fmt.Println("Example of go-term-color usage:")
	fmt.Println("=====================================")
	color.Red("This text is red")
	color.Green("This text is green")
	color.Blue("This text is blue")
	color.Yellow("This text is yellow")

	// Using multiple attributes
	c := color.New(color.FgMagenta, color.Bold, color.Underline)
	c.Println("Magenta text, bold and underlined")

	// Using strings
	text := color.CyanString("Cyan returned as string")
	fmt.Printf("We can easily interleave: %s inside a standard Printf.\n\n", text)
}
