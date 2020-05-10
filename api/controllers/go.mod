module github.com/PatrickWalker/buffStreams/controllers

go 1.13

require (
	github.com/PatrickWalker/buffStreams/helpers v0.0.1
	github.com/PatrickWalker/buffStreams/repositories v0.0.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/stretchr/testify v1.5.1
)

replace github.com/PatrickWalker/buffStreams/repositories => ../repositories

replace github.com/PatrickWalker/buffStreams/helpers => ../helpers

replace github.com/PatrickWalker/buffStreams/middleware => ../middleware
