package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

const githubRepo = "https://github.com/fathurmdr/starter-kit.git"
const templatePath = "templates"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("❌ Gunakan format:")
		fmt.Println("   starter-kit create <project-name>")
		os.Exit(1)
	}

	action := os.Args[1]
	projectName := os.Args[2]

	switch action {
	case "create":
		template := selectTemplate()
		createProject(template, projectName)
	default:
		fmt.Println("❌ Perintah tidak dikenal.")
	}
}

func selectTemplate() string {
	fmt.Println("📦 Mengambil daftar template dari repository...")

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("❌ Gagal mendapatkan lokasi program.")
		os.Exit(1)
	}

	tmpDir := filepath.Join(currentDir, "starter-kit-list")
	_ = os.RemoveAll(tmpDir)

	if err := runCommand(fmt.Sprintf("git clone --depth=1 --no-checkout %s %s", githubRepo, tmpDir)); err != nil {
		fmt.Println("❌ Gagal meng-clone repository.")
		os.Exit(1)
	}

	output, err := runCommandOutput(fmt.Sprintf("cd %s && git ls-tree --name-only HEAD %s/", tmpDir, templatePath))
	os.RemoveAll(tmpDir)

	if err != nil || output == "" {
		fmt.Println("❌ Tidak dapat mengambil daftar template.")
		os.Exit(1)
	}

	templates := strings.Split(strings.TrimSpace(output), "\n")
	for i, template := range templates {
		templates[i] = strings.TrimPrefix(template, templatePath+"/")
	}

	prompt := promptui.Select{
		Label: "📜 Pilih Template:",
		Items: templates,
	}

	_, selectedTemplate, err := prompt.Run()
	if err != nil {
		fmt.Println("❌ Pemilihan template dibatalkan.")
		os.Exit(1)
	}

	return selectedTemplate
}

func createProject(templateName, projectName string) {
	targetPath := filepath.Join(".", projectName)
	if _, err := os.Stat(targetPath); err == nil {
		fmt.Println("❌ Folder sudah ada.")
		os.Exit(1)
	}

	fmt.Printf("🚀 Membuat proyek '%s' dari template '%s'...\n", projectName, templateName)

	if err := runCommand(fmt.Sprintf("git clone --depth=1 --no-checkout %s %s", githubRepo, projectName)); err != nil {
		fmt.Println("❌ Gagal meng-clone repository.")
		os.Exit(1)
	}

	commands := []string{
		fmt.Sprintf("cd %s && git sparse-checkout init --cone", projectName),
		fmt.Sprintf("cd %s && git sparse-checkout set %s/%s", projectName, templatePath, templateName),
		fmt.Sprintf("cd %s && git checkout main", projectName),
		fmt.Sprintf("cd %s && mv %s/%s/* .", projectName, templatePath, templateName),
		fmt.Sprintf("cd %s && rm -rf .git %s", projectName, templatePath),
	}

	for _, cmd := range commands {
		if err := runCommand(cmd); err != nil {
			fmt.Println("❌ Gagal membuat proyek.")
			os.RemoveAll(targetPath)
			os.Exit(1)
		}
	}

	fmt.Println("✅ Proyek berhasil dibuat!")
}

func runCommand(cmd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	command := createCommand(ctx, cmd)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func runCommandOutput(cmd string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	command := createCommand(ctx, cmd)
	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = os.Stderr

	err := command.Run()
	return out.String(), err
}

func createCommand(ctx context.Context, cmd string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/C", cmd)
	}
	return exec.CommandContext(ctx, "bash", "-c", cmd)
}
