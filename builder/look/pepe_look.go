package look

type Skin struct {
	Color string `json:"color" datastore:"color"`
}

type Eyes struct {
	EyeColor string `json:"color" datastore:"color"`
	EyeType string `json:"type" datastore:"type"`
}

//Note; can be a hat to
type Hair struct {
	HairColor string `json:"haircolor" datastore:"haircolor"`
	HatColor string `json:"hatcolor" datastore:"hatcolor"`
	HatColor2 string `json:"hatcolor2" datastore:"hatcolor2"`
	HairType string `json:"type" datastore:"type"`
}

type Head struct {
	Hair Hair `json:"hair" datastore:"hair"`
	Eyes Eyes `json:"eyes" datastore:"eyes"`
	Mouth string `json:"mouth" datastore:"mouth"`
}

type Shirt struct {
	ShirtColor string `json:"color" datastore:"color"`
	ShirtType string `json:"type" datastore:"type"`
}

type Body struct {
	Neck string `json:"neck" datastore:"neck"`
	Shirt Shirt `json:"shirt" datastore:"shirt"`
}

type Glasses struct {
	PrimaryColor string `json:"primary" datastore:"primary"`
	SecondaryColor string `json:"secondary" datastore:"secondary"`
	GlassesType string `json:"type" datastore:"type"`
}

type Extra struct {
	Glasses Glasses `json:"glasses" datastore:"glasses"`
}

type PepeLook struct {
	Skin Skin `json:"skin" datastore:"skin"`
	Head Head `json:"head" datastore:"head"`
	Body Body `json:"body" datastore:"body"`
	Extra Extra `json:"extra" datastore:"extra"`
	BackgroundColor string `json:"background" datastore:"background"`
}

