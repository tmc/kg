package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newBackupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Create compressed, timestamped backups",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createBackup()
		},
	}
}

func createBackup() error {
	// Get the knowledge graph directory from config
	kgDir := viper.GetString("knowledge_graph_dir")
	if kgDir == "" {
		return fmt.Errorf("knowledge graph directory not set in config")
	}

	// Create a timestamped directory for the backup
	backupDir := viper.GetString("backup_dir")
	if backupDir == "" {
		backupDir = filepath.Join(kgDir, "backups")
	}
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("kg_backup_%s.zip", timestamp)
	backupPath := filepath.Join(backupDir, backupName)

	// Create a new zip file
	zipFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the knowledge graph directory and add files to the zip
	err = filepath.Walk(kgDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		relPath, err := filepath.Rel(kgDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return fmt.Errorf("failed to create file in zip: %w", err)
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(zipFile, file)
		if err != nil {
			return fmt.Errorf("failed to copy file to zip: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	fmt.Printf("Backup created: %s\n", backupPath)

	// Implement rotation of old backups
	if err := rotateBackups(backupDir); err != nil {
		return fmt.Errorf("failed to rotate backups: %w", err)
	}

	return nil
}

func rotateBackups(backupDir string) error {
	maxBackups := viper.GetInt("max_backups")
	if maxBackups <= 0 {
		maxBackups = 5 // Default to keeping 5 most recent backups
	}

	backups, err := filepath.Glob(filepath.Join(backupDir, "kg_backup_*.zip"))
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	if len(backups) <= maxBackups {
		return nil
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i] > backups[j] // Sort in descending order
	})

	for _, backup := range backups[maxBackups:] {
		if err := os.Remove(backup); err != nil {
			return fmt.Errorf("failed to remove old backup %s: %w", backup, err)
		}
		fmt.Printf("Removed old backup: %s\n", backup)
	}

	return nil
}