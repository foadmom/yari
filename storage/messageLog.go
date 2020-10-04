package storage

import (
	"fmt"
	"yari/types"
	"log"
	"os"
	"time"
)

var logFile    *os.File;
var Logger     *log.Logger;

func InitStorage (storage *types.StorageType, config *types.ConfigStorageType) error {
	var  	_err	error;

	storage.Config = config;

	Logger = log.New(os.Stdout, "" , 0)
	Logger.SetFlags(0);
	log.SetOutput(new(logWriter));
	
	storage.SystemLog = Logger;
	Logger.Println ("logger text");

	return _err;
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
    return fmt.Print(time.Now().UTC().Format("2006-01-02T15:04:05.999999Z") + " [DEBUG] " + string(bytes))
}

func Log (buffer string) {
	Logger.Println (buffer);
}

func LogMessage (key string, value []byte) error {
	return logToFile (key, value);
}

func logToFile (key string, value []byte) error {
//	var _fileName	string = common.CONFIG_LogFilePath;
	var _buffer		string = "key:"+key+" value:"+string(value)+"\n";
	var _err		error;

	_, _err = logFile.WriteString(_buffer);
	if (_err != nil) {
		log.Println(_err);
	}

	// return _err
	return nil;
}

func makeupFileName (path string) string {
	var _logFileName string;
	_now := time.Now ();
	_date := _now.Format ("20060102");
	_logFileName = path+"/YARI-"+_date+".log"
	return _logFileName;
}


func OnClose (config types.ConfigStorageType) {
	if (logFile != nil) {
		logFile.Close();
	}
}
