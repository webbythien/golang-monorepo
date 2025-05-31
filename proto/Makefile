.PHONY: fmt generate-api

tidy: ## Update dependencies
	@buf dep update
	
generate-api: ## Generate api code
	rm -rf $$(find ../backend/ ../api/ ../common/ -mindepth 1 -maxdepth 1 -type d)
	buf generate --template=templates/buf.api.gen.yaml -v
	@for file in $$(grep -iRl --include="*.pb.go" "// @gotags:" ../api/ ../backend/ ../common/); do \
		protoc-go-inject-tag -input=$$file -remove_tag_comment; \
	done

generate-ts: ## Generate typescript code
	rm -rf $$(find ../web/ -mindepth 1 -maxdepth 1 -type d)
	buf generate --template=templates/buf.ts.gen.yaml -v

generate: tidy generate-api generate-ts ## Generate code
