.PHONY: all build run clean proto proto-lint proto-format wire help

# 默认目标
all: help

# 构建可执行文件
build:
	go build -o bin/server ./cmd/server

# 运行服务
run:
	go run ./cmd/server

# 清理构建产物
clean:
	rm -rf bin/

# 生成 Proto 代码（包含 lint 检查）
proto: proto-lint
	buf generate api/

# Proto 文件 lint 检查
proto-lint:
	buf lint api/

# Proto 文件格式化
proto-format:
	buf format -w api/

# 生成 Wire 依赖注入代码
wire:
	wire ./cmd/server/

# 整理依赖
tidy:
	go mod tidy

# 运行测试
test:
	go test -v ./...

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  build        - Build the server binary"
	@echo "  run          - Run the server"
	@echo "  clean        - Remove build artifacts"
	@echo "  proto        - Generate code from proto files (with lint check)"
	@echo "  proto-lint   - Lint check proto files"
	@echo "  proto-format - Format proto files"
	@echo "  wire         - Generate dependency injection code"
	@echo "  tidy         - Run go mod tidy"
	@echo "  test         - Run tests"
