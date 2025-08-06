.PHONY: build run clean deps sample test help

# 默认目标
all: deps sample run

# 下载依赖
deps:
	@echo "正在下载依赖包..."
	go mod tidy
	go mod download

# 创建示例Excel文件
sample:
	@echo "正在创建示例Excel文件..."
	go run create_sample.go

# 运行项目
run:
	@echo "正在启动服务..."
	go run main.go

# 构建项目
build:
	@echo "正在构建项目..."
	go build -o student-lookup main.go

# 测试API
test:
	@echo "正在测试API..."
	@echo "等待服务启动..."
	@sleep 2
	@echo "测试健康检查接口..."
	curl -s http://localhost:8080/health | jq .
	@echo "\n测试查询接口..."
	curl -s -X POST http://localhost:8080/api/v1/search \
		-H "Content-Type: application/json" \
		-d '{"student_id":"202301001","name":"张三"}' | jq .
	@echo "\n测试统计接口..."
	curl -s http://localhost:8080/api/v1/stats | jq .

# 清理临时文件
clean:
	@echo "正在清理临时文件..."
	rm -f student-lookup
	rm -f create_sample

# 显示帮助信息
help:
	@echo "可用的命令:"
	@echo "  make deps    - 下载依赖包"
	@echo "  make sample  - 创建示例Excel文件"
	@echo "  make run     - 运行项目"
	@echo "  make build   - 构建项目"
	@echo "  make test    - 测试API"
	@echo "  make clean   - 清理临时文件"
	@echo "  make all     - 执行完整流程 (deps + sample + run)"
