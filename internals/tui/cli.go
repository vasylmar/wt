package tui

import (
	"errors"
	"fmt"
	"log"
	"os"
	"wt/internals/git"
	"wt/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cli struct {
	textInput textinput.Model
	worktrees git.Worktrees
	current   int
	error     error
	inputMode bool
}

const (
	title       = "Hi, there 👋 \nThis is simple cli to create and delete git worktrees."
	instruction = "q/esc: Quit | j/down: Next worktree | k/up: Previous worktree | ctrl+a: Add new worktree | ctrl+d: Delete current worktree | ctrl+l: Lock current worktree | ctrl+u: Unlock current worktree"
	inputHint   = "esc: Quit input mode | Enter: Create worktree"
)

var (
	boxStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingTop(1).
			PaddingBottom(1)
	optionsStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("252"))
	selectedOptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("204"))
	cursorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).PaddingRight(1)
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("200")).Bold(true).PaddingBottom(1)
	instructionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingBottom(1)
	inputStyle       = lipgloss.NewStyle().PaddingBottom(1)
	errorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).PaddingBottom(1)
)

func checks() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	isMain, err := git.IsMainWorktree()
	if err != nil {
		return fmt.Errorf("can't check if this is main worktree: %w", err)
	}
	if !isMain {
		return errors.New("this is not the main worktree. please run this command in the main worktree")
	}

	isGit, err := git.IsGitRepo(dir)
	if err != nil {
		return fmt.Errorf("can't check if git repo: %w", err)
	}
	if !isGit {
		return errors.New("this is not a git repository. please run this command in a git repository")
	}

	if err := git.HandleGitIgnore(); err != nil {
		return fmt.Errorf("can't handle .gitignore file: %w", err)
	}

	if err := git.HandleWorktreeDir(); err != nil {
		return fmt.Errorf("can't handle worktree directory: %w", err)
	}

	return nil
}

func NewCli() cli {
	error := checks()

	worktrees := &git.Worktrees{}
	err := worktrees.WorktreeList()
	if err != nil {
		log.Fatal(err)
	}

	ti := textinput.New()
	ti.Placeholder = "Enter new worktree name... (btw without slashes)"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 50

	return cli{
		textInput: ti,
		worktrees: *worktrees,
		current:   0,
		error:     error,
	}
}

func (m cli) Init() tea.Cmd {
	return nil
}

func (m cli) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.inputMode {
			switch msg.Type {
			case tea.KeyEsc:
				m.textInput.Reset()
				m.inputMode = false
			case tea.KeyEnter:
				worktreeName := m.textInput.Value()
				err := git.AddWorktree(worktreeName)
				if err != nil {
					m.error = err
				}
				worktreeErr := m.worktrees.WorktreeList()
				if worktreeErr != nil {
					m.error = err
				}
				m.textInput.Reset()
				m.inputMode = false
			}
		} else {

			switch msg.Type {
			case tea.KeyCtrlA:
				m.inputMode = true
			case tea.KeyCtrlD:
				worktrees := m.worktrees.Items()
				worktreeName := utils.GetWorktreeName(worktrees[m.current].Branch)
				if worktrees[m.current].Locked {
					m.error = fmt.Errorf("can't delete locked worktree, unlock it first")
				} else {
					err := git.RemoveWorktree(worktreeName)
					if err != nil {
						m.error = err
					}
					worktreeErr := m.worktrees.WorktreeList()
					if worktreeErr != nil {
						m.error = err
					}
					m.current = 0

				}

			case tea.KeyCtrlL:
				worktrees := m.worktrees.Items()
				worktreeName := utils.GetWorktreeName(worktrees[m.current].Branch)

				if !worktrees[m.current].Locked {
					if err := git.LockWorktree(worktreeName); err != nil {
						m.error = fmt.Errorf("can't lock worktree: %w", err)
					}
					worktreeErr := m.worktrees.WorktreeList()
					if worktreeErr != nil {
						m.error = worktreeErr
					}
				}

			case tea.KeyCtrlU:
				worktrees := m.worktrees.Items()
				worktreeName := utils.GetWorktreeName(worktrees[m.current].Branch)

				if worktrees[m.current].Locked {
					if err := git.UnlockWorktree(worktreeName); err != nil {
						m.error = fmt.Errorf("can't unlock worktree: %w", err)
					}
					worktreeErr := m.worktrees.WorktreeList()
					if worktreeErr != nil {
						m.error = worktreeErr
					}
				}

			}

			switch msg.String() {
			case "q", "esc":
				return m, tea.Quit
			case "up", "k":
				if m.current > 0 {
					m.current--
				}
			case "down", "j":
				if m.current < len(m.worktrees.Items())-1 {
					m.current++
				}
			}

		}
	}

	if m.inputMode {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func RenderLockEmoji(locked bool) string {
	if locked {
		return " 🔒"
	}
	return ""
}

func (m cli) View() string {
	var output string
	cursor := cursorStyle.Render(">")
	input := inputStyle.Render(m.textInput.View())

	var errorMsg string
	if m.error != nil {
		errorMsg = errorStyle.Render("Error: " + m.error.Error())
	}

	for i, wt := range m.worktrees.Items() {
		branchName := utils.GetWorktreeName(wt.Path)

		if i == m.current {
			output += cursor +
				selectedOptionStyle.Render(branchName) +
				RenderLockEmoji(wt.Locked) + "\n"
		} else {
			output += optionsStyle.Render(branchName) +
				RenderLockEmoji(wt.Locked) + "\n"
		}
	}

	var view string

	if m.inputMode {
		if m.error != nil {
			view = titleStyle.Render(title) +
				"\n" + output +
				"\n" + input +
				"\n" + errorMsg +
				"\n" + instructionStyle.Render(inputHint)
		} else {
			view = titleStyle.Render(title) +
				"\n" + output +
				"\n" + input +
				"\n" + instructionStyle.Render(inputHint)
		}
	} else {
		if m.error != nil {
			view = titleStyle.Render(title) +
				"\n" + output +
				"\n" + errorMsg +
				"\n" + instructionStyle.Render(instruction)
		} else {
			view = titleStyle.Render(title) +
				"\n" + output +
				"\n" + instructionStyle.Render(instruction)
		}
	}

	return boxStyle.Render(view)
}
