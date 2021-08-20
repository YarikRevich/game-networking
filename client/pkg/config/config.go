package config


type Config struct {
	//IP for dialer
	IP string

	//Port for dialer
	Port string
	
	//Addr which will be used for pinger
	//to count the timeout
	PingerAddr string 
}