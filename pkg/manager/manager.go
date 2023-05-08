package manager

type Manager struct {
	rootPath string
	files    []*SourceFile
}

func NewManager(rootPath string) (*Manager, error) {
	files, err := loadSourceFiles(rootPath)
	if err != nil {
		return nil, err
	}

	return &Manager{rootPath: rootPath, files: files}, nil
}

func (m *Manager) parse() error {
	for _, v := range m.files {
		if err := v.Parse(); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) Run() error {
	return m.parse()
}
