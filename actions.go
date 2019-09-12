package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/xubiosueldos/framework/configuracion"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/conexionBD/Novedad/structNovedad"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/novedad/structNovedadMin"
	_ "github.com/xubiosueldos/novedad/structNovedadMin"
)

type IdsAEliminar struct {
	Ids []int `json:"ids"`
}

var nombreMicroservicio string = "novedad"

// Sirve para controlar si el server esta OK
func Healthy(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Healthy"))
}

func NovedadList(w http.ResponseWriter, r *http.Request) {

	var legajoid = r.URL.Query()["legajoid"]

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {
		versionMicroservicio := obtenerVersionNovedad()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		var novedades []structNovedadMin.Novedad

		if legajoid != nil {
			db.Set("gorm:auto_preload", true).Where("legajoid = ?", legajoid).Find(&novedades)
		} else {
			db.Set("gorm:auto_preload", true).Find(&novedades)
		}

		framework.RespondJSON(w, http.StatusOK, novedades)
	}

}

func NovedadShow(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)
		novedad_id := params["id"]

		var novedad structNovedad.Novedad //Con &var --> lo que devuelve el metodo se le asigna a la var

		versionMicroservicio := obtenerVersionNovedad()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		//gorm:auto_preload se usa para que complete todos los struct con su informacion
		if err := db.Set("gorm:auto_preload", true).First(&novedad, "id = ?", novedad_id).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, novedad)
	}

}

func NovedadAdd(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		var novedad_data structNovedad.Novedad

		if err := decoder.Decode(&novedad_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		versionMicroservicio := obtenerVersionNovedad()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		if err := db.Create(&novedad_data).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusCreated, novedad_data)
	}
}

func NovedadUpdate(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)
		//se convirtió el string en int para poder comparar
		param_novedadid, _ := strconv.ParseInt(params["id"], 10, 64)
		p_novedadid := int(param_novedadid)

		if p_novedadid == 0 {
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroVacio)
			return
		}

		decoder := json.NewDecoder(r.Body)

		var novedad_data structNovedad.Novedad

		if err := decoder.Decode(&novedad_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		novedadid := novedad_data.ID

		if p_novedadid == novedadid || novedadid == 0 {

			novedad_data.ID = p_novedadid

			versionMicroservicio := obtenerVersionNovedad()
			tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

			db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio)

			//defer db.Close()
			defer apiclientconexionbd.CerrarDB(db)

			if err := db.Save(&novedad_data).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			framework.RespondJSON(w, http.StatusOK, novedad_data)

		} else {
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroDistintoStruct)
			return
		}
	}

}

func NovedadRemove(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		novedad_id := params["id"]

		versionMicroservicio := obtenerVersionNovedad()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		//--Borrado Fisico
		if err := db.Unscoped().Where("id = ?", novedad_id).Delete(structNovedad.Novedad{}).Error; err != nil {

			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, framework.Novedad+novedad_id+framework.MicroservicioEliminado)
	}

}

func NovedadesRemoveMasivo(w http.ResponseWriter, r *http.Request) {
	var resultadoDeEliminacion = make(map[int]string)
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		var idsEliminar IdsAEliminar
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&idsEliminar); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		versionMicroservicio := obtenerVersionNovedad()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio)

		defer apiclientconexionbd.CerrarDB(db)

		if len(idsEliminar.Ids) > 0 {
			for i := 0; i < len(idsEliminar.Ids); i++ {
				novedad_id := idsEliminar.Ids[i]
				if err := db.Unscoped().Where("id = ?", novedad_id).Delete(structNovedad.Novedad{}).Error; err != nil {
					//framework.RespondError(w, http.StatusInternalServerError, err.Error())
					resultadoDeEliminacion[novedad_id] = string(err.Error())

				} else {
					resultadoDeEliminacion[novedad_id] = "Fue eliminado con exito"
				}
			}
		} else {
			framework.RespondError(w, http.StatusInternalServerError, "Seleccione por lo menos un registro")
		}

		framework.RespondJSON(w, http.StatusOK, resultadoDeEliminacion)
	}

}

func obtenerVersionNovedad() int {
	configuracion := configuracion.GetInstance()

	return configuracion.Versionnovedad
}
