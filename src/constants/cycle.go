package constants

// --- Batas Notifikasi Worker ---

// Periode (durasi haid) dianggap 'terlalu panjang' jika
// sedang berlangsung lebih dari N hari.
const (
	CyclePeriodLongThresholdDays = 7
)

// Siklus baru dianggap 'terlambat' jika belum dimulai
// setelah N hari sejak siklus terakhir selesai.
// Ini > CycleLengthMaxNormalDays untuk memberi jeda.
const (
	CycleLateThresholdDays = 32
)

// --- Batas Kategori Normal (untuk UI, Laporan, & Handler) ---
// Digunakan untuk menentukan flag IsPeriodNormal / IsCycleNormal

const (
	CyclePeriodMinNormalDays int16 = 2 // Kurang dari ini = Hipomenorea
	CyclePeriodMaxNormalDays int16 = 7 // Lebih dari ini = Menoragia
)

const (
	CycleLengthMinNormalDays int16 = 21 // Kurang dari ini = Polimenorea
	CycleLengthMaxNormalDays int16 = 35 // Lebih dari ini = Oligomenorea
)
