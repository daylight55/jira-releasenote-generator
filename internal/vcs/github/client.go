package github

import (
    "context"
    "fmt"
    "time"

    "github.com/google/go-github/v45/github"
    "golang.org/x/oauth2"
    "github.com/daylight55/jira-changelog-generator/internal/vcs"
)

type Client struct {
    client     *github.Client
    owner      string
    repository string
    ctx        context.Context
}

func New(cfg *vcs.Config) (vcs.Client, error) {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: cfg.Token},
    )
    tc := oauth2.NewClient(ctx, ts)

    var client *github.Client
    if cfg.IsCloud {
        client = github.NewClient(tc)
    } else {
        var err error
        client, err = github.NewEnterpriseClient(cfg.ServerURL, cfg.ServerURL, tc)
        if err != nil {
            return nil, fmt.Errorf("failed to create GitHub enterprise client: %w", err)
        }
    }

    return &Client{
        client:     client,
        owner:      cfg.Owner,
        repository: cfg.Repository,
        ctx:        ctx,
    }, nil
}

func (c *Client) GetChangeRequests(fromDate, toDate time.Time) ([]vcs.ChangeRequest, error) {
    opts := &github.PullRequestListOptions{
        State: "closed",
        ListOptions: github.ListOptions{
            PerPage: 100,
        },
    }

    var allPRs []vcs.ChangeRequest
    for {
        prs, resp, err := c.client.PullRequests.List(c.ctx, c.owner, c.repository, opts)
        if err != nil {
            return nil, fmt.Errorf("failed to list pull requests: %w", err)
        }

        for _, pr := range prs {
            if pr.MergedAt != nil && pr.MergedAt.After(fromDate) && pr.MergedAt.Before(toDate) {
                // ToDo: Implement GitHub issue-PR relation check
                allPRs = append(allPRs, vcs.ChangeRequest{
                    Title:    pr.GetTitle(),
                    ID:       pr.GetNumber(),
                    MergedAt: pr.GetMergedAt(),
                    URL:      pr.GetHTMLURL(),
                    Author:   pr.User.GetLogin(),
                    Branch:   pr.GetHead().GetRef(),
                })
            }
        }

        if resp.NextPage == 0 {
            break
        }
        opts.Page = resp.NextPage
    }

    return allPRs, nil
}

func (c *Client) GetTagDate(tagName string) (time.Time, error) {
    ref, _, err := c.client.Git.GetRef(c.ctx, c.owner, c.repository, "refs/tags/"+tagName)
    if err != nil {
        return time.Time{}, fmt.Errorf("failed to get tag reference: %w", err)
    }

    tag, _, err := c.client.Git.GetTag(c.ctx, c.owner, c.repository, ref.Object.GetSHA())
    if err != nil {
        commit, _, err := c.client.Git.GetCommit(c.ctx, c.owner, c.repository, ref.Object.GetSHA())
        if err != nil {
            return time.Time{}, fmt.Errorf("failed to get commit: %w", err)
        }
        return commit.Committer.GetDate(), nil
    }

    return tag.Tagger.GetDate(), nil
}