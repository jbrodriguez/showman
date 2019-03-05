package main

import (
	"io"
	"os"
	"path/filepath"
)

// Exists - Check if File / Directory Exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// SearchFile - look for file in the specified locations
func SearchFile(name string, locations []string) string {
	for _, location := range locations {
		if b, _ := Exists(filepath.Join(location, name)); b {
			return location
		}
	}

	return ""
}

// IsDirEmpty -
func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// var lock sync.Mutex

// // Save saves a representation of v to the file at path.
// func Save(path string, v interface{}) error {
// 	lock.Lock()
// 	defer lock.Unlock()
// 	f, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	r, err := Marshal(v)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = io.Copy(f, r)
// 	return err
// }

// // Load loads the file at path into v.
// // Use os.IsNotExist() to see if the returned error is due
// // to the file being missing.
// func Load(path string, v interface{}) error {
// 	lock.Lock()
// 	defer lock.Unlock()
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	return Unmarshal(f, v)
// }

// // Marshal is a function that marshals the object into an
// // io.Reader.
// // By default, it uses the JSON marshaller.
// var Marshal = func(v interface{}) (io.Reader, error) {
// 	b, err := json.MarshalIndent(v, "", "\t")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bytes.NewReader(b), nil
// }

// // Unmarshal is a function that unmarshals the data from the
// // reader into the specified value.
// // By default, it uses the JSON unmarshaller.
// var Unmarshal = func(r io.Reader, v interface{}) error {
// 	return json.NewDecoder(r).Decode(v)
// }
