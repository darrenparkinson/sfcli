package salesforce

import (
	"context"
	"fmt"
	"net/http"
)

// DescribeResponse represents the fields for describing an object.
// Note that not all fields are provided here.  Feel free to add extra fields if you need them.
type DescribeResponse struct {
	Fields []struct {
		Aggregatable             bool        `json:"aggregatable"`
		AiPredictionField        bool        `json:"aiPredictionField"`
		AutoNumber               bool        `json:"autoNumber"`
		ByteLength               int         `json:"byteLength"`
		Calculated               bool        `json:"calculated"`
		CalculatedFormula        string      `json:"calculatedFormula"`
		CascadeDelete            bool        `json:"cascadeDelete"`
		CaseSensitive            bool        `json:"caseSensitive"`
		CompoundFieldName        string      `json:"compoundFieldName"`
		ControllerName           string      `json:"controllerName"`
		Createable               bool        `json:"createable"`
		Custom                   bool        `json:"custom"`
		DefaultValue             interface{} `json:"defaultValue"`
		DefaultValueFormula      string      `json:"defaultValueFormula"`
		DefaultedOnCreate        bool        `json:"defaultedOnCreate"`
		DependentPicklist        bool        `json:"dependentPicklist"`
		DeprecatedAndHidden      bool        `json:"deprecatedAndHidden"`
		Digits                   int         `json:"digits"`
		DisplayLocationInDecimal bool        `json:"displayLocationInDecimal"`
		Encrypted                bool        `json:"encrypted"`
		ExternalID               bool        `json:"externalId"`
		ExtraTypeInfo            string      `json:"extraTypeInfo"`
		Filterable               bool        `json:"filterable"`
		// FilteredLookupInfo           interface{} `json:"filteredLookupInfo"`
		FormulaTreatNullNumberAsZero bool   `json:"formulaTreatNullNumberAsZero"`
		Groupable                    bool   `json:"groupable"`
		HighScaleNumber              bool   `json:"highScaleNumber"`
		HTMLFormatted                bool   `json:"htmlFormatted"`
		IDLookup                     bool   `json:"idLookup"`
		InlineHelpText               string `json:"inlineHelpText"`
		Label                        string `json:"label"`
		Length                       int    `json:"length"`
		// Mask                         interface{} `json:"mask"`
		// MaskType                     interface{} `json:"maskType"`
		Name           string `json:"name"`
		NameField      bool   `json:"nameField"`
		NamePointing   bool   `json:"namePointing"`
		Nillable       bool   `json:"nillable"`
		Permissionable bool   `json:"permissionable"`
		PicklistValues []struct {
			Active       bool   `json:"active"`
			DefaultValue bool   `json:"defaultValue"`
			Label        string `json:"label"`
			ValidFor     string `json:"validFor"`
			Value        string `json:"value"`
		} `json:"picklistValues"`
		PolymorphicForeignKey bool `json:"polymorphicForeignKey"`
		Precision             int  `json:"precision"`
		QueryByDistance       bool `json:"queryByDistance"`
		// ReferenceTargetField    interface{}   `json:"referenceTargetField"`
		ReferenceTo      []string `json:"referenceTo"`
		RelationshipName string   `json:"relationshipName"`
		// RelationshipOrder       interface{} `json:"relationshipOrder"`
		RestrictedDelete        bool   `json:"restrictedDelete"`
		RestrictedPicklist      bool   `json:"restrictedPicklist"`
		Scale                   int    `json:"scale"`
		SearchPrefilterable     bool   `json:"searchPrefilterable"`
		SoapType                string `json:"soapType"`
		Sortable                bool   `json:"sortable"`
		Type                    string `json:"type"`
		Unique                  bool   `json:"unique"`
		Updateable              bool   `json:"updateable"`
		WriteRequiresMasterRead bool   `json:"writeRequiresMasterRead"`
	} `json:"fields"`
	IsSubtype       bool   `json:"isSubtype"`
	KeyPrefix       string `json:"keyPrefix"`
	Label           string `json:"label"`
	LabelPlural     string `json:"labelPlural"`
	Name            string `json:"name"`
	Queryable       bool   `json:"queryable"`
	RecordTypeInfos []struct {
		Active                   bool   `json:"active"`
		Available                bool   `json:"available"`
		DefaultRecordTypeMapping bool   `json:"defaultRecordTypeMapping"`
		DeveloperName            string `json:"developerName"`
		Master                   bool   `json:"master"`
		Name                     string `json:"name"`
		RecordTypeID             string `json:"recordTypeId"`
		Urls                     struct {
			Layout string `json:"layout"`
		} `json:"urls"`
	} `json:"recordTypeInfos"`
	Replicateable         bool   `json:"replicateable"`
	Retrieveable          bool   `json:"retrieveable"`
	SearchLayoutable      bool   `json:"searchLayoutable"`
	Searchable            bool   `json:"searchable"`
	SobjectDescribeOption string `json:"sobjectDescribeOption"`
	SupportedScopes       []struct {
		Label string `json:"label"`
		Name  string `json:"name"`
	} `json:"supportedScopes"`
	Triggerable bool `json:"triggerable"`
	Undeletable bool `json:"undeletable"`
	Updateable  bool `json:"updateable"`
	Urls        struct {
		CompactLayouts   string `json:"compactLayouts"`
		RowTemplate      string `json:"rowTemplate"`
		ApprovalLayouts  string `json:"approvalLayouts"`
		UIDetailTemplate string `json:"uiDetailTemplate"`
		UIEditTemplate   string `json:"uiEditTemplate"`
		Listviews        string `json:"listviews"`
		Describe         string `json:"describe"`
		UINewRecord      string `json:"uiNewRecord"`
		QuickActions     string `json:"quickActions"`
		Layouts          string `json:"layouts"`
		Sobject          string `json:"sobject"`
	} `json:"urls"`
}

func (s *AccountService) Describe(ctx context.Context) (*DescribeResponse, error) {
	sfurl := fmt.Sprintf("%s/services/data/%s/sobjects/Account/describe", s.client.BaseURL, s.client.Version)
	req, err := http.NewRequest("GET", sfurl, nil)
	if err != nil {
		return nil, err
	}
	var dr DescribeResponse
	if err := s.client.makeRequest(ctx, req, &dr); err != nil {
		return nil, err
	}
	return &dr, nil
}

func (s *ContactService) Describe(ctx context.Context) (*DescribeResponse, error) {
	sfurl := fmt.Sprintf("%s/services/data/%s/sobjects/Contact/describe", s.client.BaseURL, s.client.Version)
	req, err := http.NewRequest("GET", sfurl, nil)
	if err != nil {
		return nil, err
	}
	var dr DescribeResponse
	if err := s.client.makeRequest(ctx, req, &dr); err != nil {
		return nil, err
	}
	return &dr, nil
}

func (s *OpportunityService) Describe(ctx context.Context) (*DescribeResponse, error) {
	sfurl := fmt.Sprintf("%s/services/data/%s/sobjects/Opportunity/describe", s.client.BaseURL, s.client.Version)
	req, err := http.NewRequest("GET", sfurl, nil)
	if err != nil {
		return nil, err
	}
	var dr DescribeResponse
	if err := s.client.makeRequest(ctx, req, &dr); err != nil {
		return nil, err
	}
	return &dr, nil
}

func (s *UserService) Describe(ctx context.Context) (*DescribeResponse, error) {
	sfurl := fmt.Sprintf("%s/services/data/%s/sobjects/User/describe", s.client.BaseURL, s.client.Version)
	req, err := http.NewRequest("GET", sfurl, nil)
	if err != nil {
		return nil, err
	}
	var dr DescribeResponse
	if err := s.client.makeRequest(ctx, req, &dr); err != nil {
		return nil, err
	}
	return &dr, nil
}
