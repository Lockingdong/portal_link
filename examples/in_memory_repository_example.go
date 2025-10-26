package main

import (
	"context"
	"fmt"
	"log"
	"portal_link/modules/portal_page/domain"
	"portal_link/modules/portal_page/repository"
	"time"
)

// Example demonstrates how to use the InMemoryPortalPageRepository for local testing
func main() {
	// Create an in-memory repository instance
	repo := repository.NewInMemoryPortalPageRepository()
	ctx := context.Background()

	fmt.Println("=== In-Memory Portal Page Repository Example ===")

	// Example 1: Create a portal page with links
	fmt.Println("\n1. Creating a portal page with links...")
	
	portalPage := &domain.PortalPage{
		UserID:          1,
		Slug:            "john-doe",
		Title:           "John Doe's Links",
		Bio:             "Software Developer & Open Source Enthusiast",
		ProfileImageURL: "https://example.com/john-avatar.jpg",
		Theme:           domain.ThemeLight,
		Links: []*domain.Link{
			{
				Title:        "GitHub",
				URL:          "https://github.com/johndoe",
				Description:  "My GitHub Profile",
				IconURL:      "https://github.com/favicon.ico",
				DisplayOrder: 1,
			},
			{
				Title:        "Twitter",
				URL:          "https://twitter.com/johndoe",
				Description:  "Follow me on Twitter",
				IconURL:      "https://twitter.com/favicon.ico",
				DisplayOrder: 2,
			},
			{
				Title:        "Personal Website",
				URL:          "https://johndoe.dev",
				Description:  "My personal portfolio",
				DisplayOrder: 3,
			},
		},
	}

	err := repo.Create(ctx, portalPage)
	if err != nil {
		log.Fatalf("Failed to create portal page: %v", err)
	}

	fmt.Printf("✓ Created portal page with ID: %d\n", portalPage.ID)
	fmt.Printf("  - Title: %s\n", portalPage.Title)
	fmt.Printf("  - Slug: %s\n", portalPage.Slug)
	fmt.Printf("  - Number of links: %d\n", len(portalPage.Links))

	// Example 2: Find by slug
	fmt.Println("\n2. Finding portal page by slug...")
	
	found, err := repo.FindBySlug(ctx, "john-doe")
	if err != nil {
		log.Fatalf("Failed to find portal page by slug: %v", err)
	}

	fmt.Printf("✓ Found portal page: %s\n", found.Title)
	fmt.Println("  Links (sorted by display order):")
	for i, link := range found.Links {
		fmt.Printf("    %d. %s (%s)\n", i+1, link.Title, link.URL)
	}

	// Example 3: Update the portal page
	fmt.Println("\n3. Updating portal page...")
	
	// Add a new link and update existing one
	found.Title = "John Doe - Full Stack Developer"
	found.Theme = domain.ThemeDark
	
	// Update existing link
	found.Links[0].URL = "https://github.com/johndoe-updated"
	found.Links[0].Description = "Updated GitHub Profile"
	
	// Add new link
	newLink := &domain.Link{
		Title:        "LinkedIn",
		URL:          "https://linkedin.com/in/johndoe",
		Description:  "Professional Network",
		DisplayOrder: 4,
	}
	found.Links = append(found.Links, newLink)

	err = repo.Update(ctx, found)
	if err != nil {
		log.Fatalf("Failed to update portal page: %v", err)
	}

	fmt.Printf("✓ Updated portal page\n")
	fmt.Printf("  - New title: %s\n", found.Title)
	fmt.Printf("  - New theme: %s\n", found.Theme)
	fmt.Printf("  - Number of links: %d\n", len(found.Links))

	// Example 4: Create another portal page for different user
	fmt.Println("\n4. Creating portal page for another user...")
	
	anotherPage := &domain.PortalPage{
		UserID: 2,
		Slug:   "jane-smith",
		Title:  "Jane Smith's Profile",
		Bio:    "UX Designer",
		Theme:  domain.ThemeLight,
		Links: []*domain.Link{
			{
				Title:        "Dribbble",
				URL:          "https://dribbble.com/janesmith",
				DisplayOrder: 1,
			},
		},
	}

	err = repo.Create(ctx, anotherPage)
	if err != nil {
		log.Fatalf("Failed to create second portal page: %v", err)
	}

	fmt.Printf("✓ Created portal page for Jane with ID: %d\n", anotherPage.ID)

	// Example 5: List portal pages by user ID
	fmt.Println("\n5. Listing portal pages by user ID...")
	
	userPages, err := repo.ListByUserID(ctx, 1)
	if err != nil {
		log.Fatalf("Failed to list portal pages by user ID: %v", err)
	}

	fmt.Printf("✓ Found %d portal page(s) for user ID 1:\n", len(userPages))
	for _, page := range userPages {
		fmt.Printf("  - %s (slug: %s)\n", page.Title, page.Slug)
		// Note: Links are not included in list view
		fmt.Printf("    Links: %v (not loaded in list view)\n", page.Links)
	}

	// Example 6: Find by ID
	fmt.Println("\n6. Finding portal page by ID...")
	
	byID, err := repo.FindByID(ctx, portalPage.ID)
	if err != nil {
		log.Fatalf("Failed to find portal page by ID: %v", err)
	}

	fmt.Printf("✓ Found portal page by ID %d: %s\n", byID.ID, byID.Title)
	fmt.Printf("  - Created at: %s\n", byID.CreatedAt.Format(time.RFC3339))
	fmt.Printf("  - Updated at: %s\n", byID.UpdatedAt.Format(time.RFC3339))

	// Example 7: Demonstrate error handling
	fmt.Println("\n7. Demonstrating error handling...")
	
	// Try to find non-existent slug
	_, err = repo.FindBySlug(ctx, "non-existent")
	if err != nil {
		fmt.Printf("✓ Expected error for non-existent slug: %v\n", err)
	}

	// Try to create duplicate slug
	duplicatePage := &domain.PortalPage{
		UserID: 3,
		Slug:   "john-doe", // Same slug as first page
		Title:  "Duplicate Slug",
		Links:  []*domain.Link{},
	}

	err = repo.Create(ctx, duplicatePage)
	if err != nil {
		fmt.Printf("✓ Expected error for duplicate slug: %v\n", err)
	}

	fmt.Println("\n=== Example completed successfully! ===")
	fmt.Println("\nThe InMemoryPortalPageRepository is perfect for:")
	fmt.Println("• Unit testing without database dependencies")
	fmt.Println("• Local development and prototyping")
	fmt.Println("• Integration tests in CI/CD pipelines")
	fmt.Println("• Quick experimentation and demos")
}