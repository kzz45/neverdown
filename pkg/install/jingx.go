package install

type JingxImage struct {
	DiscoveryEtcd         string
	DiscoveryControlPlane string
	OpenxApiServer        string
	OpenxDashboard        string
	AuthxServer           string
	AuthxDashboard        string
	JingxApiServer        string
	JingxDashboard        string
	DiscoverApiServer     string
	DiscoveryDashboard    string
}

type NativeAppHost struct {
	DiscoveryEtcd             string
	DiscoveryControlPlane     string
	OpenxApiServer            string
	AuthxServer               string
	JingxApiServer            string
	DiscoveryAggregatorServer string
}
