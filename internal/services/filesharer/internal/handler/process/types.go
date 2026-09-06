package process

import "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"

type IDownloadProcessManager interface {
	Update(dto.IDownloadProcessKey, [2]uint64) bool

	GetList() []dto.IDownloadProcess
	Get(dto.IDownloadProcessKey) (dto.IDownloadProcess, bool)

	Unlock(dto.IDownloadProcessKey) bool
	TryLock(dto.IDownloadProcessKey, func()) bool
}
