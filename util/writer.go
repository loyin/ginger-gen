package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
)


const (
	HandlerOutput = "handler"
	ModelOutput   = "model"
	RepoOutput    = "repository"
	CacheOutput   = "cache"
)


// 输出
func OutputFile(root, module string) (io.Writer, error) {
	// 创建输出目录
	err := os.MkdirAll(root, 0755)
	if err != nil {
		// 如目录创建失败，则标准输出
		return os.Stdout, err
	}

	// 如: /handler/user_handler.go 、 /model/user_model.go 、 /repository/user_repository.go
	filename := root + "/" + module + "_" + root + ".go"

	// 创建输出的目录并创建输出的go文件
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	// 如存在，直接输出该文件
	if err == nil {
		OutputInfo("OpenFile", "File is exist,the code will append to this file.")
		return file, nil
	}

	// 不存在，创建文件并加文件头
	var std = os.Stdout
	var out io.Reader
	if os.IsNotExist(err) {
		OutputInfo("OpenFile", "File is not exist,create new file.")
		switch root {
		case HandlerOutput:
			out = addHandlerImportContent()
		case ModelOutput:
			out = addModelImportContent(module)
		case RepoOutput:
			out = addRepoImportContent(module)
		case CacheOutput:
			out = addCacheImportContent()
		}

		file, err = os.Create(filename)
		if err != nil {
			io.Copy(std, out)
			return std, err
		}

		io.Copy(file, out)
		return file, nil
	}

	// 有其它错误，则标准输出
	return std, err
}

func addHandlerImportContent() io.Reader {
	return bytes.NewBuffer([]byte(fmt.Sprintf(`package handler

import(
    "github.com/gofuncchan/ginger/common"
    "github.com/gin-gonic/gin"
)

/*
This code is generated with ginger-cli.
You must reset Request Params, and implement biz logic code.
*/

	`)))
}

func addModelImportContent(module string) io.Reader {
	return bytes.NewBuffer([]byte(fmt.Sprintf(`package model

import(
    builder "github.com/gofuncchan/ginger/dao/mysql/%s_builder"
    "github.com/gofuncchan/ginger/util/e"
)

/*
This code is generated with ginger-cli.
You should handling errors in model function,and return data or result to caller.

For example:

func CreateUserByPhone(name, phone, passwd, salt string) int64 {
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"name":   name,
		"phone":  phone,
		"password": password,
		"salt":   salt,
	})

	id, err := builder.Insert(data)
	if !e.Em(err) {
		return -1
	}
	return id
}

*/

	`,module)))
}

func addRepoImportContent(module string) io.Reader {
	return bytes.NewBuffer([]byte(fmt.Sprintf(`package repository

import(
   	mongo "github.com/gofuncchan/ginger/dao/mongodb"
	"github.com/gofuncchan/ginger/util/e"
	"gopkg.in/mgo.v2"
)

/*
This code is generated with ginger-cli.
You should handling errors in repository function,and return data or result to caller.

For example:

func InsertPost(dataMap map[string]interface{}) (b bool) {

	// Use mongodb dao common function
	err := mongo.Insert(MongoPostCollection,dataMap)
	
	if !e.Em(err) {
		return false
	}

	return true
}
*/

const Mongo%sCollection = "%s"

	`,CamelString(module),module)))
}


func addCacheImportContent() io.Reader {
	return bytes.NewBuffer([]byte(fmt.Sprintf(`package cache

import(
	redigo "github.com/garyburd/redigo/redis"
	"github.com/gofuncchan/ginger/dao/redis"
)

/*
This code is generated with ginger-cli.
You should handling errors in cache function,and return data or result to caller.

For example:

func SetKey(key, value string) bool {
	rs, _ := redigo.String(redis.R("SET", key, value, "EX", 3600))
	return rs == "OK"
}

*/

	`)))
}