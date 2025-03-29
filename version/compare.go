package version

// First returns the first non-empty version string from the provided list of
// versions.
func First(versions ...string) string {
	for _, version := range versions {
		if version != "" {
			return version
		}
	}
	return ""
}
