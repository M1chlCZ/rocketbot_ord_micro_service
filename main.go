package main

import "api/apis"

// @title Rocketbot ORD API
// @version 1.0
// @description Private API for ORD
// @termsOfService http://swagger.io/terms/

// @contact.name RocketBot
// @contact.url http://app.rocketbot.pro
// @contact.email m1chlcz18@gmail.com
// @contact.name Michal Žídek

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 89.116.25.234:7500
// @BasePath /api
func main() {
	apis.StartORDApi()
	//TestBTC()
}
