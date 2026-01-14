package envloader

import (
	"reflect"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	type Config struct {
		Name                         string `env:"NAME,required"`
		OptionSnakeCaseName          string `env:"OPTION_SNAKE_CASE_NAME"`
		OptionDifferentSnakeCaseName string `env:"OPTION_DIFF"`
		Skip                         string `env:"-"`
		OnlyLazy                     string
	}
	type args struct {
		env  map[string]string
		lazy bool
	}
	type testCase struct {
		name    string
		args    args
		wantErr bool
		want    Config
	}
	tests := []testCase{
		{
			name: "regular",
			args: args{
				env: map[string]string{
					"NAME":      "MyTest",
					"SKIP":      "should never show up",
					"ONLY_LAZY": "my lazy value",
				},
				lazy: false,
			},
			wantErr: false,
			want: Config{
				Name:                         "MyTest",
				OptionSnakeCaseName:          "",
				OptionDifferentSnakeCaseName: "",
				Skip:                         "",
				OnlyLazy:                     "",
			},
		},
		{
			name: "regular with different env name",
			args: args{
				env: map[string]string{
					"NAME":        "MyTest",
					"SKIP":        "should never show up",
					"OPTION_DIFF": "my different value",
					"ONLY_LAZY":   "my lazy value",
				},
				lazy: false,
			},
			wantErr: false,
			want: Config{
				Name:                         "MyTest",
				OptionSnakeCaseName:          "",
				OptionDifferentSnakeCaseName: "my different value",
				Skip:                         "",
				OnlyLazy:                     "",
			},
		},
		{
			name: "regular with lazy loading",
			args: args{
				env: map[string]string{
					"NAME":      "MyTest",
					"SKIP":      "should never show up",
					"ONLY_LAZY": "my lazy value",
				},
				lazy: true,
			},
			wantErr: false,
			want: Config{
				Name:                         "MyTest",
				OptionSnakeCaseName:          "",
				OptionDifferentSnakeCaseName: "",
				Skip:                         "",
				OnlyLazy:                     "my lazy value",
			},
		},
		{
			name: "lazy loading should not override option",
			args: args{
				env: map[string]string{
					"NAME":        "MyTest",
					"SKIP":        "should never show up",
					"OPTION_DIFF": "my different value",
					"ONLY_LAZY":   "my lazy value",
				},
				lazy: true,
			},
			wantErr: false,
			want: Config{
				Name:                         "MyTest",
				OptionSnakeCaseName:          "",
				OptionDifferentSnakeCaseName: "my different value",
				Skip:                         "",
				OnlyLazy:                     "my lazy value",
			},
		},
		{
			name: "lazy loading should not override option 2",
			args: args{
				env: map[string]string{
					"NAME":                             "MyTest",
					"SKIP":                             "should never show up",
					"OPTION_DIFFERENT_SNAKE_CASE_NAME": "my different value should not show up",
					"ONLY_LAZY":                        "my lazy value",
				},
				lazy: true,
			},
			wantErr: false,
			want: Config{
				Name:                         "MyTest",
				OptionSnakeCaseName:          "",
				OptionDifferentSnakeCaseName: "",
				Skip:                         "",
				OnlyLazy:                     "my lazy value",
			},
		},
		{
			name: "regular with optional",
			args: args{
				env: map[string]string{
					"NAME":                   "MyTest",
					"SKIP":                   "should never show up",
					"OPTION_SNAKE_CASE_NAME": "OptionalTest",
				},
				lazy: false,
			},
			wantErr: false,
			want: Config{
				Name:                         "MyTest",
				OptionSnakeCaseName:          "OptionalTest",
				OptionDifferentSnakeCaseName: "",
				Skip:                         "",
				OnlyLazy:                     "",
			},
		},
		{
			name: "should fail on missing required env var",
			args: args{
				env:  map[string]string{},
				lazy: false,
			},
			wantErr: true,
			want:    Config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.args.env {
				t.Setenv(k, v)
			}
			tgt := Config{}
			if tt.args.lazy {
				if err := LoadEnvLazy(&tgt); (err != nil) != tt.wantErr {
					t.Errorf("LoadEnv() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err := LoadEnv(&tgt); (err != nil) != tt.wantErr {
					t.Errorf("LoadEnv() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			if !reflect.DeepEqual(reflect.TypeOf(tgt), reflect.TypeOf(tt.want)) {
				t.Errorf("LoadEnv() got = %v, want %v", reflect.TypeOf(tgt), reflect.TypeOf(tt.want))
			}

			if !reflect.DeepEqual(tgt, tt.want) {
				t.Errorf("LoadEnv() got = %v, want %v", tgt, tt.want)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "dash",
			args: args{"my-test"},
			want: "MY_TEST",
		},
		{
			name: "pascal",
			args: args{"MyTest"},
			want: "MY_TEST",
		},
		{
			name: "spaces",
			args: args{"my Test"},
			want: "MY_TEST",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnakeCase(tt.args.str); got != tt.want {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
