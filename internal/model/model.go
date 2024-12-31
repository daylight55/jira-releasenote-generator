package model

// ChangeEntry はチェンジログの各エントリを表します
type ChangeEntry struct {
    Title     string
    ID        int
    JiraID    string
    JiraTitle string
    JiraEpic  string
    Author    string
    URL       string
}

// Changelog はリリースノートの構造を表します
type Changelog struct {
    Version    string
    FromTag    string
    ToTag      string
    Date       string
    Categories map[string][]ChangeEntry
}