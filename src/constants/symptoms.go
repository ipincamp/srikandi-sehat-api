package constants

type SymptomType string
type MoodType string

const (
	SymptomTypeBasic   SymptomType = "BASIC"
	SymptomTypeOptions SymptomType = "OPTIONS"
)

const (
	MoodTypeHappy   MoodType = "happy"
	MoodTypeNeutral MoodType = "neutral"
	MoodTypeAnxious MoodType = "anxious"
	MoodTypeSad     MoodType = "sad"
	MoodTypeAngry   MoodType = "angry"
)
