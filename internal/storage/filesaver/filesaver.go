package filesaver

import (
	"os"
	"time"
)

// FileSaver is the struct that provides methods for saving to and restoring from file RAM storage data.
type FileSaver struct {
	StoreInterval      time.Duration
	FilePath           string
	Restore            bool
	SaveTicker         *time.Ticker
	SynchronizedSaving bool // saving by ticker or synchronized with storage updating
}

// NewFileSaver creates the instance of FileSaver with given params.
func NewFileSaver(si int, fp string, r bool) FileSaver {
	return FileSaver{
		StoreInterval:      time.Duration(si) * time.Second,
		FilePath:           fp,
		Restore:            r,
		SaveTicker:         nil,
		SynchronizedSaving: false,
	}
}

// SetStoreInterval sets store interval of FileSaver. If synchronized save is not enabled, it creates time.Ticker
// to store data by tick.
func (f *FileSaver) SetStoreInterval(si int) {
	f.StoreInterval = time.Duration(si) * time.Second
	if f.StoreInterval != 0 { // set ticker to 'save storage to file' operation if not synchronized save is set
		f.SaveTicker = time.NewTicker(f.StoreInterval)
		f.SynchronizedSaving = false
	} else {
		f.SaveTicker = nil
		f.SynchronizedSaving = true
	}
}

// SetFilePath sets file path of FileSaver.
func (f *FileSaver) SetFilePath(fp string) {
	f.FilePath = fp
}

// SetRestore points to file saver whether to restore data from file.
func (f *FileSaver) SetRestore(r bool) {
	f.Restore = r
}

// Save saves input data to file.
func (f *FileSaver) Save(d []byte) error {
	file, err := os.OpenFile(f.FilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = file.Write(d)
	if err != nil {
		return err
	}

	return nil
}

// Load loads data from file set in FileSaver.
func (f *FileSaver) Load() ([]byte, error) {
	data, err := os.ReadFile(f.FilePath)
	return data, err
}

// GetRestore returns Restore flag of FileSaver.
func (f *FileSaver) GetRestore() bool {
	return f.Restore
}

// GetFilePath returns FilePath of FileSaver.
func (f *FileSaver) GetFilePath() string {
	return f.FilePath
}

// GetSynchronizedFlag returns whether saving in file is synchronized or not.
func (f *FileSaver) GetSynchronizedFlag() bool {
	return f.SynchronizedSaving
}

// GetTicker returns FileSaver ticker. If synchronized saving is enabled, returns nil.
func (f *FileSaver) GetTicker() *time.Ticker {
	return f.SaveTicker
}
