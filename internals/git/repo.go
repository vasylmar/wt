package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"wt/utils"
)

func IsMainWorktree() (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	bytes, err := cmd.Output()
	if err != nil {
		return false, err
	}

	gitDir := strings.TrimSpace(string(bytes))

	if gitDir == ".git" {
		return true, nil
	}

	return false, nil
}

func IsGitRepo(path string) (bool, error) {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}

func HandleGitIgnore() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	isMainWorktree, err := IsMainWorktree()
	if err != nil {
		return err
	}

	if isMainWorktree {
		exist, err := utils.IsFileExist(filepath.Join(dir, ".gitignore"))
		if err != nil {
			return err
		}

		if exist {
			found, _, err := utils.FindInFile(filepath.Join(dir, ".gitignore"), ".worktrees")
			if err != nil {
				return err
			}
			if found {
				return nil
			} else {
				err := utils.AppendToFile(filepath.Join(dir, ".gitignore"), "\n#Ignore .worktrees directory\n.worktrees\n")
				if err != nil {
					return err
				}
			}

		} else {
			err := utils.CreateFile(filepath.Join(dir, ".gitignore"), "#Ignore .worktrees directory\n.worktrees\n")
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func HandleWorktreeDir() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	yes, error := IsMainWorktree()
	if error != nil {
		return error
	}

	if yes {
		exist, err := utils.IsDirExist(filepath.Join(dir, ".worktrees"))
		if err != nil {
			return err
		}
		if !exist {
			err := utils.MakeDir(dir, ".worktrees")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
