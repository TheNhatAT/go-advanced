package thingsmustbeclosed

import (
	"github.com/efficientgo/core/errcapture"
	"github.com/efficientgo/core/logerrcapture"
	"github.com/efficientgo/core/merrors"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"io"
	"os"
)

func doWithFile_Wrong(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	//NOTE: a file close often flushes some of writes to disk only on Close => can lose data on an error (ex: permission error, disk full, ...)
	defer f.Close() // Wrong!

	// Use file...

	return nil
}

func doWithFile_CaptureCloseErr(fileName string) (err error) { // NOTE: return an error if Close fails.
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	//NOTE:
	defer errcapture.Do(&err, f.Close, "close file")

	// Use file...

	return nil
}

func doWithFile_LogCloseErr(logger log.Logger, fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		level.Error(logger).Log("err", err)
		return
	}
	defer logerrcapture.Do(logger, f.Close, "close file")

	// Use file...
}

func openMultiple_Wrong(fileNames ...string) ([]io.ReadCloser, error) {
	files := make([]io.ReadCloser, 0, len(fileNames))
	for _, fn := range fileNames {
		f, err := os.Open(fn)
		if err != nil {
			// NOTE: create more resources that have to be closed, or wanna close them in a different function,
			// defer can't be used => potential leaks, e.g., create three files, have an error when opening the second one,
			// => leaking resources for the first non-closed file.
			return nil, err // Leaked files!
		}
		files = append(files, f)
	}
	return files, nil
}

func openMultiple_Correct(fileNames ...string) ([]io.ReadCloser, error) {
	files := make([]io.ReadCloser, 0, len(fileNames))
	for _, fn := range fileNames {
		f, err := os.Open(fn)
		if err != nil {
			return nil, merrors.New(err, closeAll(files)).Err()
		}
		files = append(files, f)
	}
	return files, nil
}

func closeAll(closers []io.ReadCloser) error {
	errs := merrors.New()
	for _, c := range closers {
		errs.Add(c.Close())
	}
	return errs.Err()
}
