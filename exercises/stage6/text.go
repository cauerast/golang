package stage6

// ExtractHrefs returns the value of the href attribute of every <a> tag in
// html, in document order. Tag and attribute names are case-insensitive and
// values may be wrapped in single or double quotes. <a> tags without href
// are ignored.
//
// TODO: Use a compiled package-level regexp (regexp.MustCompile). Regular
// expressions are NOT a real HTML parser: for production crawlers use
// golang.org/x/net/html. Think about which inputs would break your regexp.
func ExtractHrefs(html string) []string {
	return nil
}

// ResolveLinks turns hrefs found on the page at base into absolute URLs:
//
//   - Relative references are resolved against base ("../c" on /dir/page.html is /c).
//   - Fragments (#...) are removed.
//   - Only http and https URLs are kept (drop mailto:, javascript:, ...).
//   - Hrefs that fail to parse are skipped.
//   - Duplicates are removed, keeping the first occurrence.
//   - If base itself is not a valid absolute URL, return an error.
//
// TODO: net/url: url.Parse + base.ResolveReference(ref).
func ResolveLinks(base string, hrefs []string) ([]string, error) {
	return nil, nil
}

// Tokenize lowercases text and splits it into words made of letters and
// digits. Everything else is a separator. "Go is FUN!" -> ["go" "is" "fun"].
//
// TODO: strings.FieldsFunc with unicode.IsLetter / unicode.IsDigit.
func Tokenize(text string) []string {
	return nil
}

// BuildIndex builds an inverted index: token -> ascending list of IDs of the
// documents containing it (each ID at most once per token).
//
// TODO: Iterate the docs map (random order!), then sort each posting list.
func BuildIndex(docs map[int]string) map[string][]int {
	return nil
}

// Search returns the ascending IDs of documents that contain ALL tokens of
// query. An empty query (or one with no tokens) matches nothing.
//
// TODO: Intersect the posting lists of the query tokens.
func Search(index map[string][]int, query string) []int {
	return nil
}
