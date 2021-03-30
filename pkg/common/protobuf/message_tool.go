/*
 * Copyright (c) 2018-2118
 * Author: linceyou
 * LastModified: 18-2-3 下午4:50
 */

package protobuf

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/golang/protobuf/proto"
)

// DumpMessageToFile dump proto message object to file
func DumpMessageToFile(message proto.Message, fileName string, forceOverride bool) error {
	// if file already exist and not force to override
	if _, err := os.Stat(fileName); os.IsExist(err) && !forceOverride {
		return fmt.Errorf("failed to dump to existed file[%s] without force override",
			fileName)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file[%s]", fileName)
	}
	defer file.Close()

	buffer, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal proto to buffer, err: %s", err)
	} else {
		if _, err := file.Write(buffer); err != nil {
			return fmt.Errorf("failed to write to file, err: %s", err)
		}
	}

	return nil
}

// LoadFileToMessage load proto message from file
func LoadFileToMessage(fileName string, message proto.Message) error {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read file[%s], err[%s]", fileName, err)
	}

	err = proto.Unmarshal(fileContent, message)
	if err != nil {
		return fmt.Errorf("failed to unmarshal, err[%s]", err)
	}

	return nil
}

// DumpMessageToTextFile
func DumpMessageToTextFile(message proto.Message, fileName string, forceOverride bool) error {
	// if file already exist and not force to override
	if _, err := os.Stat(fileName); os.IsExist(err) && !forceOverride {
		return fmt.Errorf("failed to dump to existed file[%s] without force override",
			fileName)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file[%s]", fileName)
	}
	defer file.Close()

	if err := proto.MarshalText(file, message); err != nil {
		return fmt.Errorf("failed to write to file, err: %s", err)
	}
	return nil
}

// LoadTextFileToMessage
func LoadTextFileToMessage(fileName string, message proto.Message) error {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read file[%s], err[%s]",
			fileName, err)
	}

	strData := string(fileContent)
	err = proto.UnmarshalText(strData, message)
	if err != nil {
		return fmt.Errorf("failed to unmarshal, err[%s]", err)
	}

	return nil
}

// CheckMarshal 检查序列化后能正常反序列化
func CheckMarshal(pb proto.Message) error {
	in := reflect.ValueOf(pb)
	if in.IsNil() {
		return fmt.Errorf("nil message")
	}

	bytes, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	out := reflect.New(in.Type().Elem())
	dst := out.Interface().(proto.Message)
	if err := proto.Unmarshal(bytes, dst); err != nil {
		return err
	}
	return nil
}
