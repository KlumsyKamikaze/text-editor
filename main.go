package main

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

var og_termios_cofig unix.Termios

func die(s error) {
	fmt.Println(s.Error())
	os.Exit(1)
}

func disableRawMode() {
	err := termios.Tcsetattr(termios.TCIFLUSH, termios.TCSAFLUSH, &og_termios_cofig)
	if err != nil {
		die(err)
	}
}

func enableRawMode() {
	err := termios.Tcgetattr(termios.TCIFLUSH, &og_termios_cofig)
	if err != nil {
		panic(err)
	}
	var raw unix.Termios = og_termios_cofig

	ECHO := uint32(unix.ECHO)
	ICANON := uint32(unix.ICANON)
	ISIG := uint32(unix.ISIG)
	raw.Iflag &= ^(uint32(unix.IXON) | uint32(unix.ICRNL) | uint32(unix.BRKINT) | uint32(unix.INPCK) | uint32(unix.ISTRIP))
	raw.Lflag &= ^(ECHO | ICANON | ISIG | uint32(unix.IEXTEN))
	raw.Cflag |= (uint32(unix.CS8))
	raw.Oflag &= ^(uint32(unix.OPOST))

	raw.Cc[unix.VMIN] = 0 // setting the timeouts for read()
	raw.Cc[unix.VTIME] = 100

	err = termios.Tcsetattr(termios.TCIFLUSH, termios.TCSAFLUSH, &raw)
	if err != nil {
		die(err)
	}
}

func CTRL_KEY(k byte) byte {
	return ((k) & 0x1f)
}

// terminal

func editorReadKey() []byte {
	c := make([]byte, 1)
	for bytesRead, err := io.ReadAtLeast(os.Stdin, c, 1); bytesRead != 1; {
		if err.(syscall.Errno) != unix.EAGAIN {
			die(err)
		}
	}
	return c
}

// input
func editorProcessKeypress() {
	switch c := editorReadKey(); c[0] {
	case CTRL_KEY('q'):
		os.Exit(0)
	}
}

// init
func main() {
	enableRawMode()
	defer disableRawMode()

	for {
		editorProcessKeypress()
	}
	// for {
	// 	c := make([]byte, 1)
	// 	if _, err := io.ReadAtLeast(os.Stdin, c, 1); err != nil {
	// 		die(err)
	// 	}
	// 	if unicode.IsControl(rune(c[0])) {
	// 		fmt.Printf("%d\r\n", c)
	// 	} else {
	// 		fmt.Printf("%d ('%c')\r\n", c, c)
	// 	}
	// 	if c[0] == CTRL_KEY('q') {
	// 		break
	// 	}
	// }
}
