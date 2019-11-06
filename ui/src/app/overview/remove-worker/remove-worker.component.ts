import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {NodeService} from '../../node/node.service';
import {Node} from '../../node/node';

@Component({
  selector: 'app-remove-worker',
  templateUrl: './remove-worker.component.html',
  styleUrls: ['./remove-worker.component.css']
})
export class RemoveWorkerComponent implements OnInit {

  constructor(private nodeService: NodeService) {
  }

  worker: string;
  opened = false;
  workers: Node[] = [];
  @Output() openedChange = new EventEmitter();
  @Output() confirm = new EventEmitter();
  @ViewChild('form', {static: true}) form: NgForm;

  ngOnInit() {
  }

  loadNodes(clusterName: string) {
    this.nodeService.listNodes(clusterName).subscribe(data => {
      this.workers = data.filter(worker => {
        return worker.roles.includes('worker');
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
