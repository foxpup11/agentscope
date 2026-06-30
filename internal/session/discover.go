package session

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DiscoverConfig 会话发现配置
type DiscoverConfig struct {
	// HomeDir 用户主目录（默认 $HOME）
	HomeDir string
	// WorkDir 当前工作目录（默认 os.Getwd）
	WorkDir string
	// SessionID 指定会话 ID（可选）
	SessionID string
}

// DiscoverResult 会话发现结果
type DiscoverResult struct {
	Path      string   // 会话文件路径
	Sessions  []string // 所有可用会话路径（当指定 SessionID 时为空）
	AgentType string   // 代理类型
}

// DiscoverSession 自动发现最近的 Agent 会话
func DiscoverSession(config DiscoverConfig) (*DiscoverResult, error) {
	if config.HomeDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("获取主目录失败: %w", err)
		}
		config.HomeDir = home
	}

	if config.WorkDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("获取工作目录失败: %w", err)
		}
		config.WorkDir = wd
	}

	// 如果指定了 SessionID，精确查找
	if config.SessionID != "" {
		return discoverByID(config)
	}

	// 按优先级发现会话
	return discoverByPriority(config)
}

func discoverByID(config DiscoverConfig) (*DiscoverResult, error) {
	claudePath := filepath.Join(config.HomeDir, ".claude", "projects")

	// 遍历所有项目目录
	projects, err := os.ReadDir(claudePath)
	if err != nil {
		return nil, fmt.Errorf("读取 Claude 项目目录失败: %w", err)
	}

	for _, project := range projects {
		if !project.IsDir() {
			continue
		}

		sessionsDir := filepath.Join(claudePath, project.Name())
		sessionFile := filepath.Join(sessionsDir, config.SessionID+".jsonl")

		if _, err := os.Stat(sessionFile); err == nil {
			return &DiscoverResult{
				Path:      sessionFile,
				AgentType: "claude-code",
			}, nil
		}
	}

	return nil, fmt.Errorf("未找到会话: %s", config.SessionID)
}

func discoverByPriority(config DiscoverConfig) (*DiscoverResult, error) {
	// 优先级 1: 当前仓库的最近会话
	if result, err := discoverInProject(config); err == nil {
		return result, nil
	}

	// 优先级 2: 所有项目中最近的会话
	if result, err := discoverLatest(config); err == nil {
		return result, nil
	}

	return nil, fmt.Errorf("未找到任何 Agent 会话")
}

func discoverInProject(config DiscoverConfig) (*DiscoverResult, error) {
	claudePath := filepath.Join(config.HomeDir, ".claude", "projects")

	// 将当前工作目录编码为 Claude 的格式
	encoded := encodeCWD(config.WorkDir)

	projectDir := filepath.Join(claudePath, encoded)
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("项目目录不存在: %s", projectDir)
	}

	return findLatestSession(projectDir)
}

func discoverLatest(config DiscoverConfig) (*DiscoverResult, error) {
	claudePath := filepath.Join(config.HomeDir, ".claude", "projects")

	projects, err := os.ReadDir(claudePath)
	if err != nil {
		return nil, fmt.Errorf("读取 Claude 项目目录失败: %w", err)
	}

	var allSessions []sessionInfo

	for _, project := range projects {
		if !project.IsDir() {
			continue
		}

		projectDir := filepath.Join(claudePath, project.Name())
		sessions, err := listSessionsInDir(projectDir)
		if err != nil {
			continue
		}

		allSessions = append(allSessions, sessions...)
	}

	if len(allSessions) == 0 {
		return nil, fmt.Errorf("未找到任何会话")
	}

	// 按修改时间排序
	sort.Slice(allSessions, func(i, j int) bool {
		return allSessions[i].ModTime.After(allSessions[j].ModTime)
	})

	latest := allSessions[0]
	return &DiscoverResult{
		Path:      latest.Path,
		AgentType: "claude-code",
	}, nil
}

type sessionInfo struct {
	Path    string
	ModTime time.Time
}

func findLatestSession(dir string) (*DiscoverResult, error) {
	sessions, err := listSessionsInDir(dir)
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("目录中没有会话: %s", dir)
	}

	// 按修改时间排序
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].ModTime.After(sessions[j].ModTime)
	})

	latest := sessions[0]
	return &DiscoverResult{
		Path:      latest.Path,
		AgentType: "claude-code",
	}, nil
}

func listSessionsInDir(dir string) ([]sessionInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var sessions []sessionInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		sessions = append(sessions, sessionInfo{
			Path:    filepath.Join(dir, entry.Name()),
			ModTime: info.ModTime(),
		})
	}

	return sessions, nil
}

// encodeCWD 将工作目录编码为 Claude Code 的格式
// 例如: /home/user/project → -home-user-project
func encodeCWD(cwd string) string {
	// 将路径分隔符替换为连字符
	encoded := strings.ReplaceAll(cwd, "/", "-")
	encoded = strings.ReplaceAll(encoded, "\\", "-")

	// 移除开头的连字符（如果有）
	encoded = strings.TrimPrefix(encoded, "-")

	return encoded
}
