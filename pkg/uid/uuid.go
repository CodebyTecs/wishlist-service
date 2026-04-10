package uid

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var uuidPattern = regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")

func NewUUID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:16],
	), nil
}

func IsValidUUID(value string) bool {
	return uuidPattern.MatchString(strings.ToLower(value))
}

func DeterministicSlotUUID(roomID string, start time.Time, end time.Time) string {
	source := roomID + "|" + start.UTC().Format(time.RFC3339Nano) + "|" + end.UTC().Format(time.RFC3339Nano)
	sum := sha1.Sum([]byte(source))

	b := sum[:16]
	b[6] = (b[6] & 0x0f) | 0x50
	b[8] = (b[8] & 0x3f) | 0x80

	hexValue := hex.EncodeToString(b)
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		hexValue[0:8],
		hexValue[8:12],
		hexValue[12:16],
		hexValue[16:20],
		hexValue[20:32],
	)
}
