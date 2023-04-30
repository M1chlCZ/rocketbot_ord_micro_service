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
	//if len(os.Args) != 2 {
	//	//apis.InitCLI()
	//	utils.WrapErrorLog("API has to be started with --ord / --launchpad")
	//	//utils.WrapErrorLog("Exiting")
	//} else {
	//	if os.Args[1] == "--ord" {
	//		utils.ReportMessage("Rest API v" + utils.VERSION + " - RocketBot API | MODE ORD")
	//		apis.StartORDApi()
	//	} else if os.Args[1] == "--launchpad" {
	//		utils.ReportMessage("Rest API v" + utils.VERSION + " - RocketBot API | MODE Launchpad")
	//		apis.StartLaunchpadApi()
	//	} else {
	//		utils.WrapErrorLog("API has to be started with --ord / --launchpad")
	//		os.Exit(0)
	//	}
	//}

}
