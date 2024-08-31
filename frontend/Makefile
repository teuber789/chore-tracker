# Packages the client-side JS into a single file for running in the browser.
.PHONY: dist
dist:
	@npx webpack ./client.js

# Serves the file built for the browser. dist must be run first.
.PHONY: serve
serve:
	python3 -m http.server 3000
