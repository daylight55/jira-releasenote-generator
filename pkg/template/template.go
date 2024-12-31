package template

import (
    "io"
    "sort"
    "text/template"

    "github.com/daylight55/jira-changelog-generator/internal/model"
)

const defaultTemplate = `# {{.Version}} ({{.Date}})

{{range $epic, $changes := .Categories}}## {{$epic}}
{{range $change := $changes}}- {{$change.JiraID}}: {{$change.Title}} ({{$change.Author}}) [#{{$change.ID}}]({{$change.URL}})
{{end}}
{{end}}
`

// Generator はチェンジログのテンプレート生成を担当します
type Generator struct {
    tmpl *template.Template
}

// New は新しいテンプレートジェネレーターを作成します
func New() (*Generator, error) {
    tmpl, err := template.New("changelog").Parse(defaultTemplate)
    if err != nil {
        return nil, err
    }

    return &Generator{
        tmpl: tmpl,
    }, nil
}

// Generate はチェンジログを生成します
func (g *Generator) Generate(w io.Writer, changelog *model.Changelog) error {
    // カテゴリーをソート
    categories := make([]string, 0, len(changelog.Categories))
    for category := range changelog.Categories {
        categories = append(categories, category)
    }
    sort.Strings(categories)

    // 各カテゴリー内の変更をソート
    for _, category := range categories {
        changes := changelog.Categories[category]
        sort.Slice(changes, func(i, j int) bool {
            return changes[i].JiraID < changes[j].JiraID
        })
    }

    sortedChangelog := &model.Changelog{
        Version:    changelog.Version,
        FromTag:    changelog.FromTag,
        ToTag:      changelog.ToTag,
        Date:       changelog.Date,
        Categories: make(map[string][]model.ChangeEntry),
    }

    for _, category := range categories {
        sortedChangelog.Categories[category] = changelog.Categories[category]
    }

    return g.tmpl.Execute(w, sortedChangelog)
}