package slices

import "fmt"

func main() {
	// Inefficient
	nums1 := []int{1, 2, 3, 4, 5} // Allocates an underlying array of size 5
	fmt.Println("nums1:", nums1)
	nums1 = append(nums1, 6) // Allocates a new underlying array of size 10
	fmt.Println("nums1:", nums1)

	// Efficient
	nums2 := []int{1, 2, 3, 4, 5} // Allocates an underlying array of size 5
	fmt.Println("nums2:", nums2)
	nums2 = append(nums2[:len(nums2):cap(nums2)], 6) // Reuses the underlying array
	fmt.Println("nums2:", nums2)
}
