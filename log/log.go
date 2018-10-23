package log

/*
*	
*	Handles logging
*	
*/

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

var (
	RESET     = "\x1b[0;37m"
	BOLD      = "\x1b[1;37m"
	ON_RED    = "\x1b[0;41;37m"
	ON_YELLOW = "\x1b[0;43;37m"
	ON_GREEN  = "\x1b[0;42;37m"
	ON_GREY   = "\x1b[0;40;37m"
	ON_WHITE  = "\x1b[0;47;37m"
)

func Success(message string) {
	fmt.Print(ON_GREEN + " " + RESET + " OK: ", message, "\n")
}

func Info(message string) {
	fmt.Print(ON_GREY + " " + RESET + " INF: ", message, "\n")
}

func Warn(message string) {
	fmt.Print(ON_YELLOW + " " + RESET + " WAR: ", message, "\n")
}

func Error(message string) {
	fmt.Print(ON_RED + " " + RESET + " ERR: ", message, "\n")
}

func Fatal(message string) {
	fmt.Print(ON_RED + " " + RESET + " ERR: ", message, "\n")
	Raw("")
	os.Exit(0)
}

func Raw(message string) {
	fmt.Print(message, "\n")
}

func Prompt(message string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(ON_WHITE + " " + RESET + " " + message + "\n" + ON_WHITE + " " + RESET + " ")

	stdin, _ := reader.ReadString('\n')
    stdin = strings.Replace(stdin, "\n", "", -1)

	return stdin
}
