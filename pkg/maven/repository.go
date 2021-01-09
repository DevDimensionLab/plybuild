package maven

import (
	"errors"
	"github.com/co-pilot-cli/co-pilot/pkg/file"
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
	Url  string
	Auth *RepositoryAuth
}

type RepositoryAuth struct {
	Username string
	Password string
}

func (settings Settings) GetRepositories() (Repositories, error) {
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

	for _, mirrors := range settings.Settings.Mirrors {
		if mirrors.Mirror.URL != "" {
			mirrorRepo := Repository{
				Url: mirrors.Mirror.URL,
			}
			if hasServer, server := settings.Settings.FindServerWith(mirrors.Mirror.ID); hasServer {
				mirrorRepo.Auth = &RepositoryAuth{
					Username: server.Username,
					Password: server.Password,
				}
			}
			repos.Mirror = append(repos.Mirror, mirrorRepo)
		}
	}

	return repos, nil
}

func (settings M2Settings) FindServerWith(id string) (bool, Server) {
	if settings.Servers == nil {
		return false, Server{}
	}

	for _, servers := range settings.Servers {
		if servers.Server.ID == id {
			return true, servers.Server
		}
	}

	return false, Server{}
}

func (repos Repositories) GetDefaultRepository() Repository {
	if len(repos.Mirror) > 0 {
		return repos.Mirror[0]
	} else {
		return repos.Fallback
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
