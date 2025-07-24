// Package fieldmask enables selective field retention on Go structs using
// dot-separated field paths. It provides a powerful and efficient way to
// manipulate struct fields while maintaining type safety and handling complex
// nested structures.
//
// Key Features:
//   - Selective field updates using dot notation paths
//   - Support for nested structs, pointers, and complex types
//   - JSON tag compatibility
//   - Built-in protection against circular references
//   - High performance through internal caching
//   - Zero external dependencies
//
// The package is particularly useful in scenarios such as:
//   - API response filtering
//   - Partial updates in REST endpoints
//   - Data transformation pipelines
//   - Selective struct serialization
//
// Basic Usage:
//
//	type Profile struct {
//	    Age  int    `json:"age"`
//	    Bio  string `json:"bio"`
//	}
//
//	type User struct {
//	    Name    string  `json:"name"`
//	    Email   string  `json:"email"`
//	    Profile Profile `json:"profile"`
//	}
//
//	user := &User{
//	    Name:  "John Doe",
//	    Email: "john@example.com",
//	    Profile: Profile{
//	        Age: 30,
//	        Bio: "Developer",
//	    },
//	}
//
//	// Create and apply a field mask
//	mask := fieldmask.New("name", "profile.age")
//	if err := mask.Apply(user); err != nil {
//	    log.Fatal(err)
//	}
//
// Path Management:
//
//	mask := fieldmask.New("user.profile.name", "user.email")
//	exists := mask.HasPath("user.email")           // Check specific path
//	exists = mask.HashAny("user.name", "user.id")  // Check multiple paths
//	mask.RemovePaths("user.email")                 // Remove paths
//	paths := mask.GetPaths()                       // Get all paths
//
// Performance Optimizations:
//
// The package implements internal caching for type descriptors and zero values,
// ensuring optimal performance for repeated operations on the same types.
// The caching is thread-safe and handles concurrent access properly.
//
// Error Handling:
//
// The Apply method returns errors in the following cases:
//   - Nil input
//   - Non-pointer input
//   - Non-struct input
//   - Circular references in nested structs
//
// Thread Safety:
//
// All operations are thread-safe and can be used concurrently.
// The internal caches use sync.Map for safe concurrent access.
package fieldmask
