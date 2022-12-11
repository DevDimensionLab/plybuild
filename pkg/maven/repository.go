package maven

import (
	"errors"
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"golang.org/x/term"
	"os/user"
)

type Settings struct {
	Path     string
	Settings M2Settings
}

type Repositories struct {
	Fallback Repository
	Profile  []Repository
	Mirror   []Repository
}

type Repository struct {
	Id   string
	Url  string
	Auth *RepositoryAuth
}

type RepositoryAuth struct {
	Username  string
	Password  string
	Encrypted bool
}

func (settings Settings) GetRepositories() (Repositories, error) {
	usr, err := user.Current()
	if err != nil {
		return Repositories{}, err
	}

	repos := Repositories{
		Fallback: Repository{
			Url: "https://repo1.maven.org/maven2",
		},
	}

	for _, profile := range settings.Settings.Profiles.Profile {
		for _, repo := range profile.Repositories.Repository {
			if repo.Releases.Enabled && repo.URL != "" {
				repos.Profile = append(repos.Profile, Repository{
					Url: repo.URL,
				})
			}
		}
	}

	for _, mirror := range settings.Settings.Mirrors.Mirror {
		if mirror.URL != "" {
			log.Debugf("Found mirror %s in m2settings\n", mirror.ID)
			mirrorRepo := Repository{
				Id:  mirror.ID,
				Url: mirror.URL,
			}
			if hasServer, server := settings.Settings.FindServerWith(mirror.ID); hasServer {
				log.Debugf("Mirror %s has a matching server, trying to copy credentials\n", mirror.ID)
				mirrorRepo.Auth = &RepositoryAuth{
					Username:  server.Username,
					Password:  server.Password,
					Encrypted: file.Exists(file.Path("%s/.m2/settings-security.xml", usr.HomeDir)),
				}
			}
			repos.Mirror = append(repos.Mirror, mirrorRepo)
		}
	}

	return repos, nil
}

func (settings M2Settings) FindServerWith(id string) (bool, Server) {
	if settings.Servers.Server == nil {
		return false, Server{}
	}

	for _, server := range settings.Servers.Server {
		if server.ID == id {
			return true, server
		}
	}

	return false, Server{}
}

func (repos Repositories) GetDefaultRepository() (Repository, error) {
	if len(repos.Mirror) > 0 {
		repo := repos.Mirror[0]
		log.Debugf("found mirrors, using the first mirror [%s] \n", repo.Id)
		if repo.Auth != nil && repo.Auth.Encrypted {
			fmt.Printf("!! Password for [%s] seems to be encrypted, please enter password: ", repo.Id)
			bytePassword, err := term.ReadPassword(0)
			fmt.Println()
			if err != nil {
				return Repository{}, err
			}
			repo.Auth.Password = string(bytePassword)
		}
		return repo, nil
	} else {
		log.Debugf("using the fallback maven repository %s\n", repos.Fallback.Url)
		return repos.Fallback, nil
	}
}

func (settings Settings) ListRepositories() error {
	repos, err := settings.GetRepositories()
	if err != nil {
		return err
	}

	for _, profileRepo := range repos.Profile {
		log.Infof("found maven profile repository: %s", profileRepo)
	}

	for _, mirrorRepo := range repos.Mirror {
		log.Infof("found maven mirror repository: %s", mirrorRepo)
	}

	log.Infof("maven repository: %s", repos.Fallback)

	return nil
}

func NewSettings() (settings Settings, err error) {
	usr, err := user.Current()
	if err != nil {
		return settings, err
	}

	home := usr.HomeDir
	m2Settings := file.Path("%s/.m2/settings.xml", home)
	confSettings := file.Path("%s/conf/settings.xml", home)

	if file.Exists(m2Settings) {
		settings.Path = m2Settings
	} else if file.Exists(confSettings) {
		settings.Path = confSettings
	} else {
		return Settings{}, errors.New("could not find settings.xml")
	}

	err = file.ReadXml(settings.Path, &settings.Settings)
	if err != nil {
		return Settings{}, err
	}
	return
}

func DefaultRepository() (Repository, error) {
	settings, _ := NewSettings()
	repos, err := settings.GetRepositories()
	if err != nil {
		return Repository{}, err
	}

	return repos.GetDefaultRepository()
}

func RepositoryFrom(url string, username string, password string) Repository {
	repo := Repository{
		Url: url,
	}

	if username != "" && password != "" {
		repo.Auth = &RepositoryAuth{
			Username: username,
			Password: password,
		}
	}

	return repo
}
