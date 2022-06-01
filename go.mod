module github.com/KubeOperator/KubeOperator

go 1.16

require (
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/Azure/azure-storage-blob-go v0.10.0
	github.com/CloudyKit/jet/v3 v3.0.1 // indirect
	github.com/KubeOperator/FusionComputeGolangSDK v0.0.4
	github.com/KubeOperator/kobe v0.1.2
	github.com/KubeOperator/kotf v0.1.4
	github.com/Shopify/goreferrer v0.0.0-20210305184658-1a4fe54f556d // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/aliyun/aliyun-oss-go-sdk v2.1.4+incompatible
	github.com/aws/aws-sdk-go v1.33.18
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/benmanns/goworker v0.1.3
	github.com/c-robinson/iplib v0.3.1
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elazarl/goproxy v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/fairwindsops/polaris v0.0.0-20210818215548-9ae4f774e98e
	github.com/ghodss/yaml v1.0.0
	github.com/go-git/go-git/v5 v5.2.0
	github.com/go-ldap/ldap v3.0.3+incompatible
	github.com/go-openapi/spec v0.19.7 // indirect
	github.com/go-ping/ping v0.0.0-20201115131931-3300c582a663
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.2.0
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gofrs/flock v0.8.0
	github.com/golang-migrate/migrate/v4 v4.12.1
	github.com/google/uuid v1.1.3
	github.com/gophercloud/gophercloud v0.12.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/icza/dyno v0.0.0-20210726202311-f1bafe5d9996
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/iris-contrib/jade v1.1.4 // indirect
	github.com/iris-contrib/middleware/jwt v0.0.0-20191219204441-78279b78a367
	github.com/iris-contrib/schema v0.0.6 // indirect
	github.com/iris-contrib/swagger/v12 v12.0.1
	github.com/jinzhu/gorm v1.9.12
	github.com/kataras/golog v0.1.7 // indirect
	github.com/kataras/iris/v12 v12.1.8
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/microcosm-cc/bluemonday v1.0.6 // indirect
	github.com/minio/minio-go/v7 v7.0.21
	github.com/mitchellh/mapstructure v1.4.1
	github.com/mojocn/base64Captcha v1.3.1
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/mozillazg/go-pinyin v0.18.0
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.11.0
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.26.0
	github.com/qri-io/jsonpointer v0.1.1 // indirect
	github.com/robfig/cron v1.1.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/russross/blackfriday v2.0.0+incompatible // indirect
	github.com/ryanuber/columnize v2.1.2+incompatible // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.8.1
	github.com/storyicon/grbac v0.0.0-20200224041032-a0461737df7e
	github.com/swaggo/swag v1.6.5
	github.com/valyala/fasthttp v1.14.0 // indirect
	github.com/vmware/govmomi v0.23.0
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0 // indirect
	github.com/youtube/vitess v2.1.1+incompatible // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/text v0.3.6
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	helm.sh/helm/v3 v3.6.1
	k8s.io/api v0.22.0
	k8s.io/apiextensions-apiserver v0.22.0
	k8s.io/apimachinery v0.22.0
	k8s.io/cli-runtime v0.22.0
	k8s.io/client-go v0.22.0
	rsc.io/letsencrypt v0.0.3 // indirect
)

replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.7.0

replace github.com/mattn/go-colorable => github.com/mattn/go-colorable v0.1.8

replace github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2

replace github.com/mattn/go-isatty => github.com/mattn/go-isatty v0.0.12

replace github.com/nats-io/nats-server/v2 => github.com/nats-io/nats-server/v2 v2.2.1

replace github.com/cespare/xxhash/v2 => github.com/cespare/xxhash/v2 v2.1.2
