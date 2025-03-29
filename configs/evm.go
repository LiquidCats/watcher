package configs

type Evm struct {
	RPC EvmRPC `yaml:"rpc" envconfig:"RPC"`
}

type EvmRPC struct {
	URL string `yaml:"url" envconfig:"URL"`
}
