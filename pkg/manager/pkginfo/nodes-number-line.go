package pkginfo

func (f *FileInfo) initNodesNumberLine() {
	// dst.Inspect(f.File, func(n dst.Node) bool {
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			fmt.Println(err, n)
	// 			panic(err)
	// 		}
	// 	}()

	// 	if n != nil && n.Pos().IsValid() {
	// 		fmt.Println(f.Fset.Position(n.Pos()).Filename, f.Fset.Position(n.Pos()).Line)
	// 	}
	// 	// f.nodesLines[n] = f.Fset.Position(n.Pos()).Line
	// 	return true
	// })
}
