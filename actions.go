package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/novedad/structNovedad"
)

func NovedadList(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := apiclientautenticacion.CheckTokenValido(r)

	if tokenError != nil {
		apiclientautenticacion.ErrorToken(w, tokenError)
		return
	} else {

		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion)
		automigrateTablasPrivadas(db)
		defer db.Close()

		var novedades []structNovedad.Novedad

		db.Find(&novedades)

		framework.RespondJSON(w, http.StatusOK, novedades)
	}

}

func NovedadShow(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := apiclientautenticacion.CheckTokenValido(r)

	if tokenError != nil {
		apiclientautenticacion.ErrorToken(w, tokenError)
		return
	} else {

		params := mux.Vars(r)
		novedad_id := params["id"]

		var novedad structNovedad.Novedad //Con &var --> lo que devuelve el metodo se le asigna a la var

		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion)
		automigrateTablasPrivadas(db)
		defer db.Close()

		//gorm:auto_preload se usa para que complete todos los struct con su informacion
		if err := db.Set("gorm:auto_preload", true).First(&novedad, "id = ?", novedad_id).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, novedad)
	}

}

func NovedadAdd(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := apiclientautenticacion.CheckTokenValido(r)

	if tokenError != nil {
		apiclientautenticacion.ErrorToken(w, tokenError)
		return
	} else {

		decoder := json.NewDecoder(r.Body)

		var novedad_data structNovedad.Novedad

		if err := decoder.Decode(&novedad_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion)
		automigrateTablasPrivadas(db)
		defer db.Close()

		if err := db.Create(&novedad_data).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusCreated, novedad_data)
	}
}

func NovedadUpdate(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := apiclientautenticacion.CheckTokenValido(r)

	if tokenError != nil {

		apiclientautenticacion.ErrorToken(w, tokenError)
		return
	} else {

		params := mux.Vars(r)
		//se convirtió el string en uint para poder comparar
		param_novedadid, _ := strconv.ParseUint(params["id"], 10, 64)
		p_novedadid := uint(param_novedadid)

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

			db := apiclientconexionbd.ObtenerDB(tokenAutenticacion)
			automigrateTablasPrivadas(db)
			defer db.Close()

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

	tokenAutenticacion, tokenError := apiclientautenticacion.CheckTokenValido(r)

	if tokenError != nil {

		apiclientautenticacion.ErrorToken(w, tokenError)
		return
	} else {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		novedad_id := params["id"]

		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion)
		automigrateTablasPrivadas(db)
		defer db.Close()

		//--Borrado Fisico
		if err := db.Unscoped().Where("id = ?", novedad_id).Delete(structNovedad.Novedad{}).Error; err != nil {

			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, framework.NovedadEliminada+novedad_id)
	}

}

func automigrateTablasPrivadas(db *gorm.DB) {

	//para actualizar tablas...agrega columnas e indices, pero no elimina
	db.AutoMigrate(&structNovedad.Novedad{})

}
