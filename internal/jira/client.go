package jira

import (
    "fmt"
    "regexp"

    "github.com/andygrunwald/go-jira"
    "github.com/daylight55/jira-changelog-generator/internal/config"
)

// Client はJiraのクライアントをラップします
type Client struct {
    client      *jira.Client
    jiraPattern *regexp.Regexp
}

// New は新しいJiraクライアントを作成します
func New(cfg *config.Config) (*Client, error) {
    tp := jira.BasicAuthTransport{
        Username: cfg.Jira.Username,
        Password: cfg.Jira.Token,
    }

    client, err := jira.NewClient(tp.Client(), cfg.GetEffectiveJiraURL())
    if err != nil {
        return nil, fmt.Errorf("failed to create JIRA client: %w", err)
    }

    jiraPattern := regexp.MustCompile(`[A-Z]+-\d+`)

    return &Client{
        client:      client,
        jiraPattern: jiraPattern,
    }, nil
}

// ExtractJiraID は文字列からJiraのチケットIDを抽出します
func (c *Client) ExtractJiraID(text string) string {
    return c.jiraPattern.FindString(text)
}

// GetIssue はJiraのチケット情報を取得します
func (c *Client) GetIssue(issueID string) (*jira.Issue, error) {
    issue, _, err := c.client.Issue.Get(issueID, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get Jira issue %s: %w", issueID, err)
    }
    return issue, nil
}

// GetEpicName はJiraのエピック名を取得します
func (c *Client) GetEpicName(issue *jira.Issue) (string, error) {
    epicLink, ok := issue.Fields.Unknowns["customfield_10014"].(string)
    if !ok {
        return "その他", nil
    }

    epic, _, err := c.client.Issue.Get(epicLink, nil)
    if err != nil {
        return "その他", fmt.Errorf("failed to get epic information: %w", err)
    }

    return epic.Fields.Summary, nil
}

// GetIssueTitle はJiraのチケットタイトルを取得します
func (c *Client) GetIssueTitle(issue *jira.Issue) string {
    return issue.Fields.Summary
}