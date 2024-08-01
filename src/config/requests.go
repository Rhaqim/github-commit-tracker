package config

import (
	ut "savannahtech/src/utils"
)

// Postgres
var (
	GithubRepoURL = ut.Env("GITHUB_REPO_URL", "https://api.github.com/repos/")
	GitHubToken   = ut.Env("GITHUB_TOKEN", "")
)
