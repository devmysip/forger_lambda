require (
	github.com/aws/aws-lambda-go v1.36.1
	github.com/aws/aws-sdk-go v1.55.5
	github.com/aws/aws-sdk-go-v2/config v1.27.32
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.9
	github.com/aws/aws-sdk-go-v2/service/sns v1.31.6
	github.com/disintegration/imaging v1.6.2
)

require (
	github.com/aws/aws-sdk-go-v2 v1.30.5 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.31 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.6 // indirect
	github.com/aws/smithy-go v1.20.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8

module forger

go 1.21

toolchain go1.21.5
