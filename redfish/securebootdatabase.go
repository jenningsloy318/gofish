//
// SPDX-License-Identifier: BSD-3-Clause
//

package redfish

import (
	"encoding/json"

	"github.com/stmcginnis/gofish/common"
)

// SecureBootDatabaseResetKeysType is
type SecureBootDatabaseResetKeysType string

const (
	// ResetAllKeysToDefaultSecureBootDatabaseResetKeysType Reset the contents of this UEFI Secure Boot key database to the default
	// values.
	ResetAllKeysToDefaultSecureBootDatabaseResetKeysType SecureBootDatabaseResetKeysType = "ResetAllKeysToDefault"
	// DeleteAllKeysSecureBootDatabaseResetKeysType Delete the contents of this UEFI Secure Boot key database.
	DeleteAllKeysSecureBootDatabaseResetKeysType SecureBootDatabaseResetKeysType = "DeleteAllKeys"
)

// SecureBootDatabase shall be used to represent a UEFI Secure Boot database for a Redfish implementation.
type SecureBootDatabase struct {
	common.Entity
	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataEtag is the odata etag.
	ODataEtag string `json:"@odata.etag"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// certificates shall be a link to a resource collection of type CertificateCollection.
	certificates []string
	// DatabaseID shall contain the name of the UEFI Secure Boot database. This property shall contain the same value
	// as the Id property. The value shall be one of the UEFI-defined Secure Boot databases: 'PK', 'KEK' 'db', 'dbx',
	// 'dbr', 'dbt', 'PKDefault', 'KEKDefault', 'dbDefault', 'dbxDefault', 'dbrDefault', or 'dbtDefault'.
	DatabaseID string
	// Description provides a description of this resource.
	Description string
	// Oem shall contain the OEM extensions. All values for properties that this object contains shall conform to the
	// Redfish Specification-described requirements.
	OEM json.RawMessage `json:"Oem"`
	// Signatures shall be a link to a resource collection of type SignatureCollection.
	signatures []string
}

// UnmarshalJSON unmarshals a SecureBootDatabase object from the raw JSON.
func (securebootdatabase *SecureBootDatabase) UnmarshalJSON(b []byte) error {
	type temp SecureBootDatabase
	var t struct {
		temp
		Certificates common.LinksCollection
		Signatures   common.LinksCollection
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	*securebootdatabase = SecureBootDatabase(t.temp)

	// Extract the links to other entities for later
	securebootdatabase.certificates = t.Certificates.ToStrings()
	securebootdatabase.signatures = t.Signatures.ToStrings()

	return nil
}

// Certificates get the certificates contained in this UEFI Secure Boot database.
func (securebootdatabase *SecureBootDatabase) Certificates() ([]*Certificate, error) {
	var result []*Certificate

	collectionError := common.NewCollectionError()
	for _, uri := range securebootdatabase.certificates {
		item, err := GetCertificate(securebootdatabase.GetClient(), uri)
		if err != nil {
			collectionError.Failures[uri] = err
		} else {
			result = append(result, item)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}

// Signatures get the certificates contained in this UEFI Secure Boot database.
func (securebootdatabase *SecureBootDatabase) Signatures() ([]*Signature, error) {
	var result []*Signature

	collectionError := common.NewCollectionError()
	for _, uri := range securebootdatabase.signatures {
		item, err := GetSignature(securebootdatabase.GetClient(), uri)
		if err != nil {
			collectionError.Failures[uri] = err
		} else {
			result = append(result, item)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}

// GetSecureBootDatabase will get a SecureBootDatabase instance from the service.
func GetSecureBootDatabase(c common.Client, uri string) (*SecureBootDatabase, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var securebootdatabase SecureBootDatabase
	err = json.NewDecoder(resp.Body).Decode(&securebootdatabase)
	if err != nil {
		return nil, err
	}

	securebootdatabase.SetClient(c)
	return &securebootdatabase, nil
}

// ListReferencedSecureBootDatabases gets the collection of SecureBootDatabase from
// a provided reference.
func ListReferencedSecureBootDatabases(c common.Client, link string) ([]*SecureBootDatabase, error) {
	var result []*SecureBootDatabase
	if link == "" {
		return result, nil
	}

	type GetResult struct {
		Item  *SecureBootDatabase
		Link  string
		Error error
	}

	ch := make(chan GetResult)
	collectionError := common.NewCollectionError()
	get := func(link string) {
		securebootdatabase, err := GetSecureBootDatabase(c, link)
		ch <- GetResult{Item: securebootdatabase, Link: link, Error: err}
	}

	go func() {
		err := common.CollectList(get, c, link)
		if err != nil {
			collectionError.Failures[link] = err
		}
		close(ch)
	}()

	for r := range ch {
		if r.Error != nil {
			collectionError.Failures[r.Link] = r.Error
		} else {
			result = append(result, r.Item)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}
