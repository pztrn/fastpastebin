// Fast Paste Bin - uberfast and easy-to-use pastebin.
//
// Copyright (c) 2018, Stanislav N. aka pztrn and Fast Paste Bin
// developers.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject
// to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
// CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package flatfiles

import (
	// stdlib
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	// local
	"github.com/pztrn/fastpastebin/pastes/model"
)

type FlatFiles struct {
	pastesIndex []*Index
	path        string
	writeMutex  sync.Mutex
}

func (ff *FlatFiles) GetDatabaseConnection() *sql.DB {
	return nil
}

func (ff *FlatFiles) GetPaste(pasteID int) (*pastesmodel.Paste, error) {
	ff.writeMutex.Lock()
	pastePath := filepath.Join(ff.path, "pastes", strconv.Itoa(pasteID)+".json")
	c.Logger.Debug().Msgf("Trying to load paste data from '%s'...", pastePath)
	pasteInBytes, err := ioutil.ReadFile(pastePath)
	if err != nil {
		c.Logger.Debug().Msgf("Failed to read paste from storage: %s", err.Error())
		return nil, err
	}
	c.Logger.Debug().Msgf("Loaded %d bytes: %s", len(pasteInBytes), string(pasteInBytes))
	ff.writeMutex.Unlock()

	paste := &pastesmodel.Paste{}
	err = json.Unmarshal(pasteInBytes, paste)
	if err != nil {
		c.Logger.Error().Msgf("Failed to parse paste: %s", err.Error())
		return nil, err
	}

	return paste, nil
}

func (ff *FlatFiles) GetPagedPastes(page int) ([]pastesmodel.Paste, error) {
	// Pagination.
	var startPagination = 0
	if page > 1 {
		startPagination = (page - 1) * c.Config.Pastes.Pagination
	}

	c.Logger.Debug().Msgf("Pastes index: %+v", ff.pastesIndex)

	// Iteration one - get only public pastes.
	var publicPastes []*Index
	for _, paste := range ff.pastesIndex {
		if !paste.Private {
			publicPastes = append(publicPastes, paste)
		}
	}

	c.Logger.Debug().Msgf("%+v", publicPastes)

	// Iteration two - get paginated pastes.
	var pastesData []pastesmodel.Paste
	for idx, paste := range publicPastes {
		if len(pastesData) == c.Config.Pastes.Pagination {
			break
		}

		if idx < startPagination {
			c.Logger.Debug().Msgf("Paste with index %d isn't in pagination query: too low index", idx)
			continue
		}

		if (idx-1 >= startPagination && page > 1 && idx > startPagination+((page-1)*c.Config.Pastes.Pagination)) || (idx-1 >= startPagination && page == 1 && idx > startPagination+(page*c.Config.Pastes.Pagination)) {
			c.Logger.Debug().Msgf("Paste with index %d isn't in pagination query: too high index", idx)
			break
		}
		c.Logger.Debug().Msgf("Getting paste data (ID: %d, index: %d)", paste.ID, idx)

		// Get paste data.
		pasteData := &pastesmodel.Paste{}
		pasteRawData, err := ioutil.ReadFile(filepath.Join(ff.path, "pastes", strconv.Itoa(paste.ID)+".json"))
		if err != nil {
			c.Logger.Error().Msgf("Failed to read paste data: %s", err.Error())
			continue
		}

		err = json.Unmarshal(pasteRawData, pasteData)
		if err != nil {
			c.Logger.Error().Msgf("Failed to parse paste data: %s", err.Error())
			continue
		}

		pastesData = append(pastesData, (*pasteData))
	}

	return pastesData, nil
}

func (ff *FlatFiles) GetPastesPages() int {
	// Get public pastes count.
	var publicPastes []*Index

	ff.writeMutex.Lock()
	for _, paste := range ff.pastesIndex {
		if !paste.Private {
			publicPastes = append(publicPastes, paste)
		}
	}
	ff.writeMutex.Unlock()

	// Calculate pages.
	pages := len(publicPastes) / c.Config.Pastes.Pagination
	// Check if we have any remainder. Add 1 to pages count if so.
	if len(publicPastes)%c.Config.Pastes.Pagination > 0 {
		pages++
	}

	return pages
}

func (ff *FlatFiles) Initialize() {
	c.Logger.Info().Msg("Initializing flatfiles storage...")

	path := c.Config.Database.Path
	// Get proper paste file path.
	if strings.Contains(c.Config.Database.Path, "~") {
		curUser, err := user.Current()
		if err != nil {
			c.Logger.Error().Msg("Failed to get current user. Will replace '~' for '/' in storage path!")
			path = strings.Replace(path, "~", "/", -1)
		}
		path = strings.Replace(path, "~", curUser.HomeDir, -1)
	}

	path, _ = filepath.Abs(path)
	ff.path = path
	c.Logger.Debug().Msgf("Storage path is now: %s", ff.path)

	// Create directory if neccessary.
	if _, err := os.Stat(ff.path); err != nil {
		c.Logger.Debug().Msgf("Directory '%s' does not exist, creating...", ff.path)
		os.MkdirAll(ff.path, os.ModePerm)
	} else {
		c.Logger.Debug().Msgf("Directory '%s' already exists", ff.path)
	}

	// Create directory for pastes.
	if _, err := os.Stat(filepath.Join(ff.path, "pastes")); err != nil {
		c.Logger.Debug().Msgf("Directory '%s' does not exist, creating...", filepath.Join(ff.path, "pastes"))
		os.MkdirAll(filepath.Join(ff.path, "pastes"), os.ModePerm)
	} else {
		c.Logger.Debug().Msgf("Directory '%s' already exists", filepath.Join(ff.path, "pastes"))
	}

	// Load pastes index.
	ff.pastesIndex = []*Index{}
	if _, err := os.Stat(filepath.Join(ff.path, "pastes", "index.json")); err != nil {
		c.Logger.Warn().Msg("Pastes index file does not exist, will create new one")
	} else {
		indexData, err := ioutil.ReadFile(filepath.Join(ff.path, "pastes", "index.json"))
		if err != nil {
			c.Logger.Fatal().Msg("Failed to read contents of index file!")
		}

		err = json.Unmarshal(indexData, &ff.pastesIndex)
		if err != nil {
			c.Logger.Error().Msgf("Failed to parse index file contents from JSON into internal structure. Will create new index file. All of your previous pastes will became unavailable. Error was: %s", err.Error())
		}

		c.Logger.Debug().Msgf("Parsed pastes index: %+v", ff.pastesIndex)
	}
}

func (ff *FlatFiles) SavePaste(p *pastesmodel.Paste) (int64, error) {
	ff.writeMutex.Lock()
	// Write paste data on disk.
	filesOnDisk, _ := ioutil.ReadDir(filepath.Join(ff.path, "pastes"))
	pasteID := len(filesOnDisk) + 1
	c.Logger.Debug().Msgf("Writing paste to disk, ID will be " + strconv.Itoa(pasteID))
	p.ID = pasteID
	data, err := json.Marshal(p)
	if err != nil {
		ff.writeMutex.Unlock()
		return 0, err
	}
	err = ioutil.WriteFile(filepath.Join(ff.path, "pastes", strconv.Itoa(pasteID)+".json"), data, 0644)
	if err != nil {
		ff.writeMutex.Unlock()
		return 0, err
	}
	// Add it to cache.
	indexData := &Index{}
	indexData.ID = pasteID
	indexData.Private = p.Private
	ff.pastesIndex = append(ff.pastesIndex, indexData)
	ff.writeMutex.Unlock()
	return int64(pasteID), nil
}

func (ff *FlatFiles) Shutdown() {
	c.Logger.Info().Msg("Saving indexes...")
	indexData, err := json.Marshal(ff.pastesIndex)
	if err != nil {
		c.Logger.Error().Msgf("Failed to encode index data into JSON: %s", err.Error())
		return
	}

	err = ioutil.WriteFile(filepath.Join(ff.path, "pastes", "index.json"), indexData, 0644)
	if err != nil {
		c.Logger.Error().Msgf("Failed to write index data to file. Pretty sure that you've lost your pastes.")
		return
	}
}
