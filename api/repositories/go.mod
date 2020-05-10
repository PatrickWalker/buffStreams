module github.com/PatrickWalker/buffStreams/repositories

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/PatrickWalker/buffStreams/helpers v0.0.1
	github.com/elastic/go-elasticsearch/v7 v7.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/stretchr/testify v1.2.2

)

replace github.com/PatrickWalker/buffStreams/helpers => ../helpers
