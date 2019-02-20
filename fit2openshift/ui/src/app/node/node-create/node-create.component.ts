import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {FormControl, FormGroup} from '@angular/forms';
import {NodeService} from '../node.service';
import {Cluster} from '../../cluster/cluster';
import {Node} from '../node';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';

@Component({
  selector: 'app-node-create',
  templateUrl: './node-create.component.html',
  styleUrls: ['./node-create.component.css']
})
export class NodeCreateComponent implements OnInit {

  @Input() currentCluster: Cluster;
  node: Node = new Node();

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


  constructor(private nodeService: NodeService, private tipService: TipService) {
  }

  ngOnInit() {
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.node.name = this.form.value.name;
    this.node.ip = this.form.value.ip;

    this.nodeService.createNode(this.currentCluster.name, this.node).subscribe(data => {
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.tipService.showTip('创建节点 ' + this.node.name + '成功', TipLevels.SUCCESS);
      this.createNodeOpened = false;
    });
  }

  newNode() {
    this.node = new Node();
    this.createNodeOpened = true;
  }


  onCancel() {
    this.createNodeOpened = false;
  }
}
