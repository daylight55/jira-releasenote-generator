package vcs

import (
    "time"
)

type ChangeRequest struct {
    Title     string    // 変更リクエストのタイトル
    ID        int      // プラットフォーム上での識別子
    MergedAt  time.Time // マージされた日時
    URL       string    // 変更リクエストへのURL
    Author    string    // 作成者
    Branch    string    // 変更元のブランチ名
    JiraIDs   []string  // 関連するJiraのissue番号のリスト
}

type Client interface {
    // GetChangeRequests は指定された期間内の変更リクエストを取得します
    GetChangeRequests(fromDate, toDate time.Time) ([]ChangeRequest, error)

    // GetTagDate はタグの作成日時を取得します
    GetTagDate(tagName string) (time.Time, error)
}

type ClientType string

const (
    GitHub ClientType = "github"
    GitLab ClientType = "gitlab"
)