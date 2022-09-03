module git.sonicoriginal.software/routes/graphql

go 1.19

require (
	git.sonicoriginal.software/server v0.0.0
	github.com/graphql-go/graphql v0.8.0
)

replace (
	git.sonicoriginal.software/server => github.com/SonicOriginalSoftware/server v0.0.0
	git.sonicoriginal.software/routes/graphql => github.com/SonicOriginalSoftware/server-routes-graphql v0.0.0
)
