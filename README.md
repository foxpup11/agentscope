# AgentScope Desktop

<p align="center">
  <img src="build/appicon.png" width="120" alt="AgentScope Logo">
</p>

<h3 align="center">可视化 AI Agent 对你的代码做了什么</h3>

<p align="center">
  <em>像 <code>git diff</code>，但懂 Agent</em>
</p>

<p align="center">
  <a href="#-features">功能</a> •
  <a href="#-installation">安装</a> •
  <a href="#-usage">使用</a> •
  <a href="#-development">开发</a> •
  <a href="#-contributing">贡献</a>
</p>

---

## 🎯 什么是 AgentScope？

AgentScope 是一个桌面应用，帮助开发者**可视化 AI Agent（如 Claude Code）对代码做了什么**。

当 AI Agent 修改了你的代码后，AgentScope 能够：

- 📊 **展示所有文件改动** - 哪些文件被创建、修改或删除
- 🔗 **关联 Agent 操作** - 每个改动对应哪个 tool call
- 🛡️ **风险评估** - 自动标注安全等级（🟢 Safe / 🟡 Review / 🔴 Danger）
- 🔍 **Diff 预览** - 语法高亮的代码差异查看
- 📈 **会话管理** - 浏览、搜索所有历史会话

## ✨ Features

### 核心功能

| 功能 | 描述 |
|------|------|
| **会话列表** | 左侧面板展示所有 Claude Code 会话，支持搜索筛选 |
| **文件表格** | 展示选中会话的所有文件改动，包含风险等级、变更类型、操作数 |
| **Diff 视图** | 右侧实时预览选中文件的代码差异，带语法高亮 |
| **风险评估** | 基于规则引擎自动评估每个改动的安全风险 |
| **实时监控** | 监控 Claude Code 会话目录，新会话自动出现 |

### 风险规则

| 等级 | 触发条件 |
|------|----------|
| 🔴 **Danger** | 修改敏感文件（.env、secrets）、执行危险命令（rm -rf）、大量代码改动 |
| 🟡 **Review** | 删除大量代码、多次编辑同一文件、修改依赖文件 |
| 🟢 **Safe** | 新增文件、小改动、只修改文档 |

## 📦 Installation

### 下载预编译版本

前往 [Releases](https://github.com/foxpup11/agentscope/releases) 下载最新版本：

- **Windows**: `agentscope-desktop-windows-amd64.exe`
- **macOS**: `agentscope-desktop-darwin-arm64.zip`
- **Linux**: `agentscope-desktop-linux-amd64.tar.gz`

### 从源码构建

**前置要求**:

- [Go](https://go.dev/dl/) 1.21+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

```bash
# 克隆仓库
git clone https://github.com/foxpup11/agentscope.git
cd agentscope

# 安装依赖
go mod tidy

# 开发模式
wails dev

# 构建生产版本
wails build
```

## 🚀 Usage

### 快速开始

1. 确保你已经使用过 Claude Code（会话数据存储在 `~/.claude/` 目录）
2. 运行 AgentScope Desktop
3. 在左侧选择一个会话
4. 右侧查看文件改动和 Diff

### 界面说明

```
┌─────────────────────────────────────────────────────────────┐
│  AgentScope Desktop                              ─ □ ×      │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌──────────────────────────────────────┐ │
│  │ 会话列表     │  │  文件改动                             │ │
│  │             │  │  Risk │ File    │ Change │ Ops        │ │
│  │ ▶ session-1 │  │  [OK] │ main.go │ Modified│ 2        │ │
│  │   session-2 │  │  [!!] │ .env    │ Deleted│ 1         │ │
│  │   session-3 │  │                                      │ │
│  │             │  │  Diff 视图                            │ │
│  │ 搜索: [___] │  │  +import "fmt"                       │ │
│  │             │  │  +func main() { ... }                │ │
│  └─────────────┘  └──────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│  会话: 1dd94fa6 │ 分支: main │ Token: 29K in / 8K out     │
└─────────────────────────────────────────────────────────────┘
```

### 快捷键

| 快捷键 | 功能 |
|--------|------|
| `↑` / `↓` | 在会话列表中导航 |
| `Enter` | 选择会话 |
| `Ctrl+F` | 搜索会话 |
| `F5` | 刷新列表 |

## 🛠️ Development

### 项目结构

```
agentscope-desktop/
├── main.go                    # Wails 入口
├── app.go                     # 后端 API (GetSessions/GetSession/GetDiff)
├── wails.json                 # Wails 配置
│
├── frontend/                  # 前端 (HTML + CSS + JS)
│   ├── index.html
│   ├── style.css
│   └── app.js
│
├── internal/                  # 业务逻辑
│   ├── session/
│   │   ├── session.go         # 数据模型
│   │   ├── discover.go        # 会话发现
│   │   └── claude/reader.go   # Claude Code 解析器
│   ├── diff/
│   │   ├── engine.go          # Git Diff 引擎
│   │   └── matcher.go         # 文件-动作匹配
│   └── risk/
│       └── engine.go          # 风险规则引擎
│
└── build/                     # 构建配置
```

### 技术栈

- **后端**: Go + Wails v2
- **前端**: Vanilla JS + CSS
- **Git 操作**: os/exec (调用 git CLI)
- **语法高亮**: 自定义 Diff 高亮

### 开发命令

```bash
# 启动开发模式（热重载）
wails dev

# 构建生产版本
wails build

# 运行测试
go test ./...

# 格式化代码
go fmt ./...
```

## 🤝 Contributing

欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详情。

### 如何贡献

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📝 License

本项目采用 MIT License - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 Acknowledgments

- [Wails](https://wails.io/) - Go 桌面应用框架
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI 框架（灵感来源）
- [Claude Code](https://claude.ai/) - AI 编程助手

## 📧 Contact

- **Issues**: [GitHub Issues](https://github.com/foxpup11/agentscope/issues)
- **Email**: sizhen02621@gmail.com

---

<p align="center">
  如果这个项目对你有帮助，请给个 ⭐ Star 支持一下！
</p>
