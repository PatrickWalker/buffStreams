module github.com/buffup/api/migrations

go 1.13

require (
	github.com/golang-migrate/migrate/v4 v4.9.1
	github.com/buffup/api/helpers v0.0.1

)

replace github.com/buffup/api/helpers => ../helpers
