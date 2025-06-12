package plugins

import "fmt"

var pluginRegistry = make(map[string]AttackPlugin)

func RegisterPlugin(p AttackPlugin) {
	pluginRegistry[p.Name()] = p
}

func GetPlugin(name string) (AttackPlugin, error) {
	p, ok := pluginRegistry[name]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", name)
	}
	return p, nil
}

func ListPlugins() []string {
	names := make([]string, 0, len(pluginRegistry))
	for k := range pluginRegistry {
		names = append(names, k)
	}
	return names
} 