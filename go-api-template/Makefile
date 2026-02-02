.PHONY: all build run clean proto proto-lint proto-format wire swagger help

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

# 生成 Swagger 文档
# 输出到 internal/swagger，包名为 swagger
swagger:
	swag init -g cmd/server/main.go -o internal/swagger --packageName swagger --parseDependency --parseInternal

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
	@echo "  run          - Run the server (press Ctrl+C to gracefully shutdown)"
	@echo "  clean        - Remove build artifacts"
	@echo "  proto        - Generate code from proto files (with lint check)"
	@echo "  proto-lint   - Lint check proto files"
	@echo "  proto-format - Format proto files"
	@echo "  wire         - Generate dependency injection code"
	@echo "  swagger      - Generate Swagger API documentation"
	@echo "  tidy         - Run go mod tidy"
	@echo "  test         - Run tests"
	@echo ""
	@echo "Graceful Shutdown:"
	@echo "  The server supports graceful shutdown. Press Ctrl+C (SIGINT) or send"
	@echo "  SIGTERM to trigger graceful shutdown. The server will:"
	@echo "    1. Stop accepting new connections"
	@echo "    2. Wait for in-flight requests to complete (with timeout)"
	@echo "    3. Close all idle connections"
