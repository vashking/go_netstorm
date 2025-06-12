package plugins

type AttackPlugin interface {
	Name() string
	Description() string
	Run(target string, threads, duration int, params map[string]string) error
}

// plugin example

type ExamplePlugin struct{}

func (e *ExamplePlugin) Name() string        { return "example" }
func (e *ExamplePlugin) Description() string { return "Пример кастомного плагина атаки" }
func (e *ExamplePlugin) Run(target string, threads, duration int, params map[string]string) error {
	
	return nil
} 