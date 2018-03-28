package mapping

type Mapping struct {
	Map map[string]string `yaml:"map"`
}

func Read(path ...string) Mapping {
	return Mapping{}
}
