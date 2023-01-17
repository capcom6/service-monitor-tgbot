package probes

type TcpSocketConfig struct {
	Host string // имя хоста
	Port uint16 // номер порта
}

type HttpGetConfig struct {
	TcpSocketConfig
	Scheme      string              // схема http/https
	Path        string              // путь
	HTTPHeaders map[string][]string // заголовки запроса
}
