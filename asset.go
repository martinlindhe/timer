package timer

func assetData(filename string) []byte {
	data, err := Asset(filename)
	if err != nil {
		panic(err)
	}
	return data
}
