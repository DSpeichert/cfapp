package main

// this file doesn't contain anything useful except for swagger annotations

// swagger only
// swagger:parameters CertificateIndex
type CertificateIndexParams struct {
	// in: query
	Offset int `json:"offset" form:"offset"`
	// in: query
	Limit int `json:"limit" form:"limit"`
	// in: query
	CustomerID int `json:"customer_id" form:"customer_id"`
	// in: query
	Active bool `json:"active" form:"active"`
}

// swagger only
// swagger:parameters CertificateCreate
type CertificateCreateParams struct {
	// in: formData
	CustomerID int `json:"customer_id" form:"customer_id"`
	// in: formData
	Active bool `json:"active" form:"customer_id"`
	// in: formData
	Certificate string `json:"certificate" form:"customer_id"`
	// in: formData
	Key string `json:"key" form:"customer_id"`
}

// swagger only
// swagger:parameters CertificateShow
type CertificateShowParams struct {
	// in: path
	ID int `json:"id"`
}

// swagger only
// swagger:parameters CertificateUpdate
type CertificateUpdateParams struct {
	// in: path
	ID int `json:"id" form:"id"`
	// in: formData
	CustomerID int `json:"customer_id" form:"customer_id"`
	// in: formData
	Active bool `json:"active" form:"customer_id"`
	// in: formData
	Certificate string `json:"certificate" form:"customer_id"`
	// in: formData
	Key string `json:"key" form:"customer_id"`
}

// swagger only
// swagger:parameters CustomerIndex
type CustomerIndexParams struct {
	// in: query
	Offset int `json:"offset"`
	// in: query
	Limit int `json:"limit"`
}

// swagger only
// swagger:parameters CustomerCreate
type CustomerCreateParams struct {
	// in: formData
	Name string `json:"name" form:"name"`
	// in: formData
	Email string `json:"email" form:"email"`
	// in: formData
	Password string `json:"password" form:"password"`
}

// swagger only
// swagger:parameters CustomerShow
type CustomerShowParams struct {
	// in: path
	ID int `json:"id"`
}

// swagger only
// swagger:parameters CustomerDelete
type CustomerDeleteParams struct {
	// in: path
	ID int `json:"id"`
}
