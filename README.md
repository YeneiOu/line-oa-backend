# LINE OA Backend

A Go Fiber backend application for LINE Official Account integration with OAuth authentication and booking system with LINE notifications.

## Features

- **LINE OAuth Integration**: Login with LINE account
- **JWT Authentication**: Secure token-based authentication
- **Booking System**: Create, read, update, and cancel bookings
- **LINE Notifications**: Automatic push messages to users via LINE
- **MongoDB Database**: Persistent data storage with official MongoDB driver
- **RESTful API**: Clean API endpoints for frontend integration

## Architecture

This backend implements the sequence diagram flow:
1. **Part 1**: LINE account connection (one-time setup)
2. **Part 2**: Booking creation and LINE notifications (recurring)

## Prerequisites

- Go 1.21 or higher
- MongoDB database (or Docker)
- LINE Developer Account with:
  - LINE Login Channel (for OAuth)
  - LINE Messaging API Channel (for notifications)

## Setup

### 1. Clone and Install Dependencies

```bash
cd line-oa-backend
go mod tidy
```

### 2. Database Setup

#### Option A: Using Docker (Recommended)

Start MongoDB using Docker Compose:

```bash
docker-compose up -d
```

This will start MongoDB with the credentials specified in `docker-compose.yml`.

#### Option B: Local MongoDB Installation

Install MongoDB locally and create a database named `line_oa_backend`.

### 3. Environment Configuration

Copy the example environment file and configure your settings:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# MongoDB Configuration
MONGO_URI=mongodb://your_root_user:your_secret_password@localhost:27017
MONGO_DATABASE=line_oa_backend

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here

# LINE OAuth Configuration (from LINE Login Channel)
LINE_CHANNEL_ID=your-line-login-channel-id
LINE_CHANNEL_SECRET=your-line-login-channel-secret
LINE_REDIRECT_URI=http://localhost:3000/callback

# LINE Messaging API Configuration (from LINE Messaging API Channel)
LINE_CHANNEL_ACCESS_TOKEN=your-line-messaging-channel-access-token

# Server Configuration
PORT=8080
FRONTEND_URL=http://localhost:3000
```

### 4. LINE Developer Console Setup

#### LINE Login Channel:
1. Go to [LINE Developers Console](https://developers.line.biz/)
2. Create a new LINE Login channel
3. Set callback URL: `http://localhost:3000/callback`
4. Copy Channel ID and Channel Secret to `.env`

#### LINE Messaging API Channel:
1. Create a new Messaging API channel
2. Generate Channel Access Token
3. Copy the token to `.env`

## Running the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /api/v1/auth/login` - Get LINE OAuth URL
- `POST /api/v1/auth/callback` - Handle LINE OAuth callback
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `GET /api/v1/me` - Get current user info (protected)

### Bookings

- `POST /api/v1/bookings` - Create new booking (protected)
- `GET /api/v1/bookings` - Get user's bookings (protected)
- `GET /api/v1/bookings/:id` - Get specific booking (protected)
- `PUT /api/v1/bookings/:id` - Update booking (protected)
- `DELETE /api/v1/bookings/:id` - Cancel booking (protected)

### Health Check

- `GET /health` - Server health status

## API Usage Examples

### 1. Login Flow

```javascript
// Step 1: Get LINE OAuth URL
const loginResponse = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST'
});
const { auth_url, state } = await loginResponse.json();

// Step 2: Redirect user to auth_url
// User will be redirected back to your frontend with code parameter

// Step 3: Exchange code for JWT token
const callbackResponse = await fetch('http://localhost:8080/api/v1/auth/callback', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ code: 'received_code', state: 'received_state' })
});
const { token, user } = await callbackResponse.json();
```

### 2. Create Booking

```javascript
const bookingResponse = await fetch('http://localhost:8080/api/v1/bookings', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    service_name: 'Hair Cut',
    booking_date: '2024-01-15T10:00:00Z',
    notes: 'Please call before arrival'
  })
});
const { booking, message } = await bookingResponse.json();
```

## Project Structure

```
line-oa-backend/
├── config/          # Configuration management
├── database/        # Database connection and setup
├── handlers/        # HTTP request handlers
├── middleware/      # HTTP middleware (auth, CORS)
├── models/          # Database models
├── services/        # Business logic services
├── main.go          # Application entry point
├── go.mod           # Go module dependencies
└── README.md        # This file
```

## Database Schema

### Users Collection
- `_id` (ObjectID, Primary Key)
- `line_user_id` (String, Unique, from LINE)
- `name` (String, Display name from LINE)
- `email` (String, Optional)
- `picture_url` (String, Profile picture from LINE)
- `created_at`, `updated_at` (DateTime)

### Bookings Collection
- `_id` (ObjectID, Primary Key)
- `user_id` (ObjectID, Reference to Users)
- `service_name` (String)
- `booking_date` (DateTime)
- `notes` (String)
- `status` (String: confirmed, cancelled, completed)
- `created_at`, `updated_at` (DateTime)

## Security Features

- JWT token authentication
- CORS protection
- Input validation
- NoSQL injection protection (via MongoDB driver)
- Secure token generation

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Error message description"
}
```

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o line-oa-backend main.go
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check MongoDB is running (`docker-compose ps` if using Docker)
   - Verify MongoDB URI in `.env`
   - Ensure MongoDB is accessible on the specified port

2. **LINE OAuth Error**
   - Verify LINE Channel ID and Secret
   - Check callback URL matches LINE Console settings
   - Ensure LINE Login channel is properly configured

3. **LINE Messaging Failed**
   - Verify Channel Access Token
   - Check if Messaging API channel is active
   - Ensure user has added your LINE Official Account as friend

### Logs

The application logs important events including:
- Database connections
- Authentication attempts
- LINE API calls
- Booking operations

## License

MIT License
# line-oa-backend
# line-oa-backend
