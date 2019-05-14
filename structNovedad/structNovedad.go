package structNovedad

import (
	"github.com/jinzhu/gorm"
)

type Novedad struct {
	gorm.Model
	Nombre      string `json:"nombre"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
	Activo      int    `json:"activo"`
}
