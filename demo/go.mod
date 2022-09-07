module demo

go 1.19

require (
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de // indirect
	github.com/hvlck/txt v0.0.0-20220828000727-b337f9b83a30 // indirect
)

require (
	github.com/gorilla/mux v1.8.0
	github.com/lacunaverse/cirrus v0.0.0
)

replace github.com/lacunaverse/cirrus => ../

replace github.com/lacunaverse/txt => ../../
