package harbor

type Registry struct {
	Status       string             `json:"status,omitempty"`
	UpdateTime   string             `json:"update_time,omitempty"`
	Name         string             `json:"name,omitempty"`
	URL          string             `json:"url,omitempty"`
	Insecure     bool               `json:"insecure,omitempty"`
	CreationTime string             `json:"creation_time,omitempty"`
	Type         string             `json:"type,omitempty"`
	ID           int64              `json:"id,omitempty"`
	Description  string             `json:"description,omitempty"`
	Credential   RegistryCredential `json:"credential,omitempty"`
}

type RegistryCredential struct {
	AccessKey    string `json:"access_key,omitempty"`
	AccessSecret string `json:"access_secret,omitempty"`
	Type         string `json:"type,omitempty"`
}

func (client *Client) GetRegistry(id string) (*Registry, error) {
	var registry *Registry

	err := client.get(id, &registry, nil)
	if err != nil {
		return nil, err
	}

	return registry, nil
}

func (client *Client) NewRegistry(registry *Registry) (string, error) {
	_, location, err := client.post("/registries", registry)
	if err != nil {
		return "", err
	}

	return location, nil
}

func (client *Client) UpdateRegistry(id string, Registry *Registry) error {
	return client.put(id, Registry)
}

func (client *Client) DeleteRegistry(id string) error {
	return client.delete(id, nil)
}
