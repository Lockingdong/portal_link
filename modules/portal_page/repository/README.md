# Portal Page Repository

This directory contains different implementations of the `PortalPageRepository` interface for the Portal Link application.

## Available Implementations

### 1. Database Repository (`portal_page_repository.go`)
The production implementation that uses PostgreSQL with SQLBoiler for persistence.

**Features:**
- Full database persistence with transactions
- Supports all CRUD operations with proper SQL handling
- Handles link relationships and cascade operations
- Production-ready with error handling

**Usage:**
```go
db := pkg.NewPG(config.GetDBConfig().DSN())
repo := repository.NewPortalPageRepository(db)
```

### 2. In-Memory Repository (`in_memory_portal_page_repository.go`)
A thread-safe in-memory implementation perfect for testing and local development.

**Features:**
- Zero dependencies (no database required)
- Thread-safe with proper locking
- Auto-increment ID generation
- Deep copying to prevent reference issues
- Proper sorting of links by DisplayOrder
- Error handling that mimics database behavior
- Slug uniqueness validation
- User indexing for efficient lookups

**Usage:**
```go
repo := repository.NewInMemoryPortalPageRepository()
```

## When to Use Each Implementation

### Use Database Repository When:
- Running in production environment
- Need persistent data storage
- Working with real user data
- Deploying to staging/production servers

### Use In-Memory Repository When:
- Writing unit tests
- Local development without database setup
- Integration testing in CI/CD pipelines
- Prototyping new features
- Running quick experiments or demos

## Interface Compatibility

Both implementations satisfy the same `domain.PortalPageRepository` interface:

```go
type PortalPageRepository interface {
    Create(ctx context.Context, portalPage *PortalPage) error
    Update(ctx context.Context, portalPage *PortalPage) error
    FindBySlug(ctx context.Context, slug string) (*PortalPage, error)
    ListByUserID(ctx context.Context, userID int) ([]*PortalPage, error)
    FindByID(ctx context.Context, id int) (*PortalPage, error)
}
```

This means you can easily swap between implementations without changing your business logic.

## Example Usage

See `examples/in_memory_repository_example.go` for a complete demonstration of the in-memory repository capabilities.

To run the example:
```bash
go run examples/in_memory_repository_example.go
```

## Testing

Both repositories have comprehensive test suites:

- `portal_page_repository_test.go` - Tests for database implementation (requires DB)
- `in_memory_portal_page_repository_test.go` - Tests for in-memory implementation

Run in-memory tests (no dependencies):
```bash
go test -v ./modules/portal_page/repository/ -run="TestInMemoryPortalPageRepository"
```

## Architecture Benefits

This dual implementation approach provides several advantages:

1. **Development Velocity**: Developers can work without database setup
2. **Test Reliability**: Tests don't depend on external services
3. **CI/CD Efficiency**: Faster test execution in pipelines
4. **Deployment Flexibility**: Easy environment-specific configuration
5. **Code Quality**: Interface-driven design promotes better architecture