package args

import (
	"fmt"
	"os"
)

const Version string = "0.2.0"

func printVersion() {
	fmt.Println(Version)
	os.Exit(0)
}
