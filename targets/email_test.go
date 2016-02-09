package targets_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/yext/revere/targets"
)

var (
	emailTargetType = Email{}
	emailId         = 0
	emailName       = "Email"
	validEmailJson  = `{
		"emails": [
			{"emailTo":"test@ex.com", "replyTo":"test2@ex.com"}
		]
	}`

	invalidEmailAddresses = []string{"", "test", "test.com", "test@ex", "@ex", "@ex.com"}
	validEmailAddresses   = []string{"test@ex.com", "a@a.a", "a+a@a.a"}
)

func validEmailTarget() (*EmailTarget, error) {
	target, err := emailTargetType.Load(validEmailJson)
	if err != nil {
		return nil, err
	}

	emailTarget, ok := target.(EmailTarget)
	if !ok {
		return nil, fmt.Errorf("Invalid target loaded for target type: %s\n", emailTargetType.Name())
	}

	if len(emailTarget.EmailAddresses) < 1 {
		return nil, errors.New("Email json contains no email addresses")
	}

	return &emailTarget, nil
}

func TestEmailId(t *testing.T) {
	if int(emailTargetType.Id()) != emailId {
		t.Errorf("Expected email target type id: %d, got %d\n", emailId, emailTargetType.Id)
	}
}

func TestEmailName(t *testing.T) {
	if emailTargetType.Name() != emailName {
		t.Errorf("Expected email target type name: %s, got %s\n", emailName, emailTargetType.Name)
	}
}

func TestLoadEmptyEmail(t *testing.T) {
	target, err := emailTargetType.Load(`{}`)
	if err != nil {
		t.Fatalf("Failed to load empty email target: %s\n", err.Error())
	}

	_, ok := target.(EmailTarget)
	if !ok {
		t.Fatalf("Invalid target loaded for target type: %s\n", emailTargetType.Name())
	}
}

func TestInvalidEmailTo(t *testing.T) {
	et, err := validEmailTarget()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, e := range invalidEmailAddresses {
		et.EmailAddresses[0].EmailTo = e
		errs := et.Validate()
		if errs == nil {
			t.Errorf("Expected error for email-to: %s\n", e)
		}
	}
}

func TestValidEmailTo(t *testing.T) {
	et, err := validEmailTarget()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, e := range validEmailAddresses {
		et.EmailAddresses[0].EmailTo = e
		errs := et.Validate()
		if errs != nil {
			t.Errorf("Unexpected error for email-to: %v\n", errs)
		}
	}
}

func TestInvalidReplyTo(t *testing.T) {
	et, err := validEmailTarget()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, e := range invalidEmailAddresses {
		et.EmailAddresses[0].ReplyTo = e
		errs := et.Validate()
		if errs == nil {
			t.Errorf("Expected error for reply-to: %s\n", e)
		}
	}
}

func TestValidReplyTo(t *testing.T) {
	et, err := validEmailTarget()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, e := range validEmailAddresses {
		et.EmailAddresses[0].ReplyTo = e
		errs := et.Validate()
		if errs != nil {
			t.Errorf("Unexpected error for reply-to: %v\n", errs)
		}
	}
}
