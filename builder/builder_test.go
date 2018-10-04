package builder

import (
	"testing"
	"bytes"
	"fmt"
	"cryptopepe.io/cryptopepe-svg/builder/look"
)

func TestSVGBuilder_ConvertToSVG(t *testing.T) {
	builder := new(SVGBuilder)
	builder.Load()

	buf := new(bytes.Buffer)
	pepeLook := look.PepeLook{
		Skin: look.Skin{
			Color: "#658749",
		},
		Head: look.Head{
			Hair: look.Hair{
				HairType: "hair>trump_hair",
				HairColor: "#ff0000",
				HatColor: "#00ff00",
				HatColor2: "#0000ff",
			},
			Eyes:  look.Eyes{
				EyeColor: "#220010",
				EyeType: "eyes>colored_eyes",
			},
			Mouth: "mouth>happy_lips",
		},
		Body: look.Body{
			Neck: "neck>dollar_necklace",
			Shirt: look.Shirt{
				ShirtColor: "#0000aa",
				ShirtType:  "shirt>basic_shirt",
			},
		},
		Extra: look.Extra{
			Glasses: look.Glasses{
				PrimaryColor:   "#225500",
				SecondaryColor: "#202050",
				GlassesType:    "glasses>sunglasses_2",
			},
		},
	}
	err := builder.ConvertToSVG(buf, &pepeLook)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())

}
