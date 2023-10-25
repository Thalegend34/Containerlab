package utils

import (
	"errors"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIsGitHubURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "github.com",
			input: "github.com",
			want:  true,
		},
		{
			name:  "github.com/containers/containerlab/blob/master/README.md",
			input: "github.com/containers/containerlab/blob/master/README.md",
			want:  true,
		},
		{
			name:  "google.com/containers",
			input: "google.com/containers",
			want:  false,
		},
		{
			name:  "google.com/containers/containerlab/blob/master/README.md",
			input: "google.com/containers/containerlab/blob/master/README.md",
			want:  false,
		},
		{
			name:  "gitlab.com/containers",
			input: "gitlab.com/containers",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if output := IsGitHubURL(tt.input); output != tt.want {
				t.Errorf("Test %q failed: want %v, but got %v", tt.name, tt.want, output)
			}
		})
	}
}

func TestGitHubGitRepoParse(t *testing.T) {
	tests := []struct {
		name           string
		ghURL          string
		expectedResult *GitHubGitRepo
		expectedError  error
	}{
		{
			name:  "bare github url without trailing slash",
			ghURL: "https://github.com/srl-labs/repo-name",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					Path:           nil,
					ProjectOwner:   "srl-labs",
					RepositoryName: "repo-name",
				},
			},
			expectedError: nil,
		},
		{
			name:  "bare github url with trailing slash",
			ghURL: "https://github.com/srl-labs/repo-name/",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					Path:           nil,
					ProjectOwner:   "srl-labs",
					RepositoryName: "repo-name",
				},
			},
			expectedError: nil,
		},
		{
			name:  "bare github.dev url with trailing slash",
			ghURL: "https://github.dev/srl-labs/repo-name/",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					Path:           nil,
					ProjectOwner:   "srl-labs",
					RepositoryName: "repo-name",
				},
			},
			expectedError: nil,
		},
		{
			name:  "bare github url with .git suffix",
			ghURL: "https://github.com/srl-labs/repo-name.git",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					Path:           nil,
					ProjectOwner:   "srl-labs",
					RepositoryName: "repo-name",
				},
			},
			expectedError: nil,
		},
		{
			name:           "invalid url with just org name",
			ghURL:          "https://github.com/srl-labs/",
			expectedResult: &GitHubGitRepo{},
			expectedError:  errInvalidURL,
		},
		{
			name:           "invalid url with no owner and no org",
			ghURL:          "https://github.com/",
			expectedResult: &GitHubGitRepo{},
			expectedError:  errInvalidURL,
		},
		{
			name:  "github url with a clab file on the main branch",
			ghURL: "https://github.com/srl-labs/repo-name/blob/main/file.clab.yml",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					Path:           []string{},
					ProjectOwner:   "srl-labs",
					RepositoryName: "repo-name",
					GitBranch:      "main",
					FileName:       "file.clab.yml",
				},
			},
			expectedError: nil,
		},
		{
			name:  "github url with a yaml file on the main branch",
			ghURL: "https://github.com/srl-labs/repo-name/blob/main/file.yaml",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					ProjectOwner:   "srl-labs",
					RepositoryName: "repo-name",
					GitBranch:      "main",
					Path:           []string{},
					FileName:       "file.yaml",
				},
			},
			expectedError: nil,
		},
		{
			name:           "url with invalid file on the main branch",
			ghURL:          "https://github.com/srl-labs/repo-name/blob/main/file.foo",
			expectedResult: &GitHubGitRepo{},
			expectedError:  errInvalidURL,
		},
		{
			name:  "github url with a specified git ref and no file",
			ghURL: "https://github.com/srl-labs/repo-name/tree/some-branch",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					ProjectOwner:   "srl-labs",
					RepositoryName: "repo-name",
					GitBranch:      "some-branch",
					Path:           []string{},
				},
			},
			expectedError: nil,
		},
		{
			name:  "github url with a specified git ref and no file and trailing slash",
			ghURL: "https://github.com/srl-labs/repo-name/tree/some-branch/",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					ProjectOwner:   "srl-labs",
					RepositoryName: "repo-name",
					GitBranch:      "some-branch",
					Path:           []string{},
				},
			},
			expectedError: nil,
		},
		{
			name:  "github url with ref to file in subdir",
			ghURL: "https://github.com/srl-labs/containerlab/blob/main/lab-examples/srl01/srl01.clab.yml",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					ProjectOwner:   "srl-labs",
					RepositoryName: "containerlab",
					GitBranch:      "main",
					Path:           []string{"lab-examples", "srl01"},
					FileName:       "srl01.clab.yml",
				},
			},
			expectedError: nil,
		},
		{
			name:  "github url with ref to subdir",
			ghURL: "https://github.com/srl-labs/containerlab/tree/main/lab-examples/srl01/",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					ProjectOwner:   "srl-labs",
					RepositoryName: "containerlab",
					GitBranch:      "main",
					Path:           []string{"lab-examples", "srl01"},
				},
			},
			expectedError: nil,
		},
		{
			name:  "github url with tree ref to repo root",
			ghURL: "https://github.com/srl-labs/containerlab/tree/main",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					ProjectOwner:   "srl-labs",
					RepositoryName: "containerlab",
					GitBranch:      "main",
					Path:           []string{},
				},
			},
			expectedError: nil,
		},
		{
			name:  "github url with tree ref to file in repo root",
			ghURL: "https://github.com/srl-labs/containerlab/blob/main/mytopo.yml",
			expectedResult: &GitHubGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "github.com",
					},
					ProjectOwner:   "srl-labs",
					RepositoryName: "containerlab",
					GitBranch:      "main",
					Path:           []string{},
					FileName:       "mytopo.yml",
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo, err := ParseGitHubRepoUrl(tt.ghURL)

			if err != nil && tt.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error: %v, but got nil", tt.expectedError)
			}

			if err != nil && tt.expectedError != nil {
				if !errors.Is(err, tt.expectedError) {
					t.Fatalf("expected error: %v, but got %v", err, tt.expectedError)
				}
				// exit the test case as we don't want to compare url structs
				// since when error is available and matches the expected error
				// we don't care about the state the struct is in
				return
			}

			if diff := cmp.Diff(repo, tt.expectedResult); diff != "" {
				t.Errorf("got result: = %v, expected %v, diff:\n%s", repo, tt.expectedResult, diff)
			}
		})
	}
}

func TestParseGitLabRepoUrl(t *testing.T) {
	tests := []struct {
		name           string
		ghURL          string
		expectedResult *GitLabGitRepo
		expectedError  error
	}{
		{
			name:  "bare gitlab url without trailing slash",
			ghURL: "https://fake.gitlab.com/user/repo-name",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					Path:           nil,
					ProjectOwner:   "user",
					RepositoryName: "repo-name",
				},
			},
			expectedError: nil,
		},
		{
			name:  "bare gitlab url with trailing slash",
			ghURL: "https://fake.gitlab.com/user/repo-name/",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					Path:           nil,
					ProjectOwner:   "user",
					RepositoryName: "repo-name",
				},
			},
			expectedError: nil,
		},
		{
			name:  "bare github url with .git suffix",
			ghURL: "https://fake.gitlab.com/user/repo-name.git",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					Path:           nil,
					ProjectOwner:   "user",
					RepositoryName: "repo-name",
				},
			},
			expectedError: nil,
		},
		{
			name:           "invalid url with just org name",
			ghURL:          "https://fake.gitlab.com/user/",
			expectedResult: &GitLabGitRepo{},
			expectedError:  errInvalidURL,
		},
		{
			name:           "invalid url with no owner and no org",
			ghURL:          "https:/fake.gitlab.com/",
			expectedResult: &GitLabGitRepo{},
			expectedError:  errInvalidURL,
		},
		{
			name:  "gitlab url with a clab file on the main branch",
			ghURL: "https://fake.gitlab.com/user/repo-name/-/blob/main/file.clab.yml",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					Path:           []string{},
					ProjectOwner:   "user",
					RepositoryName: "repo-name",
					GitBranch:      "main",
					FileName:       "file.clab.yml",
				},
			},
			expectedError: nil,
		},
		{
			name:  "gitlab url with a yaml file on the main branch",
			ghURL: "https://fake.gitlab.com/user/repo-name/-/blob/main/file.yaml",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					ProjectOwner:   "user",
					RepositoryName: "repo-name",
					GitBranch:      "main",
					Path:           []string{},
					FileName:       "file.yaml",
				},
			},
			expectedError: nil,
		},
		{
			name:           "url with invalid file on the main branch",
			ghURL:          "https://fake.gitlab.com/user/repo-name/-/blob/main/file.foo",
			expectedResult: &GitLabGitRepo{},
			expectedError:  errInvalidURL,
		},
		{
			name:  "gitlab url with a specified git ref and no file",
			ghURL: "https://fake.gitlab.com/user/repo-name/-/tree/some-branch",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					ProjectOwner:   "user",
					RepositoryName: "repo-name",
					GitBranch:      "some-branch",
					Path:           []string{},
				},
			},
			expectedError: nil,
		},
		{
			name:  "gitlab url with a specified git ref and no file and trailing slash",
			ghURL: "https://fake.gitlab.com/user/repo-name/-/tree/some-branch/",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					ProjectOwner:   "user",
					RepositoryName: "repo-name",
					GitBranch:      "some-branch",
					Path:           []string{},
				},
			},
			expectedError: nil,
		},
		{
			name:  "gitlab url with ref to file in subdir",
			ghURL: "https://fake.gitlab.com/user/containerlab/-/blob/main/lab-examples/srl01/srl01.clab.yml",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					ProjectOwner:   "user",
					RepositoryName: "containerlab",
					GitBranch:      "main",
					Path:           []string{"lab-examples", "srl01"},
					FileName:       "srl01.clab.yml",
				},
			},
			expectedError: nil,
		},
		{
			name:  "gitlab url with ref to subdir",
			ghURL: "https://fake.gitlab.com/user/containerlab/-/tree/main/lab-examples/srl01/",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					ProjectOwner:   "user",
					RepositoryName: "containerlab",
					GitBranch:      "main",
					Path:           []string{"lab-examples", "srl01"},
				},
			},
			expectedError: nil,
		},
		{
			name:  "gitlab url with tree ref to repo root",
			ghURL: "https://fake.gitlab.com/user/containerlab/-/tree/main",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					ProjectOwner:   "user",
					RepositoryName: "containerlab",
					GitBranch:      "main",
					Path:           []string{},
				},
			},
			expectedError: nil,
		},
		{
			name:  "gitlab url with tree ref to file in repo root",
			ghURL: "https://fake.gitlab.com/user/containerlab/-/blob/main/mytopo.yml",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					ProjectOwner:   "user",
					RepositoryName: "containerlab",
					GitBranch:      "main",
					Path:           []string{},
					FileName:       "mytopo.yml",
				},
			},
			expectedError: nil,
		},
		{
			name:  "gitlab url with tree ref to file in repo root and query parameters",
			ghURL: "https://fake.gitlab.com/user/containerlab/-/blob/main/mytopo.yml?foo=bar",
			expectedResult: &GitLabGitRepo{
				GitRepoStruct: GitRepoStruct{
					URLBase: url.URL{
						Scheme: "https",
						Host:   "fake.gitlab.com",
					},
					ProjectOwner:   "user",
					RepositoryName: "containerlab",
					GitBranch:      "main",
					Path:           []string{},
					FileName:       "mytopo.yml",
				},
			},
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo, err := ParseGitLabRepoUrl(tt.ghURL)

			if err != nil && tt.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error: %v, but got nil", tt.expectedError)
			}

			if err != nil && tt.expectedError != nil {
				if !errors.Is(err, tt.expectedError) {
					t.Fatalf("expected error: %v, but got %v", err, tt.expectedError)
				}
				// exit the test case as we don't want to compare url structs
				// since when error is available and matches the expected error
				// we don't care about the state the struct is in
				return
			}

			if diff := cmp.Diff(repo, tt.expectedResult); diff != "" {
				t.Errorf("got result: = %v, expected %v, diff:\n%s", repo, tt.expectedResult, diff)
			}
		})
	}
}
