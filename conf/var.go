package conf

// go build -ldflags="-s -w -X ginapp/conf.BuildDate=$(date -Iseconds) -X ginapp/conf.BuildBranch=$(git rev-parse --abbrev-ref HEAD)" -o ginapp main.go
var (
	BuildCommitId = "000"
	BuildBranch   = ""
	BuildDate     = ""
	BuildVersion  = ""
)
