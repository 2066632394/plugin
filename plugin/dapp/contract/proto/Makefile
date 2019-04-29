build:
	@protoc --go_out=plugins=grpc:../types ./*.proto
	@protoc --govalidators_out=../types ./contract.proto
	@sed -i -e '/_ "validator"/d' -e '/_ "\."/d' ../types/*.pb.go

output:
	@mkdir proto
	@cat action.proto | sed "/^import.*/d" |  sed "s/ \[.*\]//g" > proto/action.proto
	@cp common.proto transaction.proto proto/
	@tar -zcvf proto.tar.gz proto/
	@rm -rf proto

transfer:
	@cat contract.proto | sed "/^import.*/d" |  sed "s/ \[.*\]//g" > contract_transfer.proto

