# ct-cli — Code Template CLI / 代码模板命令行工具

[English](#english) | [中文](#中文)

---

## English

A CLI tool for managing code templates, scaffolding projects, and organizing code snippets.

### Installation

```bash
go install github.com/jean/codeTemplateCli@latest
```

Or build from source:

```bash
git clone https://github.com/jean/codeTemplateCli.git
cd codeTemplateCli
go build -o ct-cli .
```

### Quick Start

```bash
ct-cli init                         # Initialize ~/.codeTemplate/
ct-cli template import --git https://github.com/user/go-service-template
ct-cli new go-service-template my-project
```

### Commands

**Project scaffolding**

| Command | Description |
|---|---|
| `ct-cli new [template] [dir]` | Scaffold a new project |
| `ct-cli generate [template]` | Generate files into current directory |

Both support `--var Name=value`, `--dry-run`, and `--yes`.

**Template management**

| Command | Description |
|---|---|
| `ct-cli template list [-l lang]` | List installed templates |
| `ct-cli template info <name>` | Show template details |
| `ct-cli template import <path>\|[url]` | Import from local dir or git (`--git`) |
| `ct-cli template remove <name>` | Remove a template |
| `ct-cli template init <name>` | Scaffold a new template skeleton |
| `ct-cli template edit <name>` | Open template dir in `$EDITOR` |

**Snippet management**

| Command | Description |
|---|---|
| `ct-cli snippet create <name>` | Create a snippet |
| `ct-cli snippet list [-l lang]` | List snippets |
| `ct-cli snippet get <name>` | Display snippet content |
| `ct-cli snippet update <name>` | Edit a snippet |
| `ct-cli snippet search <query>` | Search snippets |
| `ct-cli snippet delete <name>` | Delete a snippet |

**Configuration**

| Command | Description |
|---|---|
| `ct-cli config show` | Display config |
| `ct-cli config set <key> <value>` | Set a config value |
| `ct-cli init` | Initialize data directory |

### Template Authoring

A template is a directory with `manifest.yaml` and source files under `files/`:

```yaml
apiVersion: v1
name: go-service
version: "1.0.0"
description: "Go microservice"
languages: [go]

variables:
  - name: ProjectName
    prompt: "Project name?"
    type: string
    required: true
  - name: UsePostgres
    prompt: "Include PostgreSQL?"
    type: bool
    default: "false"

files:
  - source: main.go.tmpl
    dest: "{{ .ProjectName }}/main.go"

postGenerate: "cd {{ .ProjectName }} && go mod tidy"
```

**Template functions** (plus all Go `text/template` built-ins):

| Function | Example (`"helloWorld"`) |
|---|---|
| `upper` | `HELLOWORLD` |
| `lower` | `helloworld` |
| `title` | `Helloworld` |
| `snake` | `hello_world` |
| `camel` | `helloWorld` |
| `pascal` | `HelloWorld` |
| `kebab` | `hello-world` |
| `now "layout"` | current time |

**Variable types:** `string` (default), `bool`, `choice`

### Configuration

Data stored under `~/.codeTemplate/` (override with `CT_HOME`):

| Key | Description | Default |
|---|---|---|
| `defaults.template` | Default template | — |
| `defaults.author` | Default author | — |
| `defaults.license` | Default license | `MIT` |
| `defaults.interactive` | Interactive prompts | `true` |
| `editor` | Editor command | `$EDITOR` |

**Global flags:** `-c, --config`, `-v, --verbose`, `-y, --yes`

---

## 中文

代码模板与脚手架命令行工具，支持模板管理、项目生成、代码片段收藏。

### 安装

```bash
go install github.com/jean/codeTemplateCli@latest
```

或从源码编译：

```bash
git clone https://github.com/jean/codeTemplateCli.git
cd codeTemplateCli
go build -o ct-cli .
```

### 快速开始

```bash
ct-cli init                         # 初始化 ~/.codeTemplate/ 目录
ct-cli template import --git https://github.com/user/go-service-template
ct-cli new go-service-template my-project
```

### 命令参考

**项目生成**

| 命令 | 说明 |
|---|---|
| `ct-cli new [模板] [目录]` | 从模板生成新项目 |
| `ct-cli generate [模板]` | 将模板文件生成到当前目录 |

均支持 `--var Name=value` 传递变量、`--dry-run` 预览、`--yes` 跳过交互。

**模板管理**

| 命令 | 说明 |
|---|---|
| `ct-cli template list [-l 语言]` | 列出已安装的模板 |
| `ct-cli template info <名称>` | 查看模板详情与变量 |
| `ct-cli template import <路径>\|[URL]` | 从本地目录或 git 仓库导入（`--git`） |
| `ct-cli template remove <名称>` | 删除模板 |
| `ct-cli template init <名称>` | 创建一个新的模板骨架 |
| `ct-cli template edit <名称>` | 在编辑器中打开模板目录 |

**代码片段**

| 命令 | 说明 |
|---|---|
| `ct-cli snippet create <名称>` | 创建片段（交互式或通过参数） |
| `ct-cli snippet list [-l 语言]` | 列出片段 |
| `ct-cli snippet get <名称>` | 查看片段内容 |
| `ct-cli snippet update <名称>` | 编辑片段 |
| `ct-cli snippet search <关键词>` | 按名称/描述/标签搜索 |
| `ct-cli snippet delete <名称>` | 删除片段 |

**配置**

| 命令 | 说明 |
|---|---|
| `ct-cli config show` | 查看所有配置 |
| `ct-cli config set <键> <值>` | 设置并保存配置项 |
| `ct-cli init` | 初始化数据目录 |

**全局标志：** `-c, --config`、`-v, --verbose`、`-y, --yes`（跳过所有确认）

### 编写模板

模板是一个包含 `manifest.yaml` 和 `files/` 源文件目录的文件夹：

```
my-template/
├── manifest.yaml
└── files/
    └── main.go.tmpl
```

**manifest.yaml 示例：**

```yaml
apiVersion: v1
name: go-service
version: "1.0.0"
description: "Go 微服务模板"
languages: [go]

variables:
  - name: ProjectName
    prompt: "项目名称？"
    type: string
    required: true
  - name: UsePostgres
    prompt: "是否需要 PostgreSQL？"
    type: bool
    default: "false"
  - name: Framework
    prompt: "选择 HTTP 框架？"
    type: choice
    options: [chi, gin, echo]
    default: chi

files:
  - source: main.go.tmpl
    dest: "{{ .ProjectName }}/main.go"

postGenerate: "cd {{ .ProjectName }} && go mod tidy"
```

**模板函数**（除 Go `text/template` 内置函数外还可使用）：

| 函数 | 示例（输入 `"helloWorld"`） |
|---|---|
| `upper` | `HELLOWORLD` |
| `lower` | `helloworld` |
| `title` | `Helloworld` |
| `snake` | `hello_world` |
| `camel` | `helloWorld` |
| `pascal` | `HelloWorld` |
| `kebab` | `hello-world` |
| `now "格式"` | 当前时间 |

**变量类型：** `string`（默认）、`bool`（是否）、`choice`（选择）

- `.tmpl` 文件会经过模板引擎渲染；其他文件原样复制。
- 目标路径和文件名也支持模板变量，如 `dest: "{{ .ProjectName }}/src/main.go"`。

### 配置

数据存储在 `~/.codeTemplate/`（可通过环境变量 `CT_HOME` 自定义）：

```
~/.codeTemplate/
├── config.yaml
├── templates/    # 导入的模板
└── snippets/     # 保存的片段
```

常用配置项（`ct-cli config set <键> <值>`）：

| 键 | 说明 | 默认值 |
|---|---|---|
| `defaults.template` | 默认模板 | — |
| `defaults.author` | 默认作者 | — |
| `defaults.license` | 默认开源协议 | `MIT` |
| `defaults.interactive` | 是否开启交互模式 | `true` |
| `editor` | 编辑器命令 | `$EDITOR` |

### 从模板生成项目（完整示例）

```bash
# 1. 创建自己的模板
ct-cli template init my-go-template --lang go

# 2. 编辑 manifest.yaml，放入模板文件到 files/ 目录

# 3. 导入模板
ct-cli template import ./my-go-template

# 4. 生成项目
ct-cli new my-go-template hello-app

# 5. 非交互模式
ct-cli new my-go-template hello-app -y --var ProjectName=hello-app
```

## License

MIT
