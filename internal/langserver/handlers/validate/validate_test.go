package validate

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestValidation_disabled(t *testing.T) {
	config, err := ioutil.ReadFile(fmt.Sprintf("../testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, diag := ValidateFile(config, "main.tf")
	if len(diag) != 0 {
		t.Errorf("expect no diagnostics, but got %v", diag)
	}
}

func TestValidation_whileTyping(t *testing.T) {
	config, err := ioutil.ReadFile(fmt.Sprintf("../testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, diag := ValidateFile(config, "main.tf")
	if len(diag) == 0 {
		t.Errorf("expect non-empty diagnostics, but got %v", diag)
	}
}

func TestValidation_commaAfterArrayItem(t *testing.T) {
	config, err := ioutil.ReadFile(fmt.Sprintf("../testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, diag := ValidateFile(config, "main.tf")
	if len(diag) != 0 {
		t.Errorf("expect 0 diagnostics, but got %v", diag)
	}
}

func TestValidation_missingRequiredProperty(t *testing.T) {
	config, err := ioutil.ReadFile(fmt.Sprintf("../testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, diag := ValidateFile(config, "main.tf")
	if len(diag) != 1 || diag[0].Summary != "`vaultBaseUrl` is required, but no definition was found" {
		t.Errorf("expect 1 diagnostics, but got %v", diag)
	}
}

func TestValidation_notExpectedProperty(t *testing.T) {
	config, err := ioutil.ReadFile(fmt.Sprintf("../testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, diag := ValidateFile(config, "main.tf")
	if len(diag) != 1 || diag[0].Summary != "`identity1` is not expected here. Do you mean `identity`? " {
		t.Errorf("expect 1 diagnostics, but got %v", diag)
	}
}

func TestValidation_update(t *testing.T) {
	config, err := ioutil.ReadFile(fmt.Sprintf("../testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	_, diag := ValidateFile(config, "main.tf")
	if len(diag) != 0 {
		t.Errorf("expect no diagnostics, but got %v", diag)
	}
}
