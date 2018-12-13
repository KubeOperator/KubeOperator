import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Project } from '../project';
import { NavigatorService } from '../../base/navigator/navigator.service';

@Component({
  selector: 'app-project-detail',
  templateUrl: './project-detail.component.html',
  styles: []
})
export class ProjectDetailComponent implements OnInit {
  currentProject: Project;
  projectName: string;

  constructor(
    private navService: NavigatorService,
    private route: ActivatedRoute) {

    this.currentProject = this.route.snapshot.data['project'];
    this.projectName = this.route.snapshot.params['project'];
    console.log(this.projectName);
    this.navService.showDetailNav(this.currentProject.name);
  }

  ngOnInit() {
  }
}
