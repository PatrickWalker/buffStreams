module github.com/buffup/api

go 1.13

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/buffup/api/controllers v0.0.0
	github.com/buffup/api/helpers v0.0.1
	github.com/buffup/api/migrations v0.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang-migrate/migrate/v4 v4.9.1
	github.com/gorilla/handlers v1.4.2
	github.com/julienschmidt/httprouter v1.2.0
)

replace github.com/buffup/api/modules => ./modules

replace github.com/buffup/api/controllers => ./controllers

replace github.com/buffup/api/helpers => ./helpers

replace github.com/buffup/api/migrations => ./migrations

replace github.com/buffup/api/repositories => ./repositories
