package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
)

type GoGit struct {
	gitRepo GitRepo
	r       *gogit.Repository
}

// make sure GoGit satisfies the Git interface
var _ Git = (*GoGit)(nil)

func NewGoGit(gitRepo GitRepo) (*GoGit, error) {
	// load the git repository
	r, err := gogit.PlainOpen(gitRepo.GetRepoName())
	if err != nil {
		return nil, err
	}
	return &GoGit{
		gitRepo: gitRepo,
		r:       r,
	}, nil
}

// Clone takes the given GitRepo reference and clones the repo
// with its internal implementation.
func (g *GoGit) Clone() error {
	// if the directory is not present
	if s, err := os.Stat(g.gitRepo.GetRepoName()); os.IsNotExist(err) {
		return g.cloneNonExisting()
	} else if s.IsDir() {
		return g.cloneExistingRepo()
	}
	return fmt.Errorf("error %q exists already but is a file", g.gitRepo.GetRepoName())
}

func (g *GoGit) getDefaultBranch() (string, error) {

	rem := gogit.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{g.gitRepo.GetRepoUrl().String()},
	})

	log.Debug("Fetching tags...")

	// We can then use every Remote functions to retrieve wanted information
	refs, err := rem.List(&gogit.ListOptions{})
	if err != nil {
		return "", err
	}

	for _, ref := range refs {
		if ref.Type() == plumbing.SymbolicReference && ref.Name() == plumbing.HEAD {
			return ref.Target().Short(), nil
		}
	}

	return "", fmt.Errorf("unable to determine default branch for %q", g.gitRepo.GetRepoUrl().String())

}

func (g *GoGit) cloneExistingRepo() error {
	log.Debugf("loading git repository %q", g.gitRepo.GetRepoName())

	// get the worktree reference
	tree, err := g.r.Worktree()
	if err != nil {
		return err
	}

	branch := g.gitRepo.GetBranch()
	if g.gitRepo.GetBranch() != "" {
		log.Debugf("default branch not set. determining it")
		branch, err = g.getDefaultBranch()
		if err != nil {
			return err
		}
		log.Debugf("default branch is %q", branch)
	}

	// prepare the checkout options
	checkoutOpts := &gogit.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	}

	// check if the branch already exists locally.
	// if not fetch it and check it out.

	if _, err = g.r.Reference(plumbing.NewBranchReferenceName(branch), false); err != nil {
		err = g.fetchNonExistingBranch(branch)
		if err != nil {
			return err
		}

		ref, err := g.r.Reference(plumbing.NewRemoteReferenceName("origin", branch), true)
		if err != nil {
			return err
		}

		checkoutOpts.Hash = ref.Hash()
		checkoutOpts.Create = true
	}

	log.Debugf("checking out branch %q", branch)

	// execute the checkout
	err = tree.Checkout(checkoutOpts)
	if err != nil {
		return err
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

func (g *GoGit) fetchNonExistingBranch(branch string) error {
	// init the remote
	remote, err := g.r.Remote("origin")
	if err != nil {
		return err
	}

	localRef := plumbing.NewBranchReferenceName(branch)
	remoteRef := plumbing.NewRemoteReferenceName("origin", branch)

	// init fetch options
	fetchOpts := &gogit.FetchOptions{
		Depth:    1,
		RefSpecs: []config.RefSpec{config.RefSpec(fmt.Sprintf("%s:%s", localRef, remoteRef))},
	}

	// execute the fetch
	err = remote.Fetch(fetchOpts)
	if err == gogit.NoErrAlreadyUpToDate {
		log.Debugf("git repository up to date")
	} else if err != nil {
		return err
	}
	return nil
}

func (g *GoGit) cloneNonExisting() error {
	var err error
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
		branchName, err := g.getDefaultBranch()
		if err != nil {
			return err
		}
		co.ReferenceName = plumbing.NewBranchReferenceName(branchName)
	}

	// perform clone
	_, err = gogit.PlainClone(g.gitRepo.GetRepoName(), false, co)
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
