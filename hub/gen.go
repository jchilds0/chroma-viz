package hub

func GenerateTemplateHub(geo []string, geo_count []int, filename string) {
	hub.AddTemplate(0, "left_to_right", "", "left_to_right")

	var total int
	for i := range geo {
		for j := 0; j < geo_count[i]; j++ {
			hub.AddGeometry(0, total, geo[i])
			total++
		}
	}

	ExportArchive(filename)
}
