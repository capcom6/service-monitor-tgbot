package probes

type TCPSocketConfig struct {
	Host string // имя хоста
	Port uint16 // номер порта
}

type HTTPGetConfig struct {
	TCPSocketConfig
	Scheme      string              // схема http/https
	Path        string              // путь
	HTTPHeaders map[string][]string // заголовки запроса
}
