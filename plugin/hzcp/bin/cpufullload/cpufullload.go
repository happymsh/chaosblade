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
	"github.com/sirupsen/logrus"
)

var (
	startFlag, stopFlag bool
	keyword, pid, port  string
)

func main() {
	flag.BoolVar(&startFlag, "start", false, "start hzcpcpufullload")
	flag.BoolVar(&stopFlag, "stop", false, "stop hzcpcpufullload")
	flag.StringVar(&keyword, "keyword", "", "target java process keyword")
	flag.StringVar(&port, "port", "36662", "nouse port")
	hzcpbin.ParseFlagAndInitLog()

	if keyword == "" {
		logrus.Error("hzcp-cpufullload: target Java process match error: keyword is empty")
		hzcpbin.PrintErrAndExit("target Java process match error: keyword is empty")
	}

	if port == "" {
		port = "36662"
	}

	ctx := context.Background()
	excludeStr := fmt.Sprintf("'keyword %s'", keyword)
	ctxEx := context.WithValue(ctx, cl.ExcludeProcessKey, excludeStr)
	pids, _ := cl.GetPidsByProcessName(keyword, ctxEx)
	if pids != nil && len(pids) == 1 {
		logrus.Info("hzcp-cpufullload pid:", pids)
		pid = pids[0]
	} else if len(pids) > 0 {
		hzcpbin.PrintErrAndExit("target Java process match error: Too many records")
	} else {
		hzcpbin.PrintErrAndExit("target Java process match error: no records")
	}

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

const hzcpCpuJarName string = "cpu-1.0-jar-with-dependencies.jar"

// start burn io
func start() {
	hzcpJar := path.Join(util.GetProgramPath(), "hzcp", hzcpCpuJarName)
	ctx := context.Background()
	args := fmt.Sprintf(`java -Xdebug -Xrunjdwp:transport=dt_socket,address=%s,server=y,suspend=n -Djava.ext.dirs=${JAVA_HOME}/lib -jar %s LoadAgent %s %s > %s 2>&1 &`,
		port, hzcpJar, pid, hzcpJar, logFile)
	logrus.Info("hzcp-cpufullload nohup args: ", args)
	response := channel.Run(ctx, "nohup", args)
	if !response.Success {
		hzcpbin.PrintErrAndExit(response.Err)
		return
	}
	hzcpbin.PrintOutputAndExit("success")
}

func stop() {
	ctx := context.Background()
	pids, _ := cl.GetPidsByProcessName(hzcpCpuJarName, ctx)
	if pids != nil && len(pids) > 0 {
		_ = channel.Run(ctx, "kill", fmt.Sprintf("-9 %s", strings.Join(pids, " ")))
	}
}
