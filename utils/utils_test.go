package utils

import "testing"

func TestFilter(t *testing.T) {

	t.Run("should return even numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		evenNumbers := Filter(numbers, func(n int) bool {
			return n%2 == 0
		})
		if len(evenNumbers) != 5 {
			t.Errorf("got %d even numbers, want 5", len(evenNumbers))
		}
		// Check if the even numbers are correct
		for _, n := range evenNumbers {
			if n%2 != 0 {
				t.Errorf("got %d, want an even number", n)
			}
		}
	})
}

func TestAnyMatch(t *testing.T) {

	t.Run("should return true if any number is greater than 5", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		greaterThanFive := AnyMatch(numbers, func(n int) bool {
			return n > 5
		})
		if !greaterThanFive {
			t.Error("got false, want true")
		}
	})

	t.Run("should return false if no number is greater than 10", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		greaterThanTen := AnyMatch(numbers, func(n int) bool {
			return n > 10
		})
		if greaterThanTen {
			t.Error("got true, want false")
		}
	})
}
