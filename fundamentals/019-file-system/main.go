package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// =============================================================================
// 019. Go 文件系统：从字节到文件的完整旅程
// =============================================================================
// 场景：日志分析系统
// 读取服务器日志文件，分析错误信息，生成分析报告。
// =============================================================================

// LogEntry 表示一条日志记录
type LogEntry struct {
	Time    string // 时间戳
	Level   string // 日志级别：INFO, ERROR, WARN
	Message string // 日志内容
}

// LogAnalyzer 日志分析器
type LogAnalyzer struct {
	TotalLines int        // 总行数
	InfoCount  int        // INFO 级别数量
	WarnCount  int        // WARN 级别数量
	ErrorCount int        // ERROR 级别数量
	Errors     []LogEntry // 所有 ERROR 级别的日志
}

// =============================================================================
// 第一部分：字节的基本操作
// =============================================================================

// demoByteBasics 演示字节和字符串的转换
func demoByteBasics() {
	fmt.Println("========== 字节基础操作 ==========")

	// byte 是 uint8 的别名，表示一个字节（0-255）
	var b byte = 72 // 字母 'H' 的 ASCII 码
	fmt.Printf("byte 变量 b = %d, 对应字符 = %c\n", b, b)

	// 字符串转换为字节切片
	text := "Hello"
	data := []byte(text)
	fmt.Printf("字符串 \"%s\" 转换为字节切片: %v\n", text, data)

	// 字节切片转换为字符串
	restored := string(data)
	fmt.Printf("字节切片 %v 转换为字符串: \"%s\"\n", data, restored)

	// 查看每个字节对应的字符
	fmt.Println("逐字节分析:")
	for i, b := range data {
		fmt.Printf("  位置 %d: 字节值 %d, 字符 '%c'\n", i, b, b)
	}

	fmt.Println()
}

// =============================================================================
// 第二部分：一次性读取文件
// =============================================================================

// demoReadFileAtOnce 演示一次性读取整个文件
func demoReadFileAtOnce(path string) {
	fmt.Println("========== 一次性读取文件 ==========")

	// os.ReadFile 将整个文件内容读入内存
	// 适用于小文件（几 KB 到几 MB）
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}

	// data 的类型是 []byte
	fmt.Printf("文件大小: %d 字节\n", len(data))
	fmt.Println("文件内容:")
	fmt.Println(string(data))
	fmt.Println()
}

// =============================================================================
// 第三部分：流式读取文件
// =============================================================================

// demoStreamRead 演示使用 file.Read 流式读取
func demoStreamRead(path string) {
	fmt.Println("========== 流式读取（file.Read）==========")

	// os.Open 打开文件，返回 *os.File
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("打开失败:", err)
		return
	}
	// defer 确保函数结束时关闭文件
	defer file.Close()

	// 创建一个 64 字节的缓冲区（实际使用中通常更大，如 4KB）
	buffer := make([]byte, 64)
	readCount := 0

	fmt.Println("分块读取:")
	for {
		// Read 返回实际读取的字节数
		n, err := file.Read(buffer)

		if n > 0 {
			readCount++
			// buffer[:n] 只取实际读取的部分
			fmt.Printf("第 %d 块 (%d 字节): %q\n", readCount, n, string(buffer[:n]))
		}

		// err != nil 包括 io.EOF（文件末尾）
		if err != nil {
			fmt.Printf("读取结束: %v\n", err)
			break
		}
	}

	fmt.Println()
}

// =============================================================================
// 第四部分：按行读取文件
// =============================================================================

// demoScannerRead 演示使用 bufio.Scanner 按行读取
func demoScannerRead(path string) {
	fmt.Println("========== 按行读取（bufio.Scanner）==========")

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("打开失败:", err)
		return
	}
	defer file.Close()

	// bufio.Scanner 用于按行扫描文件
	scanner := bufio.NewScanner(file)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		// Text() 返回当前行的内容（不含换行符）
		line := scanner.Text()
		fmt.Printf("第 %d 行: %s\n", lineNum, line)
	}

	// 检查扫描过程中是否有错误
	if err := scanner.Err(); err != nil {
		fmt.Println("扫描错误:", err)
	}

	fmt.Println()
}

// =============================================================================
// 第五部分：读取文件元信息
// =============================================================================

// demoFileInfo 演示获取文件信息
func demoFileInfo(path string) {
	fmt.Println("========== 文件元信息 ==========")

	// os.Stat 获取文件的元信息
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("获取信息失败:", err)
		return
	}

	fmt.Println("文件名:", info.Name())
	fmt.Println("大小:", info.Size(), "字节")
	fmt.Println("修改时间:", info.ModTime().Format("2006-01-02 15:04:05"))
	fmt.Println("是否为目录:", info.IsDir())
	fmt.Println("权限:", info.Mode())

	fmt.Println()
}

// =============================================================================
// 第六部分：写入文件
// =============================================================================

// demoWriteFile 演示一次性写入文件
func demoWriteFile() {
	fmt.Println("========== 一次性写入文件 ==========")

	content := "这是测试内容。\n第二行。\n第三行。\n"

	// os.WriteFile 一次性写入整个文件
	// 0644 = 所有者读写，其他用户只读
	err := os.WriteFile("output_simple.txt", []byte(content), 0644)
	if err != nil {
		fmt.Println("写入失败:", err)
		return
	}

	fmt.Println("文件写入成功: output_simple.txt")
	fmt.Println()
}

// demoAppendFile 演示追加写入文件
func demoAppendFile() {
	fmt.Println("========== 追加写入文件 ==========")

	// O_APPEND: 追加模式
	// O_CREATE: 文件不存在则创建
	// O_WRONLY: 只写模式
	file, err := os.OpenFile("output_simple.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("打开失败:", err)
		return
	}
	defer file.Close()

	// WriteString 直接写入字符串
	_, err = file.WriteString("追加的第四行。\n")
	if err != nil {
		fmt.Println("写入失败:", err)
		return
	}

	fmt.Println("内容已追加到: output_simple.txt")
	fmt.Println()
}

// =============================================================================
// 第七部分：带缓冲的写入
// =============================================================================

// demoBufferedWrite 演示使用 bufio.Writer 带缓冲写入
func demoBufferedWrite() {
	fmt.Println("========== 带缓冲的写入 ==========")

	// os.Create 创建文件（已存在则清空）
	file, err := os.Create("output_buffered.txt")
	if err != nil {
		fmt.Println("创建失败:", err)
		return
	}
	defer file.Close()

	// bufio.Writer 在内存中维护缓冲区
	// 数据先写入缓冲区，缓冲区满了再写入文件
	writer := bufio.NewWriter(file)

	// 写入多行数据
	for i := 1; i <= 100; i++ {
		// 数据先进入缓冲区，不会立即写入磁盘
		writer.WriteString(fmt.Sprintf("这是第 %d 行数据\n", i))
	}

	// Flush 强制将缓冲区的数据写入文件
	// 如果不调用 Flush，缓冲区中的数据可能丢失
	err = writer.Flush()
	if err != nil {
		fmt.Println("Flush 失败:", err)
		return
	}

	fmt.Println("带缓冲写入完成: output_buffered.txt")
	fmt.Println()
}

// =============================================================================
// 第八部分：日志解析函数
// =============================================================================

// parseLogLine 解析一行日志
// 输入格式: "2024-01-15 10:30:12 [ERROR] Database connection failed"
// 返回: LogEntry 结构和是否解析成功
func parseLogLine(line string) (LogEntry, bool) {
	// 查找日志级别的位置
	start := strings.Index(line, "[")
	end := strings.Index(line, "]")

	// 格式验证
	if start == -1 || end == -1 || end <= start {
		return LogEntry{}, false
	}

	return LogEntry{
		Time:    strings.TrimSpace(line[:start]),
		Level:   line[start+1 : end],
		Message: strings.TrimSpace(line[end+1:]),
	}, true
}

// =============================================================================
// 第九部分：完整的日志分析功能
// =============================================================================

// AnalyzeFile 分析日志文件
func (a *LogAnalyzer) AnalyzeFile(path string) error {
	// 先获取文件信息
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}
	fmt.Printf("正在分析: %s (%.2f KB)\n", info.Name(), float64(info.Size())/1024)

	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 使用 Scanner 按行读取
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		a.TotalLines++
		line := scanner.Text()

		entry, ok := parseLogLine(line)
		if !ok {
			continue
		}

		// 根据日志级别统计
		switch entry.Level {
		case "INFO":
			a.InfoCount++
		case "WARN":
			a.WarnCount++
		case "ERROR":
			a.ErrorCount++
			a.Errors = append(a.Errors, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	return nil
}

// WriteReport 生成分析报告
func (a *LogAnalyzer) WriteReport(path string) error {
	// 创建报告文件
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建报告文件失败: %w", err)
	}
	defer file.Close()

	// 使用带缓冲的写入器
	writer := bufio.NewWriter(file)
	// 确保函数结束时刷新缓冲区
	defer writer.Flush()

	// 写入报告内容
	writer.WriteString("═══════════════════════════════════════════\n")
	writer.WriteString("             日志分析报告\n")
	writer.WriteString("═══════════════════════════════════════════\n\n")

	writer.WriteString(fmt.Sprintf("生成时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	writer.WriteString("【统计摘要】\n")
	writer.WriteString(fmt.Sprintf("  总行数:    %d\n", a.TotalLines))
	writer.WriteString(fmt.Sprintf("  INFO:      %d\n", a.InfoCount))
	writer.WriteString(fmt.Sprintf("  WARN:      %d\n", a.WarnCount))
	writer.WriteString(fmt.Sprintf("  ERROR:     %d\n", a.ErrorCount))

	if len(a.Errors) > 0 {
		writer.WriteString("\n【错误详情】\n")
		for i, e := range a.Errors {
			writer.WriteString(fmt.Sprintf("  %d. [%s] %s\n", i+1, e.Time, e.Message))
		}
	}

	writer.WriteString("\n═══════════════════════════════════════════\n")

	return nil
}

// =============================================================================
// 第十部分：创建测试日志文件
// =============================================================================

// createSampleLogFile 创建示例日志文件
func createSampleLogFile(path string) error {
	logs := `2024-01-15 10:30:01 [INFO] Server started on port 8080
2024-01-15 10:30:05 [INFO] User login: user_123
2024-01-15 10:30:12 [ERROR] Database connection failed: timeout
2024-01-15 10:30:15 [INFO] Retry database connection
2024-01-15 10:30:18 [ERROR] Database connection failed: refused
2024-01-15 10:30:20 [WARN] Switching to backup database
2024-01-15 10:30:22 [INFO] Connected to backup database
2024-01-15 10:30:25 [INFO] User logout: user_123
2024-01-15 10:30:30 [WARN] High memory usage detected
2024-01-15 10:30:35 [ERROR] API request timeout: /api/users
`
	return os.WriteFile(path, []byte(logs), 0644)
}

// =============================================================================
// 主函数
// =============================================================================

func main() {
	// 创建测试日志文件
	logPath := "server.log"
	err := createSampleLogFile(logPath)
	if err != nil {
		fmt.Println("创建测试文件失败:", err)
		return
	}
	fmt.Println("已创建测试日志文件:", logPath)
	fmt.Println()

	// 1. 字节基础操作
	demoByteBasics()

	// 2. 一次性读取
	demoReadFileAtOnce(logPath)

	// 3. 流式读取
	demoStreamRead(logPath)

	// 4. 按行读取
	demoScannerRead(logPath)

	// 5. 文件元信息
	demoFileInfo(logPath)

	// 6. 写入文件
	demoWriteFile()
	demoAppendFile()

	// 7. 带缓冲的写入
	demoBufferedWrite()

	// 8. 完整的日志分析演示
	fmt.Println("========== 日志分析系统演示 ==========")
	analyzer := &LogAnalyzer{}

	err = analyzer.AnalyzeFile(logPath)
	if err != nil {
		fmt.Println("分析失败:", err)
		return
	}

	fmt.Printf("\n分析完成:\n")
	fmt.Printf("  总行数: %d\n", analyzer.TotalLines)
	fmt.Printf("  INFO: %d, WARN: %d, ERROR: %d\n",
		analyzer.InfoCount, analyzer.WarnCount, analyzer.ErrorCount)

	reportPath := "analysis_report.txt"
	err = analyzer.WriteReport(reportPath)
	if err != nil {
		fmt.Println("生成报告失败:", err)
		return
	}
	fmt.Printf("\n报告已生成: %s\n", reportPath)

	// 显示报告内容
	fmt.Println("\n报告内容:")
	reportData, _ := os.ReadFile(reportPath)
	fmt.Println(string(reportData))
}
