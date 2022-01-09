package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getInput(reader *bufio.Reader) string {
	t, _ := reader.ReadString('\n')
	return strings.TrimSpace(t);
}

func main() {
	reader := bufio.NewReader(os.Stdin);
	fmt.Print("Pslang 1.0.0\n")
	fmt.Print("Type exit to exit the program or press Ctrl-D.\n")

	fmt.Print(">>> ")
	line := getInput(reader)
	for ; !strings.EqualFold("exit", line); line = getInput(reader) {
		scanner Scanner = new Scanner()
		fmt.Print(">>> ")
	}
	fmt.Print("Bai bai!\n")
}