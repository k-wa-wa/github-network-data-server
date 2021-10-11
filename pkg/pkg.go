package pkg

import (
	"github-network/database/pkg/db"
	"github-network/database/pkg/fetcher"
	"github-network/database/pkg/types"
)

type pullRequests struct {
	owner        string
	repo         string
	page         int
	PullRequests *[]types.PullRequest
}

type issues struct {
	owner  string
	repo   string
	page   int
	Issues *[]types.Issue
}

type user struct {
	user_login string
	user       *types.User
}

func PullRequests(owner string, repo string, page int) *pullRequests {
	return &pullRequests{
		owner: owner,
		repo:  repo,
		page:  page,
	}
}

func Issues(owner string, repo string, page int) *issues {
	return &issues{
		owner: owner,
		repo:  repo,
		page:  page,
	}
}

func User(user_login string) *user {
	return &user{
		user_login: user_login,
	}
}

func (prs *pullRequests) Fetch() error {
	pullRequests, err := fetcher.FetchPullRequests(
		prs.owner, prs.repo, prs.page,
	)
	if err != nil {
		return err
	}
	prs.PullRequests = pullRequests
	return nil
}

func (iss *issues) Fetch() error {
	issues, err := fetcher.FetchIssues(
		iss.owner, iss.repo, iss.page,
	)
	if err != nil {
		return err
	}
	iss.Issues = issues
	return nil
}

func (user *user) Fetch() error {
	userData, err := fetcher.FetchUser(
		user.user_login,
	)
	if err != nil {
		return err
	}
	user.user = userData
	return nil
}

func (prs *pullRequests) InsertIntoDb() error {
	query := `
		INSERT INTO pull_requests (
			id, owner, repo, number, state,
			title, user_login, body, active_lock_reason,
			created_at, closed_at, merged_at,
			assignees_login_csv, requested_reviewers_login_csv,
			head_ref, base_ref, author_association
		)
		VALUES (
			?, ?, ?, ?, ?,
			?, ?, ?, ?, 
			?, ?, ?,
			?, ?,
			?, ?, ?
		)
		ON DUPLICATE KEY UPDATE
			id = VALUES(id),
			owner = VALUES(owner),
			repo = VALUES(repo),
			number = VALUES(number),
			state = VALUES(state),
			title = VALUES(title),
			user_login = VALUES(user_login),
			body = VALUES(body),
			active_lock_reason = VALUES(active_lock_reason),
			created_at = VALUES(created_at),
			closed_at = VALUES(closed_at),
			merged_at = VALUES(merged_at),
			assignees_login_csv = VALUES(assignees_login_csv),
			requested_reviewers_login_csv = VALUES(requested_reviewers_login_csv),
			head_ref = VALUES(head_ref),
			base_ref = VALUES(base_ref),
			author_association = VALUES(author_association);
	`
	owner := prs.owner
	repo := prs.repo
	for _, pr := range *prs.PullRequests {
		var assigneesLoginCsv = ""
		for _, assignee := range pr.Assignees {
			assigneesLoginCsv += assignee.Login + ","
		}
		var requestedReviewersLoginCsv = ""
		for _, requestedReviewer := range pr.RequestedReviewers {
			requestedReviewersLoginCsv += requestedReviewer.Login + ","
		}
		err := db.ExecuteQuery(query,
			pr.Id, owner, repo, pr.Number, pr.State,
			pr.Title, pr.User.Login, pr.Body, pr.ActiveLockReason,
			pr.CreatedAt, pr.ClosedAt, pr.MergedAt,
			assigneesLoginCsv, requestedReviewersLoginCsv,
			pr.Head.Ref, pr.Base.Ref, pr.AuthorAssociation,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (iss *issues) InsertIntoDb() error {
	query := `
		INSERT INTO issues (
			id, owner, repo, number, state,
			title, body, user_login,
			assignees_login_csv, active_lock_reason,
			comments, pull_request_url,
			closed_at, created_at, closed_by,
			author_association
		)
		VALUES (
			?, ?, ?, ?, ?,
			?, ?, ?,
			?, ?,
			?, ?,
			?, ?, ?,
			?
		)
		ON DUPLICATE KEY UPDATE
			id = VALUES(id),
			owner = VALUES(owner),
			repo = VALUES(repo),
			number = VALUES(number),
			state = VALUES(state),
			title = VALUES(title),
			body = VALUES(body),
			user_login = VALUES(user_login),
			assignees_login_csv = VALUES(assignees_login_csv),
			active_lock_reason = VALUES(active_lock_reason),
			comments = VALUES(comments),
			pull_request_url = VALUES(pull_request_url),
			closed_at = VALUES(closed_at),
			created_at = VALUES(created_at),
			closed_by = VALUES(closed_by),
			author_association = VALUES(author_association);
	`
	owner := iss.owner
	repo := iss.repo
	for _, is := range *iss.Issues {
		var assigneesLoginCsv = ""
		for _, assignee := range is.Assignees {
			assigneesLoginCsv += assignee.Login + ","
		}
		err := db.ExecuteQuery(query,
			is.Id, owner, repo, is.Number, is.State,
			is.Title, is.Body, is.User.Login, assigneesLoginCsv,
			is.ActiveLockReason, is.Comments,
			is.PullRequest.Url, is.ClosedAt,
			is.CreatedAt, is.ClosedBy, is.AuthorAssociation,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (user *user) InsertIntoDb() error {
	query := `
		INSERT INTO users (
			id, login, type, company,
			location, public_repos,
			created_at
		)
		VALUES (
			?, ?, ?, ?,
			?, ?,
			?
		)
		ON DUPLICATE KEY UPDATE
			id = VALUES(id),
			login = VALUES(login),
			type = VALUES(type),
			company = VALUES(company),
			location = VALUES(location),
			public_repos = VALUES(public_repos),
			created_at = VALUES(created_at);
	`
	u := user.user
	err := db.ExecuteQuery(query,
		u.Id, u.Login, u.Type,
		u.Company, u.Location,
		u.PublicRepos, u.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
