//
// SPDX-License-Identifier: BSD-3-Clause
//

package redfish

import (
	"encoding/json"

	"github.com/stmcginnis/gofish/common"
)

type CoolingEquipmentType string

const (
	// CDUCoolingEquipmentType A coolant distribution unit (CDU).
	CDUCoolingEquipmentType CoolingEquipmentType = "CDU"
	// HeatExchangerCoolingEquipmentType A heat exchanger.
	HeatExchangerCoolingEquipmentType CoolingEquipmentType = "HeatExchanger"
	// ImmersionUnitCoolingEquipmentType An immersion cooling unit.
	ImmersionUnitCoolingEquipmentType CoolingEquipmentType = "ImmersionUnit"
)

// CoolingUnit shall represent a cooling system component or unit for a Redfish implementation.
type CoolingUnit struct {
	common.Entity
	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataEtag is the odata etag.
	ODataEtag string `json:"@odata.etag"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// AllowedCoolingUnitModes shall contain the allowed values for setting the mode of this cooling unit.
	AllowedCoolingUnitModes []CoolingUnitMode
	// Assembly shall contain a link to a resource of type Assembly.
	assembly string
	// AssetTag shall contain the user-assigned asset tag, which is an identifying string that tracks the equipment for
	// inventory purposes.
	AssetTag string
	// Coolant shall contain details regarding the coolant contained or used by this unit.
	Coolant Coolant
	// CoolantConnectorRedundancy shall contain redundancy information for the set of coolant connectors attached to
	// this equipment. The values of the RedundancyGroup array shall reference resources of type CoolantConnector.
	coolantConnectorRedundancy []string
	// CoolingCapacityWatts shall contain the manufacturer-provided cooling capacity, in watt units, of this equipment.
	CoolingCapacityWatts int
	// Description provides a description of this resource.
	Description string
	// EnvironmentMetrics shall contain a link to a resource of type EnvironmentMetrics that specifies the environment
	// metrics for this equipment.
	environmentMetrics string
	// EquipmentType shall contain the type of equipment this resource represents.
	EquipmentType CoolingEquipmentType
	// FilterRedundancy shall contain redundancy information for the groups of filters in this unit.
	FilterRedundancy []RedundantGroup
	// Filters shall contain a link to a resource collection of type FilterCollection that contains the filter
	// information for this equipment.
	filters string
	// FirmwareVersion shall contain a string describing the firmware version of this equipment as provided by the
	// manufacturer.
	FirmwareVersion string
	// LeakDetection shall contain a link to a resource of type LeakDetection that contains the leak detection
	// component information for this equipment.
	leakDetection string
	// Location shall contain the location information of the associated equipment.
	Location common.Location
	// Manufacturer shall contain the name of the organization responsible for producing the equipment. This
	// organization may be the entity from which the equipment is purchased, but this is not necessarily true.
	Manufacturer string
	// Model shall contain the manufacturer-provided model information of this equipment.
	Model string
	// PartNumber shall contain the manufacturer-provided part number for the equipment.
	PartNumber string
	// PrimaryCoolantConnectors shall contain a link to a resource collection of type CoolantConnectorCollection that
	// contains the primary coolant connectors for this equipment.
	primaryCoolantConnectors string
	// ProductionDate shall contain the date of production or manufacture for this equipment.
	ProductionDate string
	// PumpRedundancy shall contain redundancy information for the groups of pumps in this unit.
	PumpRedundancy []RedundantGroup
	// Pumps shall contain a link to a resource collection of type PumpCollection that contains the pump information
	// for this equipment.
	pumps string
	// Reservoirs shall contain a link to a resource collection of type ReservoirCollection that contains the reservoir
	// information for this equipment.
	reservoirs string
	// SecondaryCoolantConnectors shall contain a link to a resource collection of type CoolantConnectorCollection that
	// contains the secondary coolant connectors for this equipment.
	secondaryCoolantConnectors string
	// SerialNumber shall contain a manufacturer-allocated number that identifies the equipment.
	SerialNumber string
	// setMode shall contain the action target for setting the mode of this cooling unit.
	setMode string
	// Status shall contain any status or health properties of the resource.
	Status common.Status
	// UserLabel shall contain a user-assigned label used to identify this resource. If a value has not been assigned
	// by a user, the value of this property shall be an empty string.
	UserLabel string
	// Version shall contain the hardware version of this equipment as determined by the vendor or supplier.
	Version string
	// rawData holds the original serialized JSON so we can compare updates.
	rawData []byte
	chassis []string
	// ChassisCount is the number of physical containers that contain this equipment.
	ChassisCount int
	facility     string
	managedBy    []string
	// ManagedByCount is the number of managers that manage this equipment.
	ManagedByCount int
}

type CoolingUnitMode string

const (
	EnabledCoolingUnitMode  CoolingUnitMode = "Enabled"
	DisabledCoolingUnitMode CoolingUnitMode = "Disabled"
)

// UnmarshalJSON unmarshals a CoolingUnit object from the raw JSON.
func (coolingunit *CoolingUnit) UnmarshalJSON(b []byte) error {
	type temp CoolingUnit
	type Links struct {
		// Chassis shall contain an array of links to resources of type Chassis that represent the physical containers that
		// contain this equipment.
		Chassis common.Links
		// Chassis@odata.count
		ChassisCount int `json:"Chassis@odata.count"`
		// Facility shall contain a link to a resource of type Facility that represents the facility that contains this
		// equipment.
		Facility common.Link
		// ManagedBy shall contain an array of links to resources of type Manager that represent the managers that manage
		// this equipment.
		ManagedBy common.Links
		// ManagedBy@odata.count
		ManagedByCount int `json:"ManagedBy@odata.count"`
	}
	type CoolingUnitActions struct {
		SetMode struct {
			AllowedCoolingUnitModes []CoolingUnitMode `json:"Mode@Redfish.AllowableValues"`
			Target                  string
		} `json:"#CoolingUnit.SetMode"`
	}

	var t struct {
		temp
		Actions                    CoolingUnitActions
		Assembly                   common.Link
		CoolantConnectorRedundancy common.Links
		EnvironmentMetrics         common.Link
		Filters                    common.Link
		LeakDetection              common.Link
		PrimaryCoolantConnectors   common.Link
		Pumps                      common.Link
		Reservoirs                 common.Link
		SecondaryCoolantConnectors common.Link
		Links                      Links
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	*coolingunit = CoolingUnit(t.temp)

	// Extract the links to other entities for later
	coolingunit.assembly = t.Assembly.String()
	coolingunit.coolantConnectorRedundancy = t.CoolantConnectorRedundancy.ToStrings()
	coolingunit.environmentMetrics = t.EnvironmentMetrics.String()
	coolingunit.filters = t.Filters.String()
	coolingunit.leakDetection = t.LeakDetection.String()
	coolingunit.primaryCoolantConnectors = t.PrimaryCoolantConnectors.String()
	coolingunit.pumps = t.Pumps.String()
	coolingunit.reservoirs = t.Reservoirs.String()
	coolingunit.secondaryCoolantConnectors = t.SecondaryCoolantConnectors.String()
	coolingunit.chassis = t.Links.Chassis.ToStrings()
	coolingunit.ChassisCount = t.Links.ChassisCount
	coolingunit.facility = t.Links.Facility.String()
	coolingunit.managedBy = t.Links.ManagedBy.ToStrings()
	coolingunit.ManagedByCount = t.Links.ManagedByCount
	coolingunit.AllowedCoolingUnitModes = t.Actions.SetMode.AllowedCoolingUnitModes
	coolingunit.setMode = t.Actions.SetMode.Target

	// This is a read/write object, so we need to save the raw object data for later
	coolingunit.rawData = b

	return nil
}

func (coolingunit *CoolingUnit) SetMode(mode CoolingUnitMode) error {
	// TODO: check if mode is in Allowable values
	properties := map[string]interface{}{
		"Mode": mode,
	}

	return coolingunit.Post(coolingunit.setMode, properties)
}

// Update commits updates to this object's properties to the running system.
func (coolingunit *CoolingUnit) Update() error {
	readWriteFields := []string{
		"AssetTag",
		"UserLabel",
	}

	return coolingunit.UpdateFromRawData(coolingunit, coolingunit.rawData, readWriteFields)
}

// GetCoolingUnit will get a CoolingUnit instance from the service.
func GetCoolingUnit(c common.Client, uri string) (*CoolingUnit, error) {
	return common.GetObject[CoolingUnit](c, uri)
}

// ListReferencedCoolingUnits gets the collection of CoolingUnit from
// a provided reference.
func ListReferencedCoolingUnits(c common.Client, link string) ([]*CoolingUnit, error) {
	return common.GetCollectionObjects[CoolingUnit](c, link)
}

// Assembly gets a collection of assemblies.
func (coolingunit *CoolingUnit) Assembly() ([]*Assembly, error) {
	return ListReferencedAssemblys(coolingunit.GetClient(), coolingunit.assembly)
}

// EnvironmentMetrics gets the environment metrics for this cooling unit.
func (coolingunit *CoolingUnit) EnvironmentMetrics() (*EnvironmentMetrics, error) {
	if coolingunit.environmentMetrics == "" {
		return nil, nil
	}
	return GetEnvironmentMetrics(coolingunit.GetClient(), coolingunit.environmentMetrics)
}

// Filters gets a collection of filters.
func (coolingunit *CoolingUnit) Filters() ([]*Filter, error) {
	return ListReferencedFilters(coolingunit.GetClient(), coolingunit.filters)
}

// LeakDetection gets the of leak detection of this cooling unit.
func (coolingunit *CoolingUnit) LeakDetection() (*LeakDetection, error) {
	if coolingunit.leakDetection == "" {
		return nil, nil
	}
	return GetLeakDetection(coolingunit.GetClient(), coolingunit.leakDetection)
}

// PrimaryCoolantConnectors gets a collection of primary coolant connectors.
func (coolingunit *CoolingUnit) PrimaryCoolantConnectors() ([]*CoolantConnector, error) {
	return ListReferencedCoolantConnectors(coolingunit.GetClient(), coolingunit.primaryCoolantConnectors)
}

// Pumps gets a collection of pumps.
func (coolingunit *CoolingUnit) Pumps() ([]*Pump, error) {
	return ListReferencedPumps(coolingunit.GetClient(), coolingunit.pumps)
}

// Reservoirs gets a collection of reservoirs.
func (coolingunit *CoolingUnit) Reservoirs() ([]*Reservoir, error) {
	return ListReferencedReservoirs(coolingunit.GetClient(), coolingunit.reservoirs)
}

// SecondaryCoolantConnectors gets a collection of secondary coolant connectors.
func (coolingunit *CoolingUnit) SecondaryCoolantConnectors() ([]*CoolantConnector, error) {
	return ListReferencedCoolantConnectors(coolingunit.GetClient(), coolingunit.secondaryCoolantConnectors)
}
