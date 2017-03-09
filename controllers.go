// cfapp
//
// A naive, insecure HTTP REST API server serving something
//
//     Schemes: http
//     BasePath: /
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//
// swagger:meta
package main

import (
	"net/http"

	"strconv"

	"github.com/labstack/echo"
)

func CertificateIndex(c echo.Context) error {
	// swagger:route GET /certificates certificates CertificateIndex
	//
	// Lists certificates filtered by some parameters.
	//
	// This will show all available certificates by default.
	var certificates []Certificate
	var total int
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	customer_id, _ := strconv.Atoi(c.QueryParam("customer_id"))
	active, activeError := strconv.ParseBool(c.QueryParam("active"))
	if limit > 100 {
		limit = 100
	} else if limit == 0 {
		limit = 10
	}

	if customer_id != 0 {
		d := db.Offset(offset).Limit(limit).Model(&Customer{ID: customer_id})
		if activeError == nil {
			d = d.Where("active = ?", active)
		}
		d.Related(&certificates)

		d = db.Model(&Certificate{}).Where("customer_id = ?", customer_id)
		if activeError == nil {
			d = d.Where("active = ?", active)
		}
		d.Count(&total)
	} else {
		d := db.Offset(offset).Limit(limit)
		if activeError == nil {
			d = d.Where("active = ?", active)
		}
		d.Find(&certificates)

		d = db.Model(&Certificate{})
		if activeError == nil {
			d = d.Where("active = ?", active)
		}
		d.Count(&total)
	}

	resp := QueryResponse{Total: total, Offset: offset, Limit: limit, Items: certificates}
	c.JSONPretty(http.StatusOK, resp, "    ")
	return nil
}

func CertificateCreate(c echo.Context) error {
	// swagger:route POST /certificates certificates CertificateCreate
	//
	// Create a certificate
	//
	// Creates a certificate
	//
	//     Consumes:
	//     - application/json
	//     - application/x-www-form-urlencoded
	//
	//     Produces:
	//     - application/json
	var certificate Certificate
	if err := c.Bind(&certificate); err != nil {
		return err
	}
	certificate.ID = 0 // just in case it's in the input, to prevent accidental update
	if db.First(&Customer{}, certificate.CustomerID).RecordNotFound() {
		return echo.NewHTTPError(404, "customer not found")
	}
	if err := db.Save(&certificate).Error; err != nil {
		return err
	}
	c.JSONPretty(http.StatusOK, certificate, "    ")
	return nil
}

func CertificateShow(c echo.Context) error {
	// swagger:route GET /certificates/{id} certificates CertificateShow
	//
	// Show a specific certificate by ID
	//
	// Show a specific certificate by ID
	//
	//     Produces:
	//     - application/json
	var certificate Certificate
	id, _ := strconv.Atoi(c.Param("id"))
	if db.First(&certificate, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}
	c.JSONPretty(http.StatusOK, certificate, "    ")
	return nil
}

func CertificateUpdate(c echo.Context) error {
	// swagger:route PUT /certificates/{id} certificates CertificateUpdate
	//
	// Update a certificate
	//
	// Updates a certificate
	//
	//     Consumes:
	//     - application/json
	//     - application/x-www-form-urlencoded
	//
	//     Produces:
	//     - application/json
	var certificate Certificate
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.Bind(&certificate); err != nil {
		return err
	}
	certificate.ID = id
	if db.First(&Certificate{}, certificate.ID).RecordNotFound() {
		return echo.NewHTTPError(404, "certificate not found")
	}
	if db.First(&Customer{}, certificate.CustomerID).RecordNotFound() {
		return echo.NewHTTPError(400, "customer not found")
	}
	if err := db.Save(&certificate).Error; err != nil {
		return err
	}
	certificate.PostUpdateHook()
	c.JSONPretty(http.StatusOK, certificate, "    ")
	return nil
}

func CustomerIndex(c echo.Context) error {
	// swagger:route GET /customers customers CustomerIndex
	//
	// Lists customers filtered by some parameters.
	//
	// This will show all available customers by default.
	//
	//     Produces:
	//     - application/json
	var customers []Customer
	var total int
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit > 100 {
		limit = 100
	} else if limit == 0 {
		limit = 10
	}

	db.Offset(offset).Limit(limit).Find(&customers)
	db.Model(&Customer{}).Count(&total)

	resp := QueryResponse{Total: total, Offset: offset, Limit: limit, Items: customers}
	c.JSONPretty(http.StatusOK, resp, "    ")
	return nil
}

func CustomerCreate(c echo.Context) error {
	// swagger:route POST /customers customers CustomerCreate
	//
	// Creates a customer
	//
	// Creates a customer
	//
	//     Consumes:
	//     - application/json
	//     - application/x-www-form-urlencoded
	//
	//     Produces:
	//     - application/json
	var customer Customer
	if err := c.Bind(&customer); err != nil {
		return err
	}
	customer.ID = 0 // just in case it's in the input, to prevent accidental update
	customer.HashPassword()
	if err := db.Save(&customer).Error; err != nil {
		return err
	}
	c.JSONPretty(http.StatusOK, customer, "    ")
	return nil
}

func CustomerShow(c echo.Context) error {
	// swagger:route GET /customers/{id} customers CustomerShow
	//
	// Show a specific customer by ID
	//
	// Show a specific customer by ID
	//
	//     Produces:
	//     - application/json
	var customer Customer
	id, _ := strconv.Atoi(c.Param("id"))
	if db.First(&customer, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}
	c.JSONPretty(http.StatusOK, customer, "    ")
	return nil
}

func CustomerDelete(c echo.Context) error {
	// swagger:route DELETE /customers/{id} customers CustomerDelete
	//
	// Delete a specific customer by ID
	//
	// Delete a specific customer by ID
	var customer Customer
	id, _ := strconv.Atoi(c.Param("id"))
	if db.First(&customer, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}
	db.Delete(&customer)
	c.JSONPretty(http.StatusOK, nil, "    ")
	return nil
}
