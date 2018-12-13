import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {FormControl, FormGroup} from '@angular/forms';
import {NodeService} from '../node.service';

@Component({
  selector: 'app-node-create',
  templateUrl: './node-create.component.html',
  styleUrls: ['./node-create.component.css']
})
export class NodeCreateComponent implements OnInit {

  form = new FormGroup({
    name: new FormControl(''),
    ip: new FormControl(''),
    username: new FormControl(''),
    password: new FormControl(''),
    comment: new FormControl('')
  });
  staticBackdrop = true;
  closable = false;
  createNodeOpened: boolean;
  isSubmitGoing = false;
  @Output() create = new EventEmitter<boolean>();

  @Input() clusterId: string;

  constructor(private nodeService: NodeService) {
  }

  ngOnInit() {
  }

  onSubmit() {
    this.create.emit();
  }

  newNode() {
    this.createNodeOpened = true;
  }

  onCancel() {
    this.createNodeOpened = false;
  }
}
