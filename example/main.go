package main

import (
	"fmt"

	"github.com/craftions/go-term-check/terminal"
	color "github.com/craftions/go-term-color"
)

func main() {
	fmt.Println("=====================================")
	fmt.Println("Ejemplo de uso de go-term-color:")
	fmt.Println("=====================================")
	fmt.Println(terminal.Hello())
	color.Red("Este texto es rojo")
	color.Green("Este texto es verde")
	color.Blue("Este texto es azul")
	color.Yellow("Este texto es amarillo")

	// Utilizando múltiples atributos
	c := color.New(color.FgMagenta, color.Bold, color.Underline)
	c.Println("Texto magenta, en negrita y subrayado")

	// Usando strings
	text := color.CyanString("Cyan devuelto como string")
	fmt.Printf("Podemos intercalar de forma sencilla: %s dentro de un Printf estándar.\n\n", text)
}
