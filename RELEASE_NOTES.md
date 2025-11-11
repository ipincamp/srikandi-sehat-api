# SRIKANDI SEHAT - Release Notes

## 📋 Table of Contents

### Quick Navigation

- [Version History](#version-history)
- [Current Release (v1.0.1)](#version-101)
- [Previous Release (v1.0.0)](#version-100)

### Version 1.0.0 Documentation

- [SRIKANDI SEHAT - Release Notes](#srikandi-sehat---release-notes)
  - [📋 Table of Contents](#-table-of-contents)
    - [Quick Navigation](#quick-navigation)
    - [Version 1.0.0 Documentation](#version-100-documentation)
  - [Version History](#version-history)
    - [Version 1.0.1](#version-101)
    - [Version 1.0.0](#version-100)
  - [Version 1.0.1](#version-101-1)
    - [📅 Release Date](#-release-date)
    - [🎯 Overview](#-overview)
  - [✨ New Features](#-new-features)
    - [1. **Email Verification System**](#1-email-verification-system)
    - [2. **Maintenance Mode Management**](#2-maintenance-mode-management)
    - [3. **Push Notification System**](#3-push-notification-system)
      - [Long Period Detection](#long-period-detection)
      - [Late Period Detection](#late-period-detection)
    - [4. **Menstrual Cycle Soft Delete**](#4-menstrual-cycle-soft-delete)
    - [5. **Report Generation Improvements**](#5-report-generation-improvements)
    - [6. **Development Environment Improvements**](#6-development-environment-improvements)
  - [🔧 Improvements \& Refactoring](#-improvements--refactoring)
    - [Code Quality](#code-quality)
    - [Performance Optimizations](#performance-optimizations)
    - [Rate Limiting](#rate-limiting)
    - [Database Management](#database-management)
    - [Authentication Improvements](#authentication-improvements)
  - [🐛 Bug Fixes](#-bug-fixes)
    - [Critical Fixes](#critical-fixes)
    - [Minor Fixes](#minor-fixes)
  - [🗄️ Database Changes](#️-database-changes)
    - [New Tables](#new-tables)
    - [Schema Modifications](#schema-modifications)
  - [📦 Dependencies](#-dependencies)
    - [New Packages](#new-packages)
    - [Updated Packages](#updated-packages)
  - [🚀 Deployment Notes](#-deployment-notes)
    - [Environment Variables](#environment-variables)
    - [Database Migrations](#database-migrations)
    - [Post-Deployment Steps](#post-deployment-steps)
  - [📊 Statistics](#-statistics)
    - [Code Changes](#code-changes)
    - [Feature Breakdown](#feature-breakdown)
  - [🔒 Security Enhancements](#-security-enhancements)
  - [🧪 Testing](#-testing)
    - [Email Testing](#email-testing)
    - [Maintenance Mode Testing](#maintenance-mode-testing)
    - [Notification Testing](#notification-testing)
  - [📝 Breaking Changes](#-breaking-changes)
    - [API Changes](#api-changes)
    - [Behavioral Changes](#behavioral-changes)
  - [🔄 Migration from v1.0.0](#-migration-from-v100)
    - [Automatic Migration](#automatic-migration)
    - [Manual Steps](#manual-steps)
  - [🐛 Known Issues](#-known-issues)
  - [📌 Version Information](#-version-information)
  - [Version 1.0.0](#version-100-1)
    - [📅 Release Date](#-release-date-1)
    - [🎯 Overview](#-overview-1)
  - [✨ Features](#-features)
    - [1. **Authentication \& Authorization System**](#1-authentication--authorization-system)
      - [Registration Process](#registration-process)
      - [Login Process](#login-process)
      - [Logout Process](#logout-process)
    - [2. **User Profile Management**](#2-user-profile-management)
      - [Get My Profile](#get-my-profile)
      - [Update or Create Profile](#update-or-create-profile)
      - [Change Password](#change-password)
    - [3. **Menstrual Cycle Tracking**](#3-menstrual-cycle-tracking)
      - [Record Cycle (Start/End)](#record-cycle-startend)
      - [Get Cycle History](#get-cycle-history)
      - [Get Cycle by ID](#get-cycle-by-id)
      - [Get Cycle Status](#get-cycle-status)
    - [4. **Symptom Tracking System**](#4-symptom-tracking-system)
      - [Get Symptoms Master Data](#get-symptoms-master-data)
      - [Log Symptoms](#log-symptoms)
      - [Get Symptom History](#get-symptom-history)
      - [Get Symptom Log by ID](#get-symptom-log-by-id)
      - [Get Recommendations by Symptoms](#get-recommendations-by-symptoms)
    - [5. **Regional Data Management**](#5-regional-data-management)
      - [Indonesian Regional Hierarchy](#indonesian-regional-hierarchy)
      - [Get All Provinces](#get-all-provinces)
      - [Get Regencies by Province](#get-regencies-by-province)
      - [Get Districts by Regency](#get-districts-by-regency)
      - [Get Villages by District](#get-villages-by-district)
    - [6. **Admin Management Features**](#6-admin-management-features)
      - [Get All Users (Admin Only)](#get-all-users-admin-only)
      - [Get User by ID (Admin Only)](#get-user-by-id-admin-only)
      - [Get User Statistics (Admin Only)](#get-user-statistics-admin-only)
      - [Download Full Report CSV (Admin Only)](#download-full-report-csv-admin-only)
  - [🛡️ Security Features](#️-security-features)
    - [Authentication \& Authorization](#authentication--authorization)
    - [Rate Limiting](#rate-limiting-1)
    - [Data Protection](#data-protection)
    - [Performance Optimizations](#performance-optimizations-1)
  - [🔧 Technical Architecture](#-technical-architecture)
    - [Technology Stack](#technology-stack)
    - [Project Structure](#project-structure)
    - [Key Design Patterns](#key-design-patterns)
    - [Database Migrations](#database-migrations-1)
  - [📊 Data Models](#-data-models)
    - [Core Entities](#core-entities)
    - [Relationships](#relationships)
  - [🚀 Deployment \& Operations](#-deployment--operations)
    - [Environment Configuration](#environment-configuration)
    - [Graceful Shutdown](#graceful-shutdown)
    - [Logging](#logging)
    - [Docker Support](#docker-support)
    - [Database Management](#database-management-1)
  - [📱 Push Notifications (FCM)](#-push-notifications-fcm)
    - [Integration](#integration)
    - [Notification Events](#notification-events)
    - [Error Handling](#error-handling)
  - [🧪 Testing \& Quality Assurance](#-testing--quality-assurance)
    - [Data Seeders](#data-seeders)
    - [Factories](#factories)
    - [Validation](#validation)
  - [🔄 API Endpoints Summary](#-api-endpoints-summary)
    - [Public Endpoints](#public-endpoints)
    - [Protected Endpoints (Authenticated Users)](#protected-endpoints-authenticated-users)
    - [Regional Data Endpoints (Public)](#regional-data-endpoints-public)
    - [Admin Endpoints (Admin Role Required)](#admin-endpoints-admin-role-required)
  - [📈 Performance Metrics](#-performance-metrics)
    - [Optimization Techniques](#optimization-techniques)
    - [Scalability Considerations](#scalability-considerations)
  - [🐛 Known Limitations \& Future Enhancements](#-known-limitations--future-enhancements)
    - [Current Limitations](#current-limitations)
    - [Planned Features](#planned-features)
  - [📝 License](#-license)
  - [🙏 Acknowledgments](#-acknowledgments)
  - [📌 Version Information](#-version-information-1)
  - [🔗 Quick Links](#-quick-links)

---

## Version History

### Version 1.0.1

**Release Date:** October 31, 2025  
**Status:** Current Release

**Highlights:**

- Email verification system for user registration
- Maintenance mode management with whitelist support
- Push notification system for menstrual cycle monitoring
- Soft delete functionality for menstrual cycles
- Report generation improvements with token-based downloads
- Enhanced rate limiting and error handling
- Performance optimizations and refactoring
- Development environment improvements with Air live reload

[View detailed release notes →](#version-101)

### Version 1.0.0

**Release Date:** September 15, 2025  
**Status:** Initial Release

**Highlights:**

- Initial production release
- Complete authentication and authorization system
- Menstrual cycle tracking with predictions
- Symptom logging with personalized recommendations
- Indonesian regional data integration
- Admin dashboard and reporting features
- Firebase Cloud Messaging integration
- Asynchronous user registration system

[View detailed release notes →](#version-100)

---

## Version 1.0.1

### 📅 Release Date

October 31, 2025

### 🎯 Overview

Version 1.0.1 introduces significant enhancements to the Srikandi Sehat API, focusing on user experience improvements, system reliability, and operational flexibility. This release includes email verification, maintenance mode management, automated menstrual cycle monitoring with push notifications, and various performance optimizations.

---

## ✨ New Features

### 1. **Email Verification System**

**Implementation:**

- Added email verification fields to user model (`is_verified`, `verification_token`, `verification_token_expires_at`)
- Added `last_otp_sent_at` field for rate limiting OTP requests
- Implemented OTP generation and email delivery system
- Created verification middleware to protect authenticated endpoints
- Integrated MailHog service for email testing in development

**Flow:**

1. User registers and receives OTP via email
2. OTP valid for 15 minutes
3. Rate limit: One OTP per minute
4. User verifies email using OTP code
5. Account becomes fully activated
6. JWT token generated upon successful verification

**Benefits:**

- Prevents fake email registrations
- Ensures valid contact information
- Improved account security
- Better user engagement tracking

### 2. **Maintenance Mode Management**

**Features:**

- System-wide maintenance mode toggle
- IP-based whitelist for admin access during maintenance
- Cached maintenance status for performance
- Middleware to block non-whitelisted requests during maintenance
- Admin endpoints for managing maintenance settings

**Endpoints:**

- `GET /api/v1/maintenance/status` - Check maintenance status
- `POST /api/v1/maintenance/status` - Toggle maintenance mode (Admin)
- `GET /api/v1/maintenance/whitelist` - List whitelisted IPs (Admin)
- `POST /api/v1/maintenance/whitelist` - Add IP to whitelist (Admin)
- `DELETE /api/v1/maintenance/whitelist/:id` - Remove IP from whitelist (Admin)

**Use Cases:**

- Scheduled system maintenance
- Emergency downtime
- Database migrations
- Testing in production environment

### 3. **Push Notification System**

**Implementation:**

- Added FCM token to users table
- Created notifications table for history tracking
- Implemented notification handlers and routes
- Added endpoint to update FCM token
- Endpoint to mark notifications as read
- Get notification history with pagination

**Automated Notifications:**

#### Long Period Detection

- Cron job runs every day at 5 AM
- Detects cycles lasting more than 7 days
- Sends notification to user
- Marks cycle as `long_period_notified` to prevent duplicates

#### Late Period Detection

- Cron job runs every day at 5 AM
- Detects when expected period is 7+ days late
- Calculates based on previous cycle average
- Sends reminder notification
- Marks cycle as `late_period_notified`

**Endpoints:**

- `POST /api/v1/users/fcm-token` - Update FCM token
- `GET /api/v1/notifications` - Get notification history
- `PUT /api/v1/notifications/:id/read` - Mark as read
- `POST /api/v1/notifications/test` - Test notification (Admin)

**Benefits:**

- Proactive health monitoring
- Improved user engagement
- Early detection of potential health issues
- Better cycle awareness

### 4. **Menstrual Cycle Soft Delete**

**Implementation:**

- Added soft delete fields to menstrual cycles table
- `deleted_at`, `deleted_by`, `deletion_reason`
- API endpoint to delete cycles with reason
- Enhanced cycle history to include deleted records
- Filter option to exclude deleted cycles

**Flow:**

1. User requests cycle deletion with reason
2. System marks cycle as deleted (soft delete)
3. Records deletion timestamp, user ID, and reason
4. Cycle remains in database for audit trail
5. Excluded from active queries by default
6. Can be viewed in history with `include_deleted=true` parameter

**Benefits:**

- Data preservation for analytics
- Audit trail for cycle management
- Ability to restore if needed
- Better data integrity

### 5. **Report Generation Improvements**

**Enhancements:**

- Token-based download system for CSV reports
- Email masking in reports for privacy
- Timestamp in filename for uniqueness
- Improved download URL generation
- Base URL configuration from environment

**Flow:**

1. Admin requests report generation
2. System generates unique token
3. Returns download link with token
4. Token valid for limited time
5. User downloads via authenticated link
6. Emails masked in CSV export

**Example Filename:**

```
full_report_2025_10_31_143022.csv
```

### 6. **Development Environment Improvements**

**Additions:**

- Air configuration for live reload during development
- Improved Makefile with separate build targets
- Binary output to `bin/` directory
- Enhanced development mode setup
- Better project structure with `cmd/` directory

**New Commands:**

```bash
make dev          # Run with Air live reload
make build-api    # Build API binary
make build-migrate # Build migration tool
make build-seed   # Build seeder tool
make install-air  # Install Air for development
```

---

## 🔧 Improvements & Refactoring

### Code Quality

- **Centralized Constants**: Menstrual cycle thresholds moved to constants package
- **Enhanced Logging**: Improved logging across all components for better debugging
- **Logger Refactoring**: Better logger initialization and file handling
- **Error Handling**: More descriptive error messages and better error tracking

### Performance Optimizations

- **Query Optimization**: Improved `GetAllUsers` subquery logic and join conditions
- **Caching Enhancements**: Maintenance status and whitelist caching
- **Symptom ID Caching**: Improved menstrual data seeding performance
- **Batch Inserts**: Enhanced simulation seeder with batch operations

### Rate Limiting

- **Login Rate Limiter**: Added email fallback for better tracking
- **Function Renaming**: Clearer rate limiter function names
- **Improved Route Formatting**: Better organization and readability

### Database Management

- **Migration Tool**: Standalone binary in `cmd/migrate/`
- **Reset Commands**: `make reset-db` with confirmation prompts
- **Drop All**: `make drop-db` for clean database reset
- **Better Seeding**: Enhanced seeders for regions and simulation data

### Authentication Improvements

- **Streamlined Registration**: Removed OTP logic from initial registration
- **Default Role Assignment**: Automatic user role assignment
- **JWT Generation**: Improved token generation process
- **Response Enhancement**: Added `is_verified` field to user responses

---

## 🐛 Bug Fixes

### Critical Fixes

- **Bloom Filter Fix**: Resolved false positive issue in duplicate detection
- **Worker Bottleneck**: Fixed I/O bottleneck in registration worker
- **Download Report Link**: Adjusted to bind with base_url correctly
- **Cron Schedule**: Restored proper cron job schedule (5 AM daily)

### Minor Fixes

- **Route Formatting**: Improved clarity in API route definitions
- **Validation**: Enhanced request body validation across endpoints
- **Error Messages**: More descriptive error responses

---

## 🗄️ Database Changes

### New Tables

1. **notifications**
   - User notifications history
   - Fields: title, message, data, read status, timestamps

2. **settings**
   - System-wide configuration
   - Maintenance mode flag

3. **maintenance_whitelist**
   - IP addresses allowed during maintenance
   - Fields: IP address, description, timestamps

### Schema Modifications

1. **users**
   - Added: `fcm_token`, `is_verified`, `verification_token`, `verification_token_expires_at`, `last_otp_sent_at`

2. **menstrual_cycles**
   - Added: `deleted_at`, `deleted_by`, `deletion_reason`
   - Added: `long_period_notified`, `late_period_notified`

---

## 📦 Dependencies

### New Packages

- Email utilities (OTP generation and delivery)
- Enhanced Firebase Cloud Messaging integration
- Improved caching mechanisms

### Updated Packages

- Go modules updated for compatibility
- Enhanced logging dependencies

---

## 🚀 Deployment Notes

### Environment Variables

**New Required Variables:**

```env
# Email Configuration (for development with MailHog)
MAILHOG_HOST=localhost:1025

# Base URL for report downloads
BASE_URL=https://api.example.com
```

### Database Migrations

Run migrations to update schema:

```bash
make migrate-up
```

Or using the standalone binary:

```bash
./bin/migrate up
```

### Post-Deployment Steps

1. Run database migrations
2. Configure MailHog or SMTP for email delivery
3. Update BASE_URL in environment configuration
4. Restart application to load new configuration
5. Test email verification flow
6. Verify cron jobs are running (check at 5 AM)

---

## 📊 Statistics

### Code Changes

- **Files Changed**: 53 files
- **Additions**: +2,235 lines
- **Deletions**: -374 lines
- **Net Change**: +1,861 lines

### Feature Breakdown

- **New Migrations**: 8 database migrations
- **New Endpoints**: 10+ new API endpoints
- **New Models**: 3 new database models
- **New Workers**: 2 cron job workers
- **Removed Components**: Registration worker (simplified)

---

## 🔒 Security Enhancements

- Email verification prevents fake accounts
- Maintenance mode whitelist for controlled access
- Token-based report downloads
- Email masking in exported data
- Rate limiting improvements
- Enhanced validation across endpoints

---

## 🧪 Testing

### Email Testing

Use MailHog for local email testing:

```bash
docker-compose up mailhog
```

Access MailHog UI: `http://localhost:8025`

### Maintenance Mode Testing

1. Enable maintenance mode via admin endpoint
2. Verify non-whitelisted requests are blocked
3. Add test IP to whitelist
4. Verify whitelisted IP can access
5. Disable maintenance mode

### Notification Testing

Use the test notification endpoint:

```bash
POST /api/v1/notifications/test
Authorization: Bearer <admin_token>
{
  "user_id": "uuid",
  "title": "Test",
  "message": "Test notification"
}
```

---

## 📝 Breaking Changes

### API Changes

**None** - This release is backward compatible with v1.0.0

### Behavioral Changes

1. **Registration Flow**: Users must verify email before accessing protected endpoints
2. **Cycle Queries**: Soft-deleted cycles excluded by default (use `include_deleted=true` to include)
3. **Rate Limiting**: Enhanced login rate limiting may affect high-frequency users

---

## 🔄 Migration from v1.0.0

### Automatic Migration

Run database migrations:

```bash
make migrate-up
```

### Manual Steps

1. Update environment variables (add new required vars)
2. Existing users are automatically marked as verified
3. No data migration required for existing cycles
4. FCM tokens can be updated via new endpoint

---

## 🐛 Known Issues

**None reported** at the time of release

---

## 📌 Version Information

**Version:** 1.0.1  
**Release Date:** October 31, 2025  
**Git Tag:** v1.0.1  
**Commit:** 5f2961243703903129fd81e3a4db6ee0f65267b5  
**Previous Version:** v1.0.0 (7eb045e91811d12914240474f9bf6d05803cdbdf)

---

## Version 1.0.0

### 📅 Release Date

September 15, 2025

### 🎯 Overview

Srikandi Sehat is a comprehensive REST API built with Go (GoFiber) designed to support menstrual health tracking and management for women in Indonesia. This system provides robust features for user management, menstrual cycle tracking, symptom logging, and data analytics with a focus on regional health data collection.

---

## ✨ Features

### 1. **Authentication & Authorization System**

#### Registration Process

**Flow:**

1. User submits registration data (name, email, password, FCM token)
2. Request is validated and checked against Bloom filter for duplicate emails
3. Registration job is queued to worker pool for asynchronous processing
4. Worker validates email uniqueness in database
5. Password is hashed using Argon2id algorithm
6. User account is created with default "user" role
7. Email is added to Bloom filter for future duplicate checks
8. FCM notification sent to user's device upon completion

**Scenarios:**

- ✅ **Success**: User receives confirmation via FCM notification and can log in
- ❌ **Duplicate Email**: Immediate rejection if email exists in Bloom filter
- ❌ **Database Error**: FCM notification sent with failure reason

**Technical Highlights:**

- Asynchronous registration using worker pool (CPU-based concurrency)
- Bloom filter for O(1) duplicate email detection
- Queue capacity: 5,000 pending registrations
- Argon2id password hashing for security

#### Login Process

**Flow:**

1. User submits credentials (email, password)
2. Rate limiter checks: max 5 attempts per minute per IP
3. User data retrieved with roles and profile preloaded
4. Password verified using Argon2id
5. JWT token generated with user UUID and role information
6. User added to frequent login Bloom filter
7. Response includes user data and JWT token

**Scenarios:**

- ✅ **Success**: Returns JWT token and user profile
- ❌ **Invalid Credentials**: Returns 401 Unauthorized
- ❌ **Rate Limited**: Returns 429 Too Many Requests after 5 failed attempts

**JWT Token Structure:**

```json
{
  "usr": "user-uuid",
  "roles": ["user"],
  "exp": 1730678400
}
```

#### Logout Process

**Flow:**

1. Extract JWT token from Authorization header
2. Parse token to get expiration time
3. Check if token is already expired
4. Return success response (token invalidation commented out for performance)

**Note:** Current implementation uses stateless JWT validation. Token blocklist feature is available but commented out to optimize performance.

---

### 2. **User Profile Management**

#### Get My Profile

**Flow:**

1. Authenticated user's UUID extracted from JWT
2. User data fetched with nested preloads:
   - Roles
   - Profile with village, district, regency, province
   - Village classification (urban/rural)
3. Returns formatted user data

**Scenarios:**

- ✅ **Profile Exists**: Full profile data returned
- ⚠️ **No Profile**: Message prompts user to complete profile
- ❌ **User Not Found**: 404 error

#### Update or Create Profile

**Flow:**

1. User submits partial or complete profile data
2. Transaction started for data consistency
3. Check if profile exists for user
4. If not exists: Create new profile with default values
5. Compare each field with existing data (only update if changed)
6. Special handling for date of birth (timezone-aware conversion)
7. Village validation by code lookup
8. Commit transaction

**Updateable Fields:**

- Name (in users table)
- Phone number
- Height (cm) and weight (kg)
- Date of birth
- Menarche age (age of first menstruation)
- Last education level
- Parent's last education
- Parent's last job
- Internet access type
- Village code (Indonesian regional data)

**Scenarios:**

- ✅ **First Time**: Profile created with provided data
- ✅ **Update**: Only changed fields are updated
- ❌ **Invalid Village Code**: 404 error
- ❌ **Transaction Failure**: Rollback ensures data consistency

#### Change Password

**Flow:**

1. User provides old password and new password
2. Verify old password matches current hash
3. Hash new password using Argon2id
4. Update user record with new password hash

**Scenarios:**

- ✅ **Success**: Password changed, user remains logged in
- ❌ **Wrong Old Password**: 401 Unauthorized
- ❌ **User Not Found**: 404 error

---

### 3. **Menstrual Cycle Tracking**

#### Record Cycle (Start/End)

**Flow:**

1. User submits start date or end date (RFC3339 format)
2. Profile validation: User must have completed profile
3. Check for active (uncompleted) cycle

**Starting a New Cycle:**

1. Validate no active cycle exists
2. Parse start date
3. Verify start date is after last completed cycle's end date
4. Create new cycle record
5. Update previous cycle's length and normality status
6. Commit transaction

**Ending Current Cycle:**

1. Verify active cycle exists
2. Parse end date
3. Validate end date is not before start date
4. Calculate period length (days between start and end)
5. Determine if period is normal (2-7 days)
6. Update cycle record
7. Commit transaction

**Scenarios:**

- ✅ **Start Success**: New cycle started
- ✅ **End Success**: Cycle completed with period length calculated
- ❌ **No Profile**: 403 Forbidden - must complete profile first
- ❌ **Active Cycle Exists**: Cannot start new cycle
- ❌ **No Active Cycle**: Cannot end non-existent cycle
- ❌ **Invalid Date Order**: Start date must be after previous end date
- ❌ **Invalid End Date**: End date cannot be before start date

**Normality Criteria:**

- **Period Length**: Normal = 2-7 days
  - < 2 days: Short (Hypomenorrhea)
  - \> 7 days: Long (Menorrhagia)
- **Cycle Length**: Normal = 21-35 days
  - < 21 days: Short (Polymenorrhea)
  - \> 35 days: Long (Oligomenorrhea)

#### Get Cycle History

**Flow:**

1. Retrieve completed cycles (end_date IS NOT NULL)
2. Apply pagination (default: page 1, limit 10)
3. Order by start date descending (newest first)
4. Format response with period/cycle length and normality flags

**Scenarios:**

- ✅ **Has History**: Paginated list of completed cycles
- ⚠️ **No History**: Empty array with message to record first cycle
- ❌ **User Not Found**: 404 error

#### Get Cycle by ID

**Flow:**

1. Retrieve specific cycle by ID
2. Verify cycle belongs to authenticated user
3. Fetch associated symptom logs within cycle date range
4. Preload symptom details and options
5. Format response with cycle data and symptom history

**Response Includes:**

- Cycle dates and lengths
- Normality flags
- Associated symptom logs grouped by date
- Each symptom with category and selected options

**Scenarios:**

- ✅ **Success**: Full cycle detail with symptoms
- ❌ **Not Found**: 404 if cycle doesn't exist or belongs to another user

#### Get Cycle Status

**Flow:**

1. Check for active (uncompleted) cycle

**Scenario A: User is Currently in a Cycle**

1. Calculate current period day (days since start)
2. Determine if current period is normal (2-7 days)
3. If previous completed cycle exists:
   - Calculate current cycle length
   - Determine if cycle length is normal (21-35 days)
4. Return status with current day information

**Scenario B: User is Not in a Cycle**

1. Retrieve up to 6 most recent completed cycles
2. If no cycles exist: Prompt user to start tracking
3. Calculate last period length
4. Calculate last cycle length (if ≥2 cycles exist)
5. Calculate average cycle length from historical data
6. Predict next period date using average cycle length
7. Calculate days until next period
8. Return prediction and statistics

**Response Scenarios:**

- ✅ **Active Cycle**: Current day, period/cycle normality
- ✅ **No Active Cycle (with history)**: Prediction for next period
- ⚠️ **No History**: Message to start tracking
- ⚠️ **Prediction Passed**: Message that predicted date has passed
- ℹ️ **Insufficient Data**: Need at least one complete cycle for prediction

**Example Predictions:**

```
"Anda sedang berada di hari ke-5 siklus menstruasi."
"Periode menstruasi Anda berikutnya diprediksi dalam 12 hari."
"Tanggal prediksi menstruasi Anda telah lewat. Silakan catat siklus baru jika sudah dimulai."
```

---

### 4. **Symptom Tracking System**

#### Get Symptoms Master Data

**Flow:**

1. Fetch all symptoms with their options
2. Group by symptom type (e.g., pain, mood, physical)
3. Return hierarchical structure: symptoms → options

**Use Case:** Provides frontend with available symptoms for logging

**Response Structure:**

```json
[
  {
    "id": 1,
    "name": "Kram Perut",
    "type": "physical",
    "options": [
      { "id": 1, "name": "Ringan" },
      { "id": 2, "name": "Sedang" },
      { "id": 3, "name": "Berat" }
    ]
  }
]
```

#### Log Symptoms

**Flow:**

1. User submits symptoms array with logged_at timestamp
2. Validate all symptom IDs exist
3. Validate symptom option IDs match their parent symptoms
4. Create symptom log record
5. Find relevant menstrual cycle (if dates overlap)
6. Link symptom log to cycle (if found)
7. Create individual symptom detail records
8. Commit transaction

**Scenarios:**

- ✅ **Success**: Symptoms logged, returns log ID
- ✅ **During Cycle**: Automatically linked to active/relevant cycle
- ✅ **Outside Cycle**: Logged but not cycle-linked
- ❌ **Invalid Symptom ID**: 400 Bad Request
- ❌ **Invalid Option ID**: 400 Bad Request
- ❌ **Transaction Failure**: Rollback ensures consistency

**Request Example:**

```json
{
  "logged_at": "2025-11-03T10:30:00Z",
  "note": "Feeling tired today",
  "symptoms": [
    {
      "symptom_id": 1,
      "symptom_option_id": 2
    },
    {
      "symptom_id": 3,
      "symptom_option_id": 5
    }
  ]
}
```

#### Get Symptom History

**Flow:**

1. Apply optional date filters:
   - Single date: Filter logs for that specific day
   - Date range: Filter logs between start_date and end_date
2. Apply pagination
3. Retrieve symptom logs with detail counts
4. Order by logged_at descending

**Scenarios:**

- ✅ **With Filters**: Returns filtered paginated results
- ✅ **No Filters**: Returns all history with pagination
- ⚠️ **No History**: Empty array with pagination metadata
- ❌ **Invalid Date Format**: 400 Bad Request

#### Get Symptom Log by ID

**Flow:**

1. Retrieve specific symptom log
2. Verify log belongs to authenticated user
3. Preload symptoms and selected options
4. Extract unique symptom IDs
5. Fetch recommendations for those symptoms
6. If linked to cycle: Calculate cycle number
7. Format detailed response

**Response Includes:**

- Logged timestamp and note
- Cycle number (if applicable)
- All symptoms with categories and selected options
- Personalized recommendations for logged symptoms

**Scenarios:**

- ✅ **Success**: Full log detail with recommendations
- ❌ **Not Found**: 404 if log doesn't exist or belongs to another user

#### Get Recommendations by Symptoms

**Flow:**

1. Define recent period: last 30 days
2. Query user's symptom logs in that period
3. Count frequency of each symptom
4. Select top 4 most frequent symptoms
5. Fetch recommendations for those symptoms
6. Return prioritized recommendations

**Scenarios:**

- ✅ **Has Recent Symptoms**: Returns relevant recommendations
- ⚠️ **No Recent Symptoms**: Empty array with message
- ❌ **User Not Found**: 404 error

**Use Case:** Provides personalized health recommendations based on user's symptom patterns

---

### 5. **Regional Data Management**

#### Indonesian Regional Hierarchy

**Structure:** Province → Regency → District → Village

**Classification:** Each village is classified as:

- Urban (Perkotaan)
- Rural (Perdesaan)

#### Get All Provinces

**Flow:**

1. Retrieve provinces from in-memory cache
2. Return list of provinces with codes and names

**Scenarios:**

- ✅ **Success**: List of all Indonesian provinces
- ❌ **No Data**: 404 error

#### Get Regencies by Province

**Flow:**

1. Validate province_code query parameter
2. Retrieve regencies from cache filtered by province code
3. Return list of regencies

**Scenarios:**

- ✅ **Success**: List of regencies in specified province
- ❌ **Invalid Province Code**: 404 error

#### Get Districts by Regency

**Flow:**

1. Validate regency_code query parameter
2. Retrieve districts from cache filtered by regency code
3. Return list of districts

**Scenarios:**

- ✅ **Success**: List of districts in specified regency
- ❌ **Invalid Regency Code**: 404 error

#### Get Villages by District

**Flow:**

1. Validate district_code query parameter
2. Retrieve villages from cache filtered by district code
3. Return list of villages with classification

**Scenarios:**

- ✅ **Success**: List of villages in specified district
- ❌ **Invalid District Code**: 404 error

**Technical Note:** All regional data is cached in memory at startup for optimal performance. This eliminates database queries for regional lookups.

---

### 6. **Admin Management Features**

#### Get All Users (Admin Only)

**Flow:**

1. Verify admin role via middleware
2. Apply optional classification filter (urban/rural)
3. Exclude admin users from results
4. Apply pagination (default: page 1, limit 10)
5. Join with profiles, villages, and classifications
6. Return user list with metadata

**Filter Options:**

- `classification`: "urban" or "rural" (maps to "Perkotaan" or "Perdesaan")

**Scenarios:**

- ✅ **Success**: Paginated user list
- ✅ **With Filter**: Filtered by classification
- ❌ **Non-Admin**: 403 Forbidden
- ❌ **Unauthorized**: 401 if no valid JWT

#### Get User by ID (Admin Only)

**Flow:**

1. Verify admin role
2. Validate user UUID parameter
3. Fetch user with nested preloads:
   - Roles
   - Profile with full regional hierarchy
   - Village classification
4. Fetch user's complete menstrual cycle history
5. Format cycle history with all relevant data
6. Combine user data with cycle history

**Response Includes:**

- Full user profile
- Complete cycle history
- Regional information
- Account creation date

**Scenarios:**

- ✅ **Success**: Complete user data with cycle history
- ❌ **User Not Found**: 404 error
- ❌ **Non-Admin**: 403 Forbidden

#### Get User Statistics (Admin Only)

**Flow:**

1. Execute 4 concurrent database queries:
   - Total rural users (Perdesaan classification)
   - Total urban users (Perkotaan classification)
   - Total active users (users with ≥2 cycles)
   - Total users (excluding admins)
2. Wait for all queries to complete
3. Aggregate results
4. Return statistics

**Concurrency:**

- Uses goroutines for parallel query execution
- Error channel captures any database failures
- Wait group ensures all queries complete before response

**Scenarios:**

- ✅ **Success**: Complete statistics dashboard data
- ❌ **Database Error**: 500 if any query fails
- ❌ **Non-Admin**: 403 Forbidden

**Response Example:**

```json
{
  "total_users": 1250,
  "total_rural_users": 780,
  "total_urban_users": 470,
  "total_active_users": 890
}
```

#### Download Full Report CSV (Admin Only)

**Flow:**

1. Verify admin role
2. Fetch all menstrual cycles (excluding admin users)
3. Preload user profiles with full regional data
4. Fetch all symptom logs linked to cycles
5. Group symptoms by cycle ID
6. Calculate additional metrics:
   - User age from date of birth
   - BMI and BMI category
   - Period category (Normal/Short/Long)
   - Cycle category (Normal/Short/Long)
7. Count cycles per user
8. Generate CSV with all columns
9. Stream CSV file to client

**CSV Columns:**

- User Information: Name, Email, Registration Date
- Demographics: Age, Phone, Education levels, Parent info
- Health Data: Height, Weight, BMI, BMI Category, Menarche Age
- Location: Province, Regency, District, Village, Classification
- Cycle Data: Start Date, End Date, Period Length, Cycle Length
- Cycle Normality: Period Category, Cycle Category, Cycle Number
- Symptoms: Comma-separated list of symptoms per cycle

**BMI Categories:**

- < 17.0: Sangat Kurus (Very Underweight)
- 17.0 - 18.4: Kurus (Underweight)
- 18.5 - 25.0: Normal
- 25.1 - 27.0: Gemuk (Overweight)
- \> 27.0: Obesitas (Obese)

**Scenarios:**

- ✅ **Success**: CSV file downloaded with all user cycle data
- ❌ **Database Error**: 500 error
- ❌ **Non-Admin**: 403 Forbidden

**Use Case:** Complete data export for research, analysis, or reporting purposes

---

## 🛡️ Security Features

### Authentication & Authorization

- **JWT-based authentication** with HMAC-SHA256 signing
- **Role-based access control (RBAC)**:
  - User role: Standard features
  - Admin role: Management and reporting features
- **Argon2id password hashing** (industry-standard secure hashing)
- **Stateless JWT validation** for scalability

### Rate Limiting

- **Login endpoint**: 5 attempts per minute per IP
- **Registration endpoint**: 10 attempts per 5 minutes per IP
- **Admin endpoints**: 100 requests per minute per IP
- **IP-based tracking** using X-Real-IP and X-Forwarded-For headers

### Data Protection

- **Database transactions** for data consistency
- **Soft deletes** for user data (UUID-based)
- **Input validation** using go-playground/validator v10
- **SQL injection prevention** via GORM ORM
- **CORS configuration** for trusted origins only

### Performance Optimizations

- **Bloom filters** for duplicate detection:
  - Registration email filter (100,000 capacity, 0.1% false positive rate)
  - Frequent login filter
  - Persistent storage with periodic saves
- **In-memory caching**:
  - Regional data (provinces, regencies, districts, villages)
  - User roles
  - Token blocklist (commented out for performance)
- **Database query optimization**:
  - Strategic preloading to prevent N+1 queries
  - Indexed columns on foreign keys and UUIDs
  - Concurrent queries for statistics

---

## 🔧 Technical Architecture

### Technology Stack

- **Language:** Go 1.24.5
- **Framework:** Fiber v2 (Express-inspired, high performance)
- **ORM:** GORM with MySQL driver
- **Database:** MySQL 8.0+
- **Validation:** go-playground/validator v10
- **Authentication:** golang-jwt/jwt v5
- **Password Hashing:** Argon2id
- **Push Notifications:** Firebase Cloud Messaging (FCM)
- **Concurrency:** Worker pool with job queue

### Project Structure

```
├── cmd/                    # CLI commands (migrate, seed)
├── config/                 # Configuration management
├── database/
│   ├── migrations/         # Database schema migrations
│   ├── seeders/           # Data seeders
│   └── factories/         # Data factories for testing
├── src/
│   ├── constants/         # Application constants (roles, enums)
│   ├── dto/              # Data Transfer Objects
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Database models
│   ├── routes/           # API route definitions
│   ├── utils/            # Utility functions
│   └── workers/          # Background job workers
├── logs/                  # Application logs
├── backups/              # Database backups
└── main.go               # Application entry point
```

### Key Design Patterns

- **Repository Pattern:** Database access abstraction
- **DTO Pattern:** Request/response data structures
- **Middleware Pattern:** Cross-cutting concerns (auth, validation, rate limiting)
- **Worker Pool Pattern:** Asynchronous job processing
- **Factory Pattern:** Test data generation

### Database Migrations

All schema changes are versioned and tracked:

1. `20250612000001_create_users_table.go`
2. `20250805064542_add_uuid_and_soft_delete_to_users.go`
3. `20250805070956_create_rbac_tables.go`
4. `20250805081129_create_invalid_tokens_table.go`
5. `20250805083354_create_region_tables.go`
6. `20250805090003_create_profiles_table.go`
7. `20250805122845_create_menstrual_tracking_tables.go`
8. `20250805124318_add_indexes_to_symptom_logs.go`
9. `20250805180836_add_cycle_id_to_symptom_logs.go`
10. `20250807210841_change_cycle_dates_to_datetime.go`
11. `20250808091852_rename_log_date_in_symptom_logs.go`

---

## 📊 Data Models

### Core Entities

- **User:** Authentication, roles, profile relationship
- **Profile:** Demographics, health info, regional data
- **Role & Permission:** RBAC implementation
- **MenstrualCycle:** Cycle tracking with dates and metrics
- **SymptomLog:** User-reported symptoms with timestamps
- **SymptomLogDetail:** Individual symptom records
- **Symptom & SymptomOption:** Master symptom data
- **Recommendation:** Health recommendations per symptom
- **Region Models:** Province, Regency, District, Village, Classification

### Relationships

- User 1:1 Profile
- User N:M Roles
- User 1:N MenstrualCycles
- User 1:N SymptomLogs
- MenstrualCycle 1:N SymptomLogs
- SymptomLog 1:N SymptomLogDetails
- Symptom 1:N SymptomOptions
- Symptom 1:N Recommendations
- Village N:1 District N:1 Regency N:1 Province
- Village N:1 Classification

---

## 🚀 Deployment & Operations

### Environment Configuration

Required environment variables:

- `API_PORT`: Server port
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `JWT_SECRET`: Secret key for JWT signing
- `JWT_EXPIRATION_HOURS`: Token validity period
- `TIMEZONE`: Application timezone (default: Asia/Jakarta)
- `CORS_ALLOWED_ORIGINS`: Trusted frontend origins
- `TRUSTED_PROXIES`: Trusted proxy IPs for rate limiting

### Graceful Shutdown

- Listens for SIGINT/SIGTERM signals
- Saves Bloom filters to disk
- 5-second timeout for in-flight requests
- Ensures data consistency on shutdown

### Logging

- **Info Logger:** Successful operations, startup messages
- **Error Logger:** Failures, exceptions, database errors
- **Access Logger:** HTTP request/response logging
- Logs stored in `logs/` directory

### Docker Support

- `Dockerfile`: Application containerization
- `docker-compose.yml`: Complete stack with MySQL
- Volume mounting for persistent data

### Database Management

```bash
# Run migrations
go run cmd/migrate/main.go

# Seed database
go run cmd/seed/main.go

# Or use Makefile shortcuts
make db-migrate
make db-seed
```

---

## 📱 Push Notifications (FCM)

### Integration

- Firebase Cloud Messaging for real-time notifications
- Credential file: `serviceAccountKey.json`
- Token collection during registration

### Notification Events

1. **Registration Success:**

   - Title: "Registration Successful!"
   - Body: "Your account has been created successfully..."
   - Data: `{"status": "success"}`

2. **Registration Failed:**
   - Title: "Registration Failed"
   - Body: Error-specific message
   - Data: `{"status": "failed", "reason": "email_exists|server_error"}`

### Error Handling

- Graceful fallback if FCM client unavailable
- Logged errors for failed notifications
- Non-blocking notification sending

---

## 🧪 Testing & Quality Assurance

### Data Seeders

- **User Seeder:** Sample users with different roles
- **Region Seeder:** Complete Indonesian regional data
- **Permission & Role Seeder:** RBAC setup
- **Menstrual Seeder:** Sample cycle data
- **Simulation Seeder:** Large dataset for testing

### Factories

- **User Factory:** Generate test users with faker
- **Region Factory:** Generate test regional data

### Validation

- Comprehensive request validation on all endpoints
- Custom validation rules for Indonesian data
- Detailed error messages for validation failures

---

## 🔄 API Endpoints Summary

### Public Endpoints

- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### Protected Endpoints (Authenticated Users)

- `POST /api/auth/logout` - User logout
- `GET /api/me` - Get profile
- `PUT /api/me/details` - Update profile
- `PATCH /api/me/password` - Change password
- `GET /api/menstrual/cycles/status` - Get cycle status
- `POST /api/menstrual/cycles` - Record cycle
- `GET /api/menstrual/cycles` - Get cycle history
- `GET /api/menstrual/cycles/:id` - Get cycle detail
- `POST /api/menstrual/symptoms/log` - Log symptoms
- `GET /api/menstrual/symptoms/master` - Get symptom options
- `GET /api/menstrual/symptoms/history` - Get symptom history
- `GET /api/menstrual/symptoms/log/:id` - Get symptom log detail
- `GET /api/menstrual/recommendations` - Get recommendations

### Regional Data Endpoints (Public)

- `GET /api/regions/provinces` - List provinces
- `GET /api/regions/regencies` - List regencies by province
- `GET /api/regions/districts` - List districts by regency
- `GET /api/regions/villages` - List villages by district

### Admin Endpoints (Admin Role Required)

- `GET /api/admin/users` - List all users
- `GET /api/admin/users/:id` - Get user detail
- `GET /api/admin/users/statistics` - Get user statistics
- `GET /api/admin/reports/csv` - Download full CSV report

---

## 📈 Performance Metrics

### Optimization Techniques

1. **Bloom Filters:** O(1) duplicate detection
2. **In-Memory Caching:** Zero-latency regional data access
3. **Database Indexing:** Fast lookups on UUID, email, foreign keys
4. **Connection Pooling:** GORM automatic connection management
5. **Concurrent Queries:** Parallel execution for statistics
6. **Worker Pool:** Non-blocking user registration
7. **Strategic Preloading:** Prevent N+1 query problems

### Scalability Considerations

- **Stateless API:** Horizontal scaling ready
- **JWT Authentication:** No session storage required
- **Job Queue:** 5,000 pending job capacity
- **Rate Limiting:** DDoS protection
- **Database Transactions:** ACID compliance

---

## 🐛 Known Limitations & Future Enhancements

### Current Limitations

1. Token blocklist disabled for performance (logout doesn't invalidate JWT)
2. FCM token not validated during registration
3. No email verification flow
4. CSV export not streaming for very large datasets
5. No pagination for symptom recommendations

### Planned Features

- Email verification with OTP
- Password reset functionality
- Profile picture upload
- Advanced analytics dashboard
- Export reports in multiple formats (PDF, Excel)
- Real-time notifications for cycle predictions
- Multi-language support
- Mobile app integration

---

## 📝 License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

---

## 🙏 Acknowledgments

Built with ❤️ for women's health in Indonesia, supporting menstrual health awareness and data-driven healthcare improvements in both urban and rural communities.

---

## 📌 Version Information

> **Current Version:** 1.0.0

> **Release Date:** September 15, 2025

> **Minimum Go Version:** 1.21

> **Database:** MySQL 8.0+

---

## 🔗 Quick Links

[🏠 Back to Top](#srikandi-sehat---release-notes)

[📋 Table of Contents](#table-of-contents)

[📚 Version History](#version-history)

[🚀 Latest Release](#version-100)

---

_Last Updated: November 3, 2025_
