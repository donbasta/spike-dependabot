package task

import "log"

func UpdateProjects(changes []Changes) {
	for i := 0; i < len(changes); i++ {
		log.Println(changes[i].project.Name)
		log.Println(changes[i].DepChanges)
	}
}
