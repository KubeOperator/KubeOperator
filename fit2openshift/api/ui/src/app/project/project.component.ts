import { Component, OnInit, ViewChild } from '@angular/core';
import { ProjectCreateComponent } from './project-create/project-create.component';
import { ProjectListComponent } from './project-list/project-list.component';

@Component({
  selector: 'app-project',
  templateUrl: './project.component.html',
  styleUrls: ['./project.component.css']
})
export class ProjectComponent implements OnInit {
  @ViewChild(ProjectCreateComponent)
  creationProject: ProjectCreateComponent;
  @ViewChild(ProjectListComponent)
  listProject: ProjectListComponent;

  constructor() { }

  ngOnInit() {
  }

  openModal() {
    this.creationProject.newProject();
  }

  createProjectTrigger(created: boolean) {
    if (created) {
      this.listProject.refresh();
    }
  }

}
