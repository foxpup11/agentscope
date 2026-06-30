//go:build !windows

package diff

import (
	"os/exec"
)

// 创建隐藏窗口的exec.Command（非Windows版本）
func newGitCommand(args ...string) *exec.Cmd {
	return exec.Command("git", args...)
}
