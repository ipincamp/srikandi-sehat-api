#!/bin/bash

# Report System Migration Script
# This script installs dependencies and sets up the new XLSX reporting system

set -e

echo "=========================================="
echo "Report System Migration - CSV to XLSX"
echo "=========================================="
echo ""

# Step 1: Install Go dependencies
echo "[1/4] Installing Go dependencies..."
go get github.com/xuri/excelize/v2
go mod tidy
echo "✓ Dependencies installed"
echo ""

# Step 2: Create reports directory
echo "[2/4] Creating reports directory..."
mkdir -p reports
chmod 755 reports
echo "✓ Reports directory created"
echo ""

# Step 3: Build application
echo "[3/4] Building application..."
if [ -f "Makefile" ]; then
    make build
else
    go build -o bin/srikandisehat cmd/api/main.go
fi
echo "✓ Application built successfully"
echo ""

# Step 4: Test compilation
echo "[4/4] Testing new code..."
go build -o /dev/null ./src/handlers/report.handler.go 2>/dev/null && echo "✓ Report handler compiles" || echo "✗ Report handler has errors"
go build -o /dev/null ./src/utils/xlsx.util.go 2>/dev/null && echo "✓ XLSX utility compiles" || echo "✗ XLSX utility has errors"
go build -o /dev/null ./src/utils/encryption.util.go 2>/dev/null && echo "✓ Encryption utility compiles" || echo "✗ Encryption utility has errors"
echo ""

echo "=========================================="
echo "Migration Complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Review REPORT_MIGRATION.md for API changes"
echo "2. Update frontend to use new endpoints:"
echo "   - POST /api/admin/reports/generate"
echo "   - POST /api/reports/download"
echo "3. Test the new flow in staging environment"
echo "4. Reload Nginx: sudo systemctl reload nginx"
echo "5. Restart application: sudo systemctl restart srikandisehat"
echo ""
echo "For troubleshooting, see REPORT_MIGRATION.md"
