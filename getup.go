package getup

import (
	"path/filepath"
	"os"
	"os/exec"
	"log"
	"io/ioutil"

	"github.com/src-d/go-git"
)

const configDir = "~/.config"
const configPath = filepath.Join(configDir, "dgellow")
const configRepo = "git@github.com:dgellow/config.git"
const ohmyzshPath = "~/.oh-my-zsh"
const ohmyzshRepo = "git@github.com:dgellow/oh-my-zsh.git"
const goPath = "~/Development/Go/gopath"

const packages = []string{
	"git",
	"tree",
	"emacs",
	"go",
	"python",
	"python3",
	"zsh",
}

const symlinkIgnore = []string{
	"boot",
	"bootstrap",
	"REAMDE.md",
	"Vagrantfile",
	"Xmodmap",
}

func main() {
	mkdirConfig()
	cloneConfig()
	// Deps
	brew()
	brewPackages()
	ohmyzsh()
	golang()
	// Symlink dotfiles
	symlinks()
}

func done() {
	log.Println("\t[done]")
}

func mkdirConfig() {
	log.Printf("Make config dir...")
	err := os.MkdirAll(configPath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	done()
}

func cloneConfig() {
	log.Printf("Cloning config...")
	_, err := git.PlainClone(configPath, false, &git.CloneOptions{
		URL: configRepo,
		Progress: os.Stdout,
	})
	done()
}

func brew() {
	log.Printf("Install brew...")
	const installBrew = `/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
	_, err := exec.Command(installBrew).Output()
	if err != nil {
		log.Fatal(err)
	}
	done()

	log.Printf("Install brew cask...")
	const installBrewCask = `/usr/local/bin/brew tap caskroom/cask`
	done()
}

func brewPackages() {
	const installCmd = `/usr/local/bin/brew install %s`
	log.Printf("Install packages...")
	for p := range packages {
		log.Printf("\t- %s\n", p)
		out, err := exec.Command(fmt.Sprintf(installCmd, p)).Output()
		if err != nil {
			log.Fatal(err)
		}
	}
	done()
}

func ohmyzsh() {
	log.Printf("Clone ohmyzsh fork...")
	_, err := git.PlainClone(ohmyzshPath, false, &git.CloneOptions{
		URL: ohmyzshRepo,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatal(err)
	}
	done()
}

func golang() {
	log.Printf("Setup go env...")
	err := os.MkdirAll(goPath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	done()
}

func symlinks() {
	log.Printf("Symlink dotfiles...")
	for f := range ioutil.ReadDir(configPath) {
		log.Printf("\t- %s\n", f)
		err := os.Symlink(filepath.Join(configPath, f), filepath.Join("~", "."+f))
		if err != nil {
			log.Fatal(err)
		}
	}
	done()
}
