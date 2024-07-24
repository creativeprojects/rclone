package api

// DirEntriesRepoToken contains a list of DirEntry
type DirEntriesRepoToken struct {
	Entries []DirEntryRepoToken `json:"dirent_list"`
}

// DirEntryRepoToken contains a directory entry
type DirEntryRepoToken struct {
	ID       string   `json:"id"`
	Type     FileType `json:"type"`
	Name     string   `json:"name"`
	Size     int64    `json:"size"`
	Path     string   `json:"parent_dir"`
	Modified string   `json:"mtime"`
}
