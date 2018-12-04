.DEFAULT_GOAL := run
run:
	@go install
	@argo-client
