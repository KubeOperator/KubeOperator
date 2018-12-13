import { Component, EventEmitter, OnInit, Output, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { NgForm } from '@angular/forms';

import { Group } from '../group';
import { GroupService } from '../group.service';
import { HostService } from '../../host/host.service';

@Component({
  selector: 'app-group-create',
  templateUrl: './group-create.component.html',
  styles: []
})
export class GroupCreateComponent implements OnInit {
  shown = false;
  @Output() created = new EventEmitter<boolean>();
  projectName: string;
  group: Group = new Group();
  hosts: string[];
  children: string[];
  @ViewChild('groupForm')
  currentForm: NgForm;

  constructor(private service: GroupService,
              private hostService: HostService,
              private route: ActivatedRoute) { }

  ngOnInit() {
    this.projectName = this.route.snapshot.parent.params['project'];
  }

  getHosts() {
    this.hosts = [];
    this.hostService.getHosts(this.projectName).subscribe(
      hosts => {
        this.hosts = hosts.map(host => host.name);
      }
    );
  }

  getGroups() {
    this.children = [];
    this.service.getGroups(this.projectName).subscribe(
      groups => {
        this.children = groups.map(group => group.name);
      }
    );
  }

  newGroup() {
    this.getHosts();
    this.getGroups();
    this.group = new Group();
    this.shown = true;
  }

  get isValid() {
    console.log(this.currentForm);
    return this.currentForm.valid;
  }

  onCancel() {
    this.shown = false;
  }

  onSubmit() {
    this.service.createGroup(this.projectName, this.group).subscribe(
      group => {
        this.created.emit(true);
        this.shown = false;
        console.log('Group created: ', group);
      }
    );
  }
}
