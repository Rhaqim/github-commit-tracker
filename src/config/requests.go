package config

import (
	ut "savannahtech/src/utils"
)

// Github
var (
	GithubRepoURL = ut.Env("GITHUB_REPO_URL", "https://api.github.com/repos/")
	GitHubToken   = ut.Env("GITHUB_TOKEN", "")
)

// Default
var (
	DefaultOwner = ut.Env("DEFAULT_OWNER", "chromium")
	DefaultRepo  = ut.Env("DEFAULT_REPO", "chromium")
)
