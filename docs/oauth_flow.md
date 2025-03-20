# OAuth Authentication Flow

## 1. OAuth Login Flow

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant Backend
    participant OAuth Provider

    User->>Frontend: Click Login with OAuth
    Frontend->>Backend: GET /api/v1/oauth/login/:provider
    Backend->>Backend: Generate State Token
    Backend->>Frontend: Set State Cookie & Redirect
    Frontend->>OAuth Provider: Redirect to Provider's Consent Page
    OAuth Provider->>User: Show Consent Screen
    User->>OAuth Provider: Grant Permissions
    OAuth Provider->>Frontend: Redirect with Auth Code
    Frontend->>Backend: GET /api/v1/oauth/callback/:provider
    Backend->>Backend: Verify State Token
    Backend->>OAuth Provider: Exchange Code for Token
    OAuth Provider->>Backend: Return Access Token
    Backend->>OAuth Provider: Get User Profile
    OAuth Provider->>Backend: Return User Data
    Backend->>Backend: Create/Update User
    Backend->>Frontend: Return Success & User Data
```

## 2. Token Refresh Flow

```mermaid
sequenceDiagram
    participant Frontend
    participant Backend
    participant OAuth Provider

    Frontend->>Backend: Request Protected Resource
    Backend->>Backend: Check Token Expiry
    Backend->>OAuth Provider: Request Token Refresh
    OAuth Provider->>Backend: Return New Access Token
    Backend->>Backend: Update Token in DB
    Backend->>Frontend: Return Protected Resource
```

## 3. Account Unlinking Flow

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant Backend
    participant OAuth Provider

    User->>Frontend: Request Unlink Account
    Frontend->>Backend: POST /api/v1/oauth/unlink/:provider
    Backend->>OAuth Provider: Revoke Access Token
    Backend->>Backend: Delete OAuth Records
    Backend->>Frontend: Return Success
```

## 4. Data Model

```mermaid
erDiagram
    User ||--o{ OAuthProfile : has
    User ||--o{ OAuthToken : has

    OAuthProfile {
        string id PK
        string providerId
        string provider
        string email
        boolean emailVerified
        string name
        string firstName
        string lastName
        string picture
        string locale
        datetime createdAt
        datetime updatedAt
    }

    OAuthToken {
        string id PK
        string userId FK
        string provider
        string accessToken
        string tokenType
        string refreshToken
        int expiresIn
        string scope
        datetime createdAt
        datetime updatedAt
    }
```

## 5. Component Architecture

```mermaid
graph TB
    subgraph Frontend
        UI[User Interface]
        Auth[Auth Module]
    end

    subgraph Backend
        API[API Layer]
        Handler[OAuth Handler]
        Service[OAuth Service]
        Repo[OAuth Repository]
    end

    subgraph External
        Google[Google OAuth]
        Facebook[Facebook OAuth]
        MongoDB[(MongoDB)]
    end

    UI --> Auth
    Auth --> API
    API --> Handler
    Handler --> Service
    Service --> Repo
    Service --> Google
    Service --> Facebook
    Repo --> MongoDB
```

## Security Considerations

1. **State Parameter**

   - Prevents CSRF attacks
   - Unique per request
   - Short expiration time

2. **Token Storage**

   - Access tokens stored encrypted
   - Refresh tokens with additional encryption
   - Regular token rotation

3. **Error Handling**

   - Failed attempts logging
   - Rate limiting
   - IP-based blocking

4. **Data Protection**
   - HTTPS everywhere
   - Minimal scope requests
   - Data encryption at rest

## Implementation Checklist

- [x] Basic OAuth Flow
- [x] Token Management
- [x] Profile Management
- [x] Error Handling
- [x] Security Measures
- [ ] Rate Limiting
- [ ] Token Rotation
- [ ] Monitoring
- [ ] Analytics

## API Endpoints

### OAuth Login

```http
GET /api/v1/oauth/login/:provider
```

### OAuth Callback

```http
GET /api/v1/oauth/callback/:provider
```

### Unlink Account

```http
POST /api/v1/oauth/unlink/:provider
```

## Environment Variables

```env
# OAuth Configuration - Google
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/oauth/callback/google

# OAuth Configuration - Facebook
FACEBOOK_CLIENT_ID=your_facebook_client_id
FACEBOOK_CLIENT_SECRET=your_facebook_client_secret
FACEBOOK_REDIRECT_URL=http://localhost:8080/api/v1/oauth/callback/facebook
```
