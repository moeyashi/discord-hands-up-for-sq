package repository

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"testing"
)

// 参考: https://www.captaincodeman.com/unit-testing-with-firestore-emulator-and-go

const FirestoreEmulatorHost = "FIRESTORE_EMULATOR_HOST"

func TestMain(m *testing.M) {
	os.Setenv("FIREBASE_PROJECT_ID", "test")
	os.Setenv(FirestoreEmulatorHost, "127.0.0.1:5000")
	// command to start firestore emulator
	cmd := exec.Command("firebase", "emulators:start", "--only", "firestore")

	// this makes it killable
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// we need to capture it's output to know when it's started
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stderr.Close()

	// start her up!
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// ensure the process is killed when we're finished, even if an error occurs
        // (thanks to Brian Moran for suggestion)
	var result int
	defer func() {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		os.Exit(result)
	}()

	// we're going to wait until it's running to start
	var wg sync.WaitGroup
	wg.Add(1)

	// by starting a separate go routine
	go func() {
		// reading it's output
		buf := make([]byte, 256, 256)
		for {
			n, err := stderr.Read(buf[:])
			if err != nil {
				// until it ends
				if err == io.EOF {
					break
				}
				log.Fatalf("reading stderr %v", err)
			}

			if n > 0 {
				d := string(buf[:n])

				// only required if we want to see the emulator output
				log.Printf("%s", d)

				// checking for the message that it's started
				if strings.Contains(d, "Dev App Server is now running") {
					wg.Done()
				}

			}
		}
	}()

	// wait until the running message has been received
	wg.Wait()

	// now it's running, we can run our unit tests
	result = m.Run()
}

func TestGetVersion(t *testing.T) {
	ctx := context.Background()
	repo, err := New(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.client.Collection("v").Doc("1").Create(ctx, map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	version, err := repo.GetVersion(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(version)
}