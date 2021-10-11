package types

type PullRequest struct {
	Id     int    `json:"id"`
	Number int    `json:"number"`
	State  string `json:"state"`
	Title  string `json:"title"`
	User   struct {
		Login string `json:"login"`
	} `json:"user"`
	Body             string `json:"body"`
	ActiveLockReason string `json:"active_lock_reason"`
	CreatedAt        string `json:"created_at"`
	ClosedAt         string `json:"closed_at"`
	MergedAt         string `json:"merged_at"`
	Assignees        []struct {
		Login string `json:"login"`
	} `json:"assignees"`
	RequestedReviewers []struct {
		Login string `json:"login"`
	} `json:"requested_reviewers"`
	Head struct {
		Ref string `json:"ref"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
	AuthorAssociation string `json:"author_association"`
}

type Issue struct {
	Id     int    `json:"id"`
	Number int    `json:"number"`
	State  string `json:"state"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	User   struct {
		Login string `json:"login"`
	} `json:"user"`
	Assignees []struct {
		Login string `json:"login"`
	} `json:"assignees"`
	ActiveLockReason string `json:"active_lock_reason"`
	Comments         int    `json:"comments"`
	PullRequest      struct {
		Url string `json:"url"`
	} `json:"pull_request"`
	ClosedAt          string `json:"closed_at"`
	CreatedAt         string `json:"created_at"`
	ClosedBy          string `json:"closed_by"`
	AuthorAssociation string `json:"author_association"`
}

type User struct {
	Login       string `json:"login"`
	Id          int    `json:"id"`
	Type        string `json:"type"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	PublicRepos int    `json:"public_repos"`
	CreatedAt   string `json:"created_at"`
}
