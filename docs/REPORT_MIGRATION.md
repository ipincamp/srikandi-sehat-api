# Report System Migration - CSV to XLSX with Password Protection

## Overview
The reporting system has been completely redesigned to use XLSX format with password-protected downloads and encrypted tokens.

## New Flow

### 1. Generate Report (Admin Only)
**Endpoint:** `POST /api/admin/reports/generate`

**Process:**
- Admin requests report generation
- System queries all user and cycle data
- Generates XLSX file and saves it to `./reports/` directory
- Creates a random 12-character password
- Encrypts the internal token with the password
- Stores metadata in cache with 30-minute expiration
- Returns encrypted token and password to admin

**Response:**
```json
{
  "success": true,
  "message": "Report generated successfully. Use the encrypted token and password to download. Password: AbC123!@#XyZ",
  "data": {
    "download_url": "https://your-domain.com/api/reports/download",
    "encrypted_token": "base64_encrypted_token_here",
    "expires_at": "2025-11-03T02:00:00Z"
  }
}
```

### 2. Download Report
**Endpoint:** `POST /api/reports/download`

**Headers:**
```
X-Report-Token: <encrypted_token>
```

**Body:**
```json
{
  "password": "AbC123!@#XyZ"
}
```

**Process:**
- Validates encrypted token exists in cache
- Checks if report has already been downloaded (one-time use)
- Verifies token hasn't expired
- Decrypts token using provided password
- Validates decrypted token matches stored token
- Sends XLSX file
- Marks token as used (prevents reuse)
- Schedules file deletion after 5 seconds

## Installation Requirements

### 1. Install Excel Library
```bash
go get github.com/xuri/excelize/v2
```

### 2. Update go.mod
Run:
```bash
go mod tidy
```

### 3. Create Reports Directory
The system will auto-create this, but you can pre-create it:
```bash
mkdir -p reports
```

### 4. Update Nginx (Already Done)
The nginx timeouts have been increased to handle report generation:
- `proxy_connect_timeout`: 10s
- `proxy_send_timeout`: 60s  
- `proxy_read_timeout`: 60s

### 5. Reload Application
```bash
# If using systemd
sudo systemctl restart your-service-name

# Or rebuild and restart
make build
./bin/srikandisehat
```

## New Files Created

1. **src/utils/encryption.util.go** - AES-256-GCM encryption/decryption
2. **src/utils/xlsx.util.go** - XLSX generation with styling
3. **src/dto/report.dto.go** - Updated with new DTOs

## Modified Files

1. **src/handlers/report.handler.go** - Complete rewrite
2. **src/routes/api.route.go** - Updated endpoints
3. **src/utils/cache.util.go** - New metadata storage structure

## Security Features

1. **One-Time Use**: Each download link can only be used once
2. **Password Protection**: Token encrypted with random password
3. **Time Expiration**: Links expire after 30 minutes
4. **Auto Cleanup**: Files deleted after successful download
5. **AES-256-GCM Encryption**: Military-grade encryption

## API Changes

### Old Endpoints (Removed)
- ❌ `POST /api/admin/reports/generate-csv-link`
- ❌ `GET /api/reports/download/:token`

### New Endpoints
- ✅ `POST /api/admin/reports/generate`
- ✅ `POST /api/reports/download`

## Frontend Integration Example

### Step 1: Generate Report (Admin)
```javascript
const generateReport = async () => {
  const response = await fetch('https://api.example.com/api/admin/reports/generate', {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer <admin_token>'
    }
  });
  
  const data = await response.json();
  
  // Store these for the download step
  const { encrypted_token, download_url } = data.data;
  const password = extractPasswordFromMessage(data.message); // e.g., "Password: AbC123!@#XyZ"
  
  return { encrypted_token, download_url, password };
};
```

### Step 2: Download Report
```javascript
const downloadReport = async (encrypted_token, password) => {
  const response = await fetch('https://api.example.com/api/reports/download', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Report-Token': encrypted_token
    },
    body: JSON.stringify({ password })
  });
  
  if (response.ok) {
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'report_srikandi-sehat.xlsx';
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
  }
};
```

## XLSX Features

The generated Excel file includes:
- **Professional Styling**: Header with blue background, white text
- **Auto-filter**: Enabled on all columns
- **Frozen Header**: Top row stays visible when scrolling
- **Bordered Cells**: Clean table appearance
- **Auto-width Columns**: Optimized for readability
- **All Previous Data**: Same data as CSV, now in Excel format

## Error Handling

The system handles:
- Invalid/expired tokens
- Already-used tokens
- Wrong passwords
- Missing files
- Database errors
- File system errors

## Testing

### Test Report Generation
```bash
curl -X POST https://your-domain.com/api/admin/reports/generate \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### Test Report Download
```bash
curl -X POST https://your-domain.com/api/reports/download \
  -H "X-Report-Token: ENCRYPTED_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"password":"PASSWORD_HERE"}' \
  --output report.xlsx
```

## Maintenance

### Cleanup Old Reports
Reports are automatically deleted after download. To manually clean up:
```bash
# Delete reports older than 1 hour
find ./reports -name "*.xlsx" -type f -mmin +60 -delete
```

### Monitor Cache
The report cache automatically cleans up expired tokens. Monitor logs:
```bash
tail -f logs/info.log | grep "Report"
```

## Troubleshooting

### Issue: "Package not found"
```bash
go get github.com/xuri/excelize/v2
go mod tidy
```

### Issue: "Reports directory not found"
```bash
mkdir -p reports
chmod 755 reports
```

### Issue: "Token expired"
- Increase expiration time in `GenerateFullReportLink` handler
- Currently set to 30 minutes

### Issue: "File not deleted"
- Check file permissions on `./reports` directory
- Review error logs for deletion failures

## Performance Considerations

- Large datasets (>10k records) may take 10-30 seconds to generate
- XLSX files are ~2-3x larger than CSV but compress well
- Memory usage increases with record count
- Consider implementing async job processing for very large reports

## Future Enhancements

Potential improvements:
1. Background job processing with progress tracking
2. Report scheduling and email delivery
3. Multiple report formats (PDF, CSV, XLSX)
4. Custom column selection
5. Data filtering before generation
6. Report history and archiving
