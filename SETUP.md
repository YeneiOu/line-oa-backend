# LINE OA Backend Setup Guide

## Prerequisites

1. Go 1.21 or later
2. MongoDB Atlas account (or local MongoDB)
3. LINE Developer account

## Environment Setup

1. **Create your `.env` file** (copy from `.env.example`):
   ```bash
   cp .env.example .env
   ```

2. **Configure MongoDB Atlas**:
   - Go to [MongoDB Atlas](https://cloud.mongodb.com/)
   - Create a cluster if you haven't already
   - Get your connection string (should look like):
     ```
     mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority
     ```
   - Update `MONGO_URI` in your `.env` file

3. **Configure LINE Developer Settings**:
   - Go to [LINE Developers Console](https://developers.line.biz/)
   - Create a new channel or use existing one
   - Get the following values and update your `.env`:
     - `LINE_CHANNEL_ID`
     - `LINE_CHANNEL_SECRET`
     - `LINE_CHANNEL_ACCESS_TOKEN`

4. **Generate JWT Secret**:
   ```bash
   # Generate a secure random string
   openssl rand -base64 32
   ```
   Update `JWT_SECRET` in your `.env` file

## MongoDB Atlas Connection Issues

If you're getting TLS errors with MongoDB Atlas:

1. **Check your connection string format**:
   ```
   mongodb+srv://username:password@cluster.mongodb.net/database?retryWrites=true&w=majority
   ```

2. **Ensure your IP is whitelisted**:
   - Go to Atlas Dashboard → Network Access
   - Add your current IP address or use `0.0.0.0/0` for development

3. **Verify credentials**:
   - Username and password are correct
   - User has proper database permissions

## Running the Application

### Local Development
```bash
# Install dependencies
go mod download

# Run the application
go run _cmd/main.go
```

### Docker
```bash
# Build the image
docker build -t line-oa-backend .

# Run with environment file
docker run --env-file .env -p 8080:8080 line-oa-backend
```

### Docker Compose (with MongoDB)
```bash
# Start services
docker-compose up -d

# For local MongoDB, use this connection string:
# MONGO_URI=mongodb://your_root_user:your_secret_password@localhost:27017/line_oa_backend?authSource=admin
```

## Health Check

Once running, test the health endpoint:
```bash
curl http://localhost:8080/health
```

## Troubleshooting

### MongoDB Connection Issues
- Verify your Atlas cluster is running
- Check network access settings in Atlas
- Ensure correct username/password
- Try connecting with MongoDB Compass first

### LINE API Issues
- Verify channel credentials
- Check webhook URL configuration
- Ensure proper permissions are set
