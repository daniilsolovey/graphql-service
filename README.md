# Simple GraphQL service

Service can send you a list of products. Register you as a new user.
For authorization, the service send SMS (to logs) and create a token for you with expiration time

## Setup steps:

Run postgres database in container:

```
docker-compose up
```

Clear all database data:

```
make db-drop
```

Run migrations:

```
make db-run-migrations
```

Add test data to database:

```
make db-add-testdata
```

Setup database_url in .env:

```
DATABASE_URL="postgres://postgres:admin@127.0.0.1:5432/graphqlservice?sslmode=disable"
```

Configure yaml file as example:
```yaml
database:
    name: "graphqlservice"
    host: "localhost"
    port: "5432"
    user: "postgres"
    password: "admin"

token:
    secret_key: "my_key"
    expiration_timer: 5

server:
    port: "8080"

sms:
    expiration_timer: 5

## all timers in minutes
```

## Usage:

##### -c --config \<path>
Read specified config file. [default: config.toml].

##### --debug
Enable debug messages.

##### -v --version
Print version.

#####  -h --help
Show this help.


## Build app:
```
make build
```

## Run app:
```
./graphql-service
```
