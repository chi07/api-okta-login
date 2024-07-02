package config

type MongoDBConfig struct {
	URI string
}

type Config struct {
	MongoDB MongoDBConfig
}
