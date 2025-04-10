package services

import (
	"sync"

	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/db/queries"
)

var allCardClassesCache map[byte]*dbModels.CardClass = nil
var allCardClassesCacheMtx sync.Mutex

func GetAllCardClasses() (map[byte]*dbModels.CardClass, error) {
	allCardClassesCacheMtx.Lock()
	defer allCardClassesCacheMtx.Unlock()

	if allCardClassesCache == nil {
		classesFromDB, err := queries.GetClasses()
		if err != nil {
			return nil, err
		}
		allCardClassesCache = make(map[byte]*dbModels.CardClass)
		for _, cardClass := range classesFromDB {
			allCardClassesCache[cardClass.ID] = cardClass
		}
	}

	return allCardClassesCache, nil
}
