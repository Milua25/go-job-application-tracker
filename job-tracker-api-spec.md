# Job Application Tracker — Complete API Specification
**Version:** 1.0.0 | **Target Stack:** Go + PostgreSQL 18 | **Deployment:** Self-hosted

---

## Table of Contents

1. [User Personas](#1-user-personas)
2. [Functional Requirements](#2-functional-requirements)
3. [Non-Functional Requirements](#3-non-functional-requirements)
4. [User Permission Matrix](#4-user-permission-matrix)
5. [Domain Model](#5-domain-model)
6. [Endpoint Catalog](#6-endpoint-catalog)
7. [REST API Specification](#7-rest-api-specification)
8. [Security Architecture](#8-security-architecture)
9. [Audit & Compliance Requirements](#9-audit--compliance-requirements)
10. [Architectural Assumptions](#10-architectural-assumptions)

---

## 1. User Personas

### P1 — Active Job Seeker (Primary)
A professional actively applying to multiple positions simultaneously. Applies to 5–30 companies per week. Primary concern is tracking pipeline status, upcoming follow-up deadlines, and avoiding duplicate applications.

**Goals:** Fast data entry, clear status pipeline, timely reminders, export capability.
**Pain points:** Losing track of where they are with each company; forgetting to follow up.

### P2 — Passive Candidate
Currently employed but open to opportunities. Applies infrequently (1–5 per month). Values notes and contact management over speed.

**Goals:** Rich per-application notes, contact tracking, preserve history.
**Pain points:** Applications get stale with no reminders to re-engage.

### P3 — Admin / Instance Owner (Self-hosted)
The person who deployed the instance. May be the only user, or may manage accounts for a household or small team. Needs to manage user accounts, reset passwords, and observe system health.

**Goals:** User management, audit visibility, system health monitoring.
**Pain points:** No visibility into what's happening at the system level.

---

## 2. Functional Requirements

### FR-AUTH: Authentication & Sessions
- FR-AUTH-01: Users can register with email and password
- FR-AUTH-02: Users can log in and receive access + refresh tokens
- FR-AUTH-03: Access tokens expire in 15 minutes; refresh tokens in 7 days
- FR-AUTH-04: Users can refresh access tokens without re-entering credentials
- FR-AUTH-05: Users can log out (invalidating the current refresh token)
- FR-AUTH-06: Users can log out of all sessions
- FR-AUTH-07: Password reset via email token (time-limited, single-use)
- FR-AUTH-08: Admin can force-revoke all sessions for any user

### FR-USER: User Management
- FR-USER-01: Users can view and update their own profile (name, email)
- FR-USER-02: Users can change their password (requires current password)
- FR-USER-03: Users can request deletion of their account (GDPR right to erasure)
- FR-USER-04: Admin can list, view, create, update, and deactivate users
- FR-USER-05: Admin can assign roles to users

### FR-APP: Job Application Management
- FR-APP-01: Users can create a job application (job title, company, URL, status, type, mode, salary, notes)
- FR-APP-02: Users can view a list of their own applications with pagination
- FR-APP-03: Users can view a single application with full details
- FR-APP-04: Users can update any field of their own application
- FR-APP-05: Users can soft-delete an application
- FR-APP-06: Users can filter applications by status, job type, work mode, company, date range
- FR-APP-07: Users can search applications by job title or company name (full-text)
- FR-APP-08: Users can move an application to a new status (triggers status history record)
- FR-APP-09: Users can view the status history of an application
- FR-APP-10: Users can restore a soft-deleted application (within 30 days)

### FR-COMPANY: Company Management
- FR-CO-01: Companies are created automatically when an application references a new company name
- FR-CO-02: Users can create, view, update, and soft-delete companies explicitly
- FR-CO-03: Users can list their own companies with application count

### FR-CONTACT: Contact Management
- FR-CON-01: Users can add, update, and delete contacts associated with an application
- FR-CON-02: Contacts hold: name, role, email, LinkedIn URL, notes

### FR-REM: Reminders
- FR-REM-01: Users can create reminders attached to an application or standalone
- FR-REM-02: Users can list pending and sent reminders
- FR-REM-03: Users can update or delete unsent reminders
- FR-REM-04: System marks reminders as sent after delivery (background worker)

### FR-RPT: Reporting
- FR-RPT-01: Users can retrieve a summary of their application pipeline (count per status)
- FR-RPT-02: Users can retrieve application trends over time (applications added per week/month)
- FR-RPT-03: Users can export all their application data as JSON or CSV

### FR-REF: Reference Data
- FR-REF-01: Application statuses, job types, and work modes are readable by all authenticated users
- FR-REF-02: Admin can create, update, and reorder application statuses (kanban column management)

### FR-AUD: Audit
- FR-AUD-01: All mutating API calls are recorded in an audit log (actor, action, entity, before/after)
- FR-AUD-02: Admin can query the audit log with filtering by user, entity type, date range

### FR-SYS: System
- FR-SYS-01: Health check endpoint returns service + database status
- FR-SYS-02: Metrics endpoint exposes Prometheus-compatible metrics
- FR-SYS-03: Admin can trigger a background job to flush expired refresh tokens

---

## 3. Non-Functional Requirements

### Performance
- NFR-P-01: All read endpoints respond in < 200ms at p99 under normal load
- NFR-P-02: List endpoints paginate to prevent unbounded result sets (max 100 per page)
- NFR-P-03: Full-text search responds in < 500ms at p99

### Security
- NFR-S-01: All communication over HTTPS/TLS 1.3 minimum
- NFR-S-02: Passwords hashed with bcrypt (cost factor ≥ 12) or Argon2id
- NFR-S-03: JWTs signed with RS256 (asymmetric); private key never leaves the server
- NFR-S-04: PII fields (email, contact email) encrypted at rest (AES-256-GCM)
- NFR-S-05: Refresh tokens stored as SHA-256 hashes only; raw token never persisted
- NFR-S-06: All inputs validated and sanitised before persistence
- NFR-S-07: SQL injection prevented by parameterised queries (sqlc enforces this)
- NFR-S-08: Rate limiting on all auth endpoints

### Reliability
- NFR-R-01: Background reminder worker retries failed deliveries with exponential backoff
- NFR-R-02: All database mutations wrapped in transactions where state spans tables

### Observability
- NFR-O-01: Structured JSON logs on all requests (correlation ID, method, path, status, duration)
- NFR-O-02: Prometheus metrics exported on `/metrics`
- NFR-O-03: Every API response includes a `X-Correlation-ID` header

### Compliance
- NFR-C-01: GDPR right-to-erasure: account deletion scrubs PII within 30 days
- NFR-C-02: Audit log retained for minimum 12 months
- NFR-C-03: Soft-deleted application data purged after 30 days

---

## 4. User Permission Matrix

### Roles

| Role | Description |
|------|-------------|
| `user` | Standard authenticated user; owns their own data |
| `admin` | Instance administrator; manages users, roles, and system config |

### Permission Matrix

| Resource | Action | `user` (own) | `user` (others) | `admin` |
|----------|--------|:---:|:---:|:---:|
| **Auth** | Register | ✓ | — | ✓ |
| | Login | ✓ | — | ✓ |
| | Refresh | ✓ | — | ✓ |
| | Logout | ✓ | — | ✓ |
| | Logout all sessions | ✓ | ✗ | ✓ |
| | Request password reset | ✓ | — | ✓ |
| | Admin force-revoke sessions | ✗ | ✗ | ✓ |
| **Users** | View own profile | ✓ | — | ✓ |
| | Update own profile | ✓ | ✗ | ✓ |
| | Change own password | ✓ | ✗ | ✓ |
| | Delete own account | ✓ | ✗ | ✓ |
| | List all users | ✗ | — | ✓ |
| | View any user | ✗ | ✗ | ✓ |
| | Create user | ✗ | ✗ | ✓ |
| | Update any user | ✗ | ✗ | ✓ |
| | Assign roles | ✗ | ✗ | ✓ |
| | Deactivate/reactivate user | ✗ | ✗ | ✓ |
| **Applications** | Create | ✓ | — | ✓ |
| | Read own | ✓ | — | ✓ |
| | Read others | ✗ | ✗ | ✓ |
| | Update own | ✓ | ✗ | ✓ |
| | Delete own (soft) | ✓ | ✗ | ✓ |
| | Restore own | ✓ | ✗ | ✓ |
| | Export own data | ✓ | ✗ | ✓ |
| **Companies** | CRUD own | ✓ | ✗ | ✓ |
| **Contacts** | CRUD (via application) | ✓ | ✗ | ✓ |
| **Reminders** | CRUD own | ✓ | ✗ | ✓ |
| **Reference data** | Read statuses/types/modes | ✓ | — | ✓ |
| | Manage statuses (create/update/reorder) | ✗ | — | ✓ |
| **Reports** | Own pipeline summary | ✓ | ✗ | ✓ |
| | Own trends | ✓ | ✗ | ✓ |
| | Any user report | ✗ | ✗ | ✓ |
| **Audit log** | Read | ✗ | ✗ | ✓ |
| **System** | Health check | ✓ | — | ✓ |
| | Metrics | ✗ | — | ✓ |
| | Admin actions | ✗ | — | ✓ |

---

## 5. Domain Model

### Core Entities

#### User
Owns all data in the system. Identified by UUID. Has a role (`user` or `admin`). PII: `email`, `full_name`. Supports soft delete (deactivation without data loss). Password stored as a bcrypt/Argon2id hash; never returned in any API response.

Lifecycle states: `active` → `inactive` (deactivated by admin) → `pending_deletion` (self-requested) → `deleted` (purged after 30-day grace period)

#### Job Application
The primary entity. Belongs to one User, optionally to one Company. Tracks the full pipeline from bookmarked through offer/rejection. Has a `status_id` FK to `application_statuses`, plus freeform `notes`. Soft-deleted (never hard-deleted except during GDPR erasure). Changing `status_id` triggers an `application_status_history` append.

Status transitions:
```
Bookmarked → Applied → Phone Screen → Technical Interview → Final Interview → Offer
Any active status → Rejected
Any active status → Withdrawn
Offer / Rejected / Withdrawn → (terminal — no further transitions)
```

#### Company
Belongs to one User (user-scoped, not global). Stores name, website, industry, headquarters location, and notes. One Company → many Applications.

#### Contact
Belongs to one Application. Stores a person at the company: name, role, email, LinkedIn URL, notes. PII: `email`.

#### Reminder
Belongs to one User; optionally linked to one Application. Has a `remind_at` timestamp and an `is_sent` boolean. Background worker polls `WHERE is_sent = false AND remind_at <= now()`.

#### application_status_history
Append-only audit table. Records every status change on an application: `application_id`, `from_status_id` (nullable on first), `to_status_id`, `changed_by`, `changed_at`. Never updated or deleted.

### Reference / Lookup Entities

| Table | Contents | Managed by |
|-------|----------|------------|
| `application_statuses` | Kanban pipeline stages with color and order | Admin |
| `job_types` | Full-time, Part-time, Contract, Internship, Freelance | Seeded / read-only |
| `work_modes` | Remote, Hybrid, On-site | Seeded / read-only |

### Relationships

```
User 1──< Company
User 1──< Job Application
User 1──< Reminder
Company 0..1──< Job Application
Job Application 1──< Contact
Job Application 1──< application_status_history
Job Application >──1 application_statuses (current)
Job Application 0..1──1 job_types
Job Application 0..1──1 work_modes
application_status_history >──1 application_statuses (from)
application_status_history >──1 application_statuses (to)
```

---

## 6. Endpoint Catalog

| # | Method | Path | Purpose | Auth | Role |
|---|--------|------|---------|------|------|
| 1 | POST | `/api/v1/auth/register` | Create a new user account | None | — |
| 2 | POST | `/api/v1/auth/login` | Authenticate and receive tokens | None | — |
| 3 | POST | `/api/v1/auth/refresh` | Exchange refresh token for new access token | Refresh token | — |
| 4 | POST | `/api/v1/auth/logout` | Revoke current refresh token | Bearer | any |
| 5 | POST | `/api/v1/auth/logout-all` | Revoke all refresh tokens for current user | Bearer | any |
| 6 | POST | `/api/v1/auth/password/reset-request` | Send password reset email | None | — |
| 7 | POST | `/api/v1/auth/password/reset` | Reset password with token | None | — |
| 8 | GET | `/api/v1/users/me` | Get current user's profile | Bearer | any |
| 9 | PATCH | `/api/v1/users/me` | Update current user's profile | Bearer | any |
| 10 | POST | `/api/v1/users/me/password` | Change password | Bearer | any |
| 11 | DELETE | `/api/v1/users/me` | Request account deletion | Bearer | any |
| 12 | GET | `/api/v1/users` | List all users | Bearer | admin |
| 13 | POST | `/api/v1/users` | Create a user (admin) | Bearer | admin |
| 14 | GET | `/api/v1/users/{userId}` | Get a specific user | Bearer | admin |
| 15 | PATCH | `/api/v1/users/{userId}` | Update a specific user | Bearer | admin |
| 16 | POST | `/api/v1/users/{userId}/roles` | Assign role to user | Bearer | admin |
| 17 | POST | `/api/v1/users/{userId}/deactivate` | Deactivate a user account | Bearer | admin |
| 18 | POST | `/api/v1/users/{userId}/reactivate` | Reactivate a user account | Bearer | admin |
| 19 | DELETE | `/api/v1/users/{userId}/sessions` | Revoke all sessions for user | Bearer | admin |
| 20 | GET | `/api/v1/applications` | List own applications | Bearer | any |
| 21 | POST | `/api/v1/applications` | Create a new application | Bearer | any |
| 22 | GET | `/api/v1/applications/{id}` | Get a single application | Bearer | any |
| 23 | PATCH | `/api/v1/applications/{id}` | Update an application | Bearer | any |
| 24 | DELETE | `/api/v1/applications/{id}` | Soft-delete an application | Bearer | any |
| 25 | POST | `/api/v1/applications/{id}/restore` | Restore a soft-deleted application | Bearer | any |
| 26 | POST | `/api/v1/applications/{id}/status` | Change application status | Bearer | any |
| 27 | GET | `/api/v1/applications/{id}/status-history` | Get status history for an application | Bearer | any |
| 28 | GET | `/api/v1/applications/{id}/contacts` | List contacts for an application | Bearer | any |
| 29 | POST | `/api/v1/applications/{id}/contacts` | Add a contact to an application | Bearer | any |
| 30 | PATCH | `/api/v1/applications/{id}/contacts/{contactId}` | Update a contact | Bearer | any |
| 31 | DELETE | `/api/v1/applications/{id}/contacts/{contactId}` | Delete a contact | Bearer | any |
| 32 | GET | `/api/v1/companies` | List own companies | Bearer | any |
| 33 | POST | `/api/v1/companies` | Create a company | Bearer | any |
| 34 | GET | `/api/v1/companies/{id}` | Get a company | Bearer | any |
| 35 | PATCH | `/api/v1/companies/{id}` | Update a company | Bearer | any |
| 36 | DELETE | `/api/v1/companies/{id}` | Soft-delete a company | Bearer | any |
| 37 | GET | `/api/v1/reminders` | List own reminders | Bearer | any |
| 38 | POST | `/api/v1/reminders` | Create a reminder | Bearer | any |
| 39 | PATCH | `/api/v1/reminders/{id}` | Update a reminder | Bearer | any |
| 40 | DELETE | `/api/v1/reminders/{id}` | Delete a reminder | Bearer | any |
| 41 | GET | `/api/v1/reference/statuses` | List application statuses | Bearer | any |
| 42 | POST | `/api/v1/reference/statuses` | Create application status | Bearer | admin |
| 43 | PATCH | `/api/v1/reference/statuses/{id}` | Update application status | Bearer | admin |
| 44 | POST | `/api/v1/reference/statuses/reorder` | Reorder kanban columns | Bearer | admin |
| 45 | GET | `/api/v1/reference/job-types` | List job types | Bearer | any |
| 46 | GET | `/api/v1/reference/work-modes` | List work modes | Bearer | any |
| 47 | GET | `/api/v1/reports/pipeline` | Pipeline summary by status | Bearer | any |
| 48 | GET | `/api/v1/reports/trends` | Application trends over time | Bearer | any |
| 49 | GET | `/api/v1/reports/export` | Export all application data | Bearer | any |
| 50 | GET | `/api/v1/audit` | Query audit log | Bearer | admin |
| 51 | GET | `/api/v1/admin/users/{userId}/applications` | List applications for any user | Bearer | admin |
| 52 | POST | `/api/v1/admin/jobs/purge-expired-tokens` | Trigger token cleanup job | Bearer | admin |
| 53 | GET | `/health` | Service health check | None | — |
| 54 | GET | `/metrics` | Prometheus metrics | Bearer (admin) or IP allowlist | admin |

---

## 7. REST API Specification

### API Standards

**Base URL:** `https://your-instance.example.com/api/v1`

**Versioning:** URI-based versioning (`/v1`). Breaking changes increment the version. Non-breaking additions are backward-compatible.

**Content-Type:** All request and response bodies use `application/json` unless otherwise noted.

**Timestamps:** ISO 8601 with UTC timezone. Example: `2026-06-17T14:30:00Z`

**IDs:** All entity IDs are UUID v4 strings.

**Pagination:** Cursor-based for large sets; offset-based for bounded sets.
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "per_page": 25,
    "total": 143,
    "total_pages": 6,
    "has_next": true,
    "has_prev": false
  }
}
```

**Sorting:** `?sort=created_at&order=desc` (default). Multi-sort: `?sort=status,created_at&order=asc,desc`

**Filtering:** Query parameter based. `?status_id=UUID&job_type_id=UUID&work_mode_id=UUID&company_id=UUID&created_after=2026-01-01&created_before=2026-06-30`

**Search:** `?q=software+engineer` — full-text search against `job_title` and company `name`.

**Soft-deleted resources:** Excluded by default. `?include_deleted=true` returns all; `?deleted_only=true` returns only deleted.

**Rate Limiting Headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 87
X-RateLimit-Reset: 1718640600
Retry-After: 60  (only on 429 responses)
```

**Correlation ID:** Every request and response carries `X-Correlation-ID`. If the client sends one, it is echoed back. Otherwise, one is generated server-side.

---

### Standard Error Response

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "details": [
      {
        "field": "email",
        "message": "Must be a valid email address"
      }
    ],
    "correlation_id": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2026-06-17T14:30:00Z"
  }
}
```

**Error Codes:**

| HTTP Status | Code | Meaning |
|-------------|------|---------|
| 400 | `VALIDATION_ERROR` | Request body or query params failed validation |
| 400 | `INVALID_TOKEN` | Malformed or expired one-time token |
| 401 | `UNAUTHORIZED` | Missing or invalid Bearer token |
| 401 | `TOKEN_EXPIRED` | Access token has expired |
| 403 | `FORBIDDEN` | Authenticated but insufficient permissions |
| 404 | `NOT_FOUND` | Resource does not exist or is not owned by caller |
| 409 | `CONFLICT` | Uniqueness violation (e.g. duplicate email) |
| 409 | `INVALID_STATE_TRANSITION` | Status move is not permitted |
| 422 | `UNPROCESSABLE` | Semantically invalid request |
| 429 | `RATE_LIMITED` | Too many requests |
| 500 | `INTERNAL_ERROR` | Unexpected server error |
| 503 | `SERVICE_UNAVAILABLE` | Database or dependency down |

---

### Authentication Endpoints

---

#### POST /api/v1/auth/register
**Purpose:** Create a new user account.
**Authentication:** None
**Idempotency:** Not idempotent (same email returns 409)

**Request Body:**
```json
{
  "full_name": "string (1–200 chars, required)",
  "email": "string (valid email, max 320 chars, required)",
  "password": "string (min 10 chars, max 128 chars, required)"
}
```

**Validation Rules:**
- `email`: RFC 5322 format; normalised to lowercase; must be unique
- `password`: Min entropy check (reject common passwords via a blocklist); not equal to email

**Response 201 Created:**
```json
{
  "data": {
    "id": "uuid",
    "full_name": "Jane Doe",
    "email": "jane@example.com",
    "role": "user",
    "created_at": "2026-06-17T14:30:00Z"
  }
}
```

**Error Responses:** 400 (validation), 409 (email already registered)

---

#### POST /api/v1/auth/login
**Purpose:** Authenticate with email and password; receive JWT access token and refresh token.
**Authentication:** None
**Rate Limiting:** 10 attempts / 15 min per IP; 5 attempts / 15 min per email

**Request Body:**
```json
{
  "email": "string (required)",
  "password": "string (required)"
}
```

**Validation Rules:**
- Both fields required; return generic 401 on any mismatch (never reveal which field is wrong)
- Account must be `active`; deactivated accounts receive 401 with code `ACCOUNT_DISABLED`

**Response 200 OK:**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900,
    "refresh_token": "opaque-random-256-bit-token",
    "user": {
      "id": "uuid",
      "full_name": "Jane Doe",
      "email": "jane@example.com",
      "role": "user"
    }
  }
}
```

**Security Notes:**
- Refresh token returned in response body AND as `HttpOnly; Secure; SameSite=Strict` cookie
- Raw refresh token never stored — only its SHA-256 hash
- JWT payload: `{ sub, email, role, iat, exp, jti }`

**Error Responses:** 400 (validation), 401 (invalid credentials / disabled), 429 (rate limited)

---

#### POST /api/v1/auth/refresh
**Purpose:** Exchange a valid refresh token for a new access token (and optionally rotate the refresh token).
**Authentication:** Refresh token (in `Authorization: Bearer <refresh_token>` header or HttpOnly cookie)
**Security:** Refresh token rotation — old token is invalidated, new one issued

**Request Body:** (empty — token from header or cookie)

**Response 200 OK:**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900,
    "refresh_token": "new-opaque-random-256-bit-token"
  }
}
```

**Error Responses:** 401 (invalid/expired/revoked refresh token)

---

#### POST /api/v1/auth/logout
**Purpose:** Revoke the current refresh token (end current session).
**Authentication:** Bearer access token required

**Request Body:** (empty)

**Response 204 No Content**

---

#### POST /api/v1/auth/logout-all
**Purpose:** Revoke all refresh tokens for the current user (force logout of all devices).
**Authentication:** Bearer + current password confirmation

**Request Body:**
```json
{
  "password": "string (required)"
}
```

**Response 204 No Content**

---

#### POST /api/v1/auth/password/reset-request
**Purpose:** Initiate a password reset by sending a time-limited token to the user's email.
**Authentication:** None
**Security:** Always returns 200 regardless of whether email exists (prevents enumeration)
**Rate Limiting:** 3 requests / 10 min per email address

**Request Body:**
```json
{
  "email": "string (required)"
}
```

**Response 200 OK:**
```json
{
  "data": {
    "message": "If that email is registered, a reset link has been sent."
  }
}
```

---

#### POST /api/v1/auth/password/reset
**Purpose:** Set a new password using a valid reset token.
**Authentication:** Reset token (single-use, 1-hour TTL)

**Request Body:**
```json
{
  "token": "string (required)",
  "password": "string (min 10 chars, required)",
  "password_confirmation": "string (must match password, required)"
}
```

**Response 204 No Content**

**Error Responses:** 400 (validation, password mismatch), 400 `INVALID_TOKEN` (expired or already used)

---

### User Management Endpoints

---

#### GET /api/v1/users/me
**Purpose:** Return the authenticated user's profile.
**Authentication:** Bearer
**Roles:** Any authenticated user

**Response 200 OK:**
```json
{
  "data": {
    "id": "uuid",
    "full_name": "Jane Doe",
    "email": "jane@example.com",
    "role": "user",
    "is_active": true,
    "created_at": "2026-06-17T14:30:00Z",
    "updated_at": "2026-06-17T14:30:00Z"
  }
}
```

---

#### PATCH /api/v1/users/me
**Purpose:** Update the current user's name or email. Email change requires re-verification.
**Authentication:** Bearer
**Roles:** Any authenticated user

**Request Body (all fields optional):**
```json
{
  "full_name": "string (1–200 chars)",
  "email": "string (valid email)"
}
```

**Validation:** If `email` changes, uniqueness check required. Email change should trigger a verification step (out of scope for v1 — treated as immediate).

**Response 200 OK:** (same schema as GET /users/me)

**Error Responses:** 400 (validation), 409 (email taken)

---

#### POST /api/v1/users/me/password
**Purpose:** Change the current user's password.
**Authentication:** Bearer

**Request Body:**
```json
{
  "current_password": "string (required)",
  "new_password": "string (min 10 chars, required)",
  "new_password_confirmation": "string (required)"
}
```

**Response 204 No Content**
**Side effects:** All existing refresh tokens for this user are revoked (force re-login on all devices)

**Error Responses:** 400 (validation), 401 (current password wrong)

---

#### DELETE /api/v1/users/me
**Purpose:** Request deletion of the current account (initiates GDPR erasure flow).
**Authentication:** Bearer + password confirmation

**Request Body:**
```json
{
  "password": "string (required)",
  "confirmation": "DELETE MY ACCOUNT"
}
```

**Behaviour:** Sets `status = pending_deletion`. A background job purges PII after a 30-day grace period. All sessions are immediately revoked.

**Response 202 Accepted:**
```json
{
  "data": {
    "message": "Your account has been scheduled for deletion. All personal data will be removed within 30 days.",
    "scheduled_deletion_at": "2026-07-17T14:30:00Z"
  }
}
```

---

#### GET /api/v1/users
**Purpose:** List all users in the system.
**Authentication:** Bearer
**Roles:** admin only

**Query Parameters:**
- `page`, `per_page` (max 100)
- `q` — search by name or email
- `role` — filter by role
- `is_active` — boolean filter

**Response 200 OK:**
```json
{
  "data": [
    {
      "id": "uuid",
      "full_name": "Jane Doe",
      "email": "jane@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2026-06-17T14:30:00Z"
    }
  ],
  "pagination": { ... }
}
```

---

#### POST /api/v1/users
**Purpose:** Admin creates a new user (bypasses self-registration).
**Authentication:** Bearer
**Roles:** admin

**Request Body:**
```json
{
  "full_name": "string (required)",
  "email": "string (required)",
  "password": "string (required, min 10 chars)",
  "role": "user | admin (default: user)"
}
```

**Response 201 Created:** (same schema as GET /users/me)
**Error Responses:** 400, 409

---

#### PATCH /api/v1/users/{userId}
**Purpose:** Admin updates a user's profile or role.
**Authentication:** Bearer
**Roles:** admin

**Request Body (all optional):**
```json
{
  "full_name": "string",
  "email": "string",
  "role": "user | admin"
}
```

**Response 200 OK**
**Error Responses:** 400, 404, 409

---

#### POST /api/v1/users/{userId}/deactivate
**Purpose:** Disable a user's access without deleting their data.
**Authentication:** Bearer
**Roles:** admin

**Request Body:** (empty)
**Response 204 No Content**
**Side effects:** All active sessions for the user are revoked immediately.

---

#### DELETE /api/v1/users/{userId}/sessions
**Purpose:** Admin force-revokes all sessions for a user (e.g. after a security incident).
**Authentication:** Bearer
**Roles:** admin

**Response 204 No Content**

---

### Application Endpoints

---

#### GET /api/v1/applications
**Purpose:** List the authenticated user's job applications.
**Authentication:** Bearer
**Roles:** any

**Query Parameters:**
- `page`, `per_page` (max 100, default 25)
- `sort` — `created_at | updated_at | job_title | applied_at` (default: `created_at`)
- `order` — `asc | desc` (default: `desc`)
- `q` — full-text search across `job_title` and company `name`
- `status_id` — UUID filter
- `job_type_id` — UUID filter
- `work_mode_id` — UUID filter
- `company_id` — UUID filter
- `created_after`, `created_before` — ISO 8601 date
- `include_deleted` — boolean (default: false)
- `deleted_only` — boolean (default: false)

**Response 200 OK:**
```json
{
  "data": [
    {
      "id": "uuid",
      "job_title": "Senior Software Engineer",
      "company": {
        "id": "uuid",
        "name": "Acme Corp"
      },
      "status": {
        "id": "uuid",
        "label": "Technical Interview",
        "color_hex": "#3B82F6"
      },
      "job_type": { "id": "uuid", "label": "Full-time" },
      "work_mode": { "id": "uuid", "label": "Remote" },
      "salary_min": 120000,
      "salary_max": 160000,
      "salary_currency": "USD",
      "job_url": "https://example.com/jobs/123",
      "applied_at": "2026-06-01T00:00:00Z",
      "notes": "Reached out via LinkedIn before applying.",
      "deleted_at": null,
      "created_at": "2026-06-01T09:00:00Z",
      "updated_at": "2026-06-10T11:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 25,
    "total": 47,
    "total_pages": 2,
    "has_next": true,
    "has_prev": false
  }
}
```

---

#### POST /api/v1/applications
**Purpose:** Create a new job application.
**Authentication:** Bearer
**Roles:** any
**Idempotency:** Clients may send `Idempotency-Key: <uuid>` header; duplicate requests within 24h return the cached 201 response.

**Request Body:**
```json
{
  "job_title": "string (1–300 chars, required)",
  "company_id": "uuid (optional — omit to create inline company)",
  "company_name": "string (required if company_id omitted; 1–300 chars)",
  "status_id": "uuid (required)",
  "job_type_id": "uuid (optional)",
  "work_mode_id": "uuid (optional)",
  "job_url": "string (valid URL, max 2048 chars, optional)",
  "salary_min": "integer ≥ 0 (optional)",
  "salary_max": "integer ≥ salary_min (optional)",
  "salary_currency": "string (ISO 4217, 3 chars, default: USD)",
  "applied_at": "ISO 8601 date (optional)",
  "notes": "string (max 10000 chars, optional)"
}
```

**Validation Rules:**
- `salary_max >= salary_min` if both provided
- `applied_at` cannot be in the future
- `status_id` must exist in `application_statuses`
- If `company_name` provided without `company_id`, a new Company is upserted by name for this user

**Response 201 Created:** Full application object (same schema as list item)

**Error Responses:** 400 (validation), 422 (invalid reference IDs)

---

#### GET /api/v1/applications/{id}
**Purpose:** Retrieve a single application with full details including contacts and current status.
**Authentication:** Bearer
**Ownership check:** `application.user_id = authenticated user's ID`

**Response 200 OK:**
```json
{
  "data": {
    "id": "uuid",
    "job_title": "Senior Software Engineer",
    "company": {
      "id": "uuid",
      "name": "Acme Corp",
      "website": "https://acme.example.com"
    },
    "status": {
      "id": "uuid",
      "label": "Technical Interview",
      "color_hex": "#3B82F6",
      "is_terminal": false
    },
    "job_type": { "id": "uuid", "label": "Full-time" },
    "work_mode": { "id": "uuid", "label": "Remote" },
    "salary_min": 120000,
    "salary_max": 160000,
    "salary_currency": "USD",
    "job_url": "https://example.com/jobs/123",
    "applied_at": "2026-06-01T00:00:00Z",
    "notes": "Reached out via LinkedIn before applying.",
    "contacts": [
      {
        "id": "uuid",
        "full_name": "Alex Rivera",
        "role": "Recruiter",
        "email": "alex@acme.example.com",
        "linkedin_url": "https://linkedin.com/in/alex-rivera",
        "notes": "Very responsive, check in every Friday."
      }
    ],
    "reminder_count": 2,
    "deleted_at": null,
    "created_at": "2026-06-01T09:00:00Z",
    "updated_at": "2026-06-10T11:00:00Z"
  }
}
```

**Error Responses:** 401, 403 (not owner), 404

---

#### PATCH /api/v1/applications/{id}
**Purpose:** Update one or more fields of an application (excluding status — use the dedicated status endpoint).
**Authentication:** Bearer
**Ownership check:** required

**Request Body (all fields optional):**
```json
{
  "job_title": "string",
  "company_id": "uuid",
  "job_type_id": "uuid",
  "work_mode_id": "uuid",
  "job_url": "string",
  "salary_min": "integer",
  "salary_max": "integer",
  "salary_currency": "string",
  "applied_at": "ISO 8601",
  "notes": "string"
}
```

**Response 200 OK:** Updated application object
**Error Responses:** 400, 403, 404

---

#### DELETE /api/v1/applications/{id}
**Purpose:** Soft-delete an application (sets `deleted_at`; excluded from default list results).
**Authentication:** Bearer
**Ownership check:** required

**Response 204 No Content**

---

#### POST /api/v1/applications/{id}/restore
**Purpose:** Restore a soft-deleted application (clears `deleted_at`).
**Authentication:** Bearer
**Constraint:** Only restorable if `deleted_at` is within the last 30 days.

**Response 200 OK:** Restored application object
**Error Responses:** 403, 404, 409 (deletion is older than 30 days)

---

#### POST /api/v1/applications/{id}/status
**Purpose:** Transition an application to a new status. Creates a status history record.
**Authentication:** Bearer
**Ownership check:** required
**Idempotency:** Re-sending the same `status_id` when it is already the current status returns 200 without creating a duplicate history record.

**Request Body:**
```json
{
  "status_id": "uuid (required)",
  "note": "string (optional — recorded on the history entry)"
}
```

**Validation:**
- `status_id` must exist
- Cannot transition from a terminal status (`is_terminal = true`)

**Response 200 OK:**
```json
{
  "data": {
    "application_id": "uuid",
    "previous_status": { "id": "uuid", "label": "Applied" },
    "current_status": { "id": "uuid", "label": "Phone Screen" },
    "changed_at": "2026-06-17T14:30:00Z"
  }
}
```

**Error Responses:** 400, 403, 404, 409 `INVALID_STATE_TRANSITION`

---

#### GET /api/v1/applications/{id}/status-history
**Purpose:** Return the full status change timeline for an application.
**Authentication:** Bearer
**Ownership check:** required

**Response 200 OK:**
```json
{
  "data": [
    {
      "id": "uuid",
      "from_status": null,
      "to_status": { "id": "uuid", "label": "Bookmarked", "color_hex": "#9CA3AF" },
      "note": null,
      "changed_by": { "id": "uuid", "full_name": "Jane Doe" },
      "changed_at": "2026-05-28T10:00:00Z"
    },
    {
      "id": "uuid",
      "from_status": { "id": "uuid", "label": "Bookmarked" },
      "to_status": { "id": "uuid", "label": "Applied", "color_hex": "#3B82F6" },
      "note": "Submitted via company portal",
      "changed_by": { "id": "uuid", "full_name": "Jane Doe" },
      "changed_at": "2026-06-01T09:15:00Z"
    }
  ]
}
```

---

### Contact Endpoints

---

#### GET /api/v1/applications/{id}/contacts
**Purpose:** List all contacts for an application.
**Authentication:** Bearer
**Ownership check:** application must belong to user

**Response 200 OK:**
```json
{
  "data": [
    {
      "id": "uuid",
      "full_name": "Alex Rivera",
      "role": "Recruiter",
      "email": "alex@acme.example.com",
      "linkedin_url": "https://linkedin.com/in/alex-rivera",
      "notes": "Very responsive",
      "created_at": "2026-06-01T09:00:00Z"
    }
  ]
}
```

---

#### POST /api/v1/applications/{id}/contacts
**Purpose:** Add a contact to an application.
**Authentication:** Bearer

**Request Body:**
```json
{
  "full_name": "string (1–200 chars, required)",
  "role": "string (max 150 chars, optional)",
  "email": "string (valid email, optional)",
  "linkedin_url": "string (valid HTTPS URL containing linkedin.com, optional)",
  "notes": "string (max 5000 chars, optional)"
}
```

**Response 201 Created:** Contact object

---

#### PATCH /api/v1/applications/{id}/contacts/{contactId}
**Purpose:** Update a contact.
**Authentication:** Bearer

**Request Body (all optional):** Same fields as POST

**Response 200 OK:** Updated contact object

---

#### DELETE /api/v1/applications/{id}/contacts/{contactId}
**Purpose:** Permanently delete a contact.
**Authentication:** Bearer

**Response 204 No Content**

---

### Company Endpoints

---

#### GET /api/v1/companies
**Purpose:** List the current user's companies.
**Authentication:** Bearer

**Query Parameters:** `page`, `per_page`, `q` (name search), `sort`, `order`

**Response 200 OK:**
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Acme Corp",
      "website": "https://acme.example.com",
      "industry": "Software",
      "headquarters": "San Francisco, CA",
      "notes": "Dream company.",
      "application_count": 2,
      "created_at": "2026-05-01T00:00:00Z"
    }
  ],
  "pagination": { ... }
}
```

---

#### POST /api/v1/companies
**Purpose:** Create a company.
**Authentication:** Bearer

**Request Body:**
```json
{
  "name": "string (1–300 chars, required)",
  "website": "string (valid URL, optional)",
  "industry": "string (max 150 chars, optional)",
  "headquarters": "string (max 200 chars, optional)",
  "notes": "string (max 5000 chars, optional)"
}
```

**Response 201 Created:** Company object

---

#### GET /api/v1/companies/{id}
**Purpose:** Get a single company with all related applications (summary).
**Authentication:** Bearer
**Ownership check:** required

**Response 200 OK:** Company object with embedded `applications` array (summary only — id, job_title, current status)

---

#### PATCH /api/v1/companies/{id}
**Purpose:** Update a company.
**Authentication:** Bearer

**Request Body (all optional):** Same fields as POST

**Response 200 OK**

---

#### DELETE /api/v1/companies/{id}
**Purpose:** Soft-delete a company.
**Authentication:** Bearer
**Note:** Applications referencing this company retain their `company_id` (FK set null on hard delete only — soft delete leaves FK intact).

**Response 204 No Content**

---

### Reminder Endpoints

---

#### GET /api/v1/reminders
**Purpose:** List the current user's reminders.
**Authentication:** Bearer

**Query Parameters:**
- `is_sent` — boolean filter (default: false to show pending)
- `application_id` — filter by linked application
- `from`, `to` — ISO 8601 date range on `remind_at`

**Response 200 OK:**
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "Follow up with Alex",
      "body": "Send a thank-you email for the interview.",
      "remind_at": "2026-06-20T09:00:00Z",
      "is_sent": false,
      "application": {
        "id": "uuid",
        "job_title": "Senior Software Engineer",
        "company_name": "Acme Corp"
      },
      "created_at": "2026-06-17T14:00:00Z"
    }
  ],
  "pagination": { ... }
}
```

---

#### POST /api/v1/reminders
**Purpose:** Create a reminder.
**Authentication:** Bearer

**Request Body:**
```json
{
  "title": "string (1–300 chars, required)",
  "body": "string (max 2000 chars, optional)",
  "remind_at": "ISO 8601 (required, must be in the future)",
  "application_id": "uuid (optional)"
}
```

**Validation:** `remind_at` must be at least 60 seconds in the future.

**Response 201 Created:** Reminder object

---

#### PATCH /api/v1/reminders/{id}
**Purpose:** Update a pending reminder.
**Authentication:** Bearer
**Constraint:** Cannot update a reminder where `is_sent = true`.

**Request Body (all optional):**
```json
{
  "title": "string",
  "body": "string",
  "remind_at": "ISO 8601 (must be in the future)",
  "application_id": "uuid | null"
}
```

**Response 200 OK:** Updated reminder

**Error Responses:** 409 (reminder already sent)

---

#### DELETE /api/v1/reminders/{id}
**Purpose:** Delete a pending reminder.
**Constraint:** Cannot delete a reminder where `is_sent = true`.

**Response 204 No Content**

---

### Reference Data Endpoints

---

#### GET /api/v1/reference/statuses
**Purpose:** Return all application statuses (kanban pipeline columns) ordered by `sort_order`.
**Authentication:** Bearer

**Response 200 OK:**
```json
{
  "data": [
    {
      "id": "uuid",
      "label": "Bookmarked",
      "color_hex": "#9CA3AF",
      "sort_order": 1,
      "is_terminal": false
    },
    {
      "id": "uuid",
      "label": "Applied",
      "color_hex": "#3B82F6",
      "sort_order": 2,
      "is_terminal": false
    }
  ]
}
```

---

#### POST /api/v1/reference/statuses
**Purpose:** Admin creates a new pipeline stage.
**Authentication:** Bearer
**Roles:** admin

**Request Body:**
```json
{
  "label": "string (1–100 chars, unique, required)",
  "color_hex": "string (#RRGGBB format, required)",
  "sort_order": "integer (required)",
  "is_terminal": "boolean (default: false)"
}
```

**Response 201 Created**

---

#### POST /api/v1/reference/statuses/reorder
**Purpose:** Admin reorders kanban columns by providing a new ordered array of status IDs.
**Authentication:** Bearer
**Roles:** admin

**Request Body:**
```json
{
  "ordered_ids": ["uuid", "uuid", "uuid"]
}
```

**Validation:** `ordered_ids` must contain every non-terminal status ID exactly once.

**Response 204 No Content**

---

### Report Endpoints

---

#### GET /api/v1/reports/pipeline
**Purpose:** Return a count of the authenticated user's applications grouped by status.
**Authentication:** Bearer

**Response 200 OK:**
```json
{
  "data": {
    "total": 47,
    "by_status": [
      { "status": { "id": "uuid", "label": "Applied", "color_hex": "#3B82F6" }, "count": 18 },
      { "status": { "id": "uuid", "label": "Phone Screen", "color_hex": "#8B5CF6" }, "count": 9 },
      { "status": { "id": "uuid", "label": "Technical Interview", "color_hex": "#EC4899" }, "count": 4 },
      { "status": { "id": "uuid", "label": "Offer", "color_hex": "#10B981" }, "count": 1 },
      { "status": { "id": "uuid", "label": "Rejected", "color_hex": "#EF4444" }, "count": 12 }
    ],
    "generated_at": "2026-06-17T14:30:00Z"
  }
}
```

---

#### GET /api/v1/reports/trends
**Purpose:** Return application counts over time for charting.
**Authentication:** Bearer

**Query Parameters:**
- `period` — `week | month` (default: `month`)
- `from` — ISO 8601 start date (default: 6 months ago)
- `to` — ISO 8601 end date (default: today)

**Response 200 OK:**
```json
{
  "data": {
    "period": "month",
    "series": [
      { "label": "2026-01", "applications_added": 12, "offers": 0, "rejections": 3 },
      { "label": "2026-02", "applications_added": 19, "offers": 1, "rejections": 7 },
      { "label": "2026-03", "applications_added": 8, "offers": 0, "rejections": 2 }
    ]
  }
}
```

---

#### GET /api/v1/reports/export
**Purpose:** Export all of the user's application data.
**Authentication:** Bearer
**Rate Limiting:** 5 exports / 24 hours per user

**Query Parameters:**
- `format` — `json | csv` (default: `json`)

**Response 200 OK:**
- JSON: `Content-Type: application/json`
- CSV: `Content-Type: text/csv; charset=utf-8` with `Content-Disposition: attachment; filename="applications-export-2026-06-17.csv"`

```json
{
  "exported_at": "2026-06-17T14:30:00Z",
  "user_id": "uuid",
  "applications": [ ... full application objects with contacts ... ]
}
```

---

### Audit Endpoints

---

#### GET /api/v1/audit
**Purpose:** Query the audit log.
**Authentication:** Bearer
**Roles:** admin

**Query Parameters:**
- `user_id` — filter by actor
- `entity_type` — `user | application | company | contact | reminder | status`
- `action` — `create | update | delete | restore | login | logout | password_reset`
- `from`, `to` — ISO 8601 timestamp range
- `page`, `per_page`

**Response 200 OK:**
```json
{
  "data": [
    {
      "id": "uuid",
      "actor": { "id": "uuid", "full_name": "Jane Doe", "email": "jane@example.com" },
      "action": "update",
      "entity_type": "application",
      "entity_id": "uuid",
      "changes": {
        "before": { "status_id": "uuid-bookmarked" },
        "after": { "status_id": "uuid-applied" }
      },
      "ip_address": "203.0.113.42",
      "user_agent": "Mozilla/5.0 ...",
      "correlation_id": "uuid",
      "occurred_at": "2026-06-17T14:30:00Z"
    }
  ],
  "pagination": { ... }
}
```

---

### System & Admin Endpoints

---

#### GET /health
**Purpose:** Service liveness and readiness check. Used by load balancers, uptime monitors, and deployment pipelines.
**Authentication:** None

**Response 200 OK (healthy):**
```json
{
  "status": "ok",
  "version": "1.0.0",
  "checks": {
    "database": { "status": "ok", "latency_ms": 3 }
  },
  "timestamp": "2026-06-17T14:30:00Z"
}
```

**Response 503 Service Unavailable (degraded):**
```json
{
  "status": "degraded",
  "checks": {
    "database": { "status": "error", "error": "connection timeout" }
  },
  "timestamp": "2026-06-17T14:30:00Z"
}
```

---

#### GET /metrics
**Purpose:** Prometheus-format metrics for monitoring.
**Authentication:** Admin Bearer token or IP allowlist (configured at the reverse proxy)
**Format:** `text/plain; version=0.0.4`

**Exported metrics include:**
- `http_requests_total` (by method, path, status code)
- `http_request_duration_seconds` (histogram)
- `active_sessions_total`
- `reminder_worker_processed_total`
- `db_query_duration_seconds`

---

#### POST /api/v1/admin/jobs/purge-expired-tokens
**Purpose:** Trigger cleanup of expired refresh tokens.
**Authentication:** Bearer
**Roles:** admin

**Response 202 Accepted:**
```json
{
  "data": {
    "message": "Token purge job enqueued.",
    "job_id": "uuid"
  }
}
```

---

## 8. Security Architecture

### Authentication Strategy

**JWT (RS256)**
- Access tokens: 15-minute TTL, RS256 signed (asymmetric key pair)
- Private key stored as environment variable / secrets manager — never in code or database
- Public key exposed at `GET /.well-known/jwks.json` for potential future OIDC compatibility
- JWT claims: `{ iss, sub, email, role, iat, exp, jti }`
- `jti` (JWT ID) allows individual token revocation if a deny-list is maintained; for v1, rely on short TTL

**Refresh Token Strategy**
- 256-bit cryptographically random tokens (not JWTs)
- Stored as `SHA-256(token)` in the database — raw token only exists in memory during issuance
- 7-day TTL with rotation on every use (detect token reuse as a security signal)
- Stored in `HttpOnly; Secure; SameSite=Strict` cookie AND returned in response body to support native clients
- Refresh token reuse (presenting a previously rotated token) triggers immediate revocation of the entire family

**Session Management**
- `refresh_tokens` table: `id, user_id, token_hash, family_id, expires_at, revoked_at, created_at, ip_address, user_agent`
- `family_id` groups related tokens so reuse detection can revoke the entire chain

**MFA:** Out of scope for v1. Architecture accommodates it: add `mfa_enabled` to users, TOTP verification endpoint, and a `mfa_verified` claim in the JWT.

---

### Rate Limiting Strategy

| Endpoint | Limit | Window | Scope |
|----------|-------|--------|-------|
| POST /auth/login | 10 req | 15 min | per IP |
| POST /auth/login | 5 req | 15 min | per email |
| POST /auth/register | 5 req | 1 hour | per IP |
| POST /auth/password/reset-request | 3 req | 10 min | per email |
| POST /auth/password/reset | 5 req | 1 hour | per IP |
| GET /reports/export | 5 req | 24 hours | per user |
| All other authenticated endpoints | 200 req | 1 min | per user |
| All unauthenticated endpoints | 60 req | 1 min | per IP |

Implementation: Token bucket algorithm. Store counters in Redis or PostgreSQL. Return 429 with `Retry-After` header.

---

### Data Encryption Requirements

| Data | At Rest | In Transit |
|------|---------|------------|
| User passwords | Argon2id (memory=64MB, time=3, threads=4) | HTTPS only |
| Refresh tokens | SHA-256 hash only | HTTPS only |
| User email | AES-256-GCM (application-level encryption) | HTTPS only |
| Contact email | AES-256-GCM (application-level encryption) | HTTPS only |
| All other fields | PostgreSQL disk encryption (if host supports it) | HTTPS only |
| JWT private key | Secrets manager / env var | Never transmitted |

**Application-level encryption for PII:**
- Encryption key derived from a master key via HKDF
- Key rotation: add new key version, re-encrypt on next write, background job migrates old records
- Encrypted column stored as `bytea`; plaintext never written to disk

---

### OWASP API Security Top 10 Mitigations

| OWASP Risk | Mitigation |
|------------|-----------|
| API1: Broken Object Level Authorization | Every query filters by `user_id = authenticated user`; ownership check before every mutation |
| API2: Broken Authentication | Short-lived JWTs, refresh token rotation, reuse detection, rate limiting on auth endpoints |
| API3: Broken Object Property Level Authorization | Allowlist-based field updates (PATCH ignores undeclared fields); passwords never returned in responses |
| API4: Unrestricted Resource Consumption | Pagination enforced (max 100), rate limiting, export rate limited, query complexity bounded |
| API5: Broken Function Level Authorization | Role check middleware on every admin endpoint; enforced in route definition, not just handler logic |
| API6: Unrestricted Access to Sensitive Business Flows | Export rate limiting; account deletion requires password + typed confirmation |
| API7: Server Side Request Forgery | `job_url` and `website` fields never fetched server-side — stored as strings only |
| API8: Security Misconfiguration | Secrets via env vars / secrets manager; debug endpoints disabled in production; CORS allowlist |
| API9: Improper Inventory Management | Single versioned API (`/v1`); deprecated endpoints return `Sunset` header before removal |
| API10: Unsafe Consumption of APIs | No third-party API calls in v1; all data user-supplied |

---

### API Gateway & Infrastructure Security

**Reverse Proxy (nginx / Caddy):**
- TLS termination (TLS 1.3 minimum; disable TLS 1.0/1.1/1.2)
- HSTS header: `Strict-Transport-Security: max-age=63072000; includeSubDomains; preload`
- Security headers: `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Content-Security-Policy`, `Referrer-Policy: strict-origin-when-cross-origin`
- CORS: allowlist of trusted origins only (no wildcard `*`)
- Request size limit: 1MB for API requests; reject larger bodies at proxy level
- `/metrics` endpoint restricted to admin subnet or IP allowlist at the proxy

**WAF Rules (if using a WAF layer):**
- Block SQL injection patterns
- Block XSS payloads in JSON bodies
- Block requests with excessively long headers
- Geo-block or challenge on anomalous request patterns

**DDoS Protection:**
- Connection rate limiting at the proxy level
- SYN flood protection at the host network level
- Consider Cloudflare or similar for public-facing deployments

---

### Secrets Management

| Secret | Storage | Rotation |
|--------|---------|----------|
| JWT private key | Environment variable or secrets manager (e.g. HashiCorp Vault) | Annually or on compromise |
| Database credentials | Environment variable | Quarterly |
| Email sender credentials | Environment variable | Quarterly |
| PII encryption master key | Secrets manager with KMS | Annually or on compromise |

**Never:** commit secrets to source control; log secrets; include secrets in error messages.

---

## 9. Audit & Compliance Requirements

### Audit Log Schema

Every write operation records an entry in the `audit_log` table:

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | PK |
| `actor_id` | UUID | FK → users (nullable for system actions) |
| `actor_email` | varchar | Snapshot of email at time of action |
| `action` | varchar | `create`, `update`, `delete`, `restore`, `login`, `logout`, `password_change`, `password_reset`, `role_change`, `account_deactivate`, `export` |
| `entity_type` | varchar | `user`, `application`, `company`, `contact`, `reminder`, `session`, `status` |
| `entity_id` | UUID | ID of the affected record |
| `changes_before` | jsonb | Snapshot of fields before mutation (PII excluded) |
| `changes_after` | jsonb | Snapshot of fields after mutation (PII excluded) |
| `ip_address` | inet | Client IP (from proxy forwarded header) |
| `user_agent` | text | Browser/client identifier |
| `correlation_id` | UUID | Request correlation ID |
| `occurred_at` | timestamptz | Timestamp of the action |

**PII handling in audit log:** Emails and other PII fields are excluded from `changes_before` / `changes_after`. Only non-PII change data is recorded.

**Retention:** Minimum 12 months. Implement PostgreSQL table partitioning by month for efficient purging.

**Immutability:** Audit log rows are never updated or deleted during the retention window. Use PostgreSQL row-level security to enforce this.

### Events That Must Be Audited

| Event | Severity |
|-------|----------|
| Successful login | Info |
| Failed login (after 3 failures on same account) | Warning |
| Password change | High |
| Password reset | High |
| Email change | High |
| Role assignment | Critical |
| Account deactivation | Critical |
| Account deletion request | Critical |
| Admin access to any user's data | High |
| Data export | High |
| All sessions revoked | High |
| Admin force-revoke sessions | Critical |
| Application status change | Info |

### GDPR Compliance

**Right to Access (Art. 15):** The `/reports/export` endpoint fulfils this — full JSON export of all user data.

**Right to Erasure (Art. 17):** `DELETE /api/v1/users/me` initiates erasure. A background job within 30 days:
1. Scrubs PII from `users` table (name → `[deleted]`, email → randomised token)
2. Scrubs PII from `contacts` table (name, email → `[deleted]`)
3. Soft-deletes all applications and companies
4. Revokes all sessions
5. Excludes the user from future audit log queries (actor displayed as `[deleted user]`)

**Data Minimisation:** Contact email and user email are the only PII fields stored. They are encrypted at rest and excluded from audit logs.

---

## 10. Architectural Assumptions

1. **Single-tenant, self-hosted deployment.** The system is designed for one instance per user or household. Multi-tenancy (multiple isolated organisations on one instance) is not in scope and has not been architected for.

2. **Go + PostgreSQL 18 stack.** All database-specific features (generated columns for FTS, `gen_random_uuid()`, jsonb, row-level security, partitioning) are written for PostgreSQL 18. No ORM — `sqlc` generates type-safe query code from plain SQL.

3. **No real-time features in v1.** WebSocket or SSE for live kanban updates is deferred. Clients poll or refresh to see changes.

4. **Email delivery is external.** Password reset and reminder notification emails are delivered via an SMTP relay (e.g. local Postfix or a transactional email service). The API triggers sends; it does not implement the transport.

5. **Background worker is in-process.** The reminder worker runs as a goroutine within the same Go process (polling on a ticker). A separate worker service can be extracted later if scale demands it.

6. **No full OAuth 2.1 / OIDC provider in v1.** The system implements its own auth (JWT + refresh tokens). OIDC compatibility (`.well-known/jwks.json`) is architected in to allow a future upgrade to an identity provider (e.g. Keycloak, Zitadel).

7. **File uploads not in scope for v1.** Resumes, cover letters, and attachments are deferred. The `job_url` field stores a link to the job posting; no binary storage is implemented.

8. **Dashboard view style (kanban vs table) is a frontend concern.** The API serves the same data regardless of view. The `/reference/statuses` endpoint with `sort_order` supports kanban column ordering without baking it into the API response shape.

9. **Soft delete retention is 30 days.** Applications and companies are hard-purged after 30 days in soft-deleted state. This is enforced by a background job, not a database trigger.

10. **Admin role is granted at deployment.** The first user registered after a fresh deployment is automatically granted `admin` role (configured via a `FIRST_USER_IS_ADMIN=true` environment variable). Subsequent self-registrations receive `user` role.
