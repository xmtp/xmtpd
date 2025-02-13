package upgrade_test

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func setup() {
	start := time.Now() // Start tracking time
	fmt.Println("=== SETUP UpgradeTestSetup")
	fmt.Println("    Setting up before all tests...")

	// Measure time for building dev image
	fmt.Println("    ⧖ Building dev image... This may take a while.")
	imageStart := time.Now()
	err := buildDevImage()
	if err != nil {
		fmt.Printf("    Error building dev image: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("    ✔ Dev image built in %v\n", time.Since(imageStart))

	// Measure time for pulling old images
	fmt.Println("    ⧖ Pulling old images...")
	pullStart := time.Now()
	for _, image := range upgradeToLatest {
		err := dockerPull(image)
		if err != nil {
			fmt.Printf("    ❌ Error pulling image %s: %v\n", image, err)
			os.Exit(1)
		}
	}
	fmt.Printf("    ✔ All images pulled in %v\n", time.Since(pullStart))

	// Print total setup time
	fmt.Printf("=== SETUP COMPLETE (%v)\n", time.Since(start))
}

// TestMain runs once for the whole test suite
func TestMain(m *testing.M) {
	skipIfNotEnabled()
	setup()
	code := m.Run()
	os.Exit(code)
}
