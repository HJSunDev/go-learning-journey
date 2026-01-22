.PHONY: all build run clean proto wire help

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

# 生成 Proto 代码 (阶段二实现)
proto:
	@echo "Proto generation not configured yet. See Phase 2."

# 生成 Wire 依赖注入代码 (阶段四实现)
wire:
	@echo "Wire generation not configured yet. See Phase 4."

# 整理依赖
tidy:
	go mod tidy

# 运行测试
test:
	go test -v ./...

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  build  - Build the server binary"
	@echo "  run    - Run the server"
	@echo "  clean  - Remove build artifacts"
	@echo "  proto  - Generate code from proto files (Phase 2)"
	@echo "  wire   - Generate dependency injection code (Phase 4)"
	@echo "  tidy   - Run go mod tidy"
	@echo "  test   - Run tests"
