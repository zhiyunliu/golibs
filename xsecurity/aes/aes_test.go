package aes

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	type args struct {
		plainText string
		key       string
		mode      string
		opt       []Option
	}
	tests := []struct {
		name           string
		args           args
		wantCipherText string
		wantErr        bool
	}{
		{name: "1.", args: args{plainText: "1234567890123456", key: "glue.xdb12345678", mode: "cbc/pkcs7"}, wantCipherText: "", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCipherText, err := Encrypt(tt.args.plainText, tt.args.key, tt.args.mode, tt.args.opt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCipherText != tt.wantCipherText {
				t.Errorf("Encrypt() = %v, want %v", gotCipherText, tt.wantCipherText)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		cipherText string
		key        string
		mode       string
		opt        []Option
	}
	tests := []struct {
		name          string
		args          args
		wantPlainText string
		wantErr       bool
	}{
		{name: "1.", args: args{cipherText: "xTumhr4RgvDUpUlpWq7WL9C7E9aeoF00PGgnqH6pxp0=", key: "glue.xdb12345678", mode: "cbc/pkcs7"}, wantPlainText: "1234567890123456x", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPlainText, err := Decrypt(tt.args.cipherText, tt.args.key, tt.args.mode, tt.args.opt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPlainText != tt.wantPlainText {
				t.Errorf("Decrypt() = %v, want %v", gotPlainText, tt.wantPlainText)
			}
		})
	}
}
