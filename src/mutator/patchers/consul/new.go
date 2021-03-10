package consul

func NewPatcher(defaultConsulInitImage string) Patcher {
	return Patcher{
		defaultConsulInitImage: defaultConsulInitImage,
	}
}
