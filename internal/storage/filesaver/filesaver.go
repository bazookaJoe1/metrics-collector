package filesaver

import (
	"os"
	"time"
)

type FileSaver struct {
	StoreInterval      time.Duration
	FilePath           string
	Restore            bool
	SaveTicker         *time.Ticker
	SynchronizedSaving bool // saving by ticker or synchronized with storage updating
}

func NewFileSaver(si int, fp string, r bool) FileSaver {
	return FileSaver{
		StoreInterval:      time.Duration(si) * time.Second,
		FilePath:           fp,
		Restore:            r,
		SaveTicker:         nil,
		SynchronizedSaving: false,
	}
}

func (f *FileSaver) SetStoreInterval(si int) {
	f.StoreInterval = time.Duration(si) * time.Second
	if f.StoreInterval != 0 { // set ticker to 'save storage to file' operation if not synchronized save
		f.SaveTicker = time.NewTicker(f.StoreInterval)
		f.SynchronizedSaving = false
	} else {
		f.SaveTicker = nil
		f.SynchronizedSaving = true
	}
}

func (f *FileSaver) SetFilePath(fp string) {
	f.FilePath = fp
}

func (f *FileSaver) SetRestore(r bool) {
	f.Restore = r
}

func (f *FileSaver) Save(d []byte) error {
	file, err := os.OpenFile(f.FilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(d)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileSaver) Load() ([]byte, error) {
	data, err := os.ReadFile(f.FilePath)
	return data, err
}
