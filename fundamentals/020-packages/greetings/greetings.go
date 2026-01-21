// Package greetings 提供问候功能
package greetings

// Hello 返回对指定名字的问候语
// 函数名首字母大写，表示这是一个"导出"的函数，其他包可以调用它
func Hello(name string) string {
	return "你好, " + name + "!"
}

// goodbye 是一个私有函数（首字母小写）
// 其他包无法调用它，只能在 greetings 包内部使用
func goodbye(name string) string {
	return "再见, " + name + "!"
}
