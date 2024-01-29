test:
	ginkgo -p --randomize-suites --randomize-all --keep-going --trace --junit-report=report.xml --cover --coverprofile=coverage.profile -covermode atomic -r