package types

import "time"

type GitCommitCommitter struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type GitCommit struct {
	Url    string `json:"url"`
	Author struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Date  string `json:"date"`
	} `json:"author"`
	Committer GitCommitCommitter `json:"committer"`
	Message   string             `json:"message"`
	Tree      struct {
		Url string `json:"url"`
		Sha string `json:"sha"`
	} `json:"tree"`
	CommentCount int `json:"comment_count"`
	Verification struct {
		Verified  bool        `json:"verified"`
		Reason    string      `json:"reason"`
		Signature interface{} `json:"signature"`
		Payload   interface{} `json:"payload"`
	} `json:"verification"`
}

type Committer struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Author struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

// Github Commit
type Commit struct {
	Url         string    `json:"url"`
	Sha         string    `json:"sha"`
	NodeId      string    `json:"node_id"`
	HtmlUrl     string    `json:"html_url"`
	CommentsUrl string    `json:"comments_url"`
	Commit      GitCommit `json:"commit"`
	Author      Author    `json:"author"`
	Committer   Committer `json:"committer"`
	Parents     []struct {
		Url string `json:"url"`
		Sha string `json:"sha"`
	} `json:"parents"`
}
