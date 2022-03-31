package http

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Xusser/SimpIndexer/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func GetEndpoint(c *gin.Context) {
	path := c.Param("any")
	log := log.With().Str("path", path).Logger()

	realPath := filepath.Join(config.Get().PathToExplore, path)

	log.Trace().Msgf("Stating %s", realPath)
	fi, err := os.Stat(realPath)
	if err != nil {
		log.Debug().Err(err).Msgf("Unable to stat: %s", realPath)
		c.AbortWithStatus(404)
		return
	}

	// 是文件
	if !fi.IsDir() {
		log.Trace().Msgf("Path %s is file", realPath)
		log.Debug().Msgf("Serving file: %s", realPath)
		c.File(realPath)
		return
	}

	// 是目录
	log.Trace().Msgf("Path %s is directory", realPath)
	items := make([]Item, 0)

	// 不是根目录
	if config.Get().PathToExplore != realPath {
		upperPath := strings.ReplaceAll(strings.Replace(filepath.Join(realPath, ".."), config.Get().PathToExplore, "", 1), "\\", "/")
		if upperPath == "" {
			upperPath = "/"
		}

		items = append(items, Item{
			Href: upperPath,
			Name: "../",
		})
	} else {
		log.Trace().Msgf("Path %s is root directory", realPath)
	}

	if err := filepath.Walk(realPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if path == realPath {
				return nil
			}

			items = append(items, Item{
				Href:    strings.ReplaceAll(strings.Replace(path, config.Get().PathToExplore, "", 1), "\\", "/") + "/",
				Name:    info.Name() + "/",
				ModDate: info.ModTime().Format(time.RFC3339),
			})

			return filepath.SkipDir
		}

		items = append(items, Item{
			Href:    strings.ReplaceAll(strings.Replace(path, config.Get().PathToExplore, "", 1), "\\", "/"),
			Name:    info.Name(),
			ModDate: info.ModTime().Format(time.RFC3339),
			Size:    strconv.FormatInt(info.Size(), 10),
		})

		return nil
	}); err != nil {
		log.Error().Err(err).Msgf("Fail when walking %s", realPath)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "page", gin.H{
		"path":  path,
		"items": items,
	})
}
