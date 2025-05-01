package build

const (
	cFileSettings = "hl_settings.yaml"
	cFileNetworks = "hl_networks.yaml"
)

func SetBuildByPath(pInputPath string) error {
	if err := setSettings(pInputPath, cFileSettings); err != nil {
		return err
	}
	if err := setNetworks(pInputPath, cFileNetworks); err != nil {
		return err
	}
	return nil
}
