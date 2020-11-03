package main

import (
	"context"
	"flag"
	"fmt"
	"path"
	"strings"

	cl "github.com/chaosblade-io/chaosblade-spec-go/channel"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
	hzcpbin "github.com/chaosblade-io/chaosblade/plugin/hzcp/bin"
)

var (
	startFlag, stopFlag bool
)

func main() {
	flag.BoolVar(&startFlag, "start", false, "start hzcpcpufullload")
	flag.BoolVar(&stopFlag, "stop", false, "stop hzcpcpufullload")
	hzcpbin.ParseFlagAndInitLog()

	if startFlag {
		start()
	} else if stopFlag {
		stop()
	} else {
		hzcpbin.PrintErrAndExit("less --start or --stop flag")
	}
}

var logFile = util.GetNohupOutput(util.Bin, "hzcpcpufullload.log")
var channel = cl.NewLocalChannel()

// start burn io
func start() {
	hzcpJar := path.Join(util.GetProgramPath(), hzcpbin.HzcpJarName)
	ctx := context.Background()
	response := channel.Run(ctx, "nohup",
		fmt.Sprintf(`java -Xdebug -Xrunjdwp:transport=dt_socket,address=8786,server=y,suspend=n 
-Djava.ext.dirs=${JAVA_HOME}/lib -jar %s LoadAgent 1 %s > %s 2>&1 &`,
			hzcpJar, hzcpJar, logFile))
	if !response.Success {
		hzcpbin.PrintErrAndExit(response.Err)
		return
	}
	hzcpbin.PrintOutputAndExit("success")
}

func stop() {
	ctx := context.Background()
	pids, _ := cl.GetPidsByProcessName(hzcpbin.HzcpJarName, ctx)
	if pids != nil && len(pids) > 0 {
		_ = channel.Run(ctx, "kill", fmt.Sprintf("-9 %s", strings.Join(pids, " ")))
	}
}
