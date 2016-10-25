package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/libgit2/git2go"
	"github.com/qhsong/golang-vm-ci/common"
	"github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func cloneRepo(repo common.Repository, commit *common.Commit) (string, error) {
	path, err := ioutil.TempDir("/tmp", "vm-test")
	if err != nil {
		return "", err
	}
	opt := &git.CloneOptions{}
	gitRepo, err := git.Clone(repo.CloneURL, path, opt)
	defer gitRepo.Free()
	if err != nil {
		return "", err
	}
	CommitId, _ := git.NewOid(commit.CommitID)
	gitCommit, err := gitRepo.LookupCommit(CommitId)
	defer gitCommit.Free()
	if err != nil {
		return "", err
	}
	//gitTree, _ := gitCommit.Tree()
	copt := git.CheckoutOpts{}
	err = gitRepo.ResetToCommit(gitCommit, git.ResetHard, &copt)
	if err != nil {
		return "", err
	}
	return path, nil
}

func triggerPush(PushInfo common.Push, db *mgo.Database, taskChan chan []common.Task) error {
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
			taskChan <- tasks
		} else {
			//Exist in DB, exit
			return errors.New("Commit have been exists")
		}

	}
	return nil
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
				Folder:   dir.Name(),
			}
			task.Name = commit.CommitID[:5] + "-" + dir.Name()
			task.UUID = uuid.NewV4().String()
			task.DiskPath = fmt.Sprintf("/home/qhsong/qemuimg/%s.qcow2", task.Name)
			tasks = append(tasks, task)
		}
	}
	return tasks
}
