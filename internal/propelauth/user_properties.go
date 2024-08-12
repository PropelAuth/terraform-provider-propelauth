package propelauth

import (
	"encoding/json"
)

// GetUserProperties - Returns current user properties settings
func (c *PropelAuthClient) GetUserProperties() (*UserProperties, error) {
	res, err := c.get("user_property_settings", nil)
	if err != nil {
		return nil, err
	}

	userProperties := UserProperties{}
	err = json.Unmarshal(res.BodyBytes, &userProperties)
	if err != nil {
		return nil, err
	}

	return &userProperties, nil
}

// UpdateUserProperties - Updates the user properties settings
func (c *PropelAuthClient) UpdateUserProperties(userProperties *UserProperties) (*UserProperties, error) {
	body, err := json.Marshal(userProperties)
	if err != nil {
		return nil, err
	}

	_, err = c.put("user_property_settings", body)
	if err != nil {
		return nil, err
	}

	return c.GetUserProperties()
}

func (up *UserProperties) defaultPropertyEnabled(propertyName string) bool {
	for i := range up.Fields {
		if up.Fields[i].Name == propertyName {
			return up.Fields[i].IsEnabled
		}
	}
	return false
}

func (up *UserProperties) disableDefaultProperty(propertyName string) {
	for i := range up.Fields {
		if up.Fields[i].Name == propertyName {
			up.Fields[i].IsEnabled = false
		}
	}
}

type NamePropertySettings struct {
	InJwt bool
}

// UpdateNameProperty - Updates the name property and sets it to enabled
func (up *UserProperties) UpdateAndEnableNameProperty(nameProperty NamePropertySettings) {
	for i := range up.Fields {
		if up.Fields[i].Name == "legacy__name" {
			up.Fields[i].IsEnabled = true
			up.Fields[i].InJwt = nameProperty.InJwt
		}
	}
}

// DisableNameProperty - Disables the name property
func (up *UserProperties) DisableNameProperty() {
	up.disableDefaultProperty("legacy__name")
}

// NamePropertyEnabled - Returns true if the name property is enabled
func (up *UserProperties) NamePropertyEnabled() bool {
	return up.defaultPropertyEnabled("legacy__name")
}

// GetNamePropertySettings - Returns the name property settings
func (up *UserProperties) GetNamePropertySettings() NamePropertySettings {
	for i := range up.Fields {
		if up.Fields[i].Name == "legacy__name" {
			return NamePropertySettings{
				InJwt: up.Fields[i].InJwt,
			}
		}
	}
	return NamePropertySettings{}
}

type MetadataPropertySettings struct {
	InJwt bool
	CollectViaSaml bool
}

// UpdateMetadataProperty - Updates the metadata property and sets it to enabled
func (up *UserProperties) UpdateAndEnableMetadataProperty(metadataProperty MetadataPropertySettings) {
	for i := range up.Fields {
		if up.Fields[i].Name == "metadata" {
			up.Fields[i].IsEnabled = true
			up.Fields[i].InJwt = metadataProperty.InJwt
			up.Fields[i].CollectViaSaml = metadataProperty.CollectViaSaml
		}
	}
}

// DisableMetadataProperty - Disables the metadata property
func (up *UserProperties) DisableMetadataProperty() {
	up.disableDefaultProperty("metadata")
}

// MetadataPropertyEnabled - Returns true if the metadata property is enabled
func (up *UserProperties) MetadataPropertyEnabled() bool {
	return up.defaultPropertyEnabled("metadata")
}

// GetMetadataPropertySettings - Returns the metadata property settings
func (up *UserProperties) GetMetadataPropertySettings() MetadataPropertySettings {
	for i := range up.Fields {
		if up.Fields[i].Name == "metadata" {
			return MetadataPropertySettings{
				InJwt: up.Fields[i].InJwt,
				CollectViaSaml: up.Fields[i].CollectViaSaml,
			}
		}
	}
	return MetadataPropertySettings{}
}

type UsernamePropertySettings struct {
	InJwt bool
	DisplayName string
}

// UpdateUsernameProperty - Updates the username property and sets it to enabled
func (up *UserProperties) UpdateAndEnableUsernameProperty(usernameProperty UsernamePropertySettings) {
	for i := range up.Fields {
		if up.Fields[i].Name == "legacy__username" {
			up.Fields[i].IsEnabled = true
			up.Fields[i].InJwt = usernameProperty.InJwt
			up.Fields[i].DisplayName = usernameProperty.DisplayName
		}
	}
}

// DisableUsernameProperty - Disables the username property
func (up *UserProperties) DisableUsernameProperty() {
	up.disableDefaultProperty("legacy__username")
}

// UsernamePropertyEnabled - Returns true if the username property is enabled
func (up *UserProperties) UsernamePropertyEnabled() bool {
	return up.defaultPropertyEnabled("legacy__username")
}

// GetUsernamePropertySettings - Returns the username property settings
func (up *UserProperties) GetUsernamePropertySettings() UsernamePropertySettings {
	for i := range up.Fields {
		if up.Fields[i].Name == "legacy__username" {
			return UsernamePropertySettings{
				InJwt: up.Fields[i].InJwt,
				DisplayName: up.Fields[i].DisplayName,
			}
		}
	}
	return UsernamePropertySettings{}
}

type PictureUrlPropertySettings struct {
	InJwt bool
}

// UpdatePictureUrlProperty - Updates the picture url property and sets it to enabled
func (up *UserProperties) UpdateAndEnablePictureUrlProperty(pictureUrlProperty PictureUrlPropertySettings) {
	for i := range up.Fields {
		if up.Fields[i].Name == "legacy__picture_url" {
			up.Fields[i].IsEnabled = true
			up.Fields[i].InJwt = pictureUrlProperty.InJwt
		}
	}
}

// DisablePictureUrlProperty - Disables the picture url property
func (up *UserProperties) DisablePictureUrlProperty() {
	up.disableDefaultProperty("legacy__picture_url")
}

// PictureUrlPropertyEnabled - Returns true if the picture url property is enabled
func (up *UserProperties) PictureUrlPropertyEnabled() bool {
	return up.defaultPropertyEnabled("legacy__picture_url")
}

// GetPictureUrlPropertySettings - Returns the picture url property settings
func (up *UserProperties) GetPictureUrlPropertySettings() PictureUrlPropertySettings {
	for i := range up.Fields {
		if up.Fields[i].Name == "legacy__picture_url" {
			return PictureUrlPropertySettings{
				InJwt: up.Fields[i].InJwt,
			}
		}
	}
	return PictureUrlPropertySettings{}
}

type PhoneNumberPropertySettings struct {
	InJwt bool
	DisplayName string
	ShowInAccount bool
	CollectViaSaml bool
	Required bool
	RequiredBy int64
	UserWritable string
}

// UpdatePhoneNumberProperty - Updates the phone number property and sets it to enabled
func (up *UserProperties) UpdateAndEnablePhoneNumberProperty(phoneNumberProperty PhoneNumberPropertySettings) {
	for i := range up.Fields {
		if up.Fields[i].Name == "phone_number" {
			up.Fields[i].IsEnabled = true
			up.Fields[i].InJwt = phoneNumberProperty.InJwt
			up.Fields[i].DisplayName = phoneNumberProperty.DisplayName
			up.Fields[i].ShowInAccount = phoneNumberProperty.ShowInAccount
			up.Fields[i].CollectViaSaml = phoneNumberProperty.CollectViaSaml
			up.Fields[i].Required = phoneNumberProperty.Required
			up.Fields[i].RequiredBy = phoneNumberProperty.RequiredBy
			up.Fields[i].UserWritable = phoneNumberProperty.UserWritable
		}
	}
}

// DisablePhoneNumberProperty - Disables the phone number property
func (up *UserProperties) DisablePhoneNumberProperty() {
	up.disableDefaultProperty("phone_number")
}

// PhoneNumberPropertyEnabled - Returns true if the phone number property is enabled
func (up *UserProperties) PhoneNumberPropertyEnabled() bool {
	return up.defaultPropertyEnabled("phone_number")
}

// GetPhoneNumberPropertySettings - Returns the phone number property settings
func (up *UserProperties) GetPhoneNumberPropertySettings() PhoneNumberPropertySettings {
	for i := range up.Fields {
		if up.Fields[i].Name == "phone_number" {
			return PhoneNumberPropertySettings{
				InJwt: up.Fields[i].InJwt,
				DisplayName: up.Fields[i].DisplayName,
				ShowInAccount: up.Fields[i].ShowInAccount,
				CollectViaSaml: up.Fields[i].CollectViaSaml,
				Required: up.Fields[i].Required,
				RequiredBy: up.Fields[i].RequiredBy,
				UserWritable: up.Fields[i].UserWritable,
			}
		}
	}
	return PhoneNumberPropertySettings{}
}

type TosPropertySettings struct {
	InJwt bool
	Required bool
	RequiredBy int64
	TosLinks []TosLink
}

// UpdateTosProperty - Updates the TOS property and sets it to enabled
func (up *UserProperties) UpdateAndEnableTosProperty(tosProperty TosPropertySettings) {
	for i := range up.Fields {
		if up.Fields[i].Name == "tos" {
			up.Fields[i].IsEnabled = true
			up.Fields[i].InJwt = tosProperty.InJwt
			up.Fields[i].Required = tosProperty.Required
			up.Fields[i].RequiredBy = tosProperty.RequiredBy
			up.Fields[i].Metadata = userPropertyMetadata{
				TosLinks: tosProperty.TosLinks,
			}
		}
	}
}

// DisableTosProperty - Disables the TOS property
func (up *UserProperties) DisableTosProperty() {
	up.disableDefaultProperty("tos")
}

// TosPropertyEnabled - Returns true if the TOS property is enabled
func (up *UserProperties) TosPropertyEnabled() bool {
	return up.defaultPropertyEnabled("tos")
}

// GetTosPropertySettings - Returns the TOS property settings
func (up *UserProperties) GetTosPropertySettings() TosPropertySettings {
	for i := range up.Fields {
		if up.Fields[i].Name == "tos" {
			return TosPropertySettings{
				InJwt: up.Fields[i].InJwt,
				Required: up.Fields[i].Required,
				RequiredBy: up.Fields[i].RequiredBy,
				TosLinks: up.Fields[i].Metadata.TosLinks,
			}
		}
	}
	return TosPropertySettings{}
}

type ReferralSourcePropertySettings struct {
	InJwt bool
	DisplayName string
	Required bool
	RequiredBy int64
	UserWritable string
	Options []string
	ShowInAccount bool
	CollectViaSaml bool
}

// UpdateReferralSourceProperty - Updates the referral source property and sets it to enabled
func (up *UserProperties) UpdateAndEnableReferralSourceProperty(referralSourceProperty ReferralSourcePropertySettings) {
	for i := range up.Fields {
		if up.Fields[i].Name == "referral_source" {
			up.Fields[i].IsEnabled = true
			up.Fields[i].InJwt = referralSourceProperty.InJwt
			up.Fields[i].DisplayName = referralSourceProperty.DisplayName
			up.Fields[i].Required = referralSourceProperty.Required
			up.Fields[i].RequiredBy = referralSourceProperty.RequiredBy
			up.Fields[i].UserWritable = referralSourceProperty.UserWritable
			up.Fields[i].Metadata = userPropertyMetadata{
				EnumValues: referralSourceProperty.Options,
			}
			up.Fields[i].ShowInAccount = referralSourceProperty.ShowInAccount
			up.Fields[i].CollectViaSaml = referralSourceProperty.CollectViaSaml
		}
	}
}

// DisableReferralSourceProperty - Disables the referral source property
func (up *UserProperties) DisableReferralSourceProperty() {
	up.disableDefaultProperty("referral_source")
}

// ReferralSourcePropertyEnabled - Returns true if the referral source property is enabled
func (up *UserProperties) ReferralSourcePropertyEnabled() bool {
	return up.defaultPropertyEnabled("referral_source")
}

// GetReferralSourcePropertySettings - Returns the referral source property settings
func (up *UserProperties) GetReferralSourcePropertySettings() ReferralSourcePropertySettings {
	for i := range up.Fields {
		if up.Fields[i].Name == "referral_source" {
			return ReferralSourcePropertySettings{
				InJwt: up.Fields[i].InJwt,
				DisplayName: up.Fields[i].DisplayName,
				Required: up.Fields[i].Required,
				RequiredBy: up.Fields[i].RequiredBy,
				UserWritable: up.Fields[i].UserWritable,
				Options: up.Fields[i].Metadata.EnumValues,
				ShowInAccount: up.Fields[i].ShowInAccount,
				CollectViaSaml: up.Fields[i].CollectViaSaml,
			}
		}
	}
	return ReferralSourcePropertySettings{}
}

type CustomPropertySettings struct {
	Name string
	DisplayName string
	FieldType string
	Required bool
	RequiredBy int64
	InJwt bool
	IsUserFacing bool
	CollectOnSignup bool
	CollectViaSaml bool
	ShowInAccount bool
	UserWritable string
	EnumValues []string
}

func (c *CustomPropertySettings) IsEqual(other CustomPropertySettings) bool {
	if c.Name != other.Name ||
		c.DisplayName != other.DisplayName ||
		c.FieldType != other.FieldType ||
		c.Required != other.Required ||
		c.RequiredBy != other.RequiredBy ||
		c.InJwt != other.InJwt ||
		c.IsUserFacing != other.IsUserFacing ||
		c.CollectOnSignup != other.CollectOnSignup ||
		c.CollectViaSaml != other.CollectViaSaml ||
		c.ShowInAccount != other.ShowInAccount ||
		c.UserWritable != other.UserWritable {
		return false
	}
	if len(c.EnumValues) != len(other.EnumValues) {
		return false
	}
	for i := range c.EnumValues {
		if c.EnumValues[i] != other.EnumValues[i] {
			return false
		}
	}
	return true
}


// UpsertCustomProperty - Upserts a custom property
func (up *UserProperties) UpsertCustomProperty(customProperty CustomPropertySettings) {
	for i := range up.Fields {
		if up.Fields[i].Name == customProperty.Name {
			up.Fields[i].DisplayName = customProperty.DisplayName
			up.Fields[i].FieldType = customProperty.FieldType
			up.Fields[i].Required = customProperty.Required
			up.Fields[i].RequiredBy = customProperty.RequiredBy
			up.Fields[i].InJwt = customProperty.InJwt
			up.Fields[i].IsUserFacing = customProperty.IsUserFacing
			up.Fields[i].CollectOnSignup = customProperty.CollectOnSignup
			up.Fields[i].CollectViaSaml = customProperty.CollectViaSaml
			up.Fields[i].ShowInAccount = customProperty.ShowInAccount
			up.Fields[i].UserWritable = customProperty.UserWritable
			up.Fields[i].Metadata = userPropertyMetadata{
				EnumValues: customProperty.EnumValues,
			}
			up.Fields[i].IsEnabled = true
			return
		}
	}
	up.Fields = append(up.Fields, UserProperty{
		Name: customProperty.Name,
		DisplayName: customProperty.DisplayName,
		FieldType: customProperty.FieldType,
		Required: customProperty.Required,
		RequiredBy: customProperty.RequiredBy,
		InJwt: customProperty.InJwt,
		IsUserFacing: customProperty.IsUserFacing,
		CollectOnSignup: customProperty.CollectOnSignup,
		CollectViaSaml: customProperty.CollectViaSaml,
		ShowInAccount: customProperty.ShowInAccount,
		UserWritable: customProperty.UserWritable,
		Metadata: userPropertyMetadata{
			EnumValues: customProperty.EnumValues,
		},
		IsEnabled: true,
	})
}

// DisableDroppedCustomProperties - Disables custom properties that are not in the provided list and are not one of the default properties
func (up *UserProperties) DisableDroppedCustomProperties(customProperties []CustomPropertySettings) {
	for i := range up.Fields {
		if !containsName(customProperties, up.Fields[i].Name) && !IsDefaultPropertyName(up.Fields[i].Name) {
			up.Fields[i].IsEnabled = false
		}
	}
}

// GetEnabledCustomProperties - Returns a list of enabled custom properties
func (up *UserProperties) GetEnabledCustomProperties() []CustomPropertySettings {
	var enabledCustomProperties []CustomPropertySettings
	for i := range up.Fields {
		if !IsDefaultPropertyName(up.Fields[i].Name) && up.Fields[i].IsEnabled {
			enabledCustomProperties = append(enabledCustomProperties, CustomPropertySettings{
				Name: up.Fields[i].Name,
				DisplayName: up.Fields[i].DisplayName,
				FieldType: up.Fields[i].FieldType,
				Required: up.Fields[i].Required,
				RequiredBy: up.Fields[i].RequiredBy,
				InJwt: up.Fields[i].InJwt,
				IsUserFacing: up.Fields[i].IsUserFacing,
				CollectOnSignup: up.Fields[i].CollectOnSignup,
				CollectViaSaml: up.Fields[i].CollectViaSaml,
				ShowInAccount: up.Fields[i].ShowInAccount,
				UserWritable: up.Fields[i].UserWritable,
				EnumValues: up.Fields[i].Metadata.EnumValues,
			})
		}
	}
	return enabledCustomProperties
}

// GetCustomPropertySettings - Returns the settings for a custom property
func (up *UserProperties) GetEnabledCustomProperty(propertyName string) (CustomPropertySettings, bool) {
	if IsDefaultPropertyName(propertyName) {
		return CustomPropertySettings{}, false
	}
	for i := range up.Fields {
		if up.Fields[i].Name == propertyName && up.Fields[i].IsEnabled {
			return CustomPropertySettings{
				Name: up.Fields[i].Name,
				DisplayName: up.Fields[i].DisplayName,
				FieldType: up.Fields[i].FieldType,
				Required: up.Fields[i].Required,
				RequiredBy: up.Fields[i].RequiredBy,
				InJwt: up.Fields[i].InJwt,
				IsUserFacing: up.Fields[i].IsUserFacing,
				CollectOnSignup: up.Fields[i].CollectOnSignup,
				CollectViaSaml: up.Fields[i].CollectViaSaml,
				ShowInAccount: up.Fields[i].ShowInAccount,
				UserWritable: up.Fields[i].UserWritable,
				EnumValues: up.Fields[i].Metadata.EnumValues,
			}, true
		}
	}
	return CustomPropertySettings{}, false
}

// GetHangingCustomProperties - Returns a list of custom properties that are enabled but not in the provided list
func (up *UserProperties) GetHangingCustomProperties(customPropertiesInState []string) []CustomPropertySettings {
	var hangingCustomProperties []CustomPropertySettings
	for i := range up.Fields {
		if !Contains(customPropertiesInState, up.Fields[i].Name) && !IsDefaultPropertyName(up.Fields[i].Name) && up.Fields[i].IsEnabled {
			hangingCustomProperties = append(hangingCustomProperties, CustomPropertySettings{
				Name: up.Fields[i].Name,
				DisplayName: up.Fields[i].DisplayName,
				FieldType: up.Fields[i].FieldType,
				Required: up.Fields[i].Required,
				RequiredBy: up.Fields[i].RequiredBy,
				InJwt: up.Fields[i].InJwt,
				IsUserFacing: up.Fields[i].IsUserFacing,
				CollectOnSignup: up.Fields[i].CollectOnSignup,
				CollectViaSaml: up.Fields[i].CollectViaSaml,
				ShowInAccount: up.Fields[i].ShowInAccount,
				UserWritable: up.Fields[i].UserWritable,
				EnumValues: up.Fields[i].Metadata.EnumValues,
			})
		}
	}
	return hangingCustomProperties
}


// internal helper functions

func containsName(slice []CustomPropertySettings, target string) bool {
    for _, s := range slice {
        if s.Name == target {
            return true
        }
    }
    return false
}

func IsDefaultPropertyName(name string) bool {
	defaultProperties := []string{
		"legacy__name",
		"metadata",
		"legacy__username",
		"legacy__picture_url",
		"phone_number",
		"tos",
		"referral_source",
	}
	for _, s := range defaultProperties {
		if s == name {
			return true
		}
	}
	return false
}