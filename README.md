# go-term-color

Go package for colouring terminals in CLI applications. Es una biblioteca ligera y sencilla de usar que detecta automáticamente si el terminal y el sistema operativo soportan la inyección de colores ANSI, activándolos o desactivándolos en consecuencia, e incluso soporta la variable de entorno `NO_COLOR`.

## Estructura del Proyecto

El proyecto se estructura de la siguiente manera:

*   **`color.go`**: Contiene la lógica principal de la librería. Define los tipos genéricos como `Attribute` y `Color`, así como constantes de colores básicos ANSI y métodos de formato rápido como `Red`, `Green`, `Blue`, `CyanString`, etc. Las detecciones base de TTY se realizan aquí y gestionan globalmente `NoColor`.
*   **`console_posix.go` y `console_windows.go`**: Archivos de compatibilidad específicos a nivel de capa del sistema operativo para verificar e inicializar configuraciones propias del terminal o de las consolas de Windows nativas (habilitando de virtual terminal processing).
*   **`example/`**: Contiene un pequeño sub-módulo para testear y visualizar rápidamente cómo se muestran los colores en la terminal al ejecutarse de manera simple.

## Uso Básico

```go
package main

import (
    "fmt"
    "github.com/craftions/go-term-color"
)

func main() {
    // Uso rápido y directo
    color.Red("Este mensaje sale en color rojo predeterminado")
    color.Green("Verde predeterminado sobre la misma línea")

    // Personalización y construcción de combinaciones complejas con atributos múltiples
    c := color.New(color.FgCyan, color.Bold, color.Underline)
    c.Println("Texto Cyan, en negrita y subrayado")

    // Manipulando para imprimir posteriormente vía strings (útil en logs en caso de encadenar)
    str := color.YellowString("Cadena string explícitamente formateada como amarilla")
    fmt.Printf("Puedo interpolar esto: %s\n", str)
}
```

## Ambientes Docker por Sistema Operativo

A fines de asegurar y mejorar el funcionamiento, prueba y desarrollo de la librería transversalmente de manera agnóstica; se incluyen Dockerfiles para poder compilar y testear la librería en entornos asilados específicos para el respectivo sistema operativo.

### Linux

Utiliza la imagen oficial de sistema ligera Alpine Linux para probar integraciones Unix POSIX.
```bash
docker build -t go-term-color-linux -f Dockerfile.linux .
docker run --rm go-term-color-linux
```

### Windows

Utiliza una imagen de Windows Server Core estricta. *Nota: Necesitas y es estrictamente mandatorio tener Docker Desktop local corriendo en su modo **Windows Containers**.*
```powershell
docker build -t go-term-color-win -f Dockerfile.windows .
docker run --rm go-term-color-win
```

### macOS (POSIX)

Docker no soporta virtualizaciones nativas del _kernel_ de OS X/macOS empaquetadas en un único contenedor. Sin embargo, para probar el _set_ de reglas POSIX de la librería que usa indirectamente macOS (a saber `console_posix.go`), se provee una imagen oficial enriquecida de Debian que refleja la completitud del ecosistema posix equivalente.
```bash
docker build -t go-term-color-mac -f Dockerfile.mac .
docker run --rm go-term-color-mac
```
