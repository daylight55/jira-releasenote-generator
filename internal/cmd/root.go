package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/daylight55/jira-changelog-generator/internal/config"
    "github.com/daylight55/jira-changelog-generator/internal/generator"
)

var (
    cfgFile string
    cfg     *config.Config
)

var rootCmd = &cobra.Command{
    Use:   "jira-changelog-generator",
    Short: "Generate changelog from VCS merge requests and Jira tickets",
    Long: `A changelog generator that creates formatted changelog entries by combining
information from GitHub/GitLab change requests and Jira tickets. 
It organizes the changes by Jira epics and includes relevant ticket information.`,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        if err := config.InitializeConfig(); err != nil {
            return err
        }

        var err error
        cfg, err = config.LoadConfig()
        if err != nil {
            return err
        }

        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        gen := generator.New(cfg)
        return gen.Generate()
    },
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jira-changelog-generator.yaml)")

    rootCmd.PersistentFlags().String("vcs-type", "github", "VCS type (github or gitlab)")
    rootCmd.PersistentFlags().String("vcs-token", "", "VCS API token")
    rootCmd.PersistentFlags().String("vcs-server-url", "", "VCS server URL (for on-premise)")
    rootCmd.PersistentFlags().Bool("vcs-cloud", true, "Use cloud version of VCS")
    rootCmd.PersistentFlags().String("vcs-owner", "", "Repository owner (for GitHub)")
    rootCmd.PersistentFlags().String("vcs-repository", "", "Repository name or project ID")

    rootCmd.PersistentFlags().String("jira-username", "", "Jira username or email")
    rootCmd.PersistentFlags().String("jira-token", "", "Jira API token")
    rootCmd.PersistentFlags().String("jira-server-url", "", "Jira server URL (for on-premise)")
    rootCmd.PersistentFlags().Bool("jira-cloud", true, "Use Jira cloud")

    rootCmd.PersistentFlags().String("from-tag", "", "Starting tag for changelog")
    rootCmd.PersistentFlags().String("to-tag", "", "Ending tag for changelog")

    viper.BindPFlags(rootCmd.PersistentFlags())

    rootCmd.MarkPersistentFlagRequired("vcs-token")
    rootCmd.MarkPersistentFlagRequired("jira-username")
    rootCmd.MarkPersistentFlagRequired("jira-token")
    rootCmd.MarkPersistentFlagRequired("from-tag")
    rootCmd.MarkPersistentFlagRequired("to-tag")

    rootCmd.AddCommand(newVersionCmd())
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    }
}