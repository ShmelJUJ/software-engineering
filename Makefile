CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d+\.\d+)/; print $$1;')

SMARTIMPORTSVER=v0.2.0
SMARTIMPORTSBIN=${BINDIR}/smartimports_${GOVER}

LINTVER=v1.57.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}

MOCKGENVER=v0.4.0
MOCKGENBIN=${BINDIR}/mockgen_${GOVER}

# ==============================================================================
# Main commands

bindir:
	mkdir -p ${BINDIR}
.PHONY: bindir

precommit: format lint
	echo "OK"
.PHONY: precommit

gen-transaction-service:
	swagger generate server \
		-f ./doc/transaction_swagger.yml \
        -t ./transaction/internal/generated -C ./transaction/swagger-templates/server.yml \
        --template-dir ./transaction/swagger-templates/templates \
        --name transaction
.PHONY: gen-transaction-service

gen-monitor-service:
	swagger generate server \
		-f ./doc/monitor_swagger.yml \
        -t ./monitor/internal/generated -C ./monitor/swagger-templates/server.yml \
        --template-dir ./monitor/swagger-templates/templates \
        --name monitor
.PHONY: gen-monitor-service

gen-monitor-client:
	swagger generate client \
		-f ./doc/monitor_swagger.yml \
		-t ./pkg/monitor_client \
		-A MonitorAPI
	${MOCKGENBIN} -package mocks -destination pkg/monitor_client/mocks/monitor_client_mocks.go github.com/ShmelJUJ/software-engineering/pkg/monitor_client/client/monitor ClientService
.PHONY: gen-monitor-client

gen-user:
	cd user
	go generate -run github.com/ogen-go/ogen/cmd/ogen@latest ./...
	cd ..
.PHONY: gen-user

# ==============================================================================
# Tools commands

install-mockgen: bindir
	test -f ${MOCKGENBIN} || \
		(GOBIN=${BINDIR} go install go.uber.org/mock/mockgen@${MOCKGENVER} && \
			go get go.uber.org/mock/mockgen/model && \
		mv ${BINDIR}/mockgen ${MOCKGENBIN})
.PHONY: install-mockgen

gen-mocks: install-mockgen
	go generate -run mockgen ./...
.PHONY: gen-mocks

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})
.PHONY: install-lint

lint: install-lint
	${LINTBIN} run --fix
.PHONY: lint

install-smartimports: bindir
	test -f ${SMARTIMPORTSBIN} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@${SMARTIMPORTSVER} && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTSBIN})
.PHONY: install-smartimports

format: install-smartimports
	${SMARTIMPORTS}
.PHONY: format

# ==============================================================================
# Tests commands

test:
	go test -v -race -count=1 ./...
.PHONY: test

test-100:
	go test -v -race -count=100 ./...
.PHONY: test-100
