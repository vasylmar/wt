package git

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

type Worktree struct {
	Path   string
	Commit string
	Branch string
}

type Worktrees struct {
	items []Worktree
}

func (w *Worktrees) Items() []Worktree {
	return w.items
}

func (w *Worktrees) WorktreeList() error {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}

	lines := strings.Split(out.String(), "\n")
	w.items = []Worktree{}
	wt := Worktree{}

	for _, line := range lines {
		if rest, ok := strings.CutPrefix(line, "worktree"); ok {
			if wt.Path != "" {
				w.items = append(w.items, wt)
			}
			wt = Worktree{Path: strings.TrimSpace(rest)}
		} else if rest, ok := strings.CutPrefix(line, "HEAD"); ok {
			wt.Commit = strings.TrimSpace(rest)
		} else if rest, ok := strings.CutPrefix(line, "branch"); ok {
			wt.Branch = strings.TrimSpace(rest)
		}
	}

	if wt.Path != "" {
		w.items = append(w.items, wt)
	}

	return nil
}

func AddWorktree(worktreeName string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "worktree", "add", dir+"/.worktrees/"+worktreeName)
	return cmd.Run()
}

func RemoveWorktree(worktreeName string) error {
	cmd := exec.Command("git", "worktree", "remove", worktreeName)
	return cmd.Run()
}

func LockWorktree(worktreeName string) error {
	cmd := exec.Command("git", "worktree", "lock", worktreeName)
	return cmd.Run()
}

func UnlockWorktree(worktreeName string) error {
	cmd := exec.Command("git", "worktree", "unlock", worktreeName)
	return cmd.Run()
}

func PruneWorktree() error {
	cmd := exec.Command("git", "worktree", "prune")
	return cmd.Run()
}
