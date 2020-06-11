package exec

import (
	"context"
	"testing"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

func TestCallExp(t *testing.T) {
	success := spec.ReturnSuccess("call exp success.")
	CRun = func(ctx context.Context, script, args string) *spec.Response {
		if args == "create cpu fullload --cpu-percent=10 --cpu-count=1 --timeout=10" {
			return success
		} else {
			return spec.Return(spec.Code[spec.ExecCommandError])
		}
	}
	result := ExecuteExp(context.TODO(), "cpu", "fullload", "--cpu-percent=10 --cpu-count=1", "10")
	if success == result {
		t.Log(result.ToString())
	} else {
		t.Error(result.ToString())
	}
}
