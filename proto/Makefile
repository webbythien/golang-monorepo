.PHONY: fmt generate-api

generate-api: ## Generate api code
	rm -rf $$(find ../backend/ ../api/ ../common/ -mindepth 1 -maxdepth 1 -type d)
	buf generate --template=templates/buf.api.gen.yaml -v
	@for file in $$(grep -iRl --include="*.pb.go" "// @gotags:" ../api/ ../backend/ ../common/); do \
		protoc-go-inject-tag -input=$$file -remove_tag_comment; \
	done