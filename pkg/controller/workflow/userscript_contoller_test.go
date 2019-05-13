package workflow

import (
	lyrav1alpha1 "github.com/lyraproj/lyra-operator/pkg/apis/lyra/v1alpha1"
	"github.com/stretchr/testify/require"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"testing"
)

func TestHandleScript(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(false))
	s := lyrav1alpha1.UserScriptSpec{
		ScriptName: "sayHi.sh",
		GitBranch:  "chaff",
		GitRepo:    "https://github.com/markfuller/k8s-play",
	}
	err := handleScript(s)
	require.NoError(t, err)
}
