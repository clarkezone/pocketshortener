package shortener

import (
	"path"
	"testing"

	"github.com/clarkezone/pocketshorten/internal"
)

func Test_LoadfromConfig(t *testing.T) {
	cpath := "testfiles/.pocketshorten.json"
	if internal.GitRoot == "" {
		t.Fatalf("GitRoot is empty, did you call SetupGitRoot() in test?")
	}
	configpath := path.Join(internal.GitRoot, cpath)
	internal.InitConfig(&configpath)
	viperLoader := &viperLoader{}
	viperLoader.Init(nil)
}
