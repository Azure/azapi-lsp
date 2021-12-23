package module

import "github.com/ms-henglu/azurerm-restapi-lsp/internal/filesystem"

func SyncWalker(fs filesystem.Filesystem, modMgr ModuleManager) *Walker {
	w := NewWalker(fs, modMgr)
	w.sync = true
	return w
}
