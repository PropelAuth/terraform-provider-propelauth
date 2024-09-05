package propelauth

func CreateApiKeyExpirationOptions(options []string) ApiKeyExpirationOptions {
	var expirationOptions ApiKeyExpirationOptions
	for _, option := range options {
		switch option {
		case "TwoWeeks":
			expirationOptions.TwoWeeks = true
		case "OneMonth":
			expirationOptions.OneMonth = true
		case "ThreeMonths":
			expirationOptions.ThreeMonths = true
		case "SixMonths":
			expirationOptions.SixMonths = true
		case "OneYear":
			expirationOptions.OneYear = true
		case "Never":
			expirationOptions.Never = true
		}
	}
	return expirationOptions
}

func (aks ApiKeyExpirationOptionSettings) GetApiKeyExpirationOptions() []string {
	var options []string
	if aks.Options.TwoWeeks {
		options = append(options, "TwoWeeks")
	}
	if aks.Options.OneMonth {
		options = append(options, "OneMonth")
	}
	if aks.Options.ThreeMonths {
		options = append(options, "ThreeMonths")
	}
	if aks.Options.SixMonths {
		options = append(options, "SixMonths")
	}
	if aks.Options.OneYear {
		options = append(options, "OneYear")
	}
	if aks.Options.Never {
		options = append(options, "Never")
	}
	return options
}
