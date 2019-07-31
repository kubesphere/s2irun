package binary

import (
	"github.com/kubesphere/s2irun/pkg/api"
	"github.com/kubesphere/s2irun/pkg/scm/git"
	testfs "github.com/kubesphere/s2irun/pkg/test/fs"
	"testing"
)

func TestDownload(t *testing.T) {
	fs := &testfs.FakeFileSystem{}
	f := &File{fs}

	config := &api.Config{
		Source:      git.MustParse("https://kubesphere.io/etcd-operator.svg"),
		IsBinaryURL: true,
	}
	info, err := f.Download(config)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if fs.CreateFile != "upload/src/etcd-operator.svg" {
		t.Errorf("Unexpected fs.CreateFile %s", fs.CreateFile)
	}
	if info.Location != config.Source.String() || info.ContextDir != config.ContextDir {
		t.Errorf("Unexpected info")
	}
}
