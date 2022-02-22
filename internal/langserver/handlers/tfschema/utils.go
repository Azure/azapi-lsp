package tfschema

func GetResourceSchema(name string) *Resource {
	for _, r := range Resources {
		if r.Name == name {
			return &r
		}
	}
	return nil
}
