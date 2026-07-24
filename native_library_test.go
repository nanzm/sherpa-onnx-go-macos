package sherpa_onnx

import (
	"bytes"
	"debug/macho"
	"os"
	"testing"
)

func TestPackagedNativeLibraries(t *testing.T) {
	tests := []struct {
		name string
		path string
		cpu  macho.Cpu
	}{
		{
			name: "arm64",
			path: "lib/aarch64-apple-darwin/libsherpa-onnx-c-api.dylib",
			cpu:  macho.CpuArm64,
		},
		{
			name: "x86_64",
			path: "lib/x86_64-apple-darwin/libsherpa-onnx-c-api.dylib",
			cpu:  macho.CpuAmd64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := macho.Open(tt.path)
			if err != nil {
				t.Fatalf("open native library: %v", err)
			}
			defer f.Close()

			if f.Cpu != tt.cpu {
				t.Fatalf("CPU = %v, want %v", f.Cpu, tt.cpu)
			}

			assertMachOSymbol(t, f, "_SherpaOnnxCreateOfflineTts")
			assertMachOSymbol(
				t,
				f,
				"_SherpaOnnxOfflineSpeakerDiarizationProcessWithCallback",
			)

			data, err := os.ReadFile(tt.path)
			if err != nil {
				t.Fatalf("read native library: %v", err)
			}
			if bytes.Contains(data, []byte("TTS is not enabled. Please rebuild sherpa-onnx")) {
				t.Fatal("native library was built without TTS support")
			}
		})
	}
}

func assertMachOSymbol(t *testing.T, f *macho.File, want string) {
	t.Helper()

	if f.Symtab == nil {
		t.Fatal("native library has no symbol table")
	}
	for _, symbol := range f.Symtab.Syms {
		if symbol.Name == want {
			return
		}
	}
	t.Fatalf("native library is missing symbol %q", want)
}
