package app

import (
	"log"
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		channelID   string
		orgID       string
		orgUser     string
		chaincodeID string
		configPath  string
		cryptoPath  string
	}
	tests := []struct {
		name    string
		args    args
		wantP   *Provider
		wantErr bool
	}{{name: "correct",
		args: args{
			channelID:   "orgschannel",
			orgID:       "CitiBank",
			orgUser:     "Admin",
			chaincodeID: "cc_gopenbanking",
			configPath:  "config.yaml",
			cryptoPath:  "../crypto-config"},
		wantErr: false}, {
		name: "wrong config",
		args: args{
			channelID:   "orgschannel",
			orgID:       "CitiBank",
			orgUser:     "Admin",
			chaincodeID: "cc_gopenbanking",
			configPath:  "fault/config.yaml",
			cryptoPath:  "../crypto-config"},
		wantP:   nil,
		wantErr: true}, {
		name: "wrong crypto-config",
		args: args{
			channelID:   "orgschannel",
			orgID:       "CitiBank",
			orgUser:     "Admin",
			chaincodeID: "cc_gopenbanking",
			configPath:  "config.yaml",
			cryptoPath:  "fault/crypto-config"},
		wantP:   nil,
		wantErr: true}, {
		name: "wrong organization",
		args: args{
			channelID:   "orgschannel",
			orgID:       "SomeBank",
			orgUser:     "Admin",
			chaincodeID: "cc_gopenbanking",
			configPath:  "config.yaml",
			cryptoPath:  "../crypto-config"},
		wantP:   nil,
		wantErr: true}, {
		name: "wrong user",
		args: args{
			channelID:   "orgschannel",
			orgID:       "CitiBank",
			orgUser:     "Someone",
			chaincodeID: "cc_gopenbanking",
			configPath:  "config.yaml",
			cryptoPath:  "../crypto-config"},
		wantP:   nil,
		wantErr: true}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP, err := New(tt.args.channelID, tt.args.orgID, tt.args.orgUser, tt.args.chaincodeID, tt.args.configPath, tt.args.cryptoPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantP != nil && !reflect.DeepEqual(gotP, tt.wantP) {
				t.Errorf("New() = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}

func TestProvider_Invoke_Once(t *testing.T) {
	type fields struct {
		channelID   string
		orgID       string
		orgUser     string
		chaincodeID string
		configPath  string
		cryptoPath  string
	}
	sharedFields := fields{
		channelID:   "orgschannel",
		orgID:       "ANZBank",
		orgUser:     "User1",
		chaincodeID: "cc_gopenbanking",
		configPath:  "config.yaml",
		cryptoPath:  "../crypto-config"}

	type args struct {
		ccFunction string
		args       []string
	}
	tests := []struct {
		name     string
		args     args
		wantResp string
		wantErr  bool
	}{{name: "test get",
		args:    args{ccFunction: "get", args: []string{"a"}},
		wantErr: false}, {
		name:    "test query",
		args:    args{ccFunction: "query", args: []string{"in", "a"}},
		wantErr: false}, {
		name:    "test add",
		args:    args{ccFunction: "add", args: []string{"a", "1"}},
		wantErr: false}, {
		name:    "test reduce",
		args:    args{ccFunction: "reduce", args: []string{"a", "1"}},
		wantErr: false}, {
		name:    "test create",
		args:    args{ccFunction: "create", args: []string{"test", "0"}},
		wantErr: false}, {
		name:    "test delete",
		args:    args{ccFunction: "delete", args: []string{"test"}},
		wantErr: false}, {
		name:    "test tranfer",
		args:    args{ccFunction: "transfer", args: []string{"a", "b", "1"}},
		wantErr: false}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap, err := New(sharedFields.channelID,
				sharedFields.orgID,
				sharedFields.orgUser,
				sharedFields.chaincodeID,
				sharedFields.configPath,
				sharedFields.cryptoPath)
			if err != nil {
				t.Errorf("Prepare Provider error: %v", err)
				return
			}

			gotResp, err := ap.Invoke(tt.args.ccFunction, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.Invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResp != "" && gotResp != tt.wantResp {
				t.Errorf("Provider.Invoke() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestProvider_Invoke_Validate_One_Account(t *testing.T) {
	type fields struct {
		channelID   string
		orgID       string
		orgUser     string
		chaincodeID string
		configPath  string
		cryptoPath  string
	}
	sharedFields := fields{
		channelID:   "orgschannel",
		orgID:       "ANZBank",
		orgUser:     "User1",
		chaincodeID: "cc_gopenbanking",
		configPath:  "config.yaml",
		cryptoPath:  "../crypto-config"}

	type args struct {
		ccFunction string
		args       []string
	}
	tests := []struct {
		name     string
		args     args
		wantResp string
		wantErr  bool
	}{{name: "velidate add",
		args:    args{ccFunction: "add", args: []string{"a", "1"}},
		wantErr: false}, {
		name:    "velidate reduce",
		args:    args{ccFunction: "reduce", args: []string{"a", "1"}},
		wantErr: false}}

	re := regexp.MustCompile(`(\d+)$`) // the num pattern in response

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap, err := New(sharedFields.channelID,
				sharedFields.orgID,
				sharedFields.orgUser,
				sharedFields.chaincodeID,
				sharedFields.configPath,
				sharedFields.cryptoPath)
			if err != nil {
				t.Errorf("Prepare Provider error: %v", err)
				return
			}

			// get original balance
			resp, err := ap.Invoke("get", []string{tt.args.args[0]})
			if err != nil {
				t.Errorf("Get original account balance error: %v", err)
				return
			}
			numRaw := re.Find([]byte(resp))
			before, err := strconv.ParseFloat(string(numRaw), 64)
			if err != nil {
				t.Errorf("Cannot parse numbers from the response: %s, error = %v", resp, err)
				return
			}

			// execute
			gotResp, err := ap.Invoke(tt.args.ccFunction, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.Invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResp != "" && gotResp != tt.wantResp {
				t.Errorf("Provider.Invoke() = %v, want %v", gotResp, tt.wantResp)
			}

			// validate
			resp, err = ap.Invoke("get", []string{tt.args.args[0]})
			if err != nil {
				t.Errorf("Get current account balance error: %v", err)
				return
			}
			numRaw = re.Find([]byte(resp))
			after, err := strconv.ParseFloat(string(numRaw), 64)
			if err != nil {
				t.Errorf("Cannot parse numbers from the response: %s, error = %v", resp, err)
				return
			}

			diff, _ := strconv.ParseFloat(tt.args.args[1], 64)
			if tt.args.ccFunction == "reduce" {
				diff = -diff
			}
			if before+diff != after {
				t.Errorf("Found inconsitency: before: %f, delta: %f, after: %f", before, diff, after)
				return
			}

			log.Printf("Account: %s, before: %f, delta: %f, after: %f\n", tt.args.args[0], before, diff, after)
		})
	}
}

func TestProvider_Invoke_Validate_Transfer(t *testing.T) {
	type fields struct {
		channelID   string
		orgID       string
		orgUser     string
		chaincodeID string
		configPath  string
		cryptoPath  string
	}
	sharedFields := fields{
		channelID:   "orgschannel",
		orgID:       "ANZBank",
		orgUser:     "User1",
		chaincodeID: "cc_gopenbanking",
		configPath:  "config.yaml",
		cryptoPath:  "../crypto-config"}

	type args struct {
		ccFunction string
		args       []string
	}
	tests := []struct {
		name     string
		args     args
		wantResp string
		wantErr  bool
	}{{name: "velidate tranfer",
		args:    args{ccFunction: "transfer", args: []string{"a", "b", "1"}},
		wantErr: false}}

	re := regexp.MustCompile(`(\d+)$`) // the num pattern in response

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap, err := New(sharedFields.channelID,
				sharedFields.orgID,
				sharedFields.orgUser,
				sharedFields.chaincodeID,
				sharedFields.configPath,
				sharedFields.cryptoPath)
			if err != nil {
				t.Errorf("Prepare Provider error: %v", err)
				return
			}

			// get original balance
			// Debit
			resp, err := ap.Invoke("get", []string{tt.args.args[0]})
			if err != nil {
				t.Errorf("Get original account balance error: %v", err)
				return
			}
			numRaw := re.Find([]byte(resp))
			debitBefore, err := strconv.ParseFloat(string(numRaw), 64)
			if err != nil {
				t.Errorf("Cannot parse numbers from the response: %s, error = %v", resp, err)
				return
			}
			// Credit
			resp, err = ap.Invoke("get", []string{tt.args.args[1]})
			if err != nil {
				t.Errorf("Get original account balance error: %v", err)
				return
			}
			numRaw = re.Find([]byte(resp))
		  creditBefore, err := strconv.ParseFloat(string(numRaw), 64)
			if err != nil {
				t.Errorf("Cannot parse numbers from the response: %s, error = %v", resp, err)
				return
			}

			// execute
			gotResp, err := ap.Invoke(tt.args.ccFunction, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.Invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResp != "" && gotResp != tt.wantResp {
				t.Errorf("Provider.Invoke() = %v, want %v", gotResp, tt.wantResp)
			}

			// validate
			// Debit
			resp, err = ap.Invoke("get", []string{tt.args.args[0]})
			if err != nil {
				t.Errorf("Get current account balance error: %v", err)
				return
			}
			numRaw = re.Find([]byte(resp))
			debitAfter, err := strconv.ParseFloat(string(numRaw), 64)
			if err != nil {
				t.Errorf("Cannot parse numbers from the response: %s, error = %v", resp, err)
				return
			}
			resp, err = ap.Invoke("get", []string{tt.args.args[1]})
			if err != nil {
				t.Errorf("Get current account balance error: %v", err)
				return
			}
			numRaw = re.Find([]byte(resp))
			creditAfter, err := strconv.ParseFloat(string(numRaw), 64)
			if err != nil {
				t.Errorf("Cannot parse numbers from the response: %s, error = %v", resp, err)
				return
			}

			diff, _ := strconv.ParseFloat(tt.args.args[2], 64)
			if debitBefore - diff != debitAfter {
				t.Errorf("Found inconsitency at Debit site: before: %f, delta: %f, after: %f", debitBefore, -diff, debitAfter)
				return
			}
			if creditBefore + diff != creditAfter {
				t.Errorf("Found inconsitency at Credit site: before: %f, delta: %f, after: %f", creditBefore, diff, creditAfter)
				return
			}

			log.Printf("Debit account: %s, before: %f, delta: %f, after: %f\n", tt.args.args[0], debitBefore, -diff, debitAfter)
			log.Printf("Credit account: %s, before: %f, delta: %f, after: %f\n", tt.args.args[0], creditBefore, diff, creditAfter)
		})
	}
}
