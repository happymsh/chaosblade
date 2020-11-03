package hzcpbin

import (
	"flag"
	"fmt"
	"os"

	"github.com/chaosblade-io/chaosblade-spec-go/util"
)

const HzcpJarName string = "myAgent-jar-with-dependencies.jar"
const ErrPrefix = "Error:"

var ExitFunc = os.Exit

func PrintAndExitWithErrPrefix(message string) {
	fmt.Fprint(os.Stderr, fmt.Sprintf("%s %s", ErrPrefix, message))
	ExitFunc(1)
}

func PrintErrAndExit(message string) {
	fmt.Fprint(os.Stderr, message)
	ExitFunc(1)
}

func PrintOutputAndExit(message string) {
	fmt.Fprintf(os.Stdout, message)
	ExitFunc(0)
}

func ParseFlagAndInitLog() {
	util.AddDebugFlag()
	flag.Parse()
	util.InitLog(util.Bin)
}
