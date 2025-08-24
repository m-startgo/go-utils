package main

import (
	"fmt"
	"time"

	"github.com/m-startgo/go-utils/m_time"
)

func main() {
	// Example demonstrates basic usage of the m_time package.
	Example()

	// ExampleNew demonstrates different ways to create a Time object.
	ExampleNew()

	// ExampleStartOf demonstrates how to get the start of different time periods.
	ExampleStartOf()

	// ExampleIsSame demonstrates how to compare times.
	ExampleIsSame()
}

// Example demonstrates basic usage of the m_time package.
func Example() {
	fmt.Println("=== Example ===")
	// Create a new Time object for the current time
	t1 := m_time.New()
	fmt.Println("Current time:", t1.Format("2006-01-02 15:04:05"))

	// Create a Time object from a string
	t2 := m_time.New("2023-10-27T10:00:00Z")
	fmt.Println("Parsed time:", t2.Format("2006-01-02 15:04:05"))

	// Add one hour to the time
	t3 := t2.Add(time.Hour)
	fmt.Println("After adding 1 hour:", t3.Format("2006-01-02 15:04:05"))

	// Check if t3 is after t2
	fmt.Println("Is t3 after t2?", t3.IsAfter(t2))

	// Calculate the difference in hours
	diffHours := t3.Diff(t2, "hour")
	fmt.Println("Difference in hours:", diffHours)

	// Create a Time object with a specific location
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t4 := m_time.NewWithLocation(loc, "2023-10-27T10:00:00Z")
	fmt.Println("Time in Asia/Shanghai:", t4.FormatInLocation("2006-01-02 15:04:05", loc))

	// Get the start of the day
	t5 := t2.StartOfDay()
	fmt.Println("Start of day:", t5.Format("2006-01-02 15:04:05"))

	fmt.Println()
}

// ExampleNew demonstrates different ways to create a Time object.
func ExampleNew() {
	fmt.Println("=== ExampleNew ===")
	// Current time
	t1 := m_time.New()
	fmt.Println("Current time:", t1.Format("2006-01-02 15:04:05"))

	// From a time.Time object
	t2 := m_time.New(time.Now())
	fmt.Println("From time.Time:", t2.Format("2006-01-02 15:04:05"))

	// From a Unix timestamp (seconds)
	t3 := m_time.New(int64(1698400800))
	fmt.Println("From Unix timestamp:", t3.Format("2006-01-02 15:04:05"))

	// From a string (multiple formats supported)
	t4 := m_time.New("2023-10-27 10:00:00")
	fmt.Println("From string:", t4.Format("2006-01-02 15:04:05"))

	// From a string with a different format
	t5 := m_time.New("10/27/2023 10:00:00")
	fmt.Println("From US date string:", t5.Format("2006-01-02 15:04:05"))

	fmt.Println()
}

// ExampleStartOf demonstrates how to get the start of different time periods.
func ExampleStartOf() {
	fmt.Println("=== ExampleStartOf ===")
	t := m_time.New("2023-10-27T10:30:45Z")

	fmt.Println("Original time:", t.Format("2006-01-02 15:04:05"))
	fmt.Println("Start of hour:", t.StartOfHour().Format("2006-01-02 15:04:05"))
	fmt.Println("Start of day:", t.StartOfDay().Format("2006-01-02 15:04:05"))
	fmt.Println("Start of week:", t.StartOfWeek().Format("2006-01-02 15:04:05"))
	fmt.Println("Start of month:", t.StartOfMonth().Format("2006-01-02 15:04:05"))
	fmt.Println("Start of year:", t.StartOfYear().Format("2006-01-02 15:04:05"))

	fmt.Println()
}

// ExampleIsSame demonstrates how to compare times.
func ExampleIsSame() {
	fmt.Println("=== ExampleIsSame ===")
	t1 := m_time.New("2023-10-27T10:00:00Z")
	t2 := m_time.New("2023-10-27T10:00:00Z")
	t3 := m_time.New("2023-10-27T11:00:00Z")

	fmt.Println("t1 and t2 are the same:", t1.IsSame(t2))
	fmt.Println("t1 and t3 are the same:", t1.IsSame(t3))
	fmt.Println("t1 and t3 are the same hour:", t1.IsSameHour(t3))
	fmt.Println("t1 and t3 are the same day:", t1.IsSameDay(t3))

	fmt.Println()
}
