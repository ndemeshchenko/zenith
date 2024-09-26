package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func MergeBSONMaps(map1, map2 bson.M) bson.M {
	result := make(bson.M)

	// Copy keys and values from the first map
	for key, value := range map1 {
		result[key] = value
	}

	// Iterate over the keys in the second map
	for key, value := range map2 {
		// Check if the key already exists in the result map
		if existingValue, ok := result[key]; ok {
			// Handle the conflict (for example, overwrite the existing value)
			// You can customize this behavior based on your needs
			fmt.Printf("Key %s already exists with value %v. Overwriting with %v.\n", key, existingValue, value)
		}

		// Add the key and value to the result map
		result[key] = value
	}

	return result
}
