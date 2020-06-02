package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

func init() {
	cacheCmd.AddCommand(cacheInitCmd)
	cacheCmd.AddCommand(cacheUpdateCmd)
	cacheCmd.AddCommand(cacheDeleteCmd)
	rootCmd.AddCommand(cacheCmd)
}

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage the offline tldr cache",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

var cacheInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize an offline tldr cache",
	Run: func(cmd *cobra.Command, args []string) {
		err := initCache()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("Successfully initialized cache")
	},
}

var cacheUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the offline tldr cache",
	Run: func(cmd *cobra.Command, args []string) {
		err := updateCache()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("Cache successfully updated")
	},
}

var cacheDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the offline tldr cache",
	Run: func(cmd *cobra.Command, args []string) {
		err := deleteCache()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("Cache successfully deleted")
	},
}

func initCache() error {
	baseDir, err := cacheBaseDir()
	if err != nil {
		return err
	}
	err = os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return err
	}
	dir, err := cacheDir()
	if err != nil {
		return err
	}
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("Cache already exists")
	}
	cmd := exec.Command("git", "clone", "https://github.com/tldr-pages/tldr.git")
	cmd.Dir = baseDir
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func updateCache() error {
	dir, err := cacheDir()
	if err != nil {
		return err
	}
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("Cache does not exist, run 'tldr cache init' to create")
	}
	cmd := exec.Command("git", "pull")
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func deleteCache() error {
	baseDir, err := cacheBaseDir()
	if err != nil {
		return err
	}
	err = os.RemoveAll(baseDir)
	if err != nil {
		return err
	}
	return nil
}

func cacheBaseDir() (string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(userDir, ".tldr", "cache"), nil
}

func cacheDir() (string, error) {
	baseDir, err := cacheBaseDir()
	if err != nil {
		return "", err
	}
	return path.Join(baseDir, "tldr"), nil
}
