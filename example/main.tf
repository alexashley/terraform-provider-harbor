provider "harbor" {
  url = "http://localhost:30002"
  username = ""
  password = ""
}

resource "harbor_project" "project" {
  name = "project"
}

resource "harbor_robot_account" "robot" {
  name = "robot$robot"
  project_id = harbor_project.project.id
  robot_account_access {
    resource = "image"
    action = "pull"
  }
}

resource "harbor_registry" "registry" {
  name = "registry"
  type = "docker-hub"
  endpoint_url = "http://hub.docker.com"
}

