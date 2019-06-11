package structNovedad

import (
	"time"

	"github.com/xubiosueldos/concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/structGormModel"
	"github.com/xubiosueldos/legajo/structLegajo"
)

type Novedad struct {
	structGormModel.GormModel
	Nombre      string                   `json:"nombre"`
	Codigo      string                   `json:"codigo"`
	Descripcion string                   `json:"descripcion"`
	Activo      int                      `json:"activo"`
	Importe     float32                  `json:"importe" sql:"type:decimal(19,4);"`
	Cantidad    int                      `json:"cantidad"`
	Fecha       time.Time                `json:"fecha"`
	Legajo      *structLegajo.Legajo     `json:"legajo" gorm:"ForeignKey:Legajoid;association_foreignkey:ID;association_autoupdate:false"`
	Legajoid    *int                     `json:"legajoid" sql:"type:int REFERENCES Legajo(ID)"`
	Concepto    *structConcepto.Concepto `json:"concepto" gorm:"ForeignKey:Conceptoid;association_foreignkey:ID;association_autoupdate:false"`
	Conceptoid  *int                     `json:"conceptoid"`
}
