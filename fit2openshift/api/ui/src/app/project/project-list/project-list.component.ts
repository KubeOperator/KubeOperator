import { Component, OnInit, EventEmitter, Output } from '@angular/core';
import { Project } from '../project';
import { ProjectService } from '../project.service';
import { NavigatorService } from '../../base/navigator/navigator.service';

@Component({
  selector: 'app-project-list',
  templateUrl: './project-list.component.html',
  styles: []
})
export class ProjectListComponent implements OnInit {
  projects:  Project[];
  selectedRow: Project[] = [];
  @Output() addProject = new EventEmitter<void>();
  loading = false;

  constructor(private service: ProjectService, private navService: NavigatorService) {
    this.navService.showProjectsNav();
  }

  ngOnInit() {
    this.getProjects();
  }

  getProjects(): void {
    this.service.getProjects()
      .subscribe(projects => this.projects = projects);
  }

  addNewProject(): void {
    this.addProject.emit();
  }

  refresh(): void {
    this.getProjects();
  }

}
