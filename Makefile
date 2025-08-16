.PHONY: build run test clean

# 构建项目
build:
	go build -o bin/oauth2-server oauth2.go

# 运行项目
run:
	go run oauth2.go

# 测试项目
test:
	go test ./...

# 清理构建文件
clean:
	rm -rf bin/

# 安装依赖
deps:
	go mod tidy

# 运行测试脚本
test-client:
	go run test/test_oauth2.go

# 数据库初始化（需要手动执行）
init-db:
	@echo "请手动执行以下命令初始化数据库："
	@echo "mysql -u root -p < scripts/init.sql" 