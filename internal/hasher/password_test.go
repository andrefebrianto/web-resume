package hasher

import "testing"

func TestGenerateFromPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "hash password longer than 72 bytes",
			args: args{
				password: "(F2r+CYG~'AZb'sPwcV#;PHkaS(tSZC4ukEV^rd7KNGN-M2l8[]&.Z8Imm%C404Pa9_Ld8Twhrex1K;HY{85}J(HS.,6+oA1a+m7",
			},
			wantErr: true,
		},
		{
			name: "hash password fewer than 72 bytes",
			args: args{
				password: "veryverysecurepassword",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateFromPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateFromPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	type args struct {
		hashedPassword string
		password       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "verify correct password",
			args: args{
				hashedPassword: "$2a$12$hWASkUwEkcS1CbsyRRwoBew5r7qwmXwH4YJyP.S149hghOg77UEQW",
				password:       "veryverysecurepassword",
			},
			wantErr: false,
		},
		{
			name: "verify incorrect password",
			args: args{
				hashedPassword: "$2a$12$hWASkUwEkcS1CbsyRRwoBew5r7qwmXwH4YJyP.S149hghOg77UEQW",
				password:       "verysafepassword",
			},
			wantErr: true,
		},
		{
			name: "verify incorrect hashed password",
			args: args{
				hashedPassword: "$2a$12$",
				password:       "veryverysecurepassword",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyPassword(tt.args.hashedPassword, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
