proto:
	kitex -module "github.com/muixstudio/clio" -gen-path "../common" -service user.rpc.service proto/user.proto
	protoc --proto_path=. --micro_out=. --go_out=:. proto/profile.proto
	protoc --proto_path=. --micro_out=. --go_out=:. proto/usergroup.proto
