# gokitelogin
Go library to login into Kite Connect API services.

To get the package,

go get github.com/sivamgr/gokitelogin


	err := gokitelogin.Login(kc.GetLoginURL(),
		os.Getenv("KITE_USERID"),
		os.Getenv("KITE_PASSWD"),
		os.Getenv("KITE_PIN"))
	if err != nil {
		log.Printf("Failed to Authenticate into Kite Connect API. %+v\n", err)
	}

