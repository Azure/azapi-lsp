package filesystem

import (
	"path/filepath"

	"github.com/Azure/azapi-lsp/internal/uri"
)

func (fs *fsystem) markDocumentAsOpen(dh DocumentHandler) error {
	if !fs.documentMetadataExists(dh) {
		return &UnknownDocumentErr{dh}
	}

	fs.docMetaMu.Lock()
	defer fs.docMetaMu.Unlock()

	fs.docMeta[dh.URI()].setOpen(true)
	return nil
}

func (fs *fsystem) HasOpenFiles(dirPath string) (bool, error) {
	files, err := fs.ReadDir(dirPath)
	if err != nil {
		return false, err
	}

	fs.docMetaMu.RLock()
	defer fs.docMetaMu.RUnlock()

	for _, fi := range files {
		u := uri.FromPath(filepath.Join(dirPath, fi.Name()))
		dm, ok := fs.docMeta[u]
		if ok && dm.IsOpen() {
			return true, nil
		}
	}

	return false, nil
}

func (fs *fsystem) createDocumentMetadata(dh DocumentHandler, langId string, text []byte) error {
	if fs.documentMetadataExists(dh) {
		return &MetadataAlreadyExistsErr{dh}
	}

	fs.docMetaMu.Lock()
	defer fs.docMetaMu.Unlock()

	fs.docMeta[dh.URI()] = NewDocumentMetadata(dh, langId, text)
	return nil
}

func (fs *fsystem) removeDocumentMetadata(dh DocumentHandler) error {
	if !fs.documentMetadataExists(dh) {
		return nil
	}

	fs.docMetaMu.Lock()
	defer fs.docMetaMu.Unlock()

	delete(fs.docMeta, dh.URI())
	return nil
}

func (fs *fsystem) documentMetadataExists(dh DocumentHandler) bool {
	fs.docMetaMu.RLock()
	defer fs.docMetaMu.RUnlock()

	_, ok := fs.docMeta[dh.URI()]
	return ok
}

func (fs *fsystem) isDocumentOpen(dh DocumentHandler) (bool, error) {
	fs.docMetaMu.RLock()
	defer fs.docMetaMu.RUnlock()

	dm, ok := fs.docMeta[dh.URI()]
	if !ok {
		return false, &UnknownDocumentErr{dh}
	}

	return dm.isOpen, nil
}

func (fs *fsystem) updateDocumentMetadataLines(dh VersionedDocumentHandler, b []byte) error {
	if !fs.documentMetadataExists(dh) {
		return &UnknownDocumentErr{dh}
	}

	fs.docMetaMu.Lock()
	defer fs.docMetaMu.Unlock()

	fs.docMeta[dh.URI()].updateLines(b)
	fs.docMeta[dh.URI()].setVersion(dh.Version())

	return nil
}

func (fs *fsystem) getDocumentMetadata(dh DocumentHandler) (*documentMetadata, error) {
	fs.docMetaMu.RLock()
	defer fs.docMetaMu.RUnlock()

	dm, ok := fs.docMeta[dh.URI()]
	if !ok {
		return nil, &UnknownDocumentErr{dh}
	}

	return dm, nil
}
