package main

import (
	"fmt"
	"os"

	"github.com/shanisharrma/tasker/internal/app/bootstrap"
	"github.com/shanisharrma/tasker/internal/app/worker/cron"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cron",
		Short: "Tasker Cron Job Runner",
		Long:  "Tasker Cron Job Runner - Execute scheduled jobs for the Tasker task management system",
	}

	// Build container ONCE
	container, err := bootstrap.BuildWorkerContainer()
	if err != nil {
		fmt.Println("failed to bootstrap worker:", err)
		os.Exit(1)
	}

	// Build registry and register wired jobs
	registry := cron.NewJobRegistry()
	registry.Register(container.DueDateJob)
	registry.Register(container.OverdueJob)
	registry.Register(container.WeeklyJob)
	registry.Register(container.AutoArchiveJob)

	// Runner (no infra inside)
	runner := cron.NewRunner(registry)

	// ---- list command ----
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available cron jobs",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(registry.Help())
		},
	}
	rootCmd.AddCommand(listCmd)

	// ---- dynamic job commands ----
	for _, jobName := range registry.List() {
		name := jobName // capture

		job, _ := registry.Get(name)

		jobCmd := &cobra.Command{
			Use:   job.Name(),
			Short: job.Description(),
			RunE: func(cmd *cobra.Command, args []string) error {
				return runner.Run(name)
			},
		}

		rootCmd.AddCommand(jobCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
