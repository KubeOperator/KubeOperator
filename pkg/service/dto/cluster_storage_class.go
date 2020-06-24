package dto

type StorageClassNFSSpec struct {
	Name      string
	Default   bool
	NFSServer string
	NFSPath   string
}
