# Branch Protection Configuration

## 推荐设置（个人开发者友好）

### Main Branch Protection

**Settings → Branches → Add rule**

```
Branch name pattern: main
```

#### 基础设置（推荐）

| 选项 | 设置 | 理由 |
|------|------|------|
| ✅ Require pull request before merging | 0 approvals | 允许自己批准 |
| ❌ Require approvals from code owners | 不勾选 | 个人项目不需要 |
| ✅ Require status checks to pass | 空列表 | 允许直接推送 |
| ❌ Require conversation resolution | 不勾选 | 简化流程 |
| ❌ Require linear history | 不勾选 | 允许 merge commit |
| ❌ Require signed commits | 不勾选 | 个人项目不需要 |
| ❌ Require branches to be up to date | 不勾选 | 简化流程 |

#### Bypass 限制

```
Allow specified actors to bypass:
  ✅ sizhen (你的用户名)
```

这样你可以直接推送 main 分支，其他人需要 PR。

---

## GitHub Ruleset 设置（新功能）

**Settings → Rules → Rulesets → New ruleset → New branch ruleset**

### 步骤 1：设置基本信息

```
Ruleset Name: main-protection
Enforcement status: Active
```

### 步骤 2：添加 Bypass（你自己）

点击 **"+ Add bypass"** 按钮：
- Role/Actor: 选择 **User**
- 搜索并选择: **sizhen**（你的用户名）
- Bypass mode: **Always**

### 步骤 3：添加 Target

向下滚动找到 **"Target branches"** 部分：
- 点击 **"Add target"**
- 选择 **"Include by pattern"**
- 输入: `refs/heads/main`
- 点击 **"Add"**

### 步骤 4：添加 Rules

点击 **"Add rule"** 按钮，添加以下规则：

1. **Pull request**
   - Required approvals: `0`
   - 其他保持默认

2. **Block force pushes**
   - 勾选启用

3. **Block deletions**
   - 勾选启用

### 步骤 5：保存

点击页面底部的 **"Create"** 或 **"Save changes"** 按钮

---

## 快速检查清单

- [ ] 创建 Branch Protection Rule
- [ ] 设置 `Required approvals: 0`
- [ ] 添加自己到 Bypass 列表
- [ ] 测试：直接 push main 分支
- [ ] 测试：创建 PR 不需要审批

---

## 注意事项

1. **0 approvals** = 自己可以批准自己的 PR
2. **Bypass** = 完全绕过所有规则
3. **个人项目**建议保持简单
4. 如果以后团队协作，再增加限制
