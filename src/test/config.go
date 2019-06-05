package wabin

type Config struct {
	ServerAddress  string
	ServerPort     uint32
	MaxReadTime    uint32
	MaxWriteTime   uint32
	MaxPackageSize uint32
}
