package profile

type Format string

const (
	FormatJSON Format = "json"
	FormatTOML Format = "toml"
	FormatYAML Format = "yaml"
)

type Mode string

const (
	ModeOverwrite Mode = "overwrite"
	ModeMerge     Mode = "merge"
)

type Field struct {
	Key      string
	Label    string
	Help     string
	Secret   bool
	Required bool
	Default  string
	Validate func(string) error
}

type PathFunc func() (string, error)

type OutputTarget struct {
	Format Format
	Path   PathFunc
	Mode   Mode
}

type Profile struct {
	ID          string
	DisplayName string
	Description string
	Fields      []Field
	Output      OutputTarget
}

func (p Profile) ResolveOutputPath() (string, error) {
	return p.Output.Path()
}
