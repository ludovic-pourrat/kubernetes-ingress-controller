package manager

// GatewayParameters contains parameters for the gateway-api Gateway that KIC
type GatewayParameters struct {
	ControllerName string `yaml:"controllerName, omitempty"`
	Name           string `yaml:"name, omitempty"`
	Namespace      string `yaml:"namespace, omitempty"`
}

func (g *GatewayParameters) Validate() error {
	// TODO Validate ControllerName, Name, NameSpace
	return nil
}
