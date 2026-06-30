//go:build windows

package diff

import (
	"os/exec"
	"syscall"
)

// 创建隐藏窗口的exec.Command（Windows版本）
func newGitCommand(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
