package maven

import "testing"

func TestIsSecondPartyGroupIdId(t *testing.T) {
	result1, err := IsSecondPartyGroupId("com.example.backend", "com.example")
	if err != nil {
		t.Errorf("%v", err)
	}
	if result1 != true {
		t.Errorf("com.example.backend and com.example.frontend should return true for sortKey")
	}

	result2, err := IsSecondPartyGroupId("com.example2.backend", "com.example")
	if err != nil {
		t.Errorf("%v", err)
	}
	if result2 {
		t.Errorf("com.example2.backend is not a secondParty com.example groupId, and should not be true")
	}

}
