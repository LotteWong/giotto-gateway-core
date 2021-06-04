module github.com/LotteWong/giotto-gateway-core

go 1.14

require (
	github.com/LotteWong/tcp-conn-server v0.0.0-20210417075150-d4ab7448a98f // direct
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/e421083458/golang_common v1.0.3
	github.com/e421083458/gorm v1.0.1
	github.com/e421083458/grpc-proxy v0.2.0
	github.com/garyburd/redigo v1.6.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-redis/redis/v8 v8.9.0
	github.com/go-redis/redis_rate v6.5.0+incompatible
	github.com/hashicorp/consul/api v1.8.1
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mwitkow/grpc-proxy v0.0.0-20181017164139-0f1106ef9c76 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.8.1
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	github.com/ugorji/go v1.2.4 // indirect
	golang.org/x/crypto v0.0.0-20210218145215-b8e89b74b9df // indirect
	golang.org/x/net v0.0.0-20210415231046-e915ea6b2b7d // indirect
	golang.org/x/sys v0.0.0-20210415045647-66c3f260301c // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	google.golang.org/genproto v0.0.0-20210416161957-9910b6c460de // indirect
	google.golang.org/grpc v1.37.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/go-playground/validator.v9 v9.29.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace github.com/gin-contrib/sse v0.1.0 => github.com/e421083458/sse v0.1.1
