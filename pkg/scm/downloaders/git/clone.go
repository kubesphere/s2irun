package git

import (
	"path/filepath"

	"github.com/golang/glog"

	"github.com/kubesphere/s2irun/pkg/api"
	"github.com/kubesphere/s2irun/pkg/api/constants"
	"github.com/kubesphere/s2irun/pkg/scm/git"
	"github.com/kubesphere/s2irun/pkg/utils/fs"
)

// Clone knows how to clone a Git repository.
type Clone struct {
	git.Git
	fs.FileSystem
}

// Download downloads the application source code from the Git repository
// and checkout the Ref specified in the config.
func (c *Clone) Download(config *api.Config) (*git.SourceInfo, error) {
	targetSourceDir := filepath.Join(config.WorkingDir, constants.Source)
	config.WorkingSourceDir = targetSourceDir

	RevisionId := config.RevisionId
	if RevisionId == "" {
		RevisionId = "HEAD"
	}

	if len(config.ContextDir) > 0 {
		targetSourceDir = filepath.Join(config.WorkingDir, constants.ContextTmp)
		glog.V(9).Infof("Downloading %q (%q) ...", config.Source, config.ContextDir)
	} else {
		glog.V(9).Infof("Downloading %q ...", config.Source)
	}

	if !config.IgnoreSubmodules {
		glog.V(2).Infof("Cloning sources into %q", targetSourceDir)
	} else {
		glog.V(2).Infof("Cloning sources (ignoring submodules) into %q", targetSourceDir)
	}

	cloneConfig := git.CloneConfig{Quiet: false}
	err := c.Clone(config.Source, targetSourceDir, cloneConfig)
	if err != nil {
		glog.V(0).Infof("error: git clone failed: %v", err)
		return nil, err
	}

	err = c.Checkout(targetSourceDir, RevisionId)
	if err != nil {
		return nil, err
	}

	glog.V(0).Infof("Checked out to %q", RevisionId)
	if !config.IgnoreSubmodules {
		err = c.SubmoduleUpdate(targetSourceDir, true, true)
		if err != nil {
			return nil, err
		}
		glog.V(0).Infof("Updated submodules for %q", RevisionId)
	}

	info := c.GetInfo(targetSourceDir)
	if len(config.ContextDir) > 0 {
		originalTargetDir := filepath.Join(config.WorkingDir, constants.Source)
		c.RemoveDirectory(originalTargetDir)
		path := filepath.Join(targetSourceDir, config.ContextDir)
		err := c.CopyContents(path, originalTargetDir)
		if err != nil {
			return nil, err
		}
		c.RemoveDirectory(targetSourceDir)
	}

	if len(config.ContextDir) > 0 {
		info.ContextDir = config.ContextDir
	}

	return info, nil
}
