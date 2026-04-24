//go:build windows && !appengine
// +build windows,!appengine

// This file is the "translator" that reads ANSI sequences and converts them into cursor movements,
// screen clears, and color changes that Windows understands.
package colorable

import (
	"bytes"
	"strconv"
	"strings"
	"unsafe"

	syscall "golang.org/x/sys/windows"
)

// `\033]0;TITLESTR\007`
func doTitleSequence(er *bytes.Reader) error {
	var c byte
	var err error

	c, err = er.ReadByte()
	if err != nil {
		return err
	}
	if c != '0' && c != '2' {
		return nil
	}
	c, err = er.ReadByte()
	if err != nil {
		return err
	}
	if c != ';' {
		return nil
	}
	title := make([]byte, 0, 80)
	for {
		c, err = er.ReadByte()
		if err != nil {
			return err
		}
		if c == 0x07 || c == '\n' {
			break
		}
		title = append(title, c)
	}
	if len(title) > 0 {
		title8, err := syscall.UTF16PtrFromString(string(title))
		if err == nil {
			procSetConsoleTitle.Call(uintptr(unsafe.Pointer(title8)))
		}
	}
	return nil
}

// returns Atoi(s) unless s == "" in which case it returns def
func atoiWithDefault(s string, def int) (int, error) {
	if s == "" {
		return def, nil
	}
	return strconv.Atoi(s)
}

// Write writes data on console
func (w *writer) Write(data []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	var csbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(w.handle), uintptr(unsafe.Pointer(&csbi)))

	handle := w.handle

	var er *bytes.Reader
	if w.rest.Len() > 0 {
		var rest bytes.Buffer
		w.rest.WriteTo(&rest)
		w.rest.Reset()
		rest.Write(data)
		er = bytes.NewReader(rest.Bytes())
	} else {
		er = bytes.NewReader(data)
	}
	var plaintext bytes.Buffer
loop:
	for {
		c1, err := er.ReadByte()
		if err != nil {
			plaintext.WriteTo(w.out)
			break loop
		}
		if c1 != 0x1b {
			plaintext.WriteByte(c1)
			continue
		}
		_, err = plaintext.WriteTo(w.out)
		if err != nil {
			break loop
		}
		c2, err := er.ReadByte()
		if err != nil {
			break loop
		}

		switch c2 {
		case '>':
			continue
		case ']':
			w.rest.WriteByte(c1)
			w.rest.WriteByte(c2)
			er.WriteTo(&w.rest)
			if bytes.IndexByte(w.rest.Bytes(), 0x07) == -1 {
				break loop
			}
			er = bytes.NewReader(w.rest.Bytes()[2:])
			err := doTitleSequence(er)
			if err != nil {
				break loop
			}
			w.rest.Reset()
			continue
		// https://github.com/mattn/go-colorable/issues/27
		case '7':
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			w.oldpos = csbi.cursorPosition
			continue
		case '8':
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&w.oldpos)))
			continue
		case 0x5b:
			// execute part after switch
		default:
			continue
		}

		w.rest.WriteByte(c1)
		w.rest.WriteByte(c2)
		er.WriteTo(&w.rest)

		var buf bytes.Buffer
		var m byte
		for i, c := range w.rest.Bytes()[2:] {
			if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '@' {
				m = c
				er = bytes.NewReader(w.rest.Bytes()[2+i+1:])
				w.rest.Reset()
				break
			}
			buf.Write([]byte(string(c)))
		}
		if m == 0 {
			break loop
		}

		switch m {
		case 'A':
			n, err = atoiWithDefault(buf.String(), 1)
			if err != nil {
				continue
			}
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.y -= short(n)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'B':
			n, err = atoiWithDefault(buf.String(), 1)
			if err != nil {
				continue
			}
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.y += short(n)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'C':
			n, err = atoiWithDefault(buf.String(), 1)
			if err != nil {
				continue
			}
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.x += short(n)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'D':
			n, err = atoiWithDefault(buf.String(), 1)
			if err != nil {
				continue
			}
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.x -= short(n)
			if csbi.cursorPosition.x < 0 {
				csbi.cursorPosition.x = 0
			}
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'E':
			n, err = atoiWithDefault(buf.String(), 1)
			if err != nil {
				continue
			}
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.x = 0
			csbi.cursorPosition.y += short(n)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'F':
			n, err = atoiWithDefault(buf.String(), 1)
			if err != nil {
				continue
			}
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.x = 0
			csbi.cursorPosition.y -= short(n)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'G':
			n, err = strconv.Atoi(buf.String())
			if err != nil {
				continue
			}
			if n < 1 {
				n = 1
			}
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			csbi.cursorPosition.x = short(n - 1)
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'H', 'f':
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			if buf.Len() > 0 {
				token := strings.Split(buf.String(), ";")
				switch len(token) {
				case 1:
					n1, err := strconv.Atoi(token[0])
					if err != nil {
						continue
					}
					csbi.cursorPosition.y = short(n1 - 1)
				case 2:
					n1, err := strconv.Atoi(token[0])
					if err != nil {
						continue
					}
					n2, err := strconv.Atoi(token[1])
					if err != nil {
						continue
					}
					csbi.cursorPosition.x = short(n2 - 1)
					csbi.cursorPosition.y = short(n1 - 1)
				}
			} else {
				csbi.cursorPosition.y = 0
			}
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&csbi.cursorPosition)))
		case 'J':
			n := 0
			if buf.Len() > 0 {
				n, err = strconv.Atoi(buf.String())
				if err != nil {
					continue
				}
			}
			var count, written dword
			var cursor coord
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			switch n {
			case 0:
				cursor = coord{x: csbi.cursorPosition.x, y: csbi.cursorPosition.y}
				count = dword(csbi.size.x) - dword(csbi.cursorPosition.x) + dword(csbi.size.y-csbi.cursorPosition.y)*dword(csbi.size.x)
			case 1:
				cursor = coord{x: csbi.window.left, y: csbi.window.top}
				count = dword(csbi.size.x) - dword(csbi.cursorPosition.x) + dword(csbi.window.top-csbi.cursorPosition.y)*dword(csbi.size.x)
			case 2:
				cursor = coord{x: csbi.window.left, y: csbi.window.top}
				count = dword(csbi.size.x) - dword(csbi.cursorPosition.x) + dword(csbi.size.y-csbi.cursorPosition.y)*dword(csbi.size.x)
			}
			procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(count), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))
			procFillConsoleOutputAttribute.Call(uintptr(handle), uintptr(csbi.attributes), uintptr(count), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))
		case 'K':
			n := 0
			if buf.Len() > 0 {
				n, err = strconv.Atoi(buf.String())
				if err != nil {
					continue
				}
			}
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			var cursor coord
			var count, written dword
			switch n {
			case 0:
				cursor = coord{x: csbi.cursorPosition.x, y: csbi.cursorPosition.y}
				count = dword(csbi.size.x - csbi.cursorPosition.x)
			case 1:
				cursor = coord{x: csbi.window.left, y: csbi.cursorPosition.y}
				count = dword(csbi.size.x - csbi.cursorPosition.x)
			case 2:
				cursor = coord{x: csbi.window.left, y: csbi.cursorPosition.y}
				count = dword(csbi.size.x)
			}
			procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(count), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))
			procFillConsoleOutputAttribute.Call(uintptr(handle), uintptr(csbi.attributes), uintptr(count), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))
		case 'X':
			n := 0
			if buf.Len() > 0 {
				n, err = strconv.Atoi(buf.String())
				if err != nil {
					continue
				}
			}
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			var cursor coord
			var written dword
			cursor = coord{x: csbi.cursorPosition.x, y: csbi.cursorPosition.y}
			procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(n), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))
			procFillConsoleOutputAttribute.Call(uintptr(handle), uintptr(csbi.attributes), uintptr(n), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&written)))

		case 'm':
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			attr := csbi.attributes
			cs := buf.String()
			if cs == "" {
				procSetConsoleTextAttribute.Call(uintptr(handle), uintptr(w.oldattr))
				continue
			}
			token := strings.Split(cs, ";")
			for i := 0; i < len(token); i++ {
				ns := token[i]
				if n, err = strconv.Atoi(ns); err == nil {
					switch {
					case n == 0 || n == 100:
						attr = w.oldattr
					case n == 4:
						attr |= commonLvbUnderscore
					case (1 <= n && n <= 3) || n == 5:
						attr |= foregroundIntensity
					case n == 7 || n == 27:
						attr =
							(attr &^ (foregroundMask | backgroundMask)) |
								((attr & foregroundMask) << 4) |
								((attr & backgroundMask) >> 4)
					case n == 22:
						attr &^= foregroundIntensity
					case n == 24:
						attr &^= commonLvbUnderscore
					case 30 <= n && n <= 37:
						attr &= backgroundMask
						if (n-30)&1 != 0 {
							attr |= foregroundRed
						}
						if (n-30)&2 != 0 {
							attr |= foregroundGreen
						}
						if (n-30)&4 != 0 {
							attr |= foregroundBlue
						}
					case n == 38: // set foreground color.
						if i < len(token)-2 && (token[i+1] == "5" || token[i+1] == "05") {
							if n256, err := strconv.Atoi(token[i+2]); err == nil {
								if n256foreAttr == nil {
									n256setup()
								}
								attr &= backgroundMask
								attr |= n256foreAttr[n256%len(n256foreAttr)]
								i += 2
							}
						} else if len(token) == 5 && token[i+1] == "2" {
							var r, g, b int
							r, _ = strconv.Atoi(token[i+2])
							g, _ = strconv.Atoi(token[i+3])
							b, _ = strconv.Atoi(token[i+4])
							i += 4
							if r > 127 {
								attr |= foregroundRed
							}
							if g > 127 {
								attr |= foregroundGreen
							}
							if b > 127 {
								attr |= foregroundBlue
							}
						} else {
							attr = attr & (w.oldattr & backgroundMask)
						}
					case n == 39: // reset foreground color.
						attr &= backgroundMask
						attr |= w.oldattr & foregroundMask
					case 40 <= n && n <= 47:
						attr &= foregroundMask
						if (n-40)&1 != 0 {
							attr |= backgroundRed
						}
						if (n-40)&2 != 0 {
							attr |= backgroundGreen
						}
						if (n-40)&4 != 0 {
							attr |= backgroundBlue
						}
					case n == 48: // set background color.
						if i < len(token)-2 && token[i+1] == "5" {
							if n256, err := strconv.Atoi(token[i+2]); err == nil {
								if n256backAttr == nil {
									n256setup()
								}
								attr &= foregroundMask
								attr |= n256backAttr[n256%len(n256backAttr)]
								i += 2
							}
						} else if len(token) == 5 && token[i+1] == "2" {
							var r, g, b int
							r, _ = strconv.Atoi(token[i+2])
							g, _ = strconv.Atoi(token[i+3])
							b, _ = strconv.Atoi(token[i+4])
							i += 4
							if r > 127 {
								attr |= backgroundRed
							}
							if g > 127 {
								attr |= backgroundGreen
							}
							if b > 127 {
								attr |= backgroundBlue
							}
						} else {
							attr = attr & (w.oldattr & foregroundMask)
						}
					case n == 49: // reset foreground color.
						attr &= foregroundMask
						attr |= w.oldattr & backgroundMask
					case 90 <= n && n <= 97:
						attr = (attr & backgroundMask)
						attr |= foregroundIntensity
						if (n-90)&1 != 0 {
							attr |= foregroundRed
						}
						if (n-90)&2 != 0 {
							attr |= foregroundGreen
						}
						if (n-90)&4 != 0 {
							attr |= foregroundBlue
						}
					case 100 <= n && n <= 107:
						attr = (attr & foregroundMask)
						attr |= backgroundIntensity
						if (n-100)&1 != 0 {
							attr |= backgroundRed
						}
						if (n-100)&2 != 0 {
							attr |= backgroundGreen
						}
						if (n-100)&4 != 0 {
							attr |= backgroundBlue
						}
					}
					procSetConsoleTextAttribute.Call(uintptr(handle), uintptr(attr))
				}
			}

		case 'h':
			var ci consoleCursorInfo
			switch cs := buf.String(); cs {
			case "5>":
				procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
				ci.visible = 0
				procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
			case "?25":
				procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
				ci.visible = 1
				procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
			case "?1049":
				if w.althandle == 0 {
					h, _, _ := procCreateConsoleScreenBuffer.Call(uintptr(genericRead|genericWrite), 0, 0, uintptr(consoleTextmodeBuffer), 0, 0)
					w.althandle = syscall.Handle(h)
					if w.althandle != 0 {
						handle = w.althandle
					}
				}
			}
		case 'l':
			var ci consoleCursorInfo
			switch cs := buf.String(); cs {
			case "5>":
				procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
				ci.visible = 1
				procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
			case "?25":
				procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
				ci.visible = 0
				procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&ci)))
			case "?1049":
				if w.althandle != 0 {
					syscall.CloseHandle(w.althandle)
					w.althandle = 0
					handle = w.handle
				}
			}
		case 's':
			procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
			w.oldpos = csbi.cursorPosition
		case 'u':
			procSetConsoleCursorPosition.Call(uintptr(handle), *(*uintptr)(unsafe.Pointer(&w.oldpos)))
		}
	}

	return len(data), nil
}
