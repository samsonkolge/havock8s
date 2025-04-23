package chaos

func init() {
	// Register all chaos injectors
	RegisterInjector("PodFailure", &PodFailureInjector{})
	RegisterInjector("DiskFailure", &DiskFailureInjector{})
	RegisterInjector("NetworkLatency", &NetworkLatencyInjector{})
	RegisterInjector("StatefulSetScaling", &StatefulSetScalingInjector{})
}
