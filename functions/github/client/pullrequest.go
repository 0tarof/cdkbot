package client

import (
	"context"
	"errors"
)

// GetPullRequestLatestCommitHash gets latest commit hash in PR
func (c *Client) GetPullRequestLatestCommitHash(
	ctx context.Context,
	owner string,
	repo string,
	number int,
) (string, error) {
	commits, _, err := c.client.PullRequests.ListCommits(ctx, owner, repo, number, nil)
	if err != nil {
		return "", err
	}
	if len(commits) == 0 {
		return "", errors.New("PR has no commits")
	}
	return commits[len(commits)-1].GetSHA(), nil
}

// GetPullRequest gets base branch of pull request
func (c *Client) GetPullRequestBaseBranch(
	ctx context.Context,
	owner string,
	repo string,
	number int,
) (string, error) {
	pr, _, err := c.client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return "", err
	}
	return pr.GetBase().GetLabel(), nil
}
