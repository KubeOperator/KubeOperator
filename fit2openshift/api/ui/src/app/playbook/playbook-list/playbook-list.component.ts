import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { Project } from '../../project/project';
import { PlaybookService } from '../playbook.service';
import { Playbook } from '../playbook';

@Component({
  selector: 'app-playbook-list',
  templateUrl: './playbook-list.component.html',
  styles: []
})
export class PlaybookListComponent implements OnInit {
  currentProject: Project;
  projectName: string;
  playbooks: Playbook[];
  selectedRow: Playbook[] = [];
  loading = false;
  @Output() executePlaybookEvt = new EventEmitter<Playbook>();

  constructor(
    private service: PlaybookService,
    private route: ActivatedRoute) {

    this.currentProject = this.route.snapshot.parent.data['project'];
    this.projectName = this.route.snapshot.parent.params['project'];
  }

  ngOnInit() {
    this.getPlaybooks();
  }

  getPlaybooks() {
    this.loading = true;
    this.service.getPlaybooks(this.projectName)
      .subscribe(playbooks => {
        this.playbooks = playbooks;
        this.loading = false;
      });
  }

  executePlaybook(playbook): void {
    this.executePlaybookEvt.emit(playbook);
  }

}
