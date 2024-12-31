# Jira Changelog Generator

Jira Changelog Generatorは、GitHubやGitLabのマージリクエストとJiraのチケット情報を組み合わせて、構造化されたチェンジログを生成するツールです。

## 特徴

- GitHubとGitLabの両方をサポート（クラウド版・オンプレミス版）
- Jiraとの統合（クラウド版・オンプレミス版）
- GitLabのrelated issuesからJira課題番号を自動抽出
- エピックベースでの変更のグループ化
- YAMLベースの設定ファイル
- 環境変数のサポート
- 柔軟なコマンドラインインターフェース

## インストール

### Go Install

```bash
go install github.com/daylight55/jira-releasenote-generator@latest
```

### ソースからビルド

```bash
# リポジトリのクローン
git clone https://github.com/daylight55/jira-releasenote-generator.git
cd jira-releasenote-generator

# Task CLIのインストール（まだの場合）
go install github.com/go-task/task/v3/cmd/task@latest

# ビルド
task build
```

## 設定

設定は以下の優先順位で読み込まれます：

1. コマンドラインフラグ
2. 環境変数
3. 設定ファイル（デフォルト: ~/.jira-releasenote-generator.yaml）

### 設定ファイルの例

```yaml
# VCS設定
vcs:
  # 使用するVCSタイプ: "github" または "gitlab"
  type: gitlab
  # APIトークン
  token: your-vcs-token
  # クラウド版を使用するかどうか
  is-cloud: true
  # オンプレミス版の場合のサーバーURL
  server-url: ""
  # GitHubの場合のリポジトリオーナー
  owner: your-org
  # GitHubの場合はリポジトリ名、GitLabの場合はプロジェクトID
  repository: your-repo-or-project-id

# Jira設定
jira:
  # Jiraのユーザー名（通常はメールアドレス）
  username: your-email@example.com
  # JiraのAPIトークン
  token: your-jira-token
  # クラウド版を使用するかどうか
  is-cloud: true
  # オンプレミス版の場合のサーバーURL
  server-url: ""
```

### 環境変数

全ての設定は環境変数でも指定できます：

```bash
export CHANGELOG_VCS_TYPE="gitlab"
export CHANGELOG_VCS_TOKEN="your-token"
export CHANGELOG_VCS_OWNER="your-org"
export CHANGELOG_VCS_REPOSITORY="your-repo"
export CHANGELOG_JIRA_USERNAME="your-email@example.com"
export CHANGELOG_JIRA_TOKEN="your-token"
```

## 使用方法

### 基本的な使用方法

```bash
# 設定ファイルを使用
jira-releasenote-generator --from-tag v1.0.0 --to-tag v1.1.0

# 全ての設定をコマンドラインで指定
jira-releasenote-generator \
  --vcs-type gitlab \
  --vcs-token "your-token" \
  --vcs-repository "123" \
  --jira-username "your-email@example.com" \
  --jira-token "your-token" \
  --from-tag "v1.0.0" \
  --to-tag "v1.1.0"

# カスタム設定ファイルを使用
jira-releasenote-generator --config ./my-config.yaml
```

### GitLabでの使用

GitLabを使用する場合、マージリクエストとJiraチケットの関連付けは以下の方法で行えます：

1. マージリクエストのタイトルにJiraチケット番号を含める（例：`PROJ-123: 機能追加`）
2. マージリクエストをGitLabのイシューと関連付け、そのイシューのタイトルまたは説明にJiraチケット番号を含める

### GitHubでの使用

GitHubを使用する場合、現在はプルリクエストのタイトルからJiraチケット番号を抽出します：

```
PROJ-123: 新機能の追加
```

## 出力例

```markdown
# v1.1.0 (2024-12-31)

## ユーザー管理機能
- PROJ-123: ユーザー登録機能の追加 (johndoe) [#456](https://gitlab.com/org/repo/-/merge_requests/456)
- PROJ-124: パスワードリセット機能の実装 (janedoe) [#457](https://gitlab.com/org/repo/-/merge_requests/457)

## セキュリティ強化
- PROJ-125: 2要素認証の導入 (bobsmith) [#458](https://gitlab.com/org/repo/-/merge_requests/458)
```

## 開発

### 必要要件

- Go 1.21以上
- Task CLI

### 開発用コマンド

```bash
# テストの実行
task test

# リンターの実行
task lint

# 開発用ビルド
task dev

# リリースビルド
task release v1.0.0
```

## ライセンス

MIT License

## 貢献

1. Forkを作成
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチをPush (`git push origin feature/amazing-feature`)
5. Pull Requestを作成

## Authors

* **daylight55** - *Initial work*