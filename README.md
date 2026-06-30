<p align="center">
  <h1 align="center">🔍 AgentScope</h1>
  <p align="center">
    <strong>AI Agent 改了你的代码？3 秒看清全貌。</strong>
  </p>
  <p align="center">
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go" alt="Go">
    <img src="https://img.shields.io/badge/Wails-v2-5C2D91?style=flat-square" alt="Wails">
    <img src="https://img.shields.io/badge/Platform-Windows-lightgrey?style=flat-square" alt="Platform">
    <img src="https://img.shields.io/badge/Version-v0.2.1-blue?style=flat-square" alt="Version">
    <img src="https://img.shields.io/github/license/foxpup11/agentscope?style=flat-square" alt="License">
    <img src="https://img.shields.io/github/stars/foxpup11/agentscope?style=social" alt="Stars">
  </p>
</p>

---

<p align="center">
  <strong>🚀 让 AI Agent 的每一次修改，都在你的掌控之中</strong>
</p>

---

## 😩 你是否遇到过这些问题？

```
$ claude "帮我重构这个模块"

# 30 分后...

$ git diff --stat
 src/api/auth.go       |  45 +++---
 src/models/user.go    | 127 ++++++++++---------
 src/utils/helper.go   |  89 ++++++------
 src/config/config.go  |  23 +++
 src/middleware/cors.go |  56 +++++----
 tests/auth_test.go    |  34 +++--
 ...
 
# 20 个文件被修改，但我完全不知道改了什么！
```

**别慌，AgentScope 帮你看清一切。**

---

## ✨ AgentScope 是什么？

一个**轻量级桌面应用**，帮你：

| 问题 | AgentScope 的解决方案 |
|------|----------------------|
| 😵 改动太多看不过来 | 📊 **一目了然的文件列表** |
| 🤔 不知道改了什么 | 🔍 **语法高亮的 Diff 预览** |
| ⚠️ 担心误删文件 | 🚨 **自动风险评估 (Safe/Review/Danger)** |
| 📁 多项目管理混乱 | 📂 **按项目智能分组** |
| 🔔 新会话没及时发现 | 👁️ **实时监控** |

---

## 🎬 30 秒上手

### 1️⃣ 下载

👉 [**点击下载**](https://github.com/foxpup11/agentscope/releases/latest) | 仅 12MB | 无需安装

### 2️⃣ 运行

双击 `agentscope-desktop.exe`，搞定！

### 3️⃣ 查看

![preview](docs/images/preview.png)

![preview2](docs/images/preview2.png)

---

## 🛡️ 风险评估引擎

AgentScope 内置智能风险评估，自动识别危险操作：

| 风险等级 | 触发条件 | 你会看到 |
|----------|----------|----------|
| 🔴 **Danger** | 删除文件、敏感文件、危险命令 | 红色警告 |
| 🟡 **Review** | 修改依赖、多次编辑、删除代码 | 黄色提醒 |
| 🟢 **Safe** | 新增文件、小改动、文档修改 | 绿色通过 |

### 检测的危险操作

- 🗑️ **删除文件** - 任何删除操作都会被标记
- 🔑 **敏感文件** - `.env`, `secret`, `password`, `token`
- ⚠️ **危险命令** - `rm -rf`, `chmod 777`, `curl | bash`
- 📦 **配置文件** - `.git/config`, `Dockerfile`, CI 配置

---

## 🎯 核心功能

<table>
<tr>
<td width="50%">

### 📊 智能会话管理

- 按项目自动分组
- 折叠/展开控制
- 实时监控新会话
- 2 分钟自动刷新

</td>
<td width="50%">

### 🔍 Diff 预览

- 语法高亮显示
- 支持会话前后对比
- 未提交改动查看
- 单文件详细 Diff

</td>
</tr>
<tr>
<td>

### ⚙️ 灵活配置

- 深色/浅色/跟随系统
- 自定义风险规则
- 文件路径模式匹配
- 规则启用/禁用

</td>
<td>

### 📤 导出报告

- HTML 格式报告
- Markdown 格式报告
- 包含完整 Diff
- Token 使用统计

</td>
</tr>
</table>

---

## 🌍 国际化

支持**中文**和**English**一键切换。

---

## 📦 下载安装

### 方式一：直接下载（推荐）

前往 [Releases](https://github.com/foxpup11/agentscope/releases/latest) 下载最新版本。

| 文件 | 说明 |
|------|------|
| `agentscope-desktop.exe` | Windows 可执行文件，约 12MB |

### 方式二：从源码构建

```bash
# 克隆仓库
git clone https://github.com/foxpup11/agentscope.git
cd agentscope

# 安装依赖
go mod tidy

# 开发模式（热重载）
wails dev

# 构建生产版本
wails build -platform windows/amd64 -o agentscope-desktop.exe
```

---

## 🗺️ Roadmap

### ✅ 已完成

- [x] Claude Code 会话解析
- [x] 文件改动列表与风险评估
- [x] Diff 语法高亮
- [x] 中英文切换
- [x] 可拖拽布局
- [x] 实时监控
- [x] 会话导出
- [x] 深色主题
- [x] 自定义风险规则
- [x] 项目分组

### 🔜 计划中

- [ ] 支持 Codex CLI
- [ ] 支持 OpenCode
- [ ] 支持 Aider
- [ ] 会话对比分析
- [ ] 团队协作模式

---

## 🤝 贡献

欢迎贡献！无论是提交 Bug、建议新功能，还是直接贡献代码。

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

详见 [CONTRIBUTING.md](CONTRIBUTING.md)。

---

## 📄 License

MIT License - 自由使用，自由分享。

---

<p align="center">
  <strong>⭐ 如果 AgentScope 帮到了你，请给个 Star 支持一下！</strong>
</p>

<p align="center">
  <sub>Made with ❤️ by <a href="https://github.com/foxpup11">foxpup11</a></sub>
</p>
