# Papo Reto - Anonymous Messaging Platform

Papo Reto is a modern web platform for receiving and managing anonymous messages through personalized links. Create custom groups, share across multiple platforms, and interact with your audience in real-time.

## About the Project

The platform allows users to:

- Create custom groups with unique links to receive anonymous messages
- Organize messages efficiently with favorites, unread markers, and more
- Share links across multiple platforms (Instagram, X, WhatsApp, etc.)
- Receive real-time notifications when new messages arrive
- Customize the appearance of forms and share pages
- Analyze message patterns with advanced statistics (Premium)

## Project Structure

This is a monorepo containing both the backend (Go) and frontend (Next.js) applications.

```
/
├── backend/       # Go backend API
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   └── README.md  # Backend-specific documentation and setup
│
├── frontend/      # Next.js frontend application
│   ├── public/
│   ├── src/
│   └── README.md  # Frontend-specific documentation and setup
│
└── README.md      # This file
```

## Tech Stack

### Backend
- **Language**: Go
- **Framework**: Gin/Echo
- **Database**: PostgreSQL/Redis
- **Real-time**: WebSockets
- **Authentication**: JWT

### Frontend
- **Framework**: Next.js 14+
- **Language**: TypeScript
- **Styling**: TailwindCSS
- **State Management**: Zustand/React Query
- **Authentication**: NextAuth.js

## Installation Guide

### Prerequisites

- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/install/)
- [Node.js](https://nodejs.org/) (v18 or later) and [npm](https://www.npmjs.com/) (for local development)
- [Go](https://golang.org/) (v1.20 or later, for local development)

### Option 1: Using Docker (Recommended for Production)

1. **Clone the repository**
   ```bash
   git clone https://github.com/ralfferreira/papo-reto.git
   cd papo-reto
   ```

2. **Start the entire stack with Docker Compose**
   ```bash
   docker-compose up -d
   ```

   This will start:
   - Frontend on http://localhost:3000
   - Backend API on http://localhost:8080
   - PostgreSQL database
   - Redis cache

3. **View logs**
   ```bash
   # All services
   docker-compose logs -f

   # Specific service
   docker-compose logs -f frontend
   docker-compose logs -f api
   ```

### Option 2: Local Development

#### Backend Setup

1. **Navigate to the backend directory**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   # Copy example env file
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start the database and Redis with Docker**
   ```bash
   docker-compose up -d postgres redis
   ```

5. **Run the backend server**
   ```bash
   go run cmd/api/main.go
   ```

   The API will be available at http://localhost:8080

#### Frontend Setup

1. **Navigate to the frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Set up environment variables**
   ```bash
   # Copy example env file
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the development server**
   ```bash
   npm run dev
   ```

   The frontend will be available at http://localhost:3000

## Features

- **Group Organization**: Create separate contexts for different types of feedback
- **Real-time Messaging**: Instant notifications when new messages arrive
- **Multi-platform Sharing**: Share your link on any social media or messaging platform
- **Identity Reveal**: Optional feature for anonymous senders to reveal identity
- **Powerful Management**: Favorite messages, mark as unread, and organize efficiently
- **Social Content Export**: Turn messages into shareable content for stories, tweets, etc.
- **Icebreakers**: Pre-defined questions to encourage engagement

## Roadmap

- [x] Basic user authentication
- [x] Message group creation
- [x] Anonymous message sending
- [ ] Real-time notifications
- [ ] Advanced sharing options
- [ ] Message exporting to social media
- [ ] Premium subscription features
- [ ] Mobile-responsive design
- [ ] Advanced analytics

## Usage Guide

### Creating an Account

1. Visit http://localhost:3000
2. Click "Register" and fill out the form
3. Verify your email (if enabled)
4. Log in with your credentials

### Creating a Message Group

1. Navigate to the Dashboard
2. Click "New Group"
3. Fill in the group details and customize settings
4. Save the group to generate a unique link

### Sharing Your Link

1. From the Dashboard, select a group
2. Click "Share" to view sharing options
3. Copy the link or use the social sharing buttons
4. Share the link on your preferred platforms

### Managing Messages

1. Messages appear in your Dashboard as they arrive
2. Use filters to sort and organize messages
3. Star important messages or mark as read/unread
4. Reply to messages (if enabled)

## Deployment

### Using Docker (Self-hosted)

Follow the Docker installation steps above on your server.

### Environment Configuration

For production deployment, make sure to set appropriate environment variables:

- Set `NODE_ENV=production`
- Configure secure database credentials
- Set up proper JWT secrets
- Configure email providers (if using email notifications)

## Contributing

For setup and development instructions, please refer to:
- [Backend README](/backend/README.md)
- [Frontend README](/frontend/README.md)

## Contributors

- [Ralf Dewrich](https://github.com/ralfferreira)
- [Ian Bittencourt](https://github.com/ianbitt)

## License

This project is licensed under the MIT License - see the LICENSE file for details.