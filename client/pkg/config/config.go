package config


type Config struct {
	IP string

	Port string
	
	//Addr which will be used for pinger
	//to count the timeout
	PingerAddr string
}