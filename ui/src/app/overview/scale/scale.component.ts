import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {Cluster} from '../../cluster/cluster';

@Component({
  selector: 'app-scale',
  templateUrl: './scale.component.html',
  styleUrls: ['./scale.component.css']
})
export class ScaleComponent implements OnInit {

  opened = false;
  worker_size = 0;

  @Input() cluster: Cluster;
  @Output() openedChange = new EventEmitter();
  @Output() confirm = new EventEmitter();
  @ViewChild('form', {static: true}) form: NgForm;

  constructor() {
  }

  ngOnInit() {
  }

  close() {
    this.opened = false;
    this.openedChange.emit(this.opened);
  }

  onConfirm() {
    this.confirm.emit();
  }

}
