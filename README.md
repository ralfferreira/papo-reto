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

## Contributing

For setup and development instructions, please refer to:
- [Backend README](/backend/README.md)
- [Frontend README](/frontend/README.md)

## Contributors

- [Ralf Dewrich](https://github.com/ralfferreira)
- [Ian Bittencourt](https://github.com/ianbitt)

## License

This project is licensed under the MIT License - see the LICENSE file for details.