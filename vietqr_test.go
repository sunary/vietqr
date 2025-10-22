package vietqr

import (
	"testing"
)

// Test helper functions
func createTestTransferInfo(bankCode, bankNo string, amount int64, message string) TransferInfo {
	return TransferInfo{
		BankCode: bankCode,
		BankNo:   bankNo,
		Amount:   amount,
		Message:  message,
	}
}

func createTestVNPAYTransferInfo(merchantID string) TransferInfo {
	return TransferInfo{
		merchantID: merchantID,
	}
}

func equalTransferInfo(t1, t2 *TransferInfo) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}
	return t1.merchantID == t2.merchantID &&
		t1.BankCode == t2.BankCode &&
		t1.BankNo == t2.BankNo &&
		t1.Amount == t2.Amount &&
		t1.Message == t2.Message
}

// Test data constants
var (
	// Valid QR codes for testing
	validQRCodeBasic       = "00020101021138540010A00000072701240006546034011009055559990208QRIBFTTA53037045802VN6304B2A0"
	validQRCodeWithAmount  = "00020101021238490010A000000727011900069704160105135790208QRIBFTTA530370454061200005802VN63049A71"
	validQRCodeWithMessage = "00020101021138510010A00000072701210006970407010797968680208QRIBFTTA53037045802VN62240820gen by sunary/vietqr6304BE74"
	validQRCodeFull        = "00020101021238490010A000000727011900069704320105193720208QRIBFTTA530370454061520005802VN62200816gen by go-vietqr63040ED4"
	validVNPAYQRCode       = "00020101021138400010A0000007750110VNP-1234560208QRIBFTTA53037045802VN6304E56C"

	// Invalid QR codes for testing
	invalidQRCodeWrongCRC   = "00020101021138540010A00000072701240006546034011009055559990208QRIBFTTA53037045802VN6304xxxx"
	invalidQRCodeMissingCRC = "00020101021138540010A00000072701240006546034011009055559990208QRIBFTTA53037045802VN6304"
	invalidQRCodeMalformed  = "00020101021138540010A00000072701240006546034011009055559990208QRIBFTTA53037045802VN"
	invalidQRCodeTooShort   = "00020101021138540010A00000072701240006546034011009055559990208QRIBFTTA53037045802VN6304"
)

// Test Encode function
func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    TransferInfo
		expected string
	}{
		// Basic bank transfer tests
		{
			name:     "Basic bank transfer - CAKE bank",
			input:    createTestTransferInfo(CAKE, "0905555999", 0, ""),
			expected: validQRCodeBasic,
		},
		{
			name:     "Bank transfer with amount - ACB bank",
			input:    createTestTransferInfo(ACB, "13579", 120000, ""),
			expected: validQRCodeWithAmount,
		},
		{
			name:     "Bank transfer with message - TECHCOMBANK",
			input:    createTestTransferInfo(TECHCOMBANK, "9796868", 0, "gen by sunary/vietqr"),
			expected: validQRCodeWithMessage,
		},
		{
			name:     "Full bank transfer - VPBANK",
			input:    createTestTransferInfo(VPBANK, "19372", 152000, "gen by go-vietqr"),
			expected: validQRCodeFull,
		},

		// VNPAY merchant tests
		{
			name:     "VNPAY merchant transfer",
			input:    createTestVNPAYTransferInfo("VNP-123456"),
			expected: validVNPAYQRCode,
		},

		// Edge cases
		{
			name:     "Zero amount",
			input:    createTestTransferInfo(BIDV, "1234567890", 0, ""),
			expected: "00020101021138540010A00000072701240006970418011012345678900208QRIBFTTA53037045802VN6304995B",
		},
		{
			name:     "Large amount",
			input:    createTestTransferInfo(VIETINBANK, "9876543210", 999999999, "Large amount transfer"),
			expected: "00020101021238540010A00000072701240006970415011098765432100208QRIBFTTA530370454099999999995802VN62250821Large amount transfer63048A8F",
		},
		{
			name:     "Long message",
			input:    createTestTransferInfo(AGRIBANK, "1111111111", 50000, "This is a very long message that tests the maximum length handling for QR code generation"),
			expected: "00020101021238540010A00000072701240006970405011011111111110208QRIBFTTA53037045405500005802VN62930889This is a very long message that tests the maximum length handling for QR code generation63043199",
		},
		{
			name:     "Special characters in message",
			input:    createTestTransferInfo(SACOMBANK, "2222222222", 75000, "Message with special chars: @#$%^&*()_+-=[]{}|;':\",./<>?"),
			expected: "00020101021238540010A00000072701240006970403011022222222220208QRIBFTTA53037045405750005802VN62600856Message with special chars: @#$%^&*()_+-=[]{}|;':\",./<>?630498F2",
		},
		{
			name:     "Empty message",
			input:    createTestTransferInfo(MBBANK, "3333333333", 100000, ""),
			expected: "00020101021238540010A00000072701240006970422011033333333330208QRIBFTTA530370454061000005802VN63047CF1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.input)
			if result != tt.expected {
				t.Errorf("Encode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test Decode function
func TestDecode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *TransferInfo
		expectError bool
	}{
		// Valid QR codes
		{
			name:  "Valid basic bank transfer - CAKE",
			input: validQRCodeBasic,
			expected: &TransferInfo{
				BankCode: CAKE,
				BankNo:   "0905555999",
				Amount:   0,
				Message:  "",
			},
			expectError: false,
		},
		{
			name:  "Valid bank transfer with amount - ACB",
			input: validQRCodeWithAmount,
			expected: &TransferInfo{
				BankCode: ACB,
				BankNo:   "13579",
				Amount:   120000,
				Message:  "",
			},
			expectError: false,
		},
		{
			name:  "Valid bank transfer with message - TECHCOMBANK",
			input: validQRCodeWithMessage,
			expected: &TransferInfo{
				BankCode: TECHCOMBANK,
				BankNo:   "9796868",
				Amount:   0,
				Message:  "gen by sunary/vietqr",
			},
			expectError: false,
		},
		{
			name:  "Valid full bank transfer - VPBANK",
			input: validQRCodeFull,
			expected: &TransferInfo{
				BankCode: VPBANK,
				BankNo:   "19372",
				Amount:   152000,
				Message:  "gen by go-vietqr",
			},
			expectError: false,
		},
		{
			name:  "Valid VNPAY merchant transfer",
			input: validVNPAYQRCode,
			expected: &TransferInfo{
				merchantID: "VNP-123456",
				BankCode:   "",
				BankNo:     "",
				Amount:     0,
				Message:    "",
			},
			expectError: false,
		},

		// Invalid QR codes
		{
			name:        "Invalid QR code - wrong CRC",
			input:       invalidQRCodeWrongCRC,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid QR code - missing CRC",
			input:       invalidQRCodeMissingCRC,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid QR code - malformed",
			input:       invalidQRCodeMalformed,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid QR code - too short",
			input:       invalidQRCodeTooShort,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid QR code - invalid bank bin",
			input:       "00020101021138540010A00000072701240006999999011009055559990208QRIBFTTA53037045802VN6304B2A0",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid QR code - missing bank number",
			input:       "00020101021138540010A0000007270124000654603401100208QRIBFTTA53037045802VN6304B2A0",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Decode(tt.input)

			if (err != nil) != tt.expectError {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.expectError)
				return
			}

			if !equalTransferInfo(result, tt.expected) {
				t.Errorf("Decode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test roundtrip encoding and decoding
func TestEncodeDecodeRoundtrip(t *testing.T) {
	testCases := []struct {
		name string
		ti   TransferInfo
	}{
		{
			name: "Basic bank transfer",
			ti:   createTestTransferInfo(BIDV, "1234567890", 0, ""),
		},
		{
			name: "Bank transfer with amount",
			ti:   createTestTransferInfo(VIETINBANK, "9876543210", 500000, ""),
		},
		{
			name: "Bank transfer with message",
			ti:   createTestTransferInfo(ACB, "5555555555", 0, "Test message"),
		},
		{
			name: "Full bank transfer",
			ti:   createTestTransferInfo(TECHCOMBANK, "1111111111", 1000000, "Full test message"),
		},
		{
			name: "VNPAY merchant",
			ti:   createTestVNPAYTransferInfo("VNP-TEST-123"),
		},
		{
			name: "Edge case - zero amount with message",
			ti:   createTestTransferInfo(MBBANK, "9999999999", 0, "Zero amount test"),
		},
		{
			name: "Edge case - maximum amount",
			ti:   createTestTransferInfo(VPBANK, "8888888888", 999999999, "Maximum amount test"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encode the transfer info
			encoded := Encode(tc.ti)

			// Decode it back
			decoded, err := Decode(encoded)
			if err != nil {
				t.Errorf("Decode failed for %s: %v", tc.name, err)
				return
			}

			// Compare the results
			if !equalTransferInfo(decoded, &tc.ti) {
				t.Errorf("Roundtrip failed for %s: original = %+v, decoded = %+v", tc.name, tc.ti, decoded)
			}
		})
	}
}

// Test specific bank codes
func TestEncodeDecodeSpecificBanks(t *testing.T) {
	bankTests := []struct {
		bankCode string
		bankNo   string
	}{
		{ACB, "1234567890"},
		{BIDV, "9876543210"},
		{VIETINBANK, "1111111111"},
		{TECHCOMBANK, "2222222222"},
		{VPBANK, "3333333333"},
		{MBBANK, "4444444444"},
		{AGRIBANK, "5555555555"},
		{SACOMBANK, "6666666666"},
		{CAKE, "7777777777"},
		{UBANK, "8888888888"},
	}

	for _, bt := range bankTests {
		t.Run("Bank_"+bt.bankCode, func(t *testing.T) {
			ti := createTestTransferInfo(bt.bankCode, bt.bankNo, 100000, "Test for "+bt.bankCode)

			encoded := Encode(ti)
			decoded, err := Decode(encoded)

			if err != nil {
				t.Errorf("Decode failed for bank %s: %v", bt.bankCode, err)
				return
			}

			if decoded.BankCode != bt.bankCode {
				t.Errorf("Bank code mismatch: expected %s, got %s", bt.bankCode, decoded.BankCode)
			}

			if decoded.BankNo != bt.bankNo {
				t.Errorf("Bank number mismatch: expected %s, got %s", bt.bankNo, decoded.BankNo)
			}
		})
	}
}

// Test error handling
func TestDecodeErrorHandling(t *testing.T) {
	errorTests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "Invalid format",
			input:         "not a qr code",
			expectedError: "invalid CRC",
		},
		{
			name:          "Missing CRC field",
			input:         "00020101021138540010A00000072701240006546034011009055559990208QRIBFTTA53037045802VN",
			expectedError: "invalid CRC",
		},
		{
			name:          "Invalid CRC value",
			input:         "00020101021138540010A00000072701240006546034011009055559990208QRIBFTTA53037045802VN6304ABCD",
			expectedError: "invalid CRC",
		},
		{
			name:          "Truncated QR code",
			input:         "00020101021138540010A00000072701240006546034011009055559990208QRIBFTTA53037045802VN6304",
			expectedError: "invalid CRC",
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decode(tt.input)
			if err == nil {
				t.Errorf("Expected error for %s, but got none", tt.name)
				return
			}

			if err.Error() != tt.expectedError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
			}
		})
	}
}
