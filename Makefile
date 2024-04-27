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

# ==============================================================================
# Tools commands

install-mockgen: bindir
	test -f ${MOCKGENBIN} || \
		(GOBIN=${BINDIR} go install go.uber.org/mock/mockgen@${MOCKGENVER} && \
			go get go.uber.org/mock/mockgen/model && \
		mv ${BINDIR}/mockgen ${MOCKGENBIN})
.PHONY: install-mockgen

gen-mocks: install-mockgen
	go generate ./...
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
