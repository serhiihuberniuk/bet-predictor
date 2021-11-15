package scanner

import (
	"bufio"
	"fmt"
	"os"
)

func ScanWithMessage(message string) string {
	fmt.Print(message)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}
