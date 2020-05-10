module github.com/PatrickWalker/buffStreams

go 1.13

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/PatrickWalker/buffStreams/controllers v0.0.0
	github.com/PatrickWalker/buffStreams/helpers v0.0.1
	github.com/PatrickWalker/buffStreams/migrations v0.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang-migrate/migrate/v4 v4.9.1
	github.com/gorilla/handlers v1.4.2
	github.com/julienschmidt/httprouter v1.2.0
)

replace github.com/PatrickWalker/buffStreams/modules => ./modules

replace github.com/PatrickWalker/buffStreams/controllers => ./controllers

replace github.com/PatrickWalker/buffStreams/helpers => ./helpers

replace github.com/PatrickWalker/buffStreams/migrations => ./migrations

replace github.com/PatrickWalker/buffStreams/repositories => ./repositories
