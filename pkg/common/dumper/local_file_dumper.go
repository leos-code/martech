package dumper

import (
	"fmt"
	"os"
)

type LocalFileDumper struct {
	filepath string
}

func (d *LocalFileDumper) Dump(obj interface{}, marshal Marshal) error {
	tmp := d.filepath + ".tmp"
	file, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("failed to create file[%s]", err)
	}

	defer func() {
		_ = file.Close()
	}()

	content, err := marshal(obj)
	if err != nil {
		return err
	}

	if _, err = file.Write(content); err != nil {
		return err
	}

	if err := os.Rename(tmp, d.filepath); err != nil {
		return err
	}

	return nil
}
