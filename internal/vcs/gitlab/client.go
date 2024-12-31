package gitlab

import (
    "fmt"
    "regexp"
    "strconv"
    "time"

    "gitlab.com/gitlab-org/go-gitlab"
    "github.com/daylight55/jira-releasenote-generator/internal/vcs"
)

// Client はGitLabのクライアントをラップします
type Client struct {
    client      *gitlab.GitLab
    projectID   int
    jiraPattern *regexp.Regexp
}

// New は新しいGitLabクライアントを作成します
func New(cfg *vcs.Config) (vcs.Client, error) {
    opts := []gitlab.ClientOptionFunc{
        gitlab.WithOAuthClient(cfg.Token),
    }
    
    if !cfg.IsCloud {
        opts = append(opts, gitlab.WithBaseURL(cfg.ServerURL))
    }

    client, err := gitlab.NewGitLab(opts...)
    if err != nil {
        return nil, fmt.Errorf("failed to create GitLab client: %w", err)
    }

    projectID, err := strconv.Atoi(cfg.Repository)
    if err != nil {
        return nil, fmt.Errorf("invalid project ID: %w", err)
    }

    jiraPattern := regexp.MustCompile(`[A-Z]+-\d+`)

    return &Client{
        client:      client,
        projectID:   projectID,
        jiraPattern: jiraPattern,
    }, nil
}

// GetChangeRequests は指定された期間内の変更リクエストを取得します
func (c *Client) GetChangeRequests(fromDate, toDate time.Time) ([]vcs.ChangeRequest, error) {
    opts := &gitlab.ListProjectMergeRequestsOptions{
        State:   gitlab.String("merged"),
        OrderBy: gitlab.String("updated_at"),
        ListOptions: gitlab.ListOptions{
            PerPage: 100,
        },
    }

    var allMRs []vcs.ChangeRequest
    page := 1

    for {
        opts.Page = page
        mrs, resp, err := c.client.MergeRequests.ListProjectMergeRequests(c.projectID, opts)
        if err != nil {
            return nil, fmt.Errorf("failed to get merge requests: %w", err)
        }

        for _, mr := range mrs {
            if mr.MergedAt != nil && mr.MergedAt.After(fromDate) && mr.MergedAt.Before(toDate) {
                jiraIDs, err := c.getJiraIDsFromMR(mr.IID)
                if err != nil {
                    fmt.Printf("Warning: failed to get related issues for MR !%d: %v\n", mr.IID, err)
                    continue
                }

                if len(jiraIDs) == 0 {
                    if jiraID := c.jiraPattern.FindString(mr.Title); jiraID != "" {
                        jiraIDs = []string{jiraID}
                    }
                }

                if len(jiraIDs) == 0 {
                    continue
                }

                for _, jiraID := range jiraIDs {
                    allMRs = append(allMRs, vcs.ChangeRequest{
                        Title:    mr.Title,
                        ID:       mr.IID,
                        MergedAt: *mr.MergedAt,
                        URL:      mr.WebURL,
                        Author:   mr.Author.Username,
                        Branch:   mr.SourceBranch,
                        JiraIDs:  jiraIDs,
                    })
                }
            }
        }

        if page >= resp.TotalPages {
            break
        }
        page++
    }

    return allMRs, nil
}

// getJiraIDsFromMR はマージリクエストに関連するイシューからJiraのissue番号を抽出します
func (c *Client) getJiraIDsFromMR(mrID int) ([]string, error) {
    opts := &gitlab.ListMergeRequestRelatedIssuesOptions{
        ListOptions: gitlab.ListOptions{
            PerPage: 100,
        },
    }

    var issues []*gitlab.Issue
    page := 1

    for {
        opts.Page = page
        pageIssues, resp, err := c.client.MergeRequests.ListMergeRequestRelatedIssues(c.projectID, mrID, opts)
        if err != nil {
            return nil, fmt.Errorf("failed to get related issues: %w", err)
        }

        issues = append(issues, pageIssues...)

        if page >= resp.TotalPages {
            break
        }
        page++
    }

    var jiraIDs []string
    seen := make(map[string]bool)

    for _, issue := range issues {
        // イシューのタイトルと説明の両方からJiraのissue番号を抽出
        matches := c.jiraPattern.FindAllString(issue.Title+" "+issue.Description, -1)
        for _, match := range matches {
            if !seen[match] {
                jiraIDs = append(jiraIDs, match)
                seen[match] = true
            }
        }
    }

    return jiraIDs, nil
}

// GetTagDate はタグの作成日時を取得します
func (c *Client) GetTagDate(tagName string) (time.Time, error) {
    tag, _, err := c.client.Tags.GetTag(c.projectID, tagName)
    if err != nil {
        return time.Time{}, fmt.Errorf("failed to get tag %s: %w", tagName, err)
    }
    return tag.Commit.CommittedDate, nil
}