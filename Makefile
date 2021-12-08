NAME = $(notdir $(shell pwd))

VERSION = $(shell printf "%s.%s" \
	$$(git rev-list --count HEAD) \
	$$(git rev-parse --short HEAD) \
)

# could be "..."
TARGET = ...

GOFLAGS = CGO_ENABLED=0

version:
	@echo $(VERSION)

test:
	$(GOFLAGS) go test -failfast -v ./$(TARGET)

get:
	$(GOFLAGS) go get -v -d

build:
	$(GOFLAGS) go build \
		 -ldflags="-s -w -X main.version=$(VERSION)" \
		 -gcflags="-trimpath=$(GOPATH)"

db-drop:
	dbmate -e DATABASE_URL drop

db-run-migrations:
	dbmate up

db-down-migrations:
	dbmate down

db-add-testdata:
	cat test_data.sql | psql postgres://postgres:admin@127.0.0.1:5432/graphqlservice

all: build


# dbmate new create_users_table