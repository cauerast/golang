package stage6

// ListGoFiles walks the directory tree rooted at root and returns the paths of
// all files ending in ".go", relative to root, using forward slashes, sorted
// ascending. Directories whose name starts with "." are skipped entirely.
//
// TODO: filepath.WalkDir, return fs.SkipDir for hidden directories,
// filepath.Rel + filepath.ToSlash for the result, then sort.
func ListGoFiles(root string) ([]string, error) {
	return nil, nil
}

// CopyFile copies the file at src to dst (creating or truncating dst) and
// returns the number of bytes copied. The destination file must be closed
// properly: a failed Close on a written file is a real error.
//
// TODO: os.Open, os.Create, io.Copy. Use defer for src; think about the right
// way to handle dst.Close() so its error is not lost (named result + defer).
func CopyFile(dst, src string) (int64, error) {
	return -1, nil
}
