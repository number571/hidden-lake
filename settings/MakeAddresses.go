package settings

func MakeAddresses(node_address map[string]string) []string {
	var (
		addresses = make([]string, len(node_address))
		index uint32
	)
	for _, address := range node_address {
        addresses[index] = address
        index++
    }
    return addresses
}
