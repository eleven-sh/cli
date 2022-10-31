package globals

type CloudProvider string

const (
	AWSCloudProvider     CloudProvider = "aws"
	HetznerCloudProvider CloudProvider = "hetzner"
)

var (
	CurrentCloudProvider     CloudProvider
	CurrentCloudProviderArgs string
)
