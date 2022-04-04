package args

import (
	"fmt"
	"os"
)

const Version string = "0.1.0"

func printVersion() {
	fmt.Println(Version)
	os.Exit(0)
}
