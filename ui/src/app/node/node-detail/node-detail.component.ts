import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Node} from '../node';

@Component({
  selector: 'app-node-detail',
  templateUrl: './node-detail.component.html',
  styleUrls: ['./node-detail.component.css']
})
export class NodeDetailComponent implements OnInit {

  node: Node = new Node();
  @Input() showInfoModal = false;
  @Output() showInfoModalChange = new EventEmitter();

  constructor() {
  }

  ngOnInit() {
  }

  cancel() {
    this.showInfoModal = false;
    this.showInfoModalChange.emit(this.showInfoModal);
  }

}
