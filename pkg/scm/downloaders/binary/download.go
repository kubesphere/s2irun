package binary

import (
	"github.com/kubesphere/s2irun/pkg/api"
	"github.com/kubesphere/s2irun/pkg/api/constants"
	"github.com/kubesphere/s2irun/pkg/scm/git"
	"github.com/kubesphere/s2irun/pkg/utils/bytefmt"
	"github.com/kubesphere/s2irun/pkg/utils/fs"
	utilglog "github.com/kubesphere/s2irun/pkg/utils/glog"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var glog = utilglog.StderrLog

// File represents a simplest possible Downloader implementation where the
// sources are just copied from local directory.
type File struct {
	fs.FileSystem
}

// Download download sources from a http link into the working directory.
// Caller guarantees that config.Source.IsLocal() is true.
func (f *File) Download(config *api.Config) (*git.SourceInfo, error) {
	_, filename := filepath.Split(config.Source.String())
	config.WorkingSourceDir = filepath.Join(config.WorkingDir, constants.Source)
	binaryPath := filepath.Join(config.WorkingSourceDir, filename)
	glog.V(0).Infof("Start Download Binary %s", filename)
	resp, err := http.Get(config.Source.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	out, err := f.Create(binaryPath)
	if err != nil {
		return nil, err
	}

	defer out.Close()

	strsize := resp.Header.Get("Content-Length")
	size, _ := strconv.ParseUint(strsize, 10, 64)

	counter := &WriteCounter{Size: size}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	glog.V(0).Infof("Finish Download Binary %s", filename)
	return &git.SourceInfo{
		Location:   config.Source.String(),
		ContextDir: config.ContextDir,
	}, nil
}

// write cycle.
type WriteCounter struct {
	Total uint64
	Size  uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	glog.V(0).Infof("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	if wc.Size > 0 {
		glog.V(0).Infof("\rDownloading... %s//%s", bytefmt.ByteSize(wc.Total), bytefmt.ByteSize(wc.Size))
	} else {
		glog.V(0).Infof("\rDownloading... %s complete", bytefmt.ByteSize(wc.Total))
	}

}
