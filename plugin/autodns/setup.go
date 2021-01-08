package autodns

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/caddy"
)

func init() {plugin.Register("autodns", setupAutoDNS)}

func setupAutoDNS(c *caddy.Controller) error {
	c.Next()
	if c.NextArg() {
		return plugin.Error("autodns", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return AutoDNS{} // set the Next field,  so the plugin chaining works
	})
	return nil
}
