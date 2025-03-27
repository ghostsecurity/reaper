.PHONY: run ui clean help

.DEFAULT_GOAL := help

run: ## Start the FastAPI app using uvicorn
	uv run python3 main.py

ui: ## Open Chrome to the ReaperBot UI
	open http://localhost:11000

clean: ## Remove SQLite database
	@rm -f utils/.*messages.sqlite

help: ## Show this help
	@echo "Usage: make [target]\n"
	@cat ${MAKEFILE_LIST} | grep "[#]# " | grep -v grep | sort | column -t -s '##' | sed -e 's/^/ /'
	@echo ""

