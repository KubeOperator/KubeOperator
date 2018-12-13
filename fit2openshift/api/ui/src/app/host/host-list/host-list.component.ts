import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { Host } from '../host';
import { HostService } from '../host.service';

@Component({
  selector: 'app-host-list',
  templateUrl: './host-list.component.html',
  styles: []
})
export class HostListComponent implements OnInit {
  hosts: Host[];
  selectedRow: Host[] = [];
  projectName: string;
  @Output() addHost = new EventEmitter<void>();
  loading = false;

  constructor(private service: HostService, private route: ActivatedRoute) { }

  ngOnInit() {
    this.projectName = this.route.snapshot.parent.params['project'];
    this.getHosts();
  }

  getHosts() {
    this.loading = true;
    this.service.getHosts(this.projectName).subscribe(
      hosts => {
        this.hosts = hosts;
        this.loading = false;
      }
    );
  }

  createHostTrigger(evt) {
    this.addHost.emit(evt);
  }

  refresh() {
    this.getHosts();
  }
}
