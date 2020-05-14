package test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneGitRepository(t *testing.T) {
	rootDir := createTempDir(t)
	originPath := filepath.Join(rootDir, "origin")
	_, err := InitGitRepository(originPath, rootDir, "")
	assert.Nil(t, err, "cannot initialze origin Git repository")
	clonedPath := filepath.Join(rootDir, "cloned")
	_, err = CloneGitRepository(originPath, clonedPath, rootDir, "")
	assert.Nil(t, err, "cannot clone repo")
	assertIsNormalGitRepo(t, clonedPath)
}

func TestInitGitRepository(t *testing.T) {
	dir := createTempDir(t)
	repo, err := InitGitRepository(dir, dir, "")
	assert.Nil(t, err, "cannot initialize normal GitRepository")
	assertIsNormalGitRepo(t, repo.Dir)
	// ensure the Git repo works, i.e. we can commit into it
	err = repo.CreateFile("test.txt", "hello")
	assert.Nil(t, err)
	err = repo.StageFiles("test.txt")
	assert.Nil(t, err)
	err = repo.CommitStagedChanges(true)
	assert.Nil(t, err)
}

func TestNewGitRepository(t *testing.T) {
	dir := createTempDir(t)
	_ = NewGitRepository(dir, dir, NewMockingShell(dir, dir, ""))
}

func TestGitRepository_AddRemote(t *testing.T) {
	repo := createTestGitTownRepo(t)
	err := repo.AddRemote("foo", "bar")
	assert.Nil(t, err)
	remotes, err := repo.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo"}, remotes)
}

func TestGitRepository_Branches(t *testing.T) {
	repo := createTestGitTownRepo(t)
	assert.Nil(t, repo.CreateFeatureBranch("branch3"))
	assert.Nil(t, repo.CreateFeatureBranch("branch2"))
	assert.Nil(t, repo.CreateFeatureBranch("branch1"))
	branches, err := repo.Branches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"main", "branch1", "branch2", "branch3", "master"}, branches)
}

func TestGitRepository_CheckoutBranch(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateBranch("branch1", "master")
	assert.Nil(t, err)
	err = repo.CheckoutBranch("branch1")
	assert.Nil(t, err)
	currentBranch, err := repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "branch1", currentBranch)
	err = repo.CheckoutBranch("master")
	assert.Nil(t, err)
	currentBranch, err = repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", currentBranch)
}

func TestGitRepository_Commits(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateCommit(Commit{
		Branch:      "master",
		FileName:    "file1",
		FileContent: "hello",
		Message:     "first commit",
	})
	assert.Nil(t, err)
	err = repo.CreateCommit(Commit{
		Branch:      "master",
		FileName:    "file2",
		FileContent: "hello again",
		Message:     "second commit",
	})
	assert.Nil(t, err)
	commits, err := repo.Commits([]string{"FILE NAME", "FILE CONTENT"})
	assert.Nil(t, err)
	assert.Len(t, commits, 2)
	assert.Equal(t, "master", commits[0].Branch)
	assert.Equal(t, "file1", commits[0].FileName)
	assert.Equal(t, "hello", commits[0].FileContent)
	assert.Equal(t, "first commit", commits[0].Message)
	assert.Equal(t, "master", commits[1].Branch)
	assert.Equal(t, "file2", commits[1].FileName)
	assert.Equal(t, "hello again", commits[1].FileContent)
	assert.Equal(t, "second commit", commits[1].Message)
}

func TestGitRepository_Configuration(t *testing.T) {
	repo := createTestRepo(t)
	config := repo.Configuration(false)
	assert.NotNil(t, config, "first path: new config")
	config = repo.Configuration(false)
	assert.NotNil(t, config, "second path: cached config")
}

func TestGitRepo_ConnectTrackingBranch(t *testing.T) {
	// replicating the situation this is used in,
	// connecting branches of repos with the same commits in them
	origin := createTestRepo(t)
	repoDir := filepath.Join(createTempDir(t), "repo") // need a non-existing directory
	err := CopyDirectory(origin.Dir, repoDir)
	assert.Nil(t, err)
	repo := NewGitRepository(repoDir, repoDir, NewMockingShell(repoDir, repoDir, ""))
	err = repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	err = repo.Fetch()
	assert.Nil(t, err)
	err = repo.ConnectTrackingBranch("master")
	assert.Nil(t, err)
	err = repo.PushBranch("master")
	assert.Nil(t, err)
}

func TestGitRepo_CreateBranch(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateBranch("branch1", "master")
	assert.Nil(t, err)
	currentBranch, err := repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", currentBranch)
	branches, err := repo.Branches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"branch1", "master"}, branches)
}

func TestGitRepo_CreateChildFeatureBranch(t *testing.T) {
	repo := createTestGitTownRepo(t)
	err := repo.CreateFeatureBranch("f1")
	assert.Nil(t, err)
	err = repo.CreateChildFeatureBranch("f1a", "f1")
	assert.Nil(t, err)
	res, err := repo.Shell.Run("git", "town", "config")
	assert.Nil(t, err)
	has := strings.Contains(res.OutputSanitized(), "Branch Ancestry:\n  main\n    f1\n      f1a")
	assert.True(t, has)
}

func TestGitRepository_CreateCommit(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateCommit(Commit{
		Branch:      "master",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "test commit",
	})
	assert.Nil(t, err)
	commits, err := repo.Commits([]string{"FILE NAME", "FILE CONTENT"})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "hello.txt", commits[0].FileName)
	assert.Equal(t, "hello world", commits[0].FileContent)
	assert.Equal(t, "test commit", commits[0].Message)
	assert.Equal(t, "master", commits[0].Branch)
}

func TestGitRepository_CreateCommit_Author(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateCommit(Commit{
		Branch:      "master",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "test commit",
		Author:      "developer <developer@example.com>",
	})
	assert.Nil(t, err)
	commits, err := repo.Commits([]string{"FILE NAME", "FILE CONTENT"})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "hello.txt", commits[0].FileName)
	assert.Equal(t, "hello world", commits[0].FileContent)
	assert.Equal(t, "test commit", commits[0].Message)
	assert.Equal(t, "master", commits[0].Branch)
	assert.Equal(t, "developer <developer@example.com>", commits[0].Author)
}

func TestGitRepository_CreateFeatureBranch(t *testing.T) {
	repo := createTestGitTownRepo(t)
	err := repo.CreateFeatureBranch("f1")
	assert.Nil(t, err)
	assert.True(t, repo.Configuration(true).IsFeatureBranch("f1"))
}

func TestGitRepository_CreateFile(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateFile("filename", "content")
	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(filepath.Join(repo.Dir, "filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}

func TestGitRepository_CreateFile_InSubFolder(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateFile("folder/filename", "content")
	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(filepath.Join(repo.Dir, "folder/filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}

func TestGitRepository_CreatePerennialBranches(t *testing.T) {
	repo := createTestGitTownRepo(t)
	err := repo.CreatePerennialBranches("p1", "p2")
	assert.Nil(t, err)
	branches, err := repo.Branches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"main", "master", "p1", "p2"}, branches)
	config := repo.Configuration(true)
	assert.True(t, config.IsPerennialBranch("p1"))
	assert.True(t, config.IsPerennialBranch("p2"))
}

func TestGitRepository_CurrentBranch(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CheckoutBranch("master")
	assert.Nil(t, err)
	err = repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = repo.CheckoutBranch("b1")
	assert.Nil(t, err)
	branch, err := repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "b1", branch)
	err = repo.CheckoutBranch("master")
	assert.Nil(t, err)
	branch, err = repo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", branch)
}

func TestGitRepository_Fetch(t *testing.T) {
	repo := createTestRepo(t)
	origin := createTestRepo(t)
	err := repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	err = repo.Fetch()
	assert.Nil(t, err)
}

func TestGitRepository_FileContentInCommit(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateCommit(Commit{
		Branch:      "master",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "commit",
	})
	assert.Nil(t, err)
	commits, err := repo.commitsInBranch("master", []string{})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	content, err := repo.FileContentInCommit(commits[0].SHA, "hello.txt")
	assert.Nil(t, err)
	assert.Equal(t, "hello world", content)
}

func TestGitRepository_FilesInCommit(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
	err = repo.CreateFile("f2.txt", "two")
	assert.Nil(t, err)
	err = repo.StageFiles("f1.txt", "f2.txt")
	assert.Nil(t, err)
	err = repo.CommitStagedChanges(true)
	assert.Nil(t, err)
	commits, err := repo.Commits([]string{})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	fileNames, err := repo.FilesInCommit(commits[0].SHA)
	assert.Nil(t, err)
	assert.Equal(t, []string{"f1.txt", "f2.txt"}, fileNames)
}

func TestGitRepository_HasGitTownConfigNow(t *testing.T) {
	repo := createTestRepo(t)
	res, err := repo.HasGitTownConfigNow()
	assert.Nil(t, err)
	assert.False(t, res)
	err = repo.CreateFeatureBranch("foo")
	assert.Nil(t, err)
	res, err = repo.HasGitTownConfigNow()
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestGitRepository_HasFile(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
	has, err := repo.HasFile("f1.txt", "one")
	assert.Nil(t, err)
	assert.True(t, has)
	_, err = repo.HasFile("f1.txt", "zonk")
	assert.Error(t, err)
	_, err = repo.HasFile("zonk.txt", "one")
	assert.Error(t, err)
}

func TestGitRepository_HasRebaseInProgress(t *testing.T) {
	repo := createTestRepo(t)
	has, err := repo.HasRebaseInProgress()
	assert.Nil(t, err)
	assert.False(t, has)
}

func TestGitRepository_LastActiveDir(t *testing.T) {
	repo := createTestRepo(t)
	dir, err := repo.LastActiveDir()
	assert.Nil(t, err)
	assert.Equal(t, repo.Dir, dir)
}

func TestGitRepository_PushBranch(t *testing.T) {
	repo := createTestRepo(t)
	origin := createTestRepo(t)
	err := repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	err = repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	err = repo.PushBranch("b1")
	assert.Nil(t, err)
	branches, err := origin.Branches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"b1", "master"}, branches)
}

func TestGitRepository_RegisterOriginalCommit(t *testing.T) {
	repo := createTestRepo(t)
	assert.Len(t, repo.originalCommits, 0)
	repo.RegisterOriginalCommit(Commit{
		FileName: "file",
	})
	assert.Len(t, repo.originalCommits, 1)
	assert.Equal(t, "file", repo.originalCommits[0].FileName)
}

func TestGitRepository_Remotes(t *testing.T) {
	repo := createTestRepo(t)
	origin := createTestRepo(t)
	err := repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	remotes, err := repo.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{"origin"}, remotes)
}

func TestGitRepository_RemoveBranch(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateBranch("b1", "master")
	assert.Nil(t, err)
	branches, err := repo.Branches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"b1", "master"}, branches)
	err = repo.RemoveBranch("b1")
	assert.Nil(t, err)
	branches, err = repo.Branches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"master"}, branches)
}

func TestGitRepository_RemoveRemote(t *testing.T) {
	repo := createTestRepo(t)
	origin := createTestRepo(t)
	err := repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	err = repo.RemoveRemote("origin")
	assert.Nil(t, err)
	remotes, err := repo.Remotes()
	assert.Nil(t, err)
	assert.Len(t, remotes, 0)
}

func TestGitRepository_SetOffline(t *testing.T) {
	repo := createTestGitTownRepo(t)
	err := repo.SetOffline(true)
	assert.Nil(t, err)
	offline, err := repo.IsOffline()
	assert.Nil(t, err)
	assert.True(t, offline)
	err = repo.SetOffline(false)
	assert.Nil(t, err)
	offline, err = repo.IsOffline()
	assert.Nil(t, err)
	assert.False(t, offline)
}

func TestGitRepository_SetRemote(t *testing.T) {
	repo := createTestRepo(t)
	remotes, err := repo.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{}, remotes)
	origin := createTestRepo(t)
	err = repo.AddRemote("origin", origin.Dir)
	assert.Nil(t, err)
	remotes, err = repo.Remotes()
	assert.Nil(t, err)
	assert.Equal(t, []string{"origin"}, remotes)
}

func TestGitRepository_ShaForCommit(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateCommit(Commit{Branch: "master", FileName: "foo", FileContent: "bar", Message: "commit"})
	assert.Nil(t, err)
	sha, err := repo.ShaForCommit("commit")
	assert.Nil(t, err)
	assert.Len(t, sha, 40)
}

func TestGitRepository_StageFile(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
}

func TestGitRepository_Stash(t *testing.T) {
	repo := createTestRepo(t)
	stashSize, err := repo.StashSize()
	assert.Nil(t, err)
	assert.Zero(t, stashSize)
	err = repo.CreateFile("f1.txt", "hello")
	assert.Nil(t, err)
	err = repo.Stash()
	assert.Nil(t, err)
	stashSize, err = repo.StashSize()
	assert.Nil(t, err)
	assert.Equal(t, 1, stashSize)
}

func TestGitRepository_UncommittedFiles(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateFile("f1.txt", "one")
	assert.Nil(t, err)
	err = repo.CreateFile("f2.txt", "two")
	assert.Nil(t, err)
	files, err := repo.UncommittedFiles()
	assert.Nil(t, err)
	assert.Equal(t, []string{".gitconfig", "f1.txt", "f2.txt"}, files)
}

// HELPERS

// createTestGitRepo creates a fully initialized Git repo including a master branch.
func createTestRepo(t *testing.T) GitRepository {
	dir := createTempDir(t)
	repo, err := InitGitRepository(dir, dir, "")
	assert.Nil(t, err, "cannot initialize Git repow")
	err = repo.Shell.RunMany([][]string{
		{"git", "commit", "--allow-empty", "-m", "initial commit"},
	})
	assert.Nil(t, err, "cannot create initial commit: %s")
	return repo
}

func createTestGitTownRepo(t *testing.T) GitRepository {
	repo := createTestRepo(t)
	err := repo.CreateBranch("main", "master")
	assert.Nil(t, err)
	err = repo.Shell.RunMany([][]string{
		{"git", "config", "git-town.main-branch-name", "main"},
		{"git", "config", "git-town.perennial-branch-names", ""},
	})
	assert.Nil(t, err)
	return repo
}
