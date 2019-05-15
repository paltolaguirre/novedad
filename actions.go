package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/novedad/structNovedad"
)

func NovedadList(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {
		errorToken(w, tokenError)
		return
	} else {

		db := obtenerDB(tokenAutenticacion)
		automigrateTablasPrivadas(db)
		defer db.Close()

		var novedades []structNovedad.Novedad

		db.Find(&novedades)

		framework.RespondJSON(w, http.StatusOK, novedades)
	}

}

func NovedadShow(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {
		errorToken(w, tokenError)
		return
	} else {

		params := mux.Vars(r)
		novedad_id := params["id"]

		var novedad structNovedad.Novedad //Con &var --> lo que devuelve el metodo se le asigna a la var

		db := obtenerDB(tokenAutenticacion)
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

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {
		errorToken(w, tokenError)
		return
	} else {

		decoder := json.NewDecoder(r.Body)

		var novedad_data structNovedad.Novedad

		if err := decoder.Decode(&novedad_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		db := obtenerDB(tokenAutenticacion)
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

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {

		errorToken(w, tokenError)
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

			db := obtenerDB(tokenAutenticacion)
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

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {

		errorToken(w, tokenError)
		return
	} else {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		novedad_id := params["id"]

		db := obtenerDB(tokenAutenticacion)
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

func obtenerDB(tokenAutenticacion *publico.Security) *gorm.DB {

	token := *tokenAutenticacion
	tenant := token.Tenant

	return conexionBD.ConnectBD(tenant)

}

func automigrateTablasPrivadas(db *gorm.DB) {

	//para actualizar tablas...agrega columnas e indices, pero no elimina
	db.AutoMigrate(&structNovedad.Novedad{})

}

func errorToken(w http.ResponseWriter, tokenError *publico.Error) {
	errorToken := *tokenError
	framework.RespondError(w, errorToken.ErrorCodigo, errorToken.ErrorNombre)

}

func checkTokenValido(r *http.Request) (*publico.Security, *publico.Error) {

	var tokenAutenticacion *publico.Security
	var tokenError *publico.Error

	url := "http://localhost:8081/check-token"

	req, _ := http.NewRequest("GET", url, nil)

	header := r.Header.Get("Authorization")

	req.Header.Add("Authorization", header)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusBadRequest {

		// tokenAutenticacion = &(TokenAutenticacion{})
		tokenAutenticacion = new(publico.Security)
		json.Unmarshal([]byte(string(body)), tokenAutenticacion)

	} else {
		tokenError = new(publico.Error)
		json.Unmarshal([]byte(string(body)), tokenError)

	}

	return tokenAutenticacion, tokenError
}
