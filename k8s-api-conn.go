
var restConfig *rest.Config
var errKubeConfig error
if config.KubeConfig != "" {
	restConfig, errKubeConfig = clientcmd.BuildConfigFromFlags("", config.KubeConfig)
} else {
	restConfig, errKubeConfig = rest.InClusterConfig()
}
if errKubeConfig != nil {
	logger.Fatal("error getting kubernetes config ", err)
}

kubeClientSet, err := kubernetes.NewForConfig(restConfig)
if err != nil {
	logger.Fatal("error getting kubernetes client ", err)
}
echov1alpha1ClientSet, err := echov1alpha1clientset.NewForConfig(restConfig)
if err != nil {
	logger.Fatal("error creating echo client ", err)
}

ctrl := controller.New(
	kubeClientSet,
	echov1alpha1ClientSet,
	config.Namespace,
	logger.WithField("type", "controller"),
)