package analyze

import "testing"

func TestGetFirstTwoPartsOfGroupId(t *testing.T) {
	firstTwoParts, err := GetFirstTwoPartsOfGroupId("com.example.application")
	if err != nil {
		t.Errorf("%v", err)
	}

	if firstTwoParts != "com.example" {
		t.Errorf("The first two parts of com.example.application is not com.example")
	}

	_, err = GetFirstTwoPartsOfGroupId("com")
	if err == nil {
		t.Errorf("com got accepted as a at-least-two part group id")
	}
}

func TestIsLocalGroupId(t *testing.T) {
	result1, err := IsLocalGroupId("com.example.backend", "com.example")
	if err != nil {
		t.Errorf("%v", err)
	}
	if result1 != true {
		t.Errorf("com.example.backend and com.example.frontend should return true for localGroupId")
	}

	result2, err := IsLocalGroupId("com.example2.backend", "com.example")
	if err != nil {
		t.Errorf("%v", err)
	}
	if result2 {
		t.Errorf("com.example2.backend is not a local com.example groupId, and should not be true")
	}

}
