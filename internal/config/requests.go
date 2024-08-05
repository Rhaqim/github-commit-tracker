package config

import (
	ut "github.com/Rhaqim/savannahtech/internal/utils"
)

// Github
var (
	GithubRepoURL = ut.Env("GITHUB_REPO_URL", "https://api.github.com/repos/")
	GitHubToken   = ut.Env("GITHUB_TOKEN", "")
)

// Default
var (
	DefaultOwner     = ut.Env("DEFAULT_OWNER", "chromium")
	DefaultRepo      = ut.Env("DEFAULT_REPO", "chromium")
	DefaultStartDate = ut.Env("DEFAULT_START_DATE", "2023-01-01:00:00:00")
	RefetchInterval  = ut.Env("REFETCH_INTERVAL", "1h")
)
