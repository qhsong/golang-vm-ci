package main

import (
	"os"
	"testing"

	"github.com/libgit2/git2go"
	common "github.com/qhsong/golang-vm-ci/common"
)

func TestCloneRepo(t *testing.T) {
	repo := common.Repository{
		CloneURL: "https://github.com/kdump-test/kdump-test.git",
	}
	commit := &common.Commit{
		CommitID: "bcded861ddf577643da282c4a186104356cef6e0",
	}

	path, err := cloneRepo(repo, commit)
	defer os.RemoveAll(path)

	gitRepo, err := git.OpenRepository(path)
	if err != nil {
		t.Error("Can not open Repository, check failed")
	}

	head, _ := gitRepo.Head()
	if head.Target().String() != commit.CommitID {
		t.Error("Git reset error, Head not suppose commit.")
	}

}
