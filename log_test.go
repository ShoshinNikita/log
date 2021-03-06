package clog

import (
	"bytes"
	stdlog "log"
	"os"
	"testing"
)

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func TestLoggerLevels(t *testing.T) {
	printFunction := func(log *Logger) {
		log.Debug("debug")
		log.Debugf("debugf %s\n", "arg")

		log.Info("info")
		log.Infof("infof %s", "arg")

		log.Warn("warn")
		log.Warnf("warnf %s %d", "arg", 15)

		log.Error("error")
		log.Errorf("errorf %s", "arg")

		// log.Fatal("fatal")
		// log.Fatalf("fatalf %s", "arg")

		log.Print("print")
		log.Printf("printf %s", "arg")

		log.Write([]byte("bytes"))
		log.WriteString("string")
	}

	tests := []struct {
		description string
		config      *Config
		output      []byte
	}{
		{
			description: "debug level",
			config: &Config{
				level:          LevelDebug,
				printColor:     false,
				printErrorLine: false,
				printTime:      false,
				timeLayout:     DefaultTimeLayout,
			},
			output: []byte(
				"[DBG] debug\n" +
					"[DBG] debugf arg\n\n" +
					"[INF] info\n" +
					"[INF] infof arg\n" +
					"[WRN] warn\n" +
					"[WRN] warnf arg 15\n" +
					"[ERR] error\n" +
					"[ERR] errorf arg\n" +
					"print\n" +
					"printf arg\n" +
					"bytesstring"),
		},
		{
			description: "info level",
			config: &Config{
				level:          LevelInfo,
				printColor:     false,
				printErrorLine: false,
				printTime:      false,
				timeLayout:     DefaultTimeLayout,
			},
			output: []byte(
				"[INF] info\n" +
					"[INF] infof arg\n" +
					"[WRN] warn\n" +
					"[WRN] warnf arg 15\n" +
					"[ERR] error\n" +
					"[ERR] errorf arg\n" +
					"print\n" +
					"printf arg\n" +
					"bytesstring"),
		},
		{
			description: "warn level",
			config: &Config{
				level:          LevelWarn,
				printColor:     false,
				printErrorLine: false,
				printTime:      false,
				timeLayout:     DefaultTimeLayout,
			},
			output: []byte(
				"[WRN] warn\n" +
					"[WRN] warnf arg 15\n" +
					"[ERR] error\n" +
					"[ERR] errorf arg\n" +
					"print\n" +
					"printf arg\n" +
					"bytesstring"),
		},
		{
			description: "error level",
			config: &Config{
				level:          LevelError,
				printColor:     false,
				printErrorLine: false,
				printTime:      false,
				timeLayout:     DefaultTimeLayout,
			},
			output: []byte(
				"[ERR] error\n" +
					"[ERR] errorf arg\n" +
					"print\n" +
					"printf arg\n" +
					"bytesstring"),
		},
		{
			description: "fatal level",
			config: &Config{
				level:          LevelFatal,
				printColor:     false,
				printErrorLine: false,
				printTime:      false,
				timeLayout:     DefaultTimeLayout,
			},
			output: []byte(
				"print\n" +
					"printf arg\n" +
					"bytesstring"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			buff := &bytes.Buffer{}
			tt.config.SetOutput(buff)

			log := tt.config.Build()

			printFunction(log)

			res := buff.Bytes()
			if !bytes.Equal(res, tt.output) {
				t.Errorf("different output")
				t.Log(string(res))
				t.Log(string(tt.output))
			}
		})
	}
}

func TestWithPrefix(t *testing.T) {
	printFunction := func(log *Logger) {
		log.Debug("debug")
		log.Debugf("debugf %s", "arg")

		log.Info("info")
		log.Infof("infof %s", "arg")

		log.Warn("warn")
		log.Warnf("warnf %s %d", "arg", 15)

		log.Error("error")
		log.Errorf("errorf %s\n", "arg")

		// log.Fatal("fatal")
		// log.Fatalf("fatalf %s", "arg")

		log.Print("print")
		log.Printf("printf %s", "arg")

		log.Write([]byte("bytes"))
		log.WriteString("string")
	}

	tests := []struct {
		log    *Logger
		output []byte
	}{
		{
			log: NewDevConfig().
				PrintColor(false).
				PrintErrorLine(false).
				PrintTime(false).
				SetPrefix("prefix").
				Build(),
			output: []byte(
				"[DBG] prefixdebug\n" +
					"[DBG] prefixdebugf arg\n" +
					"[INF] prefixinfo\n" +
					"[INF] prefixinfof arg\n" +
					"[WRN] prefixwarn\n" +
					"[WRN] prefixwarnf arg 15\n" +
					"[ERR] prefixerror\n" +
					"[ERR] prefixerrorf arg\n\n" +
					"prefixprint\n" +
					"prefixprintf arg\n" +
					"bytesstring"),
		},
		{
			log: NewDevConfig().
				PrintColor(false).
				PrintErrorLine(false).
				PrintTime(false).
				Build().WithPrefix("prefix"),
			output: []byte(
				"[DBG] prefix: debug\n" +
					"[DBG] prefix: debugf arg\n" +
					"[INF] prefix: info\n" +
					"[INF] prefix: infof arg\n" +
					"[WRN] prefix: warn\n" +
					"[WRN] prefix: warnf arg 15\n" +
					"[ERR] prefix: error\n" +
					"[ERR] prefix: errorf arg\n\n" +
					"prefix: print\n" +
					"prefix: printf arg\n" +
					"bytesstring"),
		},
		{
			log: NewDevConfig().
				PrintColor(false).
				PrintErrorLine(false).
				PrintTime(false).
				Build().WithPrefix("[first prefix]").WithPrefix("[second prefix]"),
			output: []byte(
				    "[DBG] [first prefix]: [second prefix]: debug\n" +
					"[DBG] [first prefix]: [second prefix]: debugf arg\n" +
					"[INF] [first prefix]: [second prefix]: info\n" +
					"[INF] [first prefix]: [second prefix]: infof arg\n" +
					"[WRN] [first prefix]: [second prefix]: warn\n" +
					"[WRN] [first prefix]: [second prefix]: warnf arg 15\n" +
					"[ERR] [first prefix]: [second prefix]: error\n" +
					"[ERR] [first prefix]: [second prefix]: errorf arg\n\n" +
					"[first prefix]: [second prefix]: print\n" +
					"[first prefix]: [second prefix]: printf arg\n" +
					"bytesstring"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			buff := &bytes.Buffer{}
			tt.log.output = buff

			printFunction(tt.log)

			res := buff.Bytes()
			if !bytes.Equal(res, tt.output) {
				t.Errorf("different output")
				t.Log(string(res))
				t.Log(string(tt.output))
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Benchmarks
// -----------------------------------------------------------------------------

const (
	file = "test.txt"
	msg  = "Hello, dear world!!!"
)

func BenchmarkStdLogPrintlnWithPrefixes(b *testing.B) {
	f, err := os.Create(file)
	if err != nil {
		stdlog.Fatalln(err)
	}
	defer f.Close()

	l := &stdlog.Logger{}

	// github.com/ShoshinNikita/log prints it by default
	l.SetFlags(stdlog.Lshortfile | stdlog.Ltime | stdlog.Ldate)
	l.SetOutput(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Println(msg)
	}
}

func BenchmarkDevLogPrintln(b *testing.B) {
	f, err := os.Create(file)
	if err != nil {
		stdlog.Fatalln(err)
	}
	defer f.Close()

	l := NewDevConfig().SetOutput(f).Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Print(msg)
	}
}

func BenchmarkDevLogErrorln(b *testing.B) {
	f, err := os.Create(file)
	if err != nil {
		stdlog.Fatalln(err)
	}
	defer f.Close()

	l := NewDevConfig().SetOutput(f).Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Error(msg)
	}
}

func BenchmarkProdLogPrintln(b *testing.B) {
	f, err := os.Create(file)
	if err != nil {
		stdlog.Fatalln(err)
	}
	defer f.Close()

	l := NewProdConfig().SetOutput(f).Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Print(msg)
	}
}

func BenchmarkProdLogErrorln(b *testing.B) {
	f, err := os.Create(file)
	if err != nil {
		stdlog.Fatalln(err)
	}
	defer f.Close()

	l := NewProdConfig().SetOutput(f).Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Error(msg)
	}
}
