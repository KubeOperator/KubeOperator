import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {HostService} from '../../host/host.service';
import {Host} from '../../host/host';
import {NgForm} from '@angular/forms';
import {Cluster} from '../../cluster/cluster';
import {Profile} from '../../shared/session-user';

@Component({
  selector: 'app-add-worker',
  templateUrl: './add-worker.component.html',
  styleUrls: ['./add-worker.component.css']
})
export class AddWorkerComponent implements OnInit {

  constructor(private hostService: HostService) {
  }

  hosts: Host[] = [];
  options: any[] = [];
  host_names: string[];
  opened = false;
  @Output() openedChange = new EventEmitter();
  @Output() confirm = new EventEmitter();
  @ViewChild('form', {static: true}) form: NgForm;
  @Input() currentCluster: Cluster;
  ops: any = {
    multiple: true,
    placeholder: '选择节点',
    escapeMarkup: function (markup) {
      return markup;
    },
    templateSelection: (data) => {
      return `<span class="label label-blue select2-selection__choice__remove">${data['text']}</span>`;
    },
  };


  ngOnInit() {
  }

  loadHosts() {
    this.hostService.byItem(this.currentCluster.item_name).subscribe(data => {
      this.hosts = data.filter(host => {
        return !host.cluster;
      });
      this.options = this.toOptions(this.hosts);
    });
  }

  close() {
    this.opened = false;
    this.openedChange.emit(this.opened);
  }

  onConfirm() {
    this.confirm.emit();
  }

  private toOptions(hosts: Host[]): any[] {
    const options = [];
    hosts.forEach(h => {
      options.push({'id': h.id, 'text': h.name, 'value': h.name});
    });
    return options;
  }
}
