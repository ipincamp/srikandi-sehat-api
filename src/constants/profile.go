package constants

type EducationLevel string
type InternetAccess string

const (
	EduNone    EducationLevel = "Tidak Sekolah"
	EduSD      EducationLevel = "SD"
	EduSMP     EducationLevel = "SMP"
	EduSMA     EducationLevel = "SMA"
	EduDiploma EducationLevel = "Diploma"
	EduS1      EducationLevel = "S1"
	EduS2      EducationLevel = "S2"
	EduS3      EducationLevel = "S3"
)

const (
	AccessWiFi     InternetAccess = "WiFi"
	AccessCellular InternetAccess = "Seluler"
)
