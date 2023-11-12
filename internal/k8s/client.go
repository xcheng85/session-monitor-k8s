package k8s

//go:generate mockery --name IK8sInformer
type IK8sInformer interface {
	Run()
}
