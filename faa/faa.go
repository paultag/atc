package faa

type Model struct {
	Id           string
	Manufacturer string
	Name         string
}

var (
	Glider                = "1"
	Balloon               = "2"
	Blimp                 = "3"
	FixedWingSingleEngine = "4"
	FixedWingMultiEngine  = "5"
	Rotorcraft            = "6"
	WeightShiftControl    = "7"
	PoweredParachute      = "8"
	Gryroplane            = "9"
)

type Aircraft struct {
	Id         string
	Model      Model
	TailNumber string
	Type       string
	Status     string

	RegistrantType    string
	RegistrantName    string
	RegistrantStreet  string
	RegistrantStreet2 string
	RegistrantCity    string
	RegistrantState   string
	RegistrantZipcode string
	RegistrantRegion  string
	RegistrantCounty  string
	RegistrantCountry string

	// LAST ACTION DATE
	// CERT ISSUE DATE
	// CERTIFICATION
	// FRACT OWNER
	// AIR WORTH DATE
	// EXPIRATION DATE
}
