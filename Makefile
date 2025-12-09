.PHONY: test clean

run ?= ".*"

test:
	go test -v -count=1 -timeout=1h -run "TestQA_$(run)" ./...

clean:
	find . -name "terraform.tfstate*" -delete
	find . -name ".terraform.lock.hcl" -delete
	find . -type d -name ".terraform" -exec rm -rf {} +
	rm test_results.xml test_results.raw
