package connector

type Connector interface {
	Name() string
	SetName(name string)
	SetEnabled(enabled bool)
	Enabled() bool
}

type BaseConnector struct {
	name    string
	enabled bool
}

var _ Connector = (*BaseConnector)(nil)

func (c *BaseConnector) Name() string {
	return c.name
}
func (c *BaseConnector) SetName(name string) {
	c.name = name
}
func (c *BaseConnector) Enabled() bool {
	return c.enabled
}
func (c *BaseConnector) SetEnabled(enabled bool) {
	c.enabled = enabled
}
