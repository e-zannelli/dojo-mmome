package domain

import "testing"

func TestGame_Reset(t *testing.T) {
	type fields struct {
		isRunning bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"A running game should return error",
			fields{
				isRunning: true,
			},
			true,
		},
		{
			"A not running game should not return error",
			fields{
				isRunning: false,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				isRunning: tt.fields.isRunning,
			}
			if err := g.Reset(); (err != nil) != tt.wantErr {
				t.Errorf("Reset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_Guess(t *testing.T) {
	type fields struct {
		isRunning     bool
		playableRunes []rune
		code          []rune
	}
	type args struct {
		guess []rune
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		correct   int
		misplaced int
		wantErr   bool
	}{
		{
			"success",
			fields{
				isRunning:     true,
				playableRunes: []rune("123456"),
				code:          []rune("12345"),
			},
			args{
				[]rune("12345"),
			},
			5,
			0,
			false,
		},
		{
			"misplaced are counted once",
			fields{
				isRunning:     true,
				playableRunes: []rune("123456"),
				code:          []rune("11211"),
			},
			args{
				[]rune("22322"),
			},
			0,
			1,
			false,
		},
		{
			"Correctly placed are not counted as misplaced",
			fields{
				isRunning:     true,
				playableRunes: []rune("123456"),
				code:          []rune("12211"),
			},
			args{
				[]rune("22322"),
			},
			1,
			1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				isRunning:     tt.fields.isRunning,
				playableRunes: tt.fields.playableRunes,
				code:          tt.fields.code,
			}
			gotCorrect, gotMisplaced, err := g.Guess(tt.args.guess)
			if (err != nil) != tt.wantErr {
				t.Errorf("Guess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCorrect != tt.correct {
				t.Errorf("Guess() gotCorrect = %v, correct %v", gotCorrect, tt.correct)
			}
			if gotMisplaced != tt.misplaced {
				t.Errorf("Guess() gotMisplaced = %v, misplaced %v", gotMisplaced, tt.misplaced)
			}
		})
	}
}
