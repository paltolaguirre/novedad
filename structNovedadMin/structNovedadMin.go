package structNovedadMin

import (
	"time"

	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/structGormModel"
)

type Novedad struct {
	structGormModel.GormModel
	Nombre      string                   `json:"nombre"`
	Codigo      string                   `json:"codigo"`
	Descripcion string                   `json:"descripcion"`
	Activo      int                      `json:"activo"`
	Importe     float32                  `json:"importe"`
	Cantidad    int                      `json:"cantidad"`
	Fecha       time.Time                `json:"fecha"`
	Legajo      *Legajo                  `json:"legajo" gorm:"ForeignKey:Legajoid;association_foreignkey:ID;association_autoupdate:false"`
	Legajoid    *int                     `json:"legajoid" sql:"type:int REFERENCES Legajo(ID)"`
	Concepto    *structConcepto.Concepto `json:"concepto" gorm:"ForeignKey:Conceptoid;association_foreignkey:ID;association_autoupdate:false"`
	Conceptoid  *int                     `json:"conceptoid"`
}

type Legajo struct {
	structGormModel.GormModel
	Nombre   string `json:"nombre"`
	Legajo   string `json:"legajo"`
	Apellido string `json:"apellido"`
}
