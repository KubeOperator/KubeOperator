import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {NodeService} from '../../node/node.service';
import {Node} from '../../node/node';
import {Host} from '../../host/host';

@Component({
  selector: 'app-remove-worker',
  templateUrl: './remove-worker.component.html',
  styleUrls: ['./remove-worker.component.css']
})
export class RemoveWorkerComponent implements OnInit {

  constructor(private nodeService: NodeService) {
  }

  opened = false;
  worker_names = [];
  workers: Node[] = [];
  options: any[] = [];
  @Output() openedChange = new EventEmitter();
  @Output() confirm = new EventEmitter();
  @ViewChild('form', {static: true}) form: NgForm;
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

  loadNodes(clusterName: string) {
    this.nodeService.listNodes(clusterName).subscribe(data => {
      this.workers = data.filter(worker => {
        return worker.roles.includes('worker');
      }).filter(worker => {
        return !worker.name.includes('worker1');
      });
      this.options = this.toOptions(this.workers);
    });
  }

  close() {
    this.opened = false;
    this.openedChange.emit(this.opened);
  }

  onConfirm() {
    this.confirm.emit();
  }

  private toOptions(nodes: Node[]): any[] {
    const options = [];
    nodes.forEach(n => {
      options.push({'id': n.id, 'text': n.name, 'value': n.name});
    });
    return options;
  }
}
