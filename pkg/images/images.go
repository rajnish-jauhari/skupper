package images

const (
	DefaultImageRegistry       string = "brew.registry.redhat.io/rh-osbs"
	RouterImageName            string = "service-interconnect-skupper-router-rhel9:2.7.2-2"
	ServiceControllerImageName string = "service-interconnect-service-controller-rhel9:1.8.2-2"
	ControllerPodmanImageName  string = "service-interconnect-controller-podman-rhel9:1.8.2-2"
	ConfigSyncImageName        string = "service-interconnect-config-sync-rhel9:1.8.2-2"
	FlowCollectorImageName     string = "service-interconnect-flow-collector-rhel9:1.8.2-2"
	SiteControllerImageName    string = "service-interconnect-site-controller-rhel9:1.8.2-2"
	PrometheusImageRegistry    string = "quay.io/prometheus"
	PrometheusServerImageName  string = "prometheus:v2.42.0"
	OauthProxyImageRegistry    string = "quay.io/openshift"
	OauthProxyImageName        string = "origin-oauth-proxy:4.14.0"
)
