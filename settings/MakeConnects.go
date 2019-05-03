package settings

func MakeConnects(node_address map[string]string) []string {
	var (
		connects = make([]string, len(node_address))
		index uint32
	)
	for username := range node_address {
        connects[index] = username
        index++
    }
    return connects
}
