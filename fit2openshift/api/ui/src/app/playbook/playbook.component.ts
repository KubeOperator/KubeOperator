import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { PlaybookService } from './playbook.service';
import { Project } from '../project/project';
import { PlaybookExecution } from './playbook';

@Component({
  selector: 'app-playbook',
  templateUrl: './playbook.component.html',
  styles: []
})
export class PlaybookComponent implements OnInit, OnDestroy {
  currentProject: Project;
  projectName: string;
  isLogModalShow = false;
  execution: PlaybookExecution;
  // taskId: string;

  constructor(private playbookService: PlaybookService, private route: ActivatedRoute) {
    this.currentProject = this.route.snapshot.parent.data['project'];
    this.projectName = this.route.snapshot.parent.params['project'];
  }

  ngOnInit() {
  }

  ngOnDestroy() {
  }

  executePlaybook(playbook): void {
    this.playbookService.executePlaybook(playbook).subscribe(
      (execution) => {
        this.execution = execution;
        this.isLogModalShow = true;
      }
    );
  }

  closeModal(): void {
    this.isLogModalShow = false;
  }

}
