package cli

type Parameter struct {
	Name        string
	Shortname   string
	Use         string
	Description string
}

// Aliases returns the parameter's other uses
func (p *Parameter) Aliases() []string {
	return []string{
		p.Name,
		p.Shortname,
	}
}
