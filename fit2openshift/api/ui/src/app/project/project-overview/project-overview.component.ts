import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { Project } from '../project';

@Component({
  selector: 'app-project-overview',
  templateUrl: './project-overview.component.html',
  styles: []
})
export class ProjectOverviewComponent implements OnInit {
  currentProject: Project;
  projectName: string;

  constructor(
    private route: ActivatedRoute) {

    this.currentProject = this.route.snapshot.parent.data['project'];
    this.projectName = this.route.snapshot.parent.params['project'];
  }

  ngOnInit() {
  }
}
