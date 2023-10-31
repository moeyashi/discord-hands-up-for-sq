package repository

import (
	"context"
	"net/http"
	"os"
	"testing"
)

const emulatorHost = "localhost:5000"
const projectID = "test"

func TestMain(m *testing.M) {
	os.Setenv("FIRESTORE_EMULATOR_HOST", emulatorHost)
	os.Setenv("FIREBASE_PROJECT_ID", projectID)

	code := m.Run()
	os.Exit(code)
}

func resetEmulator(t *testing.T) {
	url := "http://" + emulatorHost + "/emulator/v1/projects/" + projectID + "/databases/(default)/documents"
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetVersion(t *testing.T) {
	resetEmulator(t)
	ctx := context.Background()
	repo, err := New(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.client.Collection("v").Doc("1").Create(ctx, map[string]interface{}{"version": 1})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.client.Collection("v").Doc("2").Create(ctx, map[string]interface{}{"version": 2})
	if err != nil {
		t.Fatal(err)
	}

	version, err := repo.GetVersion(ctx)
	if err != nil {
		t.Fatal(err)
	}
	
	if version.Version != 2 {
		t.Errorf("version = %d; want 2", version.Version)
	}
}
