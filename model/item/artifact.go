package item

type artifactKey struct {
	name, class string
}

var artifacts = map[artifactKey]bool{
	{"Grayswandir", "silver saber"}: true,
}
