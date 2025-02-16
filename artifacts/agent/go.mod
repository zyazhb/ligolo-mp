module github.com/ttpreport/ligolo-mp/internal/agent

go 1.23.5

replace github.com/ttpreport/ligolo-mp => ../../

require (
	github.com/hashicorp/yamux v0.1.2
	github.com/ttpreport/ligolo-mp v2.0.0-wip+incompatible
	golang.org/x/net v0.35.0
)

require (
	github.com/go-ping/ping v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)
