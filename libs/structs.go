package utils

type ConfigStruct struct {
	Server *ServerStruct
	Beacon *BeaconStruct
}

type ServerStruct struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type BeaconStruct struct {
	Protocol   string `json:"proto"`
	Address    string `json:"address"`
	UserName   string `json:"user"`
	Password   string `json:"pass"`
	Queue      string `json:"queue"`
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routingKey"`
	Expiration int    `json:"expiration"`
}
