package files

import (
	"io"
	"lls_api/pkg/log"
	"lls_api/pkg/rerr"
)

func CloseCloser(closer io.Closer) {
	if closer == nil {
		return
	}
	if err := closer.Close(); err != nil {
		log.DefaultContext().ErrorErr(rerr.Wrap(err))
	}
}
