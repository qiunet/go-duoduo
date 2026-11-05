# AGENTS.md

> 本文件供 Cursor / 其他 Coding Agent 阅读。人类补充内容时，请保持章节结构，便于 Agent 定位。
> 未填项用 `TODO` 标记；确定后删掉 `TODO`。

## 1. 项目概览

- **名称**：go-duoduo
- **一句话说明**：TODO（例如：通用 Go 工具库 / 游戏服务基础库）
- **语言 / 版本**：Go 1.27.2（以 `go.work` 与各模块 `go.mod` 为准）
- **仓库形态**：Go Workspace（`go.work`），多模块，非单根 `go.mod`

## 2. 目录与模块

| 路径 | 模块名 | 职责 | 备注 |
|------|--------|------|------|
| `logger/` | `logger` | 日志封装 | `GetLogger(name)` 风格 |
| `utils/` | `utils` | 通用工具模块 | 子包：`bitutil`、`commonutil`、`config`、`dateutil`、`fileutil`、`httpclient`、`jsonutil`、`mathutil`、`netutil`、`pathutil`、`strutil` |
| `docs/` | — | 说明 / 设计等文档 | 非 Go 模块；必要说明文档建在此处 |
| `scripts/` | — | 脚本存放区 | 非 Go 模块；构建、工具、运维脚本等 |
| TODO | TODO | TODO | TODO |

新增模块时：

1. 建目录并写 `go.mod`
2. 在仓库根执行：`go work use ./<module>`
3. 更新本节表格

文档与脚本：

- **必要的说明文档**：放到 `docs/` 下新建（不要散落在模块根或随意路径）
- **脚本**：放到 `scripts/`；勿把可执行脚本塞进业务模块目录（除非模块自测专用）

## 3. 构建 / 测试 / 常用命令

```bash
# 列出 workspace 模块
go list -m

# 测某个模块
go test -C ./utils ./...
go test -C ./logger ./...

# 构建某个模块
go build -C ./utils ./...
go build -C ./logger ./...

# 依赖整理（在对应模块目录）
cd utils && go mod tidy
```

- **CI**：TODO
- **本地前置**：Go SDK；可选 version-fox / vfox 管理版本
- **禁止（强制）**：
  - **禁止提交密钥 / 凭证**：包括但不限于 API Key、Token、密码、私钥、证书、`.env` 真实密钥、云厂商凭据；示例可用占位符（如 `YOUR_API_KEY`）
  - 发现仓库中已有疑似密钥：不要写入 commit；提醒用户轮换并改用环境变量 / 密钥管理

## 4. 编码约定

- **风格**：遵循 `.editorconfig`；Go 代码 `gofmt` / 官方风格
- **命名**：TODO（包名、导出符号、文件名习惯）
- **错误处理**：TODO
- **日志**：优先用 `logger.GetLogger("name")`，再 `Info` / `Error` 等
- **JSON**（`utils/jsonutil`）：
  - `ToJSONString` / `FromJSONString`：对象 ↔ JSON 字符串
  - `nil` / 空串 / `"null"` 不转换
  - `time.Time` 序列化为**秒级时间戳**；零值时间为 `null`；字段可用 `omitempty` 省略
- **子包导入**：按能力引入，例如 `utils/jsonutil`、`utils/fileutil`；不要再依赖根包 `utils` 下的 `.go` 文件（根仅保留 `go.mod`）
- **注释（强制）**：
  - **每个方法 / 函数都必须有注释**，说明用途（做什么、何时用）
  - **注释必须使用中文**（含导出 / 非导出方法、类型、重要常量等）
  - 导出符号仍遵循 Go doc 习惯（注释紧挨声明；可写成 `// Foo 用于...`）
  - 仅改实现、不改语义时，同步核对注释是否仍准确
- **测试（强制）**：
  - 每新增一个「类」级单元（Go 中指 `type` / 主要能力文件，如 `jsonutil.go`），必须有对应测试文件或测试方法
  - 约定：源文件 `foo.go` → 优先 `foo_test.go`；至少覆盖该文件导出函数 / 重要行为
  - 新增或修改可测行为时，同步新增 / 更新测试；禁止只交实现不交测试
  - 测试命名：`TestXxx` / `TestXxx_场景`；表驱动优先
- **依赖**：能用标准库则优先；新增第三方依赖需说明理由（TODO：审批约定）

## 5. Agent 工作方式

- **改动范围**：只改任务相关文件；不做顺手大重构
- **注释与测试**：新增 / 改动方法时必须写**中文**用途注释；新增类型或工具文件时必须补测试；缺一不可
- **验证**：改完相关模块后跑 `go test` / `go build`；失败先修再收尾
- **文档（强制）**：
  - **尽量少改无关文档**（含 `AGENTS.md`、README、`docs/` 下已有文档等）
  - 仅在必要时改文档：须先说明原因与改动范围，**征得用户同意后再改**
  - 确有必要的说明文档：经同意后在 **`docs/`** 新建；不要在仓库根或其他目录随意堆文档
  - 脚本放到 **`scripts/`**，不与文档、业务代码混放
  - 除非用户明确要求或上述「必要且已同意」，不主动新建 Markdown
- **提交 / 推送**：
  - 仅在用户明确要求时执行；提交前检查不含密钥
  - **push 前必须测试通过（强制）**：对 workspace 内相关模块跑全量 `go test`（当前至少 `utils`、`logger`）；全部通过才允许 `git push`
  - 任一测试或编译失败：**禁止 push**，先报告失败点并等待用户指示
- **回复语言**：简体中文（若用户另有要求则跟随）
- **不确定时**：先读现有代码与本文件；仍不清楚再问，不要臆造业务规则

## 6. 架构与设计要点

TODO（可写：分层、模块依赖方向、配置方式、并发模型等）

- 模块依赖方向：TODO（例如：业务 → utils/logger，反向禁止）
- 配置：TODO
- 数据 / 存储：TODO
- 网络 / RPC：TODO

## 7. 重要文件索引

| 文件 | 说明 |
|------|------|
| `go.work` | Workspace 成员 |
| `AGENTS.md` | Agent 指南（本文件） |
| `.editorconfig` | 编辑器格式 |
| `logger/logger.go` | 命名 Logger |
| `utils/jsonutil/` | JSON 工具 |
| `utils/fileutil/` | 文件工具 |
| `utils/config/` | YAML 配置加载（default / env / `a.b.c`） |
| `utils/strutil/` | 字符串 / ToString |
| `utils/httpclient/` | HTTP 客户端 |
| `docs/` | 说明文档目录 |
| `scripts/` | 脚本目录 |
| TODO | TODO |

## 8. 常见任务清单

Agent 接到任务时可对照：

- [ ] 读本文件相关章节 + 目标模块现有代码
- [ ] 实现 / 修改（保持风格一致）
- [ ] 每个新增 / 改动的方法已写中文用途注释
- [ ] 新增类型或源文件已有对应 `*_test.go` / 测试方法
- [ ] `go test` / `go build` 验证通过
- [ ] 无密钥 / 凭证进入变更或 commit
- [ ] 无关文档未擅自修改；必要文档改动已说明并征得同意
- [ ] 若新增模块：更新 `go.work` 与「目录与模块」表（属文档/索引变更时先征得同意）
- [ ] 用户要求时再 commit
- [ ] 用户要求 push 前：workspace 模块 `go test` 全通过后再推送

## 9. 业务 / 领域笔记

TODO（术语表、状态机、协议、与 Java 侧习惯对照等）

| 术语 | 含义 |
|------|------|
| TODO | TODO |

## 10. 已知坑 / 注意事项

- Workspace 根目录没有根模块：`go list ./...` 可能报错；请用 `-C ./<module>` 或进模块目录
- `go.work` 须声明与模块一致的 `go` 版本，否则会出现 “implicitly requires go 1.18” 类错误
- TODO

## 11. 联系与决策

- **负责人**：TODO
- **决策记录**：TODO（或链到 ADR / 文档目录）
