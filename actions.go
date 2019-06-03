package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/xubiosueldos/concepto/structConcepto"

	"github.com/xubiosueldos/framework/configuracion"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/novedad/structNovedad"
)

var nombreMicroservicio string = "novedad"

func NovedadList(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {
		versionMicroservicio := obtenerVersionNovedad()
		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer db.Close()

		var novedades []structNovedad.Novedad

		db.Find(&novedades)

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
		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer db.Close()

		//gorm:auto_preload se usa para que complete todos los struct con su informacion
		if err := db.Set("gorm:auto_preload", true).First(&novedad, "id = ?", novedad_id).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		concepto := obtenerConcepto(*novedad.Conceptoid, r)
		novedad.Concepto = concepto
		framework.RespondJSON(w, http.StatusOK, novedad)
	}

}

func obtenerConcepto(conceptoid int, r *http.Request) *structConcepto.Concepto {

	var concepto structConcepto.Concepto

	config := configuracion.GetInstance()

	url := configuracion.GetUrlMicroservicio(config.Puertomicroservicioconcepto) + "concepto/conceptos/" + strconv.Itoa(conceptoid)

	//url := "http://localhost:8084/conceptos/" + strconv.Itoa(conceptoid)

	req, _ := http.NewRequest("GET", url, nil)

	header := r.Header.Get("Authorization")

	req.Header.Add("Authorization", header)

	//res, _ := http.DefaultClient.Do(req)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	str := string(body)

	json.Unmarshal([]byte(str), &concepto)

	return &concepto

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
		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer db.Close()

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
			db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

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

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		novedad_id := params["id"]

		versionMicroservicio := obtenerVersionNovedad()
		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer db.Close()

		//--Borrado Fisico
		if err := db.Unscoped().Where("id = ?", novedad_id).Delete(structNovedad.Novedad{}).Error; err != nil {

			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, framework.Novedad+novedad_id+framework.MicroservicioEliminado)
	}

}

func AutomigrateTablasPrivadas(db *gorm.DB) {

	//para actualizar tablas...agrega columnas e indices, pero no elimina
	db.AutoMigrate(&structNovedad.Novedad{})

}

func obtenerVersionNovedad() int {
	configuracion := configuracion.GetInstance()

	return configuracion.Versionnovedad
}
