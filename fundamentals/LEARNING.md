
# 🚀 Chapter 1: Windows Go 语言开发环境配置指南

## 1. 下载与版本选择

* **下载地址**：[https://golang.google.cn/dl/](https://golang.google.cn/dl/) (国内官网，速度快)
* 官网地址：https://go.dev/
* **选择版本**：`go1.25.5.windows-amd64.msi` (Installer 版本)
  * *注：推荐使用 MSI 安装包，它会自动处理注册表和部分基础变量，比 Zip 解压版省事。*

## 2. 目录结构规划 (核心)

我们在 E 盘创建了两个完全独立的目录，从物理上隔离“编译器”与“工作区”。

| 目录路径                            | Go 术语          | 本质与作用                                                                                                                                  |
| :---------------------------------- | :--------------- | :------------------------------------------------------------------------------------------------------------------------------------------ |
| **`E:\Environment\Go`**     | **GOROOT** | **SDK 安装目录**。`<br>`存放编译器 (`go.exe`)、标准库源码 (`fmt`, `net` 等)。`<br>`安装时指定，一般不需要手动修改环境变量。 |
| **`E:\Environment\GoDeps`** | **GOPATH** | **工作区/依赖缓存**。`<br>`存放第三方包 (`pkg/mod`) 和编译后的工具二进制文件 (`bin`)。`<br>`**必须**手动配置环境变量。  |

## 3. 安装过程

1. 运行下载的 `.msi` 安装包。
2. 在 **Destination Folder** 步骤，将默认的 `C:\Program Files\Go` 修改为：
   > `E:\Environment\Go`
   >
3. 完成安装。
   * *原理：此时安装程序已自动将 `GOROOT` 指向该目录，并将 `E:\Environment\Go\bin` 加入了系统 Path，确保你能运行 `go` 命令。*

## 4. 环境变量配置 (手动修正)

安装程序默认会将 GOPATH 设为用户目录下的 `go` 文件夹，我们需要将其指向我们规划的 E 盘目录。

### 步骤 A：修改 GOPATH (指定存储位置)

1. 打开 **环境变量** -> **用户变量**。
2. 找到或新建变量名 `GOPATH`。
3. 设置变量值为：
   > `E:\Environment\GoDeps`
   >

### 步骤 B：修改 Path (指定工具调用路径)

为了在终端直接运行 `dlv` (调试器)、`air` (热重载) 等 Go 工具，必须将 GOPATH 下的 bin 目录加入系统路径。

1. 在 **用户变量** 中找到 `Path`。
2. 编辑，把旧的 `%USERPROFILE%\go\bin`修改为：

   > `%GOPATH%\bin`
   > *(或者绝对路径 `E:\Environment\GoDeps\bin`)*
   >

### 步骤 C：配置 GOPROXY (网络加速)

Go 默认使用 Google 的模块代理，国内无法访问。需要通过命令行修改为国内镜像。

打开终端 (CMD/PowerShell)，执行：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

## 5. 验证环境

关闭所有旧终端，打开新的终端，输入 `go env`，检查以下关键项：

```bash
# 1. 检查 SDK 位置
GOROOT=E:\Environment\Go

# 2. 检查依赖存储位置 (必须是 E 盘那个)
GOPATH=E:\Environment\GoDeps

# 3. 检查代理 (必须是 goproxy.cn)
GOPROXY=https://goproxy.cn,direct
```

## 6. 常见问题排查

* **问题**：安装了工具 (如 `go install github.com/air-verse/air@latest`) 但终端提示 "command not found"。
* **原因**：**步骤 4-B** 中的 `Path` 没配对，或者没重启终端。
* **解决**：检查 `Path` 中是否包含 `E:\Environment\GoDeps\bin`。
