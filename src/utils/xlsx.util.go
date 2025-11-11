package utils

import (
	"crypto/rand"
	"fmt"
	"ipincamp/srikandi-sehat/src/dto"
	"time"

	"github.com/xuri/excelize/v2"
)

// GenerateReportXLSX creates an XLSX file from report records
func GenerateReportXLSX(records []dto.FullExportRecord, filePath string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			ErrorLogger.Printf("Failed to close Excel file: %v", err)
		}
	}()

	sheetName := "Report"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create sheet: %w", err)
	}

	// Set active sheet
	f.SetActiveSheet(index)

	// Define headers
	headers := []string{
		"Nama Pengguna", "Email", "Tanggal Registrasi", "Umur", "No. Telepon",
		"Tinggi (cm)", "Berat (kg)", "IMT", "Kategori IMT", "Usia Menarche",
		"Pendidikan Terakhir", "Pendidikan Ortu", "Pekerjaan Ortu", "Akses Internet",
		"Desa/Kelurahan", "Kecamatan", "Kabupaten/Kota", "Provinsi", "Klasifikasi Alamat",
		"Siklus Ke-", "Tanggal Mulai", "Tanggal Selesai", "Lama Haid (Hari)",
		"Kategori Lama Haid", "Panjang Siklus (Hari)", "Kategori Panjang Siklus", "Gejala yang Dirasakan",
	}

	// Create header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   11,
			Color:  "FFFFFF",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create header style: %w", err)
	}

	// Write headers
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Set column widths
	columnWidths := map[string]float64{
		"A": 20, "B": 25, "C": 20, "D": 8, "E": 15,
		"F": 12, "G": 12, "H": 8, "I": 20, "J": 15,
		"K": 20, "L": 20, "M": 20, "N": 15, "O": 20,
		"P": 20, "Q": 20, "R": 15, "S": 20, "T": 10,
		"U": 15, "V": 15, "W": 18, "X": 22, "Y": 18,
		"Z": 25, "AA": 40,
	}
	for col, width := range columnWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	// Create data style
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "D3D3D3", Style: 1},
			{Type: "top", Color: "D3D3D3", Style: 1},
			{Type: "bottom", Color: "D3D3D3", Style: 1},
			{Type: "right", Color: "D3D3D3", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: false,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create data style: %w", err)
	}

	// Write data
	for i, rec := range records {
		row := i + 2
		values := []interface{}{
			rec.UserName,
			rec.UserEmail,
			rec.UserRegisteredAt.Format("2006-01-02 15:04:05"),
			rec.Age,
			rec.PhoneNumber,
			rec.HeightCM,
			rec.WeightKG,
			rec.BMI,
			rec.BMICategory,
			rec.MenarcheAge,
			rec.LastEducation,
			rec.ParentLastEducation,
			rec.ParentLastJob,
			rec.InternetAccess,
			rec.Village,
			rec.District,
			rec.Regency,
			rec.Province,
			rec.Classification,
			rec.CycleNumber,
			rec.StartDate,
			rec.EndDate,
			rec.PeriodLength,
			rec.PeriodCategory,
			rec.CycleLength,
			rec.CycleCategory,
			rec.Symptoms,
		}

		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheetName, cell, value)
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}
	}

	// Set row height for header
	f.SetRowHeight(sheetName, 1, 30)

	// Enable auto-filter
	lastCol, _ := excelize.ColumnNumberToName(len(headers))
	f.AutoFilter(sheetName, fmt.Sprintf("A1:%s1", lastCol), []excelize.AutoFilterOptions{})

	// Freeze first row
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// Delete default Sheet1 if it exists and is not our sheet
	if sheetName != "Sheet1" {
		f.DeleteSheet("Sheet1")
	}

	// Save the file
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}

// GenerateRandomPassword generates a random 12-character password
func GenerateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, 12)
	for i := range b {
		// Use crypto/rand for secure random generation
		randomByte := make([]byte, 1)
		_, err := rand.Read(randomByte)
		if err != nil {
			// Fallback to timestamp-based if crypto/rand fails
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		} else {
			b[i] = charset[int(randomByte[0])%len(charset)]
		}
	}
	return string(b)
}
