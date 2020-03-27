import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {HostService} from '../../host/host.service';
import {Host} from '../../host/host';
import {NgForm} from '@angular/forms';
import {Cluster} from '../../cluster/cluster';

@Component({
  selector: 'app-add-worker',
  templateUrl: './add-worker.component.html',
  styleUrls: ['./add-worker.component.css']
})
export class AddWorkerComponent implements OnInit {

  constructor(private hostService: HostService) {
  }

  hosts: Host[] = [];
  host: string;
  opened = false;
  @Output() openedChange = new EventEmitter();
  @Output() confirm = new EventEmitter();
  @ViewChild('form', {static: true}) form: NgForm;
  @Input() currentCluster: Cluster;


  ngOnInit() {
  }

  loadHosts() {
    this.hostService.byItem(this.currentCluster.item_name).subscribe(data => {
      this.hosts = data.filter(host => {
        return !host.cluster;
      });
    });
  }

  close() {
    this.opened = false;
    this.openedChange.emit(this.opened);
  }

  onConfirm() {
    this.confirm.emit();
  }
}
