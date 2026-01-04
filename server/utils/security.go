package utils

import (
	"bytes"
	"mime/multipart"
	"strings"
)

// FileSecurityScanner 文件安全扫描器
type FileSecurityScanner struct{}

// 恶意文件特征签名
var maliciousSignatures = [][]byte{
	// PHP 标签
	[]byte("<?php"),
	[]byte("<%"),
	[]byte("<script"),

	// 可执行文件魔术字节
	[]byte{0x4D, 0x5A},             // PE executable (Windows .exe)
	[]byte{0x7F, 0x45, 0x4C, 0x46}, // ELF executable (Linux)
	[]byte{0xCA, 0xFE, 0xBA, 0xBE}, // Mach-O executable (macOS)
	[]byte{0xFE, 0xED, 0xFA, 0xCE}, // Mach-O executable (macOS, different endian)

	// 脚本文件标识
	[]byte("#!/bin/sh"),
	[]byte("#!/bin/bash"),
	[]byte("#!/usr/bin/env"),

	// VBS/PowerShell
	[]byte("WScript"),
	[]byte("CreateObject"),
	[]byte("powershell"),
	[]byte("cmd.exe"),

	// SQL注入特征
	[]byte("UNION SELECT"),
	[]byte("DROP TABLE"),
	[]byte("DELETE FROM"),
	[]byte("INSERT INTO"),
}

// 恶意文件内容模式
var maliciousPatterns = []string{
	"eval(",
	"exec(",
	"system(",
	"shell_exec(",
	"passthru(",
	"base64_decode(",
	"file_get_contents(",
	"file_put_contents(",
	"fopen(",
	"fwrite(",
	"include(",
	"require(",
	"javascript:",
	"vbscript:",
	"data:text/html",
	"data:application/",
}

// ScanFile 扫描文件是否包含恶意内容
func (s *FileSecurityScanner) ScanFile(file *multipart.FileHeader) error {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 读取文件前8KB进行检测
	buffer := make([]byte, 8192)
	n, _ := src.Read(buffer)
	content := buffer[:n]

	// 检查文件签名
	if err := s.checkFileSignatures(content); err != nil {
		return err
	}

	// 检查恶意模式
	if err := s.checkMaliciousPatterns(content); err != nil {
		return err
	}

	return nil
}

// checkFileSignatures 检查文件魔术字节签名
func (s *FileSecurityScanner) checkFileSignatures(content []byte) error {
	for _, signature := range maliciousSignatures {
		if bytes.HasPrefix(content, signature) || bytes.Contains(content, signature) {
			return NewSecurityError("检测到潜在的恶意文件签名")
		}
	}
	return nil
}

// checkMaliciousPatterns 检查恶意代码模式
func (s *FileSecurityScanner) checkMaliciousPatterns(content []byte) error {
	contentStr := strings.ToLower(string(content))

	for _, pattern := range maliciousPatterns {
		if strings.Contains(contentStr, strings.ToLower(pattern)) {
			return NewSecurityError("检测到潜在的恶意代码模式")
		}
	}
	return nil
}

// SecurityError 安全错误类型
type SecurityError struct {
	Message string
}

func (e *SecurityError) Error() string {
	return e.Message
}

func NewSecurityError(message string) *SecurityError {
	return &SecurityError{Message: message}
}

// IsSecurityError 判断是否为安全错误
func IsSecurityError(err error) bool {
	_, ok := err.(*SecurityError)
	return ok
}
