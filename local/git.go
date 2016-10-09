package main

import (
	"errors"
	"io/ioutil"

	"github.com/libgit2/git2go"
	"github.com/qhsong/golang-vm-ci/common"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func cloneRepo(repo common.Repository, commit *common.Commit) (string, error) {
	path, err := ioutil.TempDir("/tmp", "vm-test")
	if err != nil {
		return "", err
	}
	gitRepo, err := git.Clone(repo.CloneURL, path, nil)
	if err != nil {
		return "", err
	}
	CommitId, _ := git.NewOid(commit.CommitID)
	gitCommit, err := gitRepo.LookupCommit(CommitId)
	if err != nil {
		return "", err
	}
	gitTree, _ := gitCommit.Tree()
	err = gitRepo.CheckoutTree(gitTree, nil)
	if err != nil {
		return "", err
	}
	return path, nil
}

func triggerPush(PushInfo common.Push, db *mgo.Database) error {
	for _, commit := range PushInfo.Commits {
		n, _ := db.C(common.C_COMMIT).Find(bson.M{"commitid": commit.CommitID}).Count()
		if n == 0 {
			path, err := cloneRepo(PushInfo.Repo, commit)
			if err != nil {
				return err
			}
			db.C(common.C_COMMIT).Insert(commit)
			tasks := getTaskList(commit, path)
			for _, task := range tasks {
				err := db.C(common.C_TASK).Insert(task)
				if err != nil {
					return err
				}
			}

		} else {
			//Exist in DB, exit
			return errors.New("Commit have been exists")
		}

	}
}

func inArray(target string, array []string) bool {
	for _, tmp := range array {
		if target == tmp {
			return true
		}
	}
	return false
}

func getTaskList(commit *common.Commit, path string) []common.Task {
	dirList, _ := ioutil.ReadDir(path)
	var tasks []common.Task
	ignoreDir := []string{".git", "include"}
	for _, dir := range dirList {
		if inArray(dir.Name(), ignoreDir) {
			continue
		} else {
			task := common.Task{
				ID:       bson.NewObjectId(),
				CommitID: commit.CommitID,
				Status:   common.TaskNewStatus,
				Name:     dir.Name(),
				Folder:   dir.Name(),
			}
			tasks = append(tasks, task)
		}
	}
	return tasks
}