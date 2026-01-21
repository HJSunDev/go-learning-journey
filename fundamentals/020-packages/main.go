package main

import (
	"fmt"

	// 导入自定义包
	// 格式：模块名/包目录名
	// myapp 是 go.mod 中定义的模块名
	// greetings 是包所在的目录名
	"myapp/greetings"

	// 导入第三方包
	// 格式：域名/用户名/仓库名
	"github.com/fatih/color"
)

func main() {
	// 调用 greetings 包的 Hello 函数
	// 语法：包名.函数名(参数)
	message := greetings.Hello("世界")
	fmt.Println(message)

	// 使用第三方包
	color.Green("成功：自定义包调用成功")
	color.Yellow("提示：第三方包也能正常使用")
	color.Cyan("信息：包路径就是 GitHub 地址")

	// 尝试调用私有函数会编译失败（取消注释可验证）：
	// greetings.goodbye("世界")
	// 错误：cannot refer to unexported name greetings.goodbye
}
