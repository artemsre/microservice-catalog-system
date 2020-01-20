# microservice-catalog-system
Information system for storing micro service additional data, dependences, owner SLA SLO SLI

Load /migration/init.sql to postgresDB

export GOOGLE_OAUTH2_CLIENT_ID=xxxxxxxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com
export GOOGLE_OAUTH2_CLIENT_SECRET=XXXXXXXXXXXXXXXXXXXXXXXXXXxxx
export DATABASE_URL=postgres://user:pass@127.0.0.1/dbname?sslmode=disable
export VALID_DOMAIN=mycompany_domain
go run main.go
