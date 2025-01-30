module github.com/trysourcetool/sourcetool-go

go 1.22.1

replace github.com/trysourcetool/sourcetool-proto => ../sourcetool-proto

require (
	github.com/gofrs/uuid/v5 v5.3.0
	github.com/gorilla/websocket v1.5.3
	github.com/samber/lo v1.47.0
	github.com/trysourcetool/sourcetool-proto v0.0.0
	google.golang.org/protobuf v1.32.0
)

require golang.org/x/text v0.16.0 // indirect
