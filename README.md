# go-term-color

Go package for colouring terminals in CLI applications. It is a lightweight and easy-to-use library that automatically detects whether the terminal and operating system support ANSI color injection, enabling or disabling them accordingly. It also supports the `NO_COLOR` environment variable.

## Project Structure

The project is structured as follows:

*   **`color.go`**: Contains the main logic of the library. It defines generic types such as `Attribute` and `Color`, as well as basic ANSI color constants and quick formatting methods like `Red`, `Green`, `Blue`, `CyanString`, etc. The base TTY detections are performed here, globally managing `NoColor`.
*   **`console_posix.go` and `console_windows.go`**: Specific OS-layer compatibility files to verify and initialize native terminal settings or Windows consoles (enabling virtual terminal processing).
*   **`internal/colorable/`**: Contains the internal implementation to process and render ANSI escape sequences natively on Windows environments that do not support virtual terminal processing by default, interpreting the sequences and translating them into Windows console API calls.
*   **`example/`**: Contains a small sub-module to quickly test and visualize how colors are displayed in the terminal when executed.

## Basic Usage

```go
package main

import (
    "fmt"
    "github.com/craftions/go-term-color"
)

func main() {
    // Quick and direct usage
    color.Red("This message prints in default red color")
    color.Green("Default green on the same line")

    // Customization and building complex combinations with multiple attributes
    c := color.New(color.FgCyan, color.Bold, color.Underline)
    c.Println("Cyan text, bold and underlined")

    // Manipulating to print later via strings (useful for logging if chaining)
    str := color.YellowString("String explicitly formatted as yellow")
    fmt.Printf("I can interpolate this: %s\n", str)
}
```

## Docker Environments per Operating System

To ensure and improve the functionality, testing, and cross-platform development of the library agnostically, Dockerfiles are included to compile and test the library in isolated environments specific to the respective operating system.

### Linux

Uses the official lightweight Alpine Linux image to test Unix POSIX integrations.
```bash
docker build -t go-term-color-linux -f Dockerfile.linux .
docker run --rm go-term-color-linux
```

### Windows

Uses a strict Windows Server Core image. *Note: You strictly need to have Docker Desktop running locally in its **Windows Containers** mode.*
```powershell
docker build -t go-term-color-win -f Dockerfile.windows .
docker run --rm go-term-color-win
```

### macOS (POSIX)

Docker does not support native virtualization of the OS X/macOS _kernel_ packaged in a single container. However, to test the POSIX rule set of the library indirectly used by macOS (namely `console_posix.go`), an enriched official Debian image is provided that reflects the completeness of the equivalent posix ecosystem.
```bash
docker build -t go-term-color-mac -f Dockerfile.mac .
docker run --rm go-term-color-mac
```

## Testing and Coverage

This project emphasizes code quality and reliability through rigorous unit testing, achieving over **80% global code coverage**.

- **Table-Driven Tests:** Tests are implemented using the table-driven pattern to ensure multiple scenarios and attributes are evaluated systematically.
- **Cross-Platform Mocks:** The testing suite leverages invalid file descriptors and environment variable manipulation (`TERM=dumb`, `NO_COLOR=1`) to emulate different terminal states.
- **Windows Console Evaluation:** The `colorable` tests simulate both native Windows consoles (using `CONOUT$`) and non-interactive pipes to guarantee that the `writer` correctly initializes or passes through based on the underlying output stream.

To run tests and check coverage:
```bash
make test
make coverage
```
