import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { Group } from '../group';
import { GroupService } from '../group.service';

@Component({
  selector: 'app-group-list',
  templateUrl: './group-list.component.html',
  styles: []
})
export class GroupListComponent implements OnInit {
  @Output() addGroup = new EventEmitter<void>();
  groups: Group[];
  selectedRow: Group[] = [];
  projectName: string;
  loading = false;

  constructor(private service: GroupService, private route: ActivatedRoute) { }

  ngOnInit() {
    this.projectName = this.route.snapshot.parent.params['project'];
    this.getGroups();
  }

  getGroups(): void {
    this.loading = true;
    this.service.getGroups(this.projectName).subscribe(
      groups => {
        this.groups = groups;
        this.loading = false;
      }
    );
  }

  addGroupTrigger() {
    this.addGroup.emit();
  }

  refresh()  {
    this.getGroups();
  }

}
