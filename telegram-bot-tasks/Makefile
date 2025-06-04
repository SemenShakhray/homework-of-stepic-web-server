test:
	go test -v -race
	
YAML_FILE=./config.yaml	

MIGRATION_DIR:=$(shell yq eval '.migration_dir' $(YAML_FILE))
STORAGE_PATH:=$(shell yq eval '.storage_path' $(YAML_FILE))

migration-add:
	goose -dir ${MIGRATION_DIR} create $(name) sql

migration-up:
	goose -dir ${MIGRATION_DIR} sqlite3 ${STORAGE_PATH} up -v

migration-down:
	goose -dir ${MIGRATION_DIR} sqlite3 ${STORAGE_PATH} down -v	