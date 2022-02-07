package explorer

type FabricExplorerChart struct {
	PostgreSQL PostgreSQLConfig `json:"postgresql"`
}

type PostgreSQLConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Database string `json:"database"`
}
