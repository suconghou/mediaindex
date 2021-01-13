package ebml

type eleType struct {
	name  string
	vtype string
}

var types = map[string]*eleType{
	"1c53bb6b": {"Cues", "m"},               // lvl. 1
	"bb":       {"CuePoint", "m"},           // lvl. 2
	"b3":       {"CueTime", "u"},            // lvl. 3
	"b7":       {"CueTrackPositions", "m"},  // lvl. 3
	"f7":       {"CueTrack", "u"},           // lvl. 4
	"f1":       {"CueClusterPosition", "u"}, // lvl. 4
}

func getOne(id string) *eleType {
	return types[id]
}
