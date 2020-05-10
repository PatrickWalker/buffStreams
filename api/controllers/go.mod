module github.com/buffup/api/controllers

go 1.13

require (
	github.com/buffup/api/helpers v0.0.1
	github.com/buffup/api/repositories v0.0.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/stretchr/testify v1.5.1
)

replace github.com/buffup/api/repositories => ../repositories

replace github.com/buffup/api/helpers => ../helpers

replace github.com/buffup/api/middleware => ../middleware
