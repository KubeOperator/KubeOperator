module github.com/KubeOperator/KubeOperator

go 1.14

require (
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/Azure/azure-storage-blob-go v0.10.0
	github.com/KubeOperator/FusionComputeGolangSDK v0.0.2
	github.com/KubeOperator/kobe v0.0.8
	github.com/KubeOperator/kotf v0.0.8
	github.com/ajg/form v1.5.1 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/aliyun/aliyun-oss-go-sdk v2.1.4+incompatible
	github.com/aws/aws-sdk-go v1.33.18
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/benmanns/goworker v0.1.3
	github.com/c-robinson/iplib v0.3.1
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/spdystream v0.0.0-20170912183627-bc6354cbbc29 // indirect
	github.com/elazarl/goproxy v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/elazarl/goproxy/ext v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/fairwindsops/polaris v0.0.0-20201005190522-9cce9fcec757
	github.com/fasthttp-contrib/websocket v0.0.0-20160511215533-1f3b11f56072 // indirect
	github.com/fluxcd/flux v1.17.2-0.20200121140732-3903cf8e71c3
	github.com/fluxcd/helm-operator v1.0.0-rc6
	github.com/ghodss/yaml v1.0.0
	github.com/go-git/go-git/v5 v5.2.0
	github.com/go-ldap/ldap v3.0.3+incompatible
	github.com/go-openapi/spec v0.19.7 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/go-playground/validator/v10 v10.2.0
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gofrs/flock v0.7.1
	github.com/golang-migrate/migrate/v4 v4.12.1
	github.com/gophercloud/gophercloud v0.12.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/iris-contrib/middleware/jwt v0.0.0-20191219204441-78279b78a367
	github.com/iris-contrib/swagger/v12 v12.0.1
	github.com/jinzhu/gorm v1.9.12
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/golog v0.0.18 // indirect
	github.com/kataras/iris/v12 v12.1.8
	github.com/klauspost/compress v1.10.8 // indirect
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/microcosm-cc/bluemonday v1.0.3 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mojocn/base64Captcha v1.3.1
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/mozillazg/go-pinyin v0.18.0
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.11.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.0
	github.com/storyicon/grbac v0.0.0-20200224041032-a0461737df7e
	github.com/swaggo/swag v1.6.5
	github.com/valyala/fasthttp v1.14.0 // indirect
	github.com/vmware/govmomi v0.23.0
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0 // indirect
	github.com/youtube/vitess v2.1.1+incompatible // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	github.com/yudai/pp v2.0.1+incompatible // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sys v0.0.0-20201029080932-201ba4db2418 // indirect
	golang.org/x/text v0.3.3
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	helm.sh/helm/v3 v3.2.3
	k8s.io/api v0.18.8
	k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery v0.18.8
	k8s.io/cli-runtime v0.18.0
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kubernetes v1.13.0
	rsc.io/letsencrypt v0.0.3 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.18.6

replace github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2

replace (
	github.com/fluxcd/flux => ./thirdparty/flux
	github.com/fluxcd/helm-operator => ./thirdparty/helm-operator
	github.com/fluxcd/helm-operator/pkg/install => ./thirdparty/helm-operator/pkg/install
)
