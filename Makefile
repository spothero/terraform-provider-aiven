all: dep plugin

dep:
	dep ensure

plugin:
	go build -o terraform-provider-aiven .
