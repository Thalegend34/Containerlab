package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
)

type GoGit struct {
	gitRepo GitRepo
}

// make sure GoGit satisfies the Git interface
var _ Git = (*GoGit)(nil)

func NewGoGit(gitRepo GitRepo) *GoGit {
	return &GoGit{
		gitRepo: gitRepo,
	}
}

// Clone takes the given GitRepo reference and clones the repo
// with its internal implementation.
func (g *GoGit) Clone() error {
	// if the directory is not present
	if s, err := os.Stat(g.gitRepo.GetRepoName()); os.IsNotExist(err) {
		return g.cloneNonExisting()
	} else if s.IsDir() {
		return g.cloneExisting()
	}
	return fmt.Errorf("error %q exists already but is a file", g.gitRepo.GetRepoName())
}

func (g *GoGit) cloneExisting() error {
	log.Debugf("loading git repository %q", g.gitRepo.GetRepoName())
	// load the git repository
	r, err := gogit.PlainOpen(g.gitRepo.GetRepoName())
	if err != nil {
		return err
	}
	// get the worktree reference
	tree, err := r.Worktree()
	if err != nil {
		return err
	}

	if g.gitRepo.GetBranch() != "" {
		log.Debugf("checking out branch %q", g.gitRepo.GetBranch())
		// prepare the checkout options
		checkoutOpts := &gogit.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(g.gitRepo.GetBranch()),
			Keep:   true,
			Force:  true,
		}
		// execute the checkout
		err = tree.Checkout(checkoutOpts)
		if err != nil {
			return err
		}
	}

	log.Debug("pulling latest repo data")
	// init the pull options
	pullOpts := &gogit.PullOptions{
		Depth:        1,
		SingleBranch: true,
		Force:        true,
	}
	// execute the pull
	err = tree.Pull(pullOpts)
	if err == gogit.NoErrAlreadyUpToDate {
		log.Debugf("git repository up to date")
		err = nil
	}

	return err
}

func (g *GoGit) cloneNonExisting() error {
	// init clone options
	co := &gogit.CloneOptions{
		Depth:        1,
		URL:          g.gitRepo.GetRepoUrl().String(),
		SingleBranch: true,
	}
	// set brach reference if set
	if g.gitRepo.GetBranch() != "" {
		co.ReferenceName = plumbing.NewBranchReferenceName(g.gitRepo.GetBranch())
	} else {
		co.ReferenceName = plumbing.HEAD
	}
	// perform clone
	_, err := gogit.PlainClone(g.gitRepo.GetRepoName(), false, co)
	return err
}

type ExecGit struct {
	gitRepo GitRepo
}

// make sure ExecGit satisfies the Git interface
var _ Git = (*ExecGit)(nil)

func NewExecGit(gitRepo GitRepo) *ExecGit {
	return &ExecGit{
		gitRepo: gitRepo,
	}
}

// Clone takes the given GitRepo reference and clones the repo
// with its internal implementation.
func (g *ExecGit) Clone() error {
	// build the URL with owner and repo name
	repoUrl := g.gitRepo.GetRepoUrl().String()
	cloneArgs := []string{"clone", repoUrl, "--depth", "1"}
	if g.gitRepo.GetBranch() != "" {
		cloneArgs = append(cloneArgs, []string{"--branch", g.gitRepo.GetBranch()}...)
	}

	cmd := exec.Command("git", cloneArgs...)

	log.Infof("cloning %q", repoUrl)

	cmd.Stdout = log.New().Writer()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Errorf("failed to clone %q: %v", repoUrl, err)
		log.Error(stderr.String())
		return err
	}

	return nil
}

type Git interface {
	// Clone takes the given GitRepo reference and clones the repo
	// with its internal implementation.
	Clone() error
}
