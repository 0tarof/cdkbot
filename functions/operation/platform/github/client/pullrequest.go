package client

import (
	"context"
	"fmt"
	"github.com/google/go-github/v26/github"
	"github.com/sambaiz/cdkbot/functions/operation/constant"
	"github.com/sambaiz/cdkbot/functions/operation/platform"
	"strings"
)

// GetPullRequest gets a PR
func (c *Client) GetPullRequest(ctx context.Context) (*platform.PullRequest, error) {
	pr, _, err := c.client.PullRequests.Get(ctx, c.owner, c.repo, c.number)
	if err != nil {
		return nil, err
	}
	refParts := strings.Split(pr.GetBase().GetLabel(), ":")
	labels := map[string]constant.Label{}
	for _, label := range pr.Labels {
		if lb, ok := constant.NameToLabel[label.GetName()]; ok {
			labels[lb.Name] = constant.NameToLabel[lb.Name]
		}
	}
	return &platform.PullRequest{
		BaseBranch:     refParts[len(refParts)-1],
		BaseCommitHash: pr.GetBase().GetSHA(),
		HeadCommitHash: pr.GetHead().GetSHA(),
		Labels:         labels,
	}, nil
}

// GetOpenPullRequestNumbersByLabel gets open PRs having the label
func (c *Client) GetOpenPullRequestNumbersByLabel(
	ctx context.Context,
	label constant.Label,
	excludeMySelf bool,
) ([]int, error) {
	page := 1
	prs := []*github.PullRequest{}
	for true {
		paging, _, err := c.client.PullRequests.List(ctx, c.owner, c.repo, &github.PullRequestListOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			return nil, err
		}
		if len(paging) == 0 {
			break
		}
		prs = append(prs, paging...)
		page++
		if page > maxPage {
			return nil, fmt.Errorf("Too many PRs")
		}
	}
	var ret []int
	for _, pr := range prs {
		if excludeMySelf && pr.GetNumber() == c.number {
			continue
		}
		for _, lbl := range pr.Labels {
			if lbl.GetName() == label.Name {
				ret = append(ret, pr.GetNumber())
			}
		}
	}
	return ret, nil
}

// MergePullRequest merges PR
func (c *Client) MergePullRequest(ctx context.Context, message string) error {
	_, _, err := c.client.PullRequests.Merge(ctx, c.owner, c.repo, c.number, message, nil)
	return err
}

func (c *Client) getOpenPullRequestNumbers(
	ctx context.Context,
) ([]int, error) {
	page := 1
	prs := []*github.PullRequest{}
	for true {
		paging, _, err := c.client.PullRequests.List(ctx, c.owner, c.repo, &github.PullRequestListOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			return nil, err
		}
		if len(paging) == 0 {
			break
		}
		prs = append(prs, paging...)
		page++
		if page > maxPage {
			return nil, fmt.Errorf("Too many PRs")
		}
	}
	numbers := make([]int, 0, len(prs))
	for _, pr := range prs {
		numbers = append(numbers, pr.GetNumber())
	}
	return numbers, nil
}
