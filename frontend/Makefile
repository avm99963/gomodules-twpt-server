.PHONY: serve deploy build_js

WEBPACK := ./node_modules/webpack-cli/bin/cli.js

serve:
	$(WEBPACK) serve --mode development

node_deps:
	npm ci --no-save

clean_deps:
	rm -rf node_modules

deps: node_deps
	rm -rf dist

build_js:
	$(WEBPACK) --mode production

deploy: deps build_js
	firebase deploy --only hosting:twpt-dashboard
