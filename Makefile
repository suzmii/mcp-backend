DIST := dist

up-core:
	@echo "Updating core..."
	protoc \
		--proto_path=./core/proto \
		--go_out=.. \
		--go-grpc_out=.. \
		core/proto/*.proto

up-auth:
	@echo "Updating auth..."
	goctl rpc protoc ./applications/auth/proto/auth.proto \
		--go_out=./applications/auth/pb \
		--go-grpc_out=./applications/auth/pb \
		--zrpc_out=./applications/auth \
		--proto_path=. \
		--proto_path=./applications/auth/proto \
		--style go_zero

up-user:
	@echo "Updating user..."
	goctl rpc protoc ./applications/user/proto/user.proto \
		--go_out=./applications/user/pb \
		--go-grpc_out=./applications/user/pb \
		--zrpc_out=./applications/user \
		--proto_path=. \
		--proto_path=./applications/user/proto \
		--style go_zero

up-mcp:
	@echo "Updating mcp..."
	goctl rpc protoc ./applications/mcp/proto/mcp.proto \
		--go_out=./applications/mcp/pb \
		--go-grpc_out=./applications/mcp/pb \
		--zrpc_out=./applications/mcp \
		--proto_path=. \
		--proto_path=./applications/mcp/proto \
		--style go_zero