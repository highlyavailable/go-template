# Web Application Features

This document describes the new web application features added to the Go template project.

## GORM Models and Migrations

### Models Structure
Created a comprehensive model structure in `internal/models/`:
- **BaseModel**: Common fields (ID, timestamps, soft delete)
- **User**: User management with authentication fields
- **Post**: Blog/article content management
- **Comment**: Nested commenting system
- **Tag**: Tag system with many-to-many relationships

### Migration System
- Automated migration system in `internal/db/migrations/`
- Auto-migration for all models on startup
- Custom index creation for performance optimization
- Support for soft deletes

### Usage Example
```go
// Run migrations
migrator := migrations.NewMigrator(db.DB())
if err := migrator.AutoMigrate(); err != nil {
    log.Fatal("Failed to run migrations:", err)
}
```

## Web UI with Gin/HTMX/Templ/Tailwind

### Template Structure
Created a modern web UI using:
- **Templ**: Type-safe Go templates
- **HTMX**: Dynamic updates without JavaScript
- **Tailwind CSS**: Utility-first CSS framework
- **Alpine.js**: Lightweight JavaScript for interactivity

### Directory Structure
```
web/
├── templates/
│   ├── layout.templ         # Base layout with common structure
│   ├── components/          # Reusable components
│   │   ├── navbar.templ     # Top navigation
│   │   ├── sidebar.templ    # Side navigation
│   │   └── footer.templ     # Footer component
│   ├── pages/               # Full page templates
│   │   ├── home.templ       # Dashboard/home page
│   │   └── posts.templ      # Posts listing page
│   └── partials/            # HTMX partial templates
│       ├── activity_feed.templ
│       ├── notifications.templ
│       └── user_menu.templ
└── static/
    ├── css/style.css        # Custom styles
    └── js/app.js            # JavaScript utilities
```

### Key Features

#### 1. Dynamic Layout System
- Base layout with navbar, sidebar, main content, and footer
- Responsive design with mobile support
- Dynamic content areas for HTMX updates

#### 2. HTMX Integration
- Partial template loading without page refresh
- Activity feed with real-time updates
- Notification system with mark-as-read functionality
- User menu dropdown

#### 3. Component-Based Architecture
- Reusable Templ components
- Type-safe template rendering
- Server-side rendering with dynamic updates

### Routes Configuration
```go
// Web routes
router.GET("/", homeHandler.Index)
router.GET("/posts", postsHandler.Index)

// Partial routes for HTMX
partials := router.Group("/partials")
{
    partials.GET("/activity-feed", partialsHandler.ActivityFeed)
    partials.GET("/notifications", partialsHandler.Notifications)
    partials.GET("/user-menu", partialsHandler.UserMenu)
    partials.POST("/notifications/:id/read", partialsHandler.MarkNotificationRead)
}

// Static files
router.Static("/static", "./web/static")
```

## Running the Application

1. Generate Templ files:
```bash
templ generate
```

2. Run the application:
```bash
go run cmd/goapp/main.go
```

3. Access the web UI:
- Home: http://localhost:8080/
- Posts: http://localhost:8080/posts

## Testing

All components include comprehensive tests:
- Model tests with SQLite in-memory database
- Web handler tests with test containers
- Template rendering tests

Run tests:
```bash
go test ./...
```

## Development Workflow

1. **Adding New Models**:
   - Create model in `internal/models/`
   - Add to migrator in `internal/db/migrations/migrator.go`
   - Create tests

2. **Adding New Pages**:
   - Create template in `web/templates/pages/`
   - Create handler in `api/handlers/web/`
   - Add route in `api/routes/routes.go`
   - Run `templ generate`

3. **Adding HTMX Partials**:
   - Create partial template in `web/templates/partials/`
   - Add handler method in `partials_handler.go`
   - Add route to partials group
   - Use `hx-get` in parent template

## Best Practices

1. **Models**:
   - Use GORM hooks for validation
   - Implement soft deletes
   - Create appropriate indexes

2. **Templates**:
   - Keep components small and reusable
   - Use type-safe Templ syntax
   - Leverage HTMX for dynamic updates

3. **Handlers**:
   - Use dependency injection via container
   - Proper error handling and logging
   - Write comprehensive tests

4. **Frontend**:
   - Use Tailwind utility classes
   - Implement loading states for HTMX
   - Ensure accessibility

## Next Steps

Consider adding:
- User authentication and sessions
- Form handling with validation
- File upload capabilities
- WebSocket support for real-time features
- API endpoints for the web UI
- Production build process for assets