module demo-person

go 1.13

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/robfig/cron v1.2.0
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.94+incompatible
	github.com/unidoc/unioffice v1.2.1
	google.golang.org/appengine v1.6.5 // indirect
)

replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => github.com/golang/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190603091049-60506f45cf65 => github.com/golang/net v0.0.0-20190603091049-60506f45cf65
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a => github.com/golang/sys v0.0.0-20190215142949-d0b11bdaac8a
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
	golang.org/x/text v0.3.2 => github.com/golang/text v0.3.2
	golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e => github.com/golang/tools v0.0.0-20180917221912-90fa682c2a6e
)
