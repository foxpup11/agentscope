package diff

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"agentscope-desktop/internal/session"
)

// Engine Git Diff 引擎
type Engine struct {
	WorkDir string // git 仓库根目录
}

// NewEngine 创建新的 Diff 引擎
func NewEngine(workDir string) *Engine {
	return &Engine{WorkDir: workDir}
}

// DiffResult 单个文件的 diff 结果
type DiffResult struct {
	FilePath   string // 文件路径
	ChangeType session.ChangeType
	AddedLines int
	RemovedLines int
	Patch      string // unified diff 内容
}

// GetDiff 获取指定范围的 git diff
func (e *Engine) GetDiff(from, to string) ([]DiffResult, error) {
	args := []string{"diff"}
	if from != "" {
		args = append(args, from)
	}
	if to != "" {
		args = append(args, to)
	}
	args = append(args, "--numstat")

	cmd := exec.Command("git", args...)
	cmd.Dir = e.WorkDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("执行 git diff 失败: %w", err)
	}

	return e.parseNumstat(string(output))
}

// GetDiffBetweenRefs 获取两个引用之间的 diff
func (e *Engine) GetDiffBetweenRefs(fromRef, toRef string) ([]DiffResult, error) {
	args := []string{"diff", fromRef, toRef, "--numstat"}

	cmd := exec.Command("git", args...)
	cmd.Dir = e.WorkDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("执行 git diff 失败: %w", err)
	}

	return e.parseNumstat(string(output))
}

// GetUncommittedDiff 获取未提交的改动（包括已跟踪和未跟踪的文件）
func (e *Engine) GetUncommittedDiff() ([]DiffResult, error) {
	var results []DiffResult

	// 1. 获取已跟踪文件的改动 (git diff)
	cmd := exec.Command("git", "diff", "--numstat")
	cmd.Dir = e.WorkDir
	output, err := cmd.Output()
	if err == nil {
		diffs, _ := e.parseNumstat(string(output))
		results = append(results, diffs...)
	}

	// 2. 获取已暂存文件的改动 (git diff --cached)
	cmd = exec.Command("git", "diff", "--cached", "--numstat")
	cmd.Dir = e.WorkDir
	output, err = cmd.Output()
	if err == nil {
		diffs, _ := e.parseNumstat(string(output))
		results = append(results, diffs...)
	}

	// 3. 获取未跟踪的新文件 (git ls-files --others --exclude-standard)
	cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard")
	cmd.Dir = e.WorkDir
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// 检查是否已经在结果中
			exists := false
			for _, r := range results {
				if r.FilePath == line {
					exists = true
					break
				}
			}
			if !exists {
				results = append(results, DiffResult{
					FilePath:   line,
					ChangeType: session.ChangeCreated,
					AddedLines: 0,
					RemovedLines: 0,
					Patch:      "",
				})
			}
		}
	}

	return results, nil
}

// GetStagedDiff 获取已暂存的改动
func (e *Engine) GetStagedDiff() ([]DiffResult, error) {
	cmd := exec.Command("git", "diff", "--cached", "--numstat")
	cmd.Dir = e.WorkDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("执行 git diff 失败: %w", err)
	}

	return e.parseNumstat(string(output))
}

// GetFilePatch 获取单个文件的完整 patch
func (e *Engine) GetFilePatch(filePath string) (string, error) {
	cmd := exec.Command("git", "diff", "--", filePath)
	cmd.Dir = e.WorkDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取文件 patch 失败: %w", err)
	}

	return string(output), nil
}

// GetDiffWithActions 获取 diff 并关联到 Agent 的 actions
func (e *Engine) GetDiffWithActions(actions []session.Action) ([]DiffResult, error) {
	// 1. 获取所有文件的 diff
	diffs, err := e.GetUncommittedDiff()
	if err != nil {
		return nil, err
	}

	// 2. 获取已暂存的 diff
	stagedDiffs, err := e.GetStagedDiff()
	if err == nil {
		diffs = append(diffs, stagedDiffs...)
	}

	// 3. 对于每个 diff，获取完整 patch
	for i := range diffs {
		patch, err := e.GetFilePatch(diffs[i].FilePath)
		if err != nil {
			continue
		}
		diffs[i].Patch = patch
	}

	return diffs, nil
}

// GetDiffBetweenSession 获取会话前后的 diff
func (e *Engine) GetDiffBetweenSession(sess *session.Session) ([]DiffResult, error) {
	// 获取会话开始前的 HEAD
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = e.WorkDir
	headOutput, err := cmd.Output()
	if err != nil {
		// 如果没有 git 仓库，尝试获取未提交的 diff
		return e.GetUncommittedDiff()
	}
	headRef := strings.TrimSpace(string(headOutput))

	// 获取所有 reflog，找到会话开始前的 commit
	cmd = exec.Command("git", "reflog", "--format=%H %ci")
	cmd.Dir = e.WorkDir
	reflogOutput, err := cmd.Output()
	if err != nil {
		return e.GetUncommittedDiff()
	}

	// 解析 reflog，找到会话开始前的 commit
	refBeforeSession := findRefBeforeTime(string(reflogOutput), sess.StartedAt.Format(time.RFC3339))
	if refBeforeSession == "" {
		return e.GetUncommittedDiff()
	}

	return e.GetDiffBetweenRefs(refBeforeSession, headRef)
}

func findRefBeforeTime(reflog string, t string) string {
	// 简化实现：返回第一个 ref
	lines := strings.Split(reflog, "\n")
	if len(lines) > 0 {
		parts := strings.SplitN(lines[0], " ", 2)
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}

// parseNumstat 解析 git diff --numstat 的输出
func (e *Engine) parseNumstat(output string) ([]DiffResult, error) {
	var results []DiffResult

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// numstat 格式: "added\tremoved\tfilename"
		// 删除的文件: "deleted\t0\tfilename"
		// 新增的文件: "0\tdeleted\tfilename"
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}

		addedStr := parts[0]
		removedStr := parts[1]
		filePath := parts[2]

		// 跳过二进制文件
		if addedStr == "-" || removedStr == "-" {
			continue
		}

		var added, removed int
		fmt.Sscanf(addedStr, "%d", &added)
		fmt.Sscanf(removedStr, "%d", &removed)

		changeType := session.ChangeModified
		if added > 0 && removed == 0 {
			changeType = session.ChangeCreated
		} else if added == 0 && removed > 0 {
			changeType = session.ChangeDeleted
		}

		results = append(results, DiffResult{
			FilePath:     filePath,
			ChangeType:   changeType,
			AddedLines:   added,
			RemovedLines: removed,
		})
	}

	return results, nil
}

// FindGitRoot 查找 git 仓库根目录
func FindGitRoot(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	for {
		// 检查 .git 目录是否存在
		gitDir := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			return dir, nil
		}

		// 检查父目录
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("未找到 git 仓库")
}

// GetStatus 获取 git 状态
func (e *Engine) GetStatus() (map[string]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = e.WorkDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("执行 git status 失败: %w", err)
	}

	status := make(map[string]string)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// porcelain 格式: "XY filename"
		if len(line) < 3 {
			continue
		}

		statusCode := line[:2]
		filePath := strings.TrimSpace(line[3:])

		status[filePath] = statusCode
	}

	return status, nil
}
